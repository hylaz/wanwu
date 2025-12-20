package response

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/UnicomAI/wanwu/internal/agent-service/model/request"
	"github.com/UnicomAI/wanwu/pkg/log"
	"github.com/cloudwego/eino/schema"
)

const (
	toolStartTitle        = `<tool>`
	toolStartTitleFormat  = `工具名：%s`
	toolParamsStartFormat = "\n\n```工具参数：\n"
	toolParamsEndFormat   = "\n```\n\n"
	toolEndFormat         = "\n\n```工具%s调用结果：\n %s \n```\n\n"
	toolEndTitle          = `</tool>`
	endLine               = "\n\n"
	agentSuccessCode      = 0
	agentFailCode         = 1
	finish                = 1
	notFinish             = 0
)

type AgentChatRespContext struct {
	HasTool            bool // 是否包含工具
	ToolStart          bool // 是否工具已开始
	ToolEnd            bool // 是否工具已结束
	ToolIndex          int  // 工具索引
	ToolCountMap       map[string]int
	ReplaceContent     strings.Builder // 替换内容，如果出现相同内则则进行替换
	ReplaceContentStr  string          // 替换内容，如果出现相同内则则进行替换
	ReplaceContentDone bool            //替换内容准备完成
}

func NewAgentChatRespContext() *AgentChatRespContext {
	return &AgentChatRespContext{
		ToolCountMap: make(map[string]int),
		ToolIndex:    -1,
	}
}

type AgentChatResp struct {
	Code           int             `json:"code"`
	Message        string          `json:"message"`
	Response       string          `json:"response"`
	GenFileUrlList []interface{}   `json:"gen_file_url_list"`
	History        []interface{}   `json:"history"`
	Finish         int             `json:"finish"`
	Usage          *AgentChatUsage `json:"usage"`
	SearchList     []interface{}   `json:"search_list"`
	QaType         int             `json:"qa_type"`
}

type AgentChatUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

func NewAgentChatRespWithTool(chatMessage *schema.Message, respContext *AgentChatRespContext, req *request.AgentChatContext) ([]string, error) {
	contentList := buildContentWithTool(chatMessage, respContext)
	var outputList = make([]string, 0)
	for _, content := range contentList {
		var agentChatResp = &AgentChatResp{
			Code:           agentSuccessCode,
			Message:        "success",
			Response:       content,
			GenFileUrlList: []interface{}{},
			History:        []interface{}{},
			QaType:         buildQaType(req),
			SearchList:     buildSearchList(req),
			Finish:         buildFinish(chatMessage),
			Usage:          buildUsage(chatMessage),
		}
		respString, err := buildRespString(agentChatResp)
		if err != nil {
			return nil, err
		}
		outputList = append(outputList, respString)
	}
	return outputList, nil
}

func AgentChatFailResp() string {
	var agentChatResp = &AgentChatResp{
		Code:     agentFailCode,
		Message:  "智能体处理异常，请稍后重试",
		Response: "智能体处理异常，请稍后重试",
		Finish:   finish,
	}
	respString, err := buildRespString(agentChatResp)
	if err != nil {
		log.Errorf("buildRespString error: %v", err)
		return ""
	}
	return respString
}

func buildRespString(agentChatResp *AgentChatResp) (string, error) {
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false) // 关键：禁用 HTML 转义

	if err := encoder.Encode(agentChatResp); err != nil {
		return "", err
	}
	return "data:" + buf.String(), nil
}

