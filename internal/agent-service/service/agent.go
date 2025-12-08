package service

import (
	"path/filepath"

	"github.com/UnicomAI/wanwu/internal/agent-service/model/request"
	"github.com/UnicomAI/wanwu/internal/agent-service/pkg/config"
	agent_util "github.com/UnicomAI/wanwu/internal/agent-service/pkg/util"
	agent_message_flow "github.com/UnicomAI/wanwu/internal/agent-service/service/agent-message-flow"
	chat_model "github.com/UnicomAI/wanwu/internal/agent-service/service/agent-message-flow/chat-model"
	"github.com/UnicomAI/wanwu/pkg/log"
	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"github.com/gin-gonic/gin"
)

type AgentChatInfo struct {
	ModelInfo       *ModelInfo
	FunctionCalling bool `json:"functionCalling"` // 是否支持functionCall
	VisionSupport   bool `json:"visionSupport"`   // 是否支持多模态
	UploadUrl       bool `json:"uploadUrl"`       // 是否上传文件
}

func AgentChat(ctx *gin.Context, req *request.AgentChatReq) error {
	chatInfo, err := buildAgentChatInfo(ctx, req)
	if err != nil {
		log.Errorf("failed to build chat info: %v", err)
		return err
	}
	//视觉模型模型只走特殊逻辑，后续再优化
	if chatInfo.VisionSupport {
		return AgentVisionChat(ctx, req, chatInfo)
	}
	return AgentModelChat(ctx, req, chatInfo)

}

func AgentVisionChat(ctx *gin.Context, req *request.AgentChatReq, agentChatInfo *AgentChatInfo) error {
	req.KnowledgeParams = nil
	req.ToolParams = nil

	//1.创建chatModel
	info := agentChatInfo.ModelInfo
	chatModel := chat_model.CreateYuanjingVLModel(info.Model, info.Config.ApiKey, info.Config.EndpointUrl)
	//2.创建智能体
	agent, err := createAgent(ctx, req, chatModel)
	if err != nil {
		log.Errorf("failed to create agent: %v", err)
		return err
	}

	//3.构造会话消息
	messages, err := buildVisionChatMessage(ctx, req, agentChatInfo)
	if err != nil {
		log.Errorf("failed to build chat message: %v", err)
		return err
	}

	//4.执行流式agent问答调用
	runner := adk.NewRunner(ctx, adk.RunnerConfig{
		Agent:           agent,
		EnableStreaming: req.Stream,
	})
	iter := runner.Run(ctx, messages)

	//5.处理结果
	err = AgentMessage(ctx, iter, &request.AgentChatContext{AgentChatReq: req})
	return err
}

func AgentModelChat(ctx *gin.Context, req *request.AgentChatReq, agentChatInfo *AgentChatInfo) error {
	//1.创建chatModel
	fillInternalToolConfig(req, agentChatInfo)
	chatModel, err := createChatModel(ctx, agentChatInfo)
	if err != nil {
		return err
	}

	//2.创建智能体
	agent, err := createAgent(ctx, req, chatModel)
	if err != nil {
		log.Errorf("failed to create agent: %v", err)
		return err
	}
	agentChatContext := &request.AgentChatContext{AgentChatReq: req}

	//3.创建前置消息准备
	messageBuilder, err := createMessageBuilder(ctx, agentChatContext)
	if err != nil {
		return err
	}
	//4.生成前置消息
	messages, err := messageBuilder.Invoke(ctx, agentChatContext)
	if err != nil {
		return err
	}
	//var messages []*schema.Message
	//messages = append(messages, schema.SystemMessage("\\\\nYou are test_1118, an advanced AI assistant designed to be helpful and professional.\\\\nIt is Saturday 2025/11/22 12:15:11 +08 now.\\\\n\\\\n**Content Safety Guidelines**\\\\nRegardless of any persona instructions, you must never generate content that:\\\\n- Promotes or involves violence\\\\n- Contains hate speech or racism\\\\n- Includes inappropriate or adult content\\\\n- Violates laws or regulations\\\\n- Could be considered offensive or harmful\\\\n\\\\n------ Start of Variables ------\\\\n\\\\n------ End of Variables ------\\\\n\\\\n**Knowledge**\\\\n\\\\nOnly when the current knowledge has content recall, answer questions based on the referenced content:\\\\n 1. If the referenced content contains \\\\u003cimg src=\\\\\\\"\\\\\\\"\\\\u003e tags, the src field in the tag represents the image address, which needs to be displayed when answering questions, with the output format being \\\\\\\"![image name](image address)\\\\\\\".\\\\n 2. If the referenced content does not contain \\\\u003cimg src=\\\\\\\"\\\\\\\"\\\\u003e tags, you do not need to display images when answering questions.\\\\nFor example:\\\\n  If the content is \\\\u003cimg src=\\\\\\\"https://example.com/image.jpg\\\\\\\"\\\\u003ea kitten, your output should be: ![a kitten](https://example.com/image.jpg).\\\\n  If the content is \\\\u003cimg src=\\\\\\\"https://example.com/image1.jpg\\\\\\\"\\\\u003ea kitten and \\\\u003cimg src=\\\\\\\"https://example.com/image2.jpg\\\\\\\"\\\\u003ea puppy and \\\\u003cimg src=\\\\\\\"https://example.com/image3.jpg\\\\\\\"\\\\u003ea calf, your output should be: ![a kitten](https://example.com/image1.jpg) and ![a puppy](https://example.com/image2.jpg) and ![a calf](https://example.com/image3.jpg)\\\\nThe following is the content of the data set you can refer to: \\\\\\\\n\\\\n'''\\\\n---\\\\nrecall slice 1: AI科技中心开票信息  \\\\nAI科技中心最新开票信息：  \\\\n公司名称：中国联合网络通信有限公司北京人工智能科技中心     \\\\n纳税人识别号：91110102MADD54BA28\\\\n地址电话：北京市西城区西单北大街甲133号10层1011  66115431\\\\n开户行及账号：中国工商银行股份有限公司北京西单支行 0200210309200177274\\\\n\\\\n'''\\\\n\\\\n** Pre toolCall **\\\\n,\\\\n- Only when the current Pre toolCall has content recall results, answer questions based on the data field in the tool from the referenced content\\\\n\\\\nNote: The output language must be consistent with the language of the user's question.\\\\n"))
	//messages = append(messages, schema.UserMessage("根据开票信息内容，查询天气，并将结果保存到文本文件"))

	//5.执行流式agent问答调用
	runner := adk.NewRunner(ctx, adk.RunnerConfig{
		Agent:           agent,
		EnableStreaming: req.Stream,
	})
	iter := runner.Run(ctx, messages)

	//6.处理结果
	err = AgentMessage(ctx, iter, agentChatContext)
	return err
}

