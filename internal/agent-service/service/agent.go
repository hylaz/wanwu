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

// createChatModel 创建chatmodel
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

// createMessageBuilder 创建消息
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
