package chat_model

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"time"

	agent_http_client "github.com/UnicomAI/wanwu/internal/agent-service/pkg/http"
	http_client "github.com/UnicomAI/wanwu/pkg/http-client"
	"github.com/UnicomAI/wanwu/pkg/log"
	"github.com/cloudwego/eino/callbacks"
	"github.com/cloudwego/eino/components"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
	"github.com/meguminnnnnnnnn/go-openai"
)

const (
	Text     ContentType = "text"
	ImageUrl ContentType = "image_url"
)

var (
	headerData  = regexp.MustCompile(`^data:\s*`)
	errorPrefix = regexp.MustCompile(`^data:\s*{"error":`)
)

type YuanjingVLModel struct {
	// Model specifies the ID of the model to use
	// Required
	Model string `json:"model"`
	// APIKey is your authentication key
	// Use OpenAI API key or Azure API key depending on the service
	// Required
	APIKey string `json:"api_key"`
	// BaseURL is the Azure OpenAI endpoint URL
	// Format: https://{YOUR_RESOURCE_NAME}.openai.azure.com. YOUR_RESOURCE_NAME is the name of your resource that you have created on Azure.
	// Required for Azure
	BaseURL string `json:"base_url"`
}

func CreateYuanjingVLModel(model, apiKey, baseURL string) *YuanjingVLModel {
	return &YuanjingVLModel{
		Model:   model,
		APIKey:  apiKey,
		BaseURL: baseURL,
	}
}

type ContentType string

type YuanjingVLParams struct {
	Model     string                 `json:"model"`
	Messages  []*YuanjingVLMessage   `json:"messages"`
	Stream    bool                   `json:"stream"`
	ExtraBody map[string]interface{} `json:"extra_body"`
}

type YuanjingVLMessage struct {
	Role    string         `json:"role"`
	Content []*ContentInfo `json:"content"`
}

type ContentInfo struct {
	Type     ContentType   `json:"type"`
	Text     string        `json:"text,omitempty"`
	ImageUrl *ContentImage `json:"image_url,omitempty"`
}

type ContentImage struct {
	Url string `json:"url"`
}

func (cm *YuanjingVLModel) Generate(ctx context.Context, in []*schema.Message, opts ...model.Option) (
	outMsg *schema.Message, err error) {
	return nil, errors.New("generate not implement")
}

func (cm *YuanjingVLModel) Stream(ctx context.Context, in []*schema.Message, opts ...model.Option) (outStream *schema.StreamReader[*schema.Message], err error) {
	ctx = callbacks.EnsureRunInfo(ctx, "yuanjingVL", components.ComponentOfChatModel)
	out, err := streamChat(ctx, in, cm, opts...)
	if err != nil {
		return nil, errors.New("yuanjing vl stream error")
	}
	return out, nil
}

func (cm *YuanjingVLModel) WithTools(tools []*schema.ToolInfo) (model.ToolCallingChatModel, error) {
	return &YuanjingVLModel{}, nil
}