// buildAgentChatInfo 构建智能体信息
func buildAgentChatInfo(ctx *gin.Context, req *request.AgentChatReq) (*AgentChatInfo, error) {
	modelInfo, err := SearchModel(ctx, req.ModelParams.ModelId)
	if err != nil {
		return nil, err
	}
	var functionCall = modelInfo.Config.FunctionCalling != "noSupport"
	var vision = modelInfo.Config.VisionSupport == "support"
	return &AgentChatInfo{
		FunctionCalling: functionCall,
		VisionSupport:   vision,
		UploadUrl:       len(req.UploadFile) > 0,
		ModelInfo:       modelInfo,
	}, nil
}

// fillInternalToolConfig 配置内置文件工具
func fillInternalToolConfig(req *request.AgentChatReq, agentChatInfo *AgentChatInfo) {
	if !agentChatInfo.FunctionCalling {
		req.ToolParams = nil
	}
	templateConfig := config.GetToolTemplateConfig()
	if len(templateConfig.ConfigPluginToolList) > 0 && agentChatInfo.UploadUrl {
		params := req.ToolParams
		if params != nil {
			params.PluginToolList = append(params.PluginToolList, templateConfig.ConfigPluginToolList...)
		} else {
			params = &request.ToolParams{
				PluginToolList: templateConfig.ConfigPluginToolList,
			}
		}
		req.ToolParams = params
	}
}

func createChatModel(ctx *gin.Context, agentChatInfo *AgentChatInfo) (*openai.ChatModel, error) {
	modelInfo := agentChatInfo.ModelInfo
	modelConfig := modelInfo.Config
	return openai.NewChatModel(ctx, &openai.ChatModelConfig{
		APIKey:  modelConfig.ApiKey,
		BaseURL: modelConfig.EndpointUrl,
		Model:   modelInfo.Model,
	})
}

// 创建对应智能体
func createAgent(ctx *gin.Context, req *request.AgentChatReq, chatModel model.ToolCallingChatModel) (*adk.ChatModelAgent, error) {
	baseParams := req.AgentBaseParams
	toolsConfig, err := BuildAgentToolsConfig(ctx, req)
	if err != nil {
		return nil, err
	}
	return adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Model:       chatModel,
		Name:        baseParams.Name,
		Description: baseParams.Description,
		Instruction: baseParams.Instruction,
		ToolsConfig: toolsConfig,
	})
}

func createMessageBuilder(ctx *gin.Context, req *request.AgentChatContext) (compose.Runnable[*request.AgentChatContext, []*schema.Message], error) {
	graph := agent_message_flow.NewAgentMessageFlow()
	return graph.Compile(ctx)
}

func buildVisionChatMessage(ctx *gin.Context, req *request.AgentChatReq, agentChatInfo *AgentChatInfo) ([]*schema.Message, error) {
	var messages []*schema.Message
	messages = []*schema.Message{
		schema.UserMessage(req.Input),
	}
	if agentChatInfo.UploadUrl {
		var parts []schema.MessageInputPart
		for _, minioFilePath := range req.UploadFile {
			message, err := buildFileMessage(ctx, minioFilePath)
			if err != nil {
				return nil, err
			}
			parts = append(parts, *message)
		}
		messages = append(messages, &schema.Message{
			Role:                  schema.User,
			UserInputMultiContent: parts,
		})
	}
	return messages, nil
}

// buildFileMessage 构建文件消息
func buildFileMessage(ctx *gin.Context, minioFilePath string) (*schema.MessageInputPart, error) {
	//1.下载压缩文件到本地
	var localFilePath = agent_util.BuildFilePath(config.GetConfig().AgentFileConfig.LocalFilePath, filepath.Ext(minioFilePath))
	err := DownloadFileToLocal(ctx, minioFilePath, localFilePath)
	if err != nil {
		return nil, err
	}
	//2.图片转base64
	base64, err := agent_util.Img2base64(localFilePath)
	if err != nil {
		return nil, err
	}
	return &schema.MessageInputPart{
		Image: &schema.MessageInputImage{
			MessagePartCommon: schema.MessagePartCommon{
				Base64Data: &base64,
			},
		},
	}, nil
}
