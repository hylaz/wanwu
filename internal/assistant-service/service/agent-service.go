package service

import (
	"github.com/UnicomAI/wanwu/internal/assistant-service/client/model"
	"github.com/UnicomAI/wanwu/internal/assistant-service/config"
)

const (
	maxHistory = 5
)

type AgentChatReq struct {
	Input      string   `json:"input"`
	Stream     bool     `json:"stream"`
	UploadFile []string `json:"uploadFile"`

	AgentBaseParams *AgentBaseParams `json:"agentBaseParams"` // 智能体基础参数
	ModelParams     *ModelParams     `json:"modelParams"`     // 模型参数
	KnowledgeParams *KnowledgeParams `json:"knowledgeParams"` // 知识库参数，如果后续需要增加透传，理论上只需要修改此KnowledgeParams即可
	ToolParams      *ToolParams      `json:"toolParams"`      // 工具相关参数，mcp tool plugin tool
}

type AgentBaseParams struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Instruction string `json:"instruction"`
}

type ModelParams struct {
	ModelId          string                              `json:"modelId"`
	History          []config.AssistantConversionHistory `json:"history,omitempty"`
	MaxHistory       int                                 `json:"maxHistory"`
	Temperature      *float32                            `json:"temperature,omitempty"`      //温度
	TopP             *float32                            `json:"topP,omitempty"`             //topP
	FrequencyPenalty *float32                            `json:"frequencyPenalty,omitempty"` //频率惩罚
	PresencePenalty  *float32                            `json:"presence_penalty,omitempty"` //存在惩罚
	MaxTokens        *int                                `json:"max_tokens,omitempty"`       //模型输出最大token数，这个字段暂时不设置，因为模型可能触发接口调用不确定是否会超，先不传
}

type KnowledgeParams struct {
	UserId               string                        `json:"userId"`          // 用户id
	KnowledgeIdList      []string                      `json:"knowledgeIdList"` // 知识库id列表
	Question             string                        `json:"question"`
	Threshold            float32                       `json:"threshold"` // Score阈值
	TopK                 int32                         `json:"topK"`
	Stream               bool                          `json:"stream"`
	Chichat              bool                          `json:"chichat"` // 当知识库召回结果为空时是否使用默认话术（兜底），默认为true
	RerankModelId        string                        `json:"rerank_model_id"`
	CustomModelInfo      *CustomModelInfo              `json:"custom_model_info"`
	MaxHistory           int32                         `json:"max_history"`
	RewriteQuery         bool                          `json:"rewrite_query"`   // 是否query改写
	RerankMod            string                        `json:"rerank_mod"`      // rerank_model:重排序模式，weighted_score：权重搜索
	RetrieveMethod       string                        `json:"retrieve_method"` // hybrid_search:混合搜索， semantic_search:向量搜索， full_text_search：文本搜索
	Weight               *config.WeightParams          `json:"weights"`         // 权重搜索下的权重配置
	Temperature          float32                       `json:"temperature,omitempty"`
	TopP                 float32                       `json:"top_p,omitempty"`               // 多样性
	RepetitionPenalty    float32                       `json:"repetition_penalty,omitempty"`  // 重复惩罚/频率惩罚
	ReturnMeta           bool                          `json:"return_meta,omitempty"`         // 是否返回元数据
	AutoCitation         bool                          `json:"auto_citation"`                 // 是否自动角标
	TermWeight           float32                       `json:"term_weight_coefficient"`       // 关键词系数
	MetaFilter           bool                          `json:"metadata_filtering"`            // 元数据过滤开关
	MetaFilterConditions []*config.MetadataFilterParam `json:"metadata_filtering_conditions"` // 元数据过滤条件
	UseGraph             bool                          `json:"use_graph"`                     // 是否启动知识图谱查询
}
type CustomModelInfo struct {
	LlmModelID string `json:"llm_model_id"`
}

