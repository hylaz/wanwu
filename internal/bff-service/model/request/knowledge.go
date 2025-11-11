package request

import (
	"errors"
	"fmt"

	"github.com/UnicomAI/wanwu/pkg/util"
)

type KnowledgeSelectReq struct {
	Name      string   `json:"name" form:"name" `
	TagIdList []string `json:"tagId" form:"tagId" `
	CommonCheck
}

type KnowledgeBatchSelectReq struct {
	KnowledgeIdList []string `json:"knowledgeIdList" form:"knowledgeIdList" `
	UserId          string   `json:"userId" form:"userId" `
	CommonCheck
}

type CreateKnowledgeReq struct {
	Name           string          `json:"name"  validate:"required"`
	Description    string          `json:"description"`
	EmbeddingModel *EmbeddingModel `json:"embeddingModelInfo" validate:"required"`
	KnowledgeGraph *KnowledgeGraph `json:"knowledgeGraph" validate:"required"`
}

type UpdateKnowledgeReq struct {
	KnowledgeId string `json:"knowledgeId"   validate:"required"`
	Name        string `json:"name"   validate:"required"`
	Description string `json:"description"`
	CommonCheck
}

type KnowledgeHitReq struct {
	KnowledgeList        []*AppKnowledgeBase   `json:"knowledgeList"`
	Question             string                `json:"question"   validate:"required"`
	KnowledgeMatchParams *KnowledgeMatchParams `json:"knowledgeMatchParams"   validate:"required"`
	CommonCheck
}

type KnowledgeMatchParams struct {
	MatchType         string  `json:"matchType"  validate:"required"` //matchType：vector（向量检索）、text（文本检索）、mix（混合检索：向量+文本）
	RerankModelId     string  `json:"rerankModelId"`                  //rerank模型id
	PriorityMatch     int32   `json:"priorityMatch"`                  // 权重匹配，只有在混合检索模式下，选择权重设置后，这个才设置为1
	SemanticsPriority float32 `json:"semanticsPriority"`              // 语义权重
	KeywordPriority   float32 `json:"keywordPriority"`                // 关键词权重
	TopK              int32   `json:"topK"`                           //topK 获取最高的几行
	Threshold         float32 `json:"threshold"`                      //threshold 过滤分数阈值
	TermWeight        float32 `json:"termWeight"`                     // 关键词系数
	TermWeightEnable  bool    `json:"termWeightEnable"`               // 关键词系数开关
	CommonCheck
}

type EmbeddingModel struct {
	ModelId string `json:"modelId"  validate:"required"`
}

// KnowledgeGraph 知识图谱信息
type KnowledgeGraph struct {
	Switch     bool   `json:"switch"`     //知识图谱开关
	LLMModelId string `json:"llmModelId"` //大模型id，开关为true必填
	SchemaUrl  string `json:"schemaUrl"`  //模型schema文件地址，可以为空
}

type DeleteKnowledge struct {
	KnowledgeId string `json:"knowledgeId" validate:"required"`
	CommonCheck
}

type GetKnowledgeReq struct {
	KnowledgeId string `json:"knowledgeId" validate:"required"`
	CommonCheck
}

type CallbackUpdateDocStatusReq struct {
	DocId        string              `json:"id" validate:"required"`
	Status       int32               `json:"status" validate:"required"`
	MetaDataList []*CallbackMetaData `json:"metaDataList"`
	CommonCheck
}

type CallbackUpdateKnowledgeStatusReq struct {
	KnowledgeId  string `json:"knowledgeId" validate:"required"`
	ReportStatus int32  `json:"reportStatus" validate:"required"` //此状态不会是0
	CommonCheck
}

type CallbackMetaData struct {
	Key    string `json:"key"`
	MetaId string `json:"metaId" validate:"required"`
	Value  string `json:"value" validate:"required"`
}

type DocMetaData struct {
	MetaId        string `json:"metaId"`        // 元数据id
	MetaKey       string `json:"metaKey"`       // key
	MetaValue     string `json:"metaValue"`     // 确定值
	MetaValueType string `json:"metaValueType"` // string，number，time
	MetaRule      string `json:"metaRule"`      // 正则表达式
	Option        string `json:"option"`        // option:add(新增)、update(更新)、delete(删除),update 和delete 的时候metaId 不能为空
}

type SearchKnowledgeInfoReq struct {
	KnowledgeName string `json:"categoryName" form:"categoryName" validate:"required"`
	UserId        string `json:"userId" form:"userId" validate:"required"`
	OrgId         string `json:"orgId"`
	CommonCheck
}

type GetKnowledgeMetaSelectReq struct {
	KnowledgeId string `json:"knowledgeId"  form:"knowledgeId" validate:"required"`
	CommonCheck
}

type KnowledgeMetaValueListReq struct {
	KnowledgeId string   `json:"knowledgeId"  form:"knowledgeId" validate:"required"`
	DocIdList   []string `json:"docIdList" form:"docIdList" validate:"required" `
	CommonCheck
}

type UpdateMetaValueReq struct {
	KnowledgeId     string         `json:"knowledgeId"  form:"knowledgeId" validate:"required"`
	DocIdList       []string       `json:"docIdList"  validate:"required"`
	MetaValueList   []*DocMetaData `json:"metaValueList"`
	ApplyToSelected bool           `json:"applyToSelected"`
}