func streamChat(ctx context.Context, in []*schema.Message, modelInfo *YuanjingVLModel, opts ...model.Option) (outStream *schema.StreamReader[*schema.Message], err error) {
	defer func() {
		if err != nil {
			callbacks.OnError(ctx, err)
		}
	}()

	vlMessage := YuanjingVLMessage{
		Role: "user",
	}
	var contentInfoList []*ContentInfo
	for _, message := range in {
		if len(message.Content) > 0 {
			contentInfoList = append(contentInfoList, &ContentInfo{
				Type: Text,
				Text: message.Content,
			})
		}
		if len(message.UserInputMultiContent) > 0 {
			for _, ImageInfo := range message.UserInputMultiContent {
				contentInfoList = append(contentInfoList, &ContentInfo{
					Type:     ImageUrl,
					ImageUrl: &ContentImage{Url: *ImageInfo.Image.Base64Data},
				})
			}
		}
	}

	vlMessage.Content = contentInfoList

	cbInput := &model.CallbackInput{
		Messages:   in,
		Tools:      nil,
		ToolChoice: nil,
		Config:     &model.Config{},
	}

	ctx = callbacks.OnStart(ctx, cbInput)

	// 执行请求
	resp, err := requestYuanjingVLStreamChat(ctx, &vlMessage, modelInfo)
	if err != nil {
		log.Errorf("requestYuanjingVLStreamChat error: 请求异常: %v", err)
		return nil, err
	}
	sr, sw := schema.Pipe[*model.CallbackOutput](1)

	go func() {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				log.Errorf("error: 响应体关闭异常: %v", err)
			}
		}(resp.Body) // 确保响应体关闭
		defer func() {
			panicErr := recover()

			if panicErr != nil {
				_ = sw.Send(nil, errors.New("panic error"))
			}

			sw.Close()
		}()

		reader := bufio.NewReader(resp.Body)

		for {
			select {
			case <-ctx.Done():
				// 客户端断开连接
				log.Errorf("writeSSE: ctx canceled")
				return
			default:
				bytes, err := reader.ReadBytes('\n')

				if len(bytes) > 0 {
					lines, err := processLines(bytes)
					if err != nil {
						break
					}
					if len(lines) == 0 {
						continue
					}
					log.Infof("yuanjing vl stream result, %s", string(lines))

					responseData := openai.ChatCompletionStreamResponse{}
					err = json.Unmarshal(lines, &responseData)
					if err != nil {
						log.Errorf("error: 解析响应数据异常: %v", err)
						continue
					}
					msg, found, buildErr := build(responseData)
					if buildErr != nil {
						_ = sw.Send(nil, fmt.Errorf("failed to build message from stream chunk: %w", buildErr))
						return
					}
					if !found {
						continue
					}

					closed := sw.Send(&model.CallbackOutput{
						Message:    msg,
						Config:     cbInput.Config,
						TokenUsage: toModelCallbackUsage(msg.ResponseMeta),
					}, nil)

					if closed {
						return
					}
				}

				if err != nil {
					if err == io.EOF {
						return // 正常结束
					}
					if err != nil {
						_ = sw.Send(nil, fmt.Errorf("failed to receive stream chunk from OpenAI: %w", err))
						return
					}
				}
			}
		}

	}()

	ctx, nsr := callbacks.OnEndWithStreamOutput(ctx, schema.StreamReaderWithConvert(sr,
		func(src *model.CallbackOutput) (callbacks.CallbackOutput, error) {
			return src, nil
		}))

	outStream = schema.StreamReaderWithConvert(nsr,
		func(src callbacks.CallbackOutput) (*schema.Message, error) {
			s := src.(*model.CallbackOutput)
			if s.Message == nil {
				return nil, schema.ErrNoValue
			}

			return s.Message, nil
		},
	)

	return outStream, nil
}

func processLines(rawLine []byte) ([]byte, error) {
	var (
		hasErrorPrefix bool
	)

	noSpaceLine := bytes.TrimSpace(rawLine)
	if errorPrefix.Match(noSpaceLine) {
		hasErrorPrefix = true
	}
	if !headerData.Match(noSpaceLine) || hasErrorPrefix {
		if hasErrorPrefix {
			noSpaceLine = headerData.ReplaceAll(noSpaceLine, nil)
		}
		return noSpaceLine, nil
	}

	noPrefixLine := headerData.ReplaceAll(noSpaceLine, nil)
	if string(noPrefixLine) == "[DONE]" {
		return nil, io.EOF
	}

	return noPrefixLine, nil
}

