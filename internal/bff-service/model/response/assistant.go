package response

import (
	"github.com/UnicomAI/wanwu/internal/bff-service/model/request"
)

type Assistant struct {
	request.AppBriefConfig                                // 基本信息
	AssistantId            string                         `json:"assistantId"  validate:"required"`
	Prologue               string                         `json:"prologue"`            // 开场白
	Instructions           string                         `json:"instructions"`        // 系统提示词
	RecommendQuestion      []string                       `json:"recommendQuestion"`   // 推荐问题
	ModelConfig            request.AppModelConfig         `json:"modelConfig"`         // 模型
	KnowledgeBaseConfig    request.AppKnowledgebaseConfig `json:"knowledgeBaseConfig"` // 知识库
	RerankConfig           request.AppModelConfig         `json:"rerankConfig"`        // Rerank模型
	OnlineSearchConfig     request.OnlineSearchConfig     `json:"onlineSearchConfig"`  // 在线搜索配置
	SafetyConfig           request.AppSafetyConfig        `json:"safetyConfig"`        // 敏感词表配置
	VisionConfig           VisionConfig                   `json:"visionConfig"`        // 视觉配置
	Scope                  int32                          `json:"scope"`               // 作用域
	WorkFlowInfos          []*WorkFlowInfos               `json:"workFlowInfos"`       // 工作流信息
	MCPInfos               []*MCPInfos                    `json:"mcpInfos"`            // MCP信息
	CustomInfos            []*CustomInfos                 `json:"customInfos"`         // 自定义工具信息
	CreatedAt              string                         `json:"createdAt"`           // 创建时间
	UpdatedAt              string                         `json:"updatedAt"`           // 更新时间
}

type WorkFlowInfos struct {
	UniqueId     string `json:"uniqueId"` // 随机unique id(每次动态生成)
	WorkFlowId   string `json:"workFlowId"`
	ApiName      string `json:"apiName"`
	Enable       bool   `json:"enable"`
	WorkFlowName string `json:"name"`
	WorkFlowDesc string `json:"workFlowDesc"`
}

type MCPInfos struct {
	UniqueId      string `json:"uniqueId"` // 随机unique id(每次动态生成)
	MCPId         string `json:"mcpId"`
	MCPSquareId   string `json:"mcpSquareId"`
	Enable        bool   `json:"enable"`
	MCPName       string `json:"name"`
	MCPDesc       string `json:"mcpDesc"`
	MCPServerFrom string `json:"mcpServerFrom"`
	MCPServerUrl  string `json:"mcpServerUrl"`
	Valid         bool   `json:"valid"`
}

type CustomInfos struct {
	UniqueId   string `json:"uniqueId"` // 随机unique id(每次动态生成)
	CustomId   string `json:"customId"`
	Enable     bool   `json:"enable"`
	CustomName string `json:"name"`
	CustomDesc string `json:"customDesc"`
	Valid      bool   `json:"valid"`
}

type ConversationInfo struct {
	ConversationId string `json:"conversationId"`
	AssistantId    string `json:"assistantId"`
	Title          string `json:"title"`
	CreatedAt      string `json:"createdAt"`
}

type ConversationDetailInfo struct {
	Id             string        `json:"id"`
	AssistantId    string        `json:"assistantId"`
	ConversationId string        `json:"conversationId"`
	Prompt         string        `json:"prompt"`
	SysPrompt      string        `json:"sysPrompt"`
	Response       string        `json:"response"`
	SearchList     interface{}   `json:"searchList"`
	QaType         int32         `json:"qa_type"`
	CreatedBy      string        `json:"createdBy"`
	CreatedAt      int64         `json:"createdAt"`
	UpdatedAt      int64         `json:"updatedAt"`
	RequestFiles   []RequestFile `json:"requestFiles"`
	FileSize       int64         `json:"fileSize"`
	FileName       string        `json:"fileName"`
}
type RequestFile struct {
	FileName string `json:"name"`
	FileSize int64  `json:"size"`
	FileUrl  string `json:"fileUrl"`
}
type ConversationCreateResp struct {
	ConversationId string `json:"conversationId"`
}

type AssistantCreateResp struct {
	AssistantId string `json:"assistantId"`
}

type AssistantTemplateInfo struct {
	AssistantTemplateId string `json:"assistantTemplateId"` // 智能体模板Id
	AppType             string `json:"appType"`             // 应用类型(固定值: agentTemplate)
	Category            string `json:"category"`            // 种类(gov:政务,industry:工业,edu:文教,medical:医疗)
	request.AppBriefConfig
	Prologue                  string   `json:"prologue"`            // 开场白
	Instructions              string   `json:"instructions"`        // 系统提示词
	RecommendQuestion         []string `json:"recommendQuestion"`   // 推荐问题
	Summary                   string   `json:"summary"`             // 使用概述
	Feature                   string   `json:"feature"`             // 特性说明
	Scenario                  string   `json:"scenario"`            // 应用场景
	WorkFlowConfigInstruction string   `json:"workFlowInstruction"` // 工作流配置说明
}