// RagSearchKnowledgeBaseReq rag知识库查询请求
type RagSearchKnowledgeBaseReq struct {
	UserId               string                         `json:"userId" validate:"required"`
	Question             string                         `json:"question" validate:"required"`
	KnowledgeIdList      []string                       `json:"knowledgeIdList,omitempty" validate:"required"`
	KnowledgeUser        map[string][]*RagKnowledgeInfo `json:"knowledge_base_info"`
	Threshold            float64                        `json:"threshold"`
	TopK                 int32                          `json:"topK"`
	RerankModelId        string                         `json:"rerank_model_id"`               // rerankId
	RerankMod            string                         `json:"rerank_mod"`                    // rerank_model:重排序模式，weighted_score：权重搜索
	RetrieveMethod       string                         `json:"retrieve_method"`               // hybrid_search:混合搜索， semantic_search:向量搜索， full_text_search：文本搜索
	Weight               *WeightParams                  `json:"weights"`                       // 权重搜索下的权重配置
	TermWeight           float32                        `json:"term_weight_coefficient"`       // 关键词系数
	MetaFilter           bool                           `json:"metadata_filtering"`            // 元数据过滤开关
	MetaFilterConditions []*MetadataFilterItem          `json:"metadata_filtering_conditions"` // 元数据过滤条件
	CommonCheck
}

type RagKnowledgeChatReq struct {
	UserId               string                         `json:"userId"`
	KnowledgeUser        map[string][]*RagKnowledgeInfo `json:"knowledge_base_info"`
	KnowledgeIdList      []string                       `json:"knowledgeIdList"` // 知识库id列表
	Question             string                         `json:"question"`
	Threshold            float32                        `json:"threshold"` // Score阈值
	TopK                 int32                          `json:"topK"`
	Stream               bool                           `json:"stream"`
	Chichat              bool                           `json:"chichat"` // 当知识库召回结果为空时是否使用默认话术（兜底），默认为true
	RerankModelId        string                         `json:"rerank_model_id"`
	CustomModelInfo      *CustomModelInfo               `json:"custom_model_info"`
	History              []*HistoryItem                 `json:"history"`
	MaxHistory           int32                          `json:"max_history"`
	RewriteQuery         bool                           `json:"rewrite_query"`   // 是否query改写
	RerankMod            string                         `json:"rerank_mod"`      // rerank_model:重排序模式，weighted_score：权重搜索
	RetrieveMethod       string                         `json:"retrieve_method"` // hybrid_search:混合搜索， semantic_search:向量搜索， full_text_search：文本搜索
	Weight               *WeightParams                  `json:"weights"`         // 权重搜索下的权重配置
	Temperature          float32                        `json:"temperature,omitempty"`
	TopP                 float32                        `json:"top_p,omitempty"`               // 多样性
	RepetitionPenalty    float32                        `json:"repetition_penalty,omitempty"`  // 重复惩罚/频率惩罚
	ReturnMeta           bool                           `json:"return_meta,omitempty"`         // 是否返回元数据
	AutoCitation         bool                           `json:"auto_citation"`                 // 是否自动角标
	TermWeight           float32                        `json:"term_weight_coefficient"`       // 关键词系数
	MetaFilter           bool                           `json:"metadata_filtering"`            // 元数据过滤开关
	MetaFilterConditions []*MetadataFilterItem          `json:"metadata_filtering_conditions"` // 元数据过滤条件
	CommonCheck
}

type CustomModelInfo struct {
	LlmModelID string `json:"llm_model_id"`
}

type HistoryItem struct {
	Query       string `json:"query"`
	Response    string `json:"response"`
	NeedHistory bool   `json:"needHistory"`
}

type RagKnowledgeInfo struct {
	KnowledgeId   string `json:"kb_id"`
	KnowledgeName string `json:"kb_name"`
}

type WeightParams struct {
	VectorWeight float32 `json:"vector_weight"` //语义权重
	TextWeight   float32 `json:"text_weight"`   //关键字权重
}

type MetadataFilterItem struct {
	FilterKnowledgeName string      `json:"filtering_kb_name"`
	LogicalOperator     string      `json:"logical_operator"`
	Conditions          []*MetaItem `json:"conditions"`
}

type MetaItem struct {
	MetaName           string      `json:"meta_name"`           // 元数据名称
	MetaType           string      `json:"meta_type"`           // 元数据类型
	ComparisonOperator string      `json:"comparison_operator"` // 比较运算符
	Value              interface{} `json:"value,omitempty"`     // 用于过滤的条件值
}

func (c *UpdateMetaValueReq) Check() error {
	for _, v := range c.MetaValueList {
		if v.Option == "" {
			return errors.New("option为空")
		}
	}
	return nil
}
func (c *CreateKnowledgeReq) Check() error {
	if !util.IsAlphanumeric(c.Name) {
		errMsg := fmt.Sprintf("知识库名称只能包含中文、数字、小写英文，符号之只能包含下划线和减号 参数(%v)", c.Name)
		return errors.New(errMsg)
	}
	if c.KnowledgeGraph == nil {
		return errors.New("knowledge graph can not be nil")
	}
	if c.KnowledgeGraph.Switch && c.KnowledgeGraph.LLMModelId == "" {
		return errors.New("knowledge graph llmModelId can not be empty")
	}
	return nil
}