func toModelCallbackUsage(respMeta *schema.ResponseMeta) *model.TokenUsage {
	if respMeta == nil {
		return nil
	}
	usage := respMeta.Usage
	if usage == nil {
		return nil
	}
	return &model.TokenUsage{
		PromptTokens: usage.PromptTokens,
		PromptTokenDetails: model.PromptTokenDetails{
			CachedTokens: usage.PromptTokenDetails.CachedTokens,
		},
		CompletionTokens: usage.CompletionTokens,
		TotalTokens:      usage.TotalTokens,
	}
}

func build(resp openai.ChatCompletionStreamResponse) (msg *schema.Message, found bool, err error) {
	for _, choice := range resp.Choices {
		// take 0 index as response, rewrite if needed
		if choice.Index != 0 {
			continue
		}

		found = true

		msg = &schema.Message{
			Role:    toMessageRole(choice.Delta.Role),
			Content: choice.Delta.Content,
			ResponseMeta: &schema.ResponseMeta{
				FinishReason: string(choice.FinishReason),
				Usage:        toEinoTokenUsage(resp.Usage),
			},
		}

		break
	}

	if resp.Usage != nil && !found {
		msg = &schema.Message{
			ResponseMeta: &schema.ResponseMeta{
				Usage: toEinoTokenUsage(resp.Usage),
			},
		}
		found = true
	}

	return msg, found, nil
}

func toMessageRole(role string) schema.RoleType {
	switch role {
	case openai.ChatMessageRoleUser:
		return schema.User
	case openai.ChatMessageRoleAssistant:
		return schema.Assistant
	case openai.ChatMessageRoleSystem:
		return schema.System
	case openai.ChatMessageRoleTool:
		return schema.Tool
	case "":
		// When the role field is an empty string, populate it with the schema.Assistant.
		return schema.Assistant
	default:
		return schema.RoleType(role)
	}
}

func toEinoTokenUsage(usage *openai.Usage) *schema.TokenUsage {
	if usage == nil {
		return nil
	}

	promptTokenDetails := schema.PromptTokenDetails{}
	if usage.PromptTokensDetails != nil {
		promptTokenDetails.CachedTokens = usage.PromptTokensDetails.CachedTokens
	}

	return &schema.TokenUsage{
		PromptTokens:       usage.PromptTokens,
		PromptTokenDetails: promptTokenDetails,
		CompletionTokens:   usage.CompletionTokens,
		TotalTokens:        usage.TotalTokens,
	}
}

func requestYuanjingVLStreamChat(ctx context.Context, req *YuanjingVLMessage, modelInfo *YuanjingVLModel) (*http.Response, error) {
	params, err := buildYuanjingVLChatHttpParams(req, modelInfo)
	if err != nil {
		log.Errorf("build http params fail %s", err.Error())
		return nil, err
	}
	// 捕获 panic 并记录日志（不重新抛出，避免崩溃）
	defer func() {
		if r := recover(); r != nil {
			log.Errorf("RagStreamChat panic: %v", r)
		}
	}()

	resp, err := agent_http_client.GetClient().PostJsonOriResp(ctx, params)
	if err != nil {
		errMsg := fmt.Sprintf("error: 调用下游服务异常: %v", err)
		log.Errorf(errMsg)
		return nil, err
	}
	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error: 调用下游服务异常: %s", resp.Status)
	}
	return resp, nil
}

func buildYuanjingVLChatHttpParams(req *YuanjingVLMessage, modelInfo *YuanjingVLModel) (*http_client.HttpRequestParams, error) {
	params := &YuanjingVLParams{
		Model:     modelInfo.Model,
		Messages:  []*YuanjingVLMessage{req},
		Stream:    true,
		ExtraBody: map[string]interface{}{"api_option": "general"},
	}
	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	return &http_client.HttpRequestParams{
		Url:        modelInfo.BaseURL + "/chat/completions",
		Body:       body,
		Headers:    map[string]string{"Authorization": "Bearer " + modelInfo.APIKey},
		Timeout:    time.Minute * 10,
		MonitorKey: "yuanjing_vl_chat_service",
		LogLevel:   http_client.LogAll,
	}, nil
}