type ToolParams struct {
	PluginToolList []config.PluginListAlgRequest `json:"pluginTool,omitempty"`
	McpToolList    []*MCPToolInfo                `json:"mcpToolList,omitempty"`
}

type APIAuth struct {
	Type  string `json:"type"`
	In    string `json:"in"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

type MCPToolInfo struct {
	URL          string   `json:"url"`
	Transport    string   `json:"transport"`
	ToolNameList []string `json:"toolNameList"` // MCP工具方法列表,会根据此方法名的列表进行mcp方法的过滤，如果此列为空，则标识不进行过滤
}

func BuildAgentChatReq(sseRequest *config.AgentSSERequest, assistant *model.Assistant) *AgentChatReq {
	var req = &AgentChatReq{
		Input:           sseRequest.Input,
		Stream:          true,
		UploadFile:      sseRequest.UploadFileUrl,
		AgentBaseParams: buildAgentBaseParams(assistant),
		ModelParams:     buildAgentModelParams(sseRequest),
		KnowledgeParams: buildKnowledgeParams(sseRequest, assistant),
		ToolParams:      buildToolParams(sseRequest),
	}
	return req
}

func buildAgentBaseParams(assistant *model.Assistant) *AgentBaseParams {
	return &AgentBaseParams{
		Name:        assistant.Name,
		Description: assistant.Desc,
		Instruction: assistant.Instructions,
	}
}

func buildAgentModelParams(sseRequest *config.AgentSSERequest) *ModelParams {
	modelParamsReq := sseRequest.ModelParams
	modelParams := &ModelParams{
		ModelId:    sseRequest.ModelId,
		History:    sseRequest.History,
		MaxHistory: maxHistory,
	}
	return buildModelParams(modelParamsReq, modelParams)
}

func buildModelParams(params map[string]interface{}, modelParams *ModelParams) *ModelParams {
	if len(params) == 0 {
		return modelParams
	}
	modelParams.Temperature = toFloat(params["temperature"])
	modelParams.TopP = toFloat(params["top_p"])
	modelParams.FrequencyPenalty = toFloat(params["frequency_penalty"])
	modelParams.PresencePenalty = toFloat(params["presence_penalty"])
	return modelParams
}

func toFloat(data interface{}) *float32 {
	if data == nil {
		return nil
	}
	f, ok := data.(float32)
	if !ok {
		return nil
	}
	return &f
}

func buildKnowledgeParams(sseRequest *config.AgentSSERequest, assistant *model.Assistant) *KnowledgeParams {
	if !sseRequest.UseKnow || sseRequest.KnParams == nil {
		return nil
	}
	params := sseRequest.KnParams
	return &KnowledgeParams{
		UserId:               assistant.UserId,
		KnowledgeIdList:      params.KnowledgeIdList,
		Threshold:            params.Threshold,
		TopK:                 params.TopK,
		Stream:               sseRequest.Stream,
		RerankModelId:        toString(params.RerankId),
		MaxHistory:           params.MaxHistory,
		RewriteQuery:         params.RewriteQuery,
		RerankMod:            params.RerankMod,
		RetrieveMethod:       params.RetrieveMethod,
		Weight:               params.Weights,
		TermWeight:           params.TermWeight,
		MetaFilter:           params.MetaFilter,
		MetaFilterConditions: params.MetaFilterConditions,
		UseGraph:             params.UseGraph,
		AutoCitation:         true,
	}
}

func toString(data interface{}) string {
	if data != nil {
		return data.(string)
	}
	return ""
}

func buildToolParams(sseRequest *config.AgentSSERequest) *ToolParams {
	var mcpToolList []*MCPToolInfo
	if len(sseRequest.McpTools) > 0 {
		for _, info := range sseRequest.McpTools {
			mcpToolList = append(mcpToolList, &MCPToolInfo{
				Transport:    info.Transport,
				URL:          info.URL,
				ToolNameList: info.ToolNameList,
			})
		}
	}
	return &ToolParams{
		PluginToolList: sseRequest.PluginList,
		McpToolList:    mcpToolList,
	}
}
