package request

import (
	"github.com/UnicomAI/wanwu/internal/agent-service/model"
	"github.com/getkin/kin-openapi/openapi3"
)

type AgentChatContext struct {
	AgentChatReq     *AgentChatReq
	KnowledgeHitData *model.KnowledgeHitData //  rag命中数据
}

type AgentChatReq struct {
	Input           string           `json:"input" validate:"required"`
	UploadFile      []string         `json:"uploadFile"`
	Stream          bool             `json:"stream"`
	AgentBaseParams *AgentBaseParams `json:"agentBaseParams" validate:"required"` // 智能体基础参数
	ModelParams     *ModelParams     `json:"modelParams" validate:"required"`     // 模型参数
	KnowledgeParams *KnowledgeParams `json:"knowledgeParams"`                     // 知识库参数，如果后续需要增加透传，理论上只需要修改此KnowledgeParams即可
	ToolParams      *ToolParams      `json:"toolParams"`                          // 工具相关参数，mcp tool plugin tool
}

type AgentBaseParams struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Instruction string `json:"instruction"`
}

type ModelParams struct {
	ModelId    string                       `json:"modelId" validate:"required"`
	History    []AssistantConversionHistory `json:"history,omitempty"`
	MaxHistory int                          `json:"maxHistory"`
}

type KnowledgeParams struct {
	UserId               string                 `json:"userId"`          // 用户id
	KnowledgeIdList      []string               `json:"knowledgeIdList"` // 知识库id列表
	Question             string                 `json:"question"`
	Threshold            float32                `json:"threshold"` // Score阈值
	TopK                 int32                  `json:"topK"`
	Stream               bool                   `json:"stream"`
	Chichat              bool                   `json:"chichat"` // 当知识库召回结果为空时是否使用默认话术（兜底），默认为true
	RerankModelId        string                 `json:"rerank_model_id"`
	CustomModelInfo      *CustomModelInfo       `json:"custom_model_info"`
	MaxHistory           int32                  `json:"max_history"`
	RewriteQuery         bool                   `json:"rewrite_query"`   // 是否query改写
	RerankMod            string                 `json:"rerank_mod"`      // rerank_model:重排序模式，weighted_score：权重搜索
	RetrieveMethod       string                 `json:"retrieve_method"` // hybrid_search:混合搜索， semantic_search:向量搜索， full_text_search：文本搜索
	Weight               *WeightParams          `json:"weights"`         // 权重搜索下的权重配置
	Temperature          float32                `json:"temperature,omitempty"`
	TopP                 float32                `json:"top_p,omitempty"`               // 多样性
	RepetitionPenalty    float32                `json:"repetition_penalty,omitempty"`  // 重复惩罚/频率惩罚
	ReturnMeta           bool                   `json:"return_meta,omitempty"`         // 是否返回元数据
	AutoCitation         bool                   `json:"auto_citation"`                 // 是否自动角标
	TermWeight           float32                `json:"term_weight_coefficient"`       // 关键词系数
	MetaFilter           bool                   `json:"metadata_filtering"`            // 元数据过滤开关
	MetaFilterConditions []*MetadataFilterParam `json:"metadata_filtering_conditions"` // 元数据过滤条件
	UseGraph             bool                   `json:"use_graph"`                     // 是否启动知识图谱查询
}

type CustomModelInfo struct {
	LlmModelID string `json:"llm_model_id"`
}

type ToolParams struct {
	PluginToolList []*PluginToolInfo `json:"pluginTool,omitempty"`
	McpToolList    []*MCPToolInfo    `json:"mcpToolList,omitempty"`
}

type NetSearchParams struct {
	SearchUrl string `json:"search_url,omitempty"`
	SearchKey string `json:"search_key,omitempty"`
	UseSearch bool   `json:"use_search,omitempty"`
}

type AssistantConversionHistory struct {
	Query         string   `json:"query"`
	UploadFileUrl []string `json:"upload_file_url,omitempty"`
	Response      string   `json:"response"`
}

type PluginToolInfo struct {
	APISchema *openapi3.T `json:"api_schema"`
	APIAuth   *APIAuth    `json:"api_auth,omitempty"`
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

type MetadataFilterParam struct {
	FilterKnowledgeName string                `json:"filtering_kb_name"`
	LogicalOperator     string                `json:"logical_operator"`
	MetaList            []*MetadataFilterItem `json:"conditions"` // 元数据过滤列表
}

type MetadataFilterItem struct {
	MetaName           string      `json:"meta_name"`           // 元数据名称
	MetaType           string      `json:"meta_type"`           // 元数据类型
	ComparisonOperator string      `json:"comparison_operator"` // 比较运算符
	Value              interface{} `json:"value,omitempty"`     // 用于过滤的条件值
}

type WeightParams struct {
	VectorWeight float32 `json:"vector_weight"` //语义权重
	TextWeight   float32 `json:"text_weight"`   //关键字权重
}

func (c *AgentChatReq) Check() error {
	return nil
}