func buildContentWithTool(chatMessage *schema.Message, respContext *AgentChatRespContext) []string {
	if toolStart(chatMessage) {
		respContext.ToolStart = true
		respContext.HasTool = true
		var retList []string
		//if len(respContext.ToolCountMap) == 0 {
		//	retList = append(retList, toolStartTitle)
		//}

		for _, tool := range chatMessage.ToolCalls {
			if len(tool.Function.Arguments) > 0 {
				retList = append(retList, tool.Function.Arguments)
				continue
			}
			if tool.Type == "function" {
				if respContext.ToolIndex == -1 {
					respContext.ToolIndex = *tool.Index
				} else if *tool.Index != respContext.ToolIndex { //模型触发并发请求工具的bad case
					respContext.ToolIndex = *tool.Index
					retList = append(retList, toolParamsEndFormat)
				}
				if len(tool.Function.Name) > 0 {
					toolName := fmt.Sprintf(toolStartTitleFormat, tool.Function.Name)
					retList = append(retList, toolName)
				}

				if isNewTool(tool, respContext) {
					retList = append(retList, toolStartTitle)
					retList = append(retList, toolParamsStartFormat)
					respContext.ToolCountMap[tool.ID] = 1
				}
			}
		}

		return retList
	} else if toolParamsEnd(chatMessage) {
		return []string{toolParamsEndFormat}
	} else if toolEnd(chatMessage) {
		respContext.ToolEnd = true
		toolResult := fmt.Sprintf(toolEndFormat, chatMessage.ToolName, chatMessage.Content)
		respContext.ToolCountMap[chatMessage.ToolCallID] = 0
		var allStop = true
		for _, tool := range respContext.ToolCountMap {
			if tool > 0 {
				allStop = false
				break
			}
		}
		if !allStop {
			return []string{toolResult}
		}
		return []string{toolResult, toolEndTitle}
	} else {
		//在工具期间，不输出任何content内容
		if respContext.ToolStart && !respContext.ToolEnd {
			return []string{}
		}
		//替换内容准备(工具未开始，但是输出了内容)
		if !respContext.ToolStart {
			if utf8.RuneCountInString(chatMessage.Content) > 10 {
				var replaceContent = respContext.ReplaceContentStr
				if len(replaceContent) == 0 {
					replaceContent = respContext.ReplaceContent.String()
				}
				if replaceContent == chatMessage.Content {
					respContext.ReplaceContentDone = true
					respContext.ReplaceContentStr = replaceContent
					return []string{}
				}
			}
			if !respContext.ReplaceContentDone {
				respContext.ReplaceContent.WriteString(chatMessage.Content)
			}

			//if !respContext.ReplaceContentDone {
			//	respContext.ReplaceContent.WriteString(chatMessage.Content)
			//	if strings.HasSuffix(chatMessage.Content, endLine) {
			//		respContext.ReplaceContentDone = true
			//		respContext.ReplaceContentStr = respContext.ReplaceContent.String()
			//	}
			//} else {
			//	if respContext.ReplaceContentStr == chatMessage.Content {
			//		return []string{}
			//	}
			//}
		}
		return []string{chatMessage.Content}
	}
}

func isNewTool(tool schema.ToolCall, respContext *AgentChatRespContext) bool {
	return len(tool.ID) > 0 && respContext.ToolCountMap[tool.ID] == 0
}

func toolStart(chatMessage *schema.Message) bool {
	//responseMeta := chatMessage.ResponseMeta
	//if responseMeta == nil {
	//	return false
	//}
	//return responseMeta.FinishReason == "tool_calls"
	return len(chatMessage.ToolCalls) > 0
}

func toolParamsEnd(chatMessage *schema.Message) bool {
	responseMeta := chatMessage.ResponseMeta
	if responseMeta == nil {
		return false
	}
	return responseMeta.FinishReason == "tool_calls"
}

func toolEnd(chatMessage *schema.Message) bool {
	return chatMessage.Role == schema.Tool
}

func buildFinish(chatMessage *schema.Message) int {
	if chatMessage.ResponseMeta != nil && chatMessage.ResponseMeta.FinishReason == "stop" {
		return finish
	}
	return notFinish
}

func buildUsage(chatMessage *schema.Message) *AgentChatUsage {
	if chatMessage.ResponseMeta != nil && chatMessage.ResponseMeta.Usage != nil {
		usage := chatMessage.ResponseMeta.Usage
		return &AgentChatUsage{
			PromptTokens:     usage.PromptTokens,
			CompletionTokens: usage.CompletionTokens,
			TotalTokens:      usage.TotalTokens,
		}
	}
	return &AgentChatUsage{}
}

func buildSearchList(req *request.AgentChatContext) []interface{} {
	if req.KnowledgeHitData == nil {
		return []interface{}{}
	}
	list := req.KnowledgeHitData.SearchList
	var retList = make([]interface{}, 0)
	if len(list) > 0 {
		for _, item := range list {
			retList = append(retList, item)
		}
	}
	return retList
}

func buildQaType(req *request.AgentChatContext) int {
	if req.KnowledgeHitData == nil {
		return 0
	}
	return 1
}
