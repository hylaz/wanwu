package service

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/UnicomAI/wanwu/internal/knowledge-service/pkg/config"
	"github.com/UnicomAI/wanwu/internal/knowledge-service/pkg/http"
	http_client "github.com/UnicomAI/wanwu/pkg/http-client"
	"github.com/UnicomAI/wanwu/pkg/log"
)

type RagQACreateParams struct {
	UserId           string `json:"userId"`             // 用户id
	QABase           string `json:"QABase"`             // 问答库名称
	QAId             string `json:"QAId"`               // 问答id
	EmbeddingModelId string `json:"embedding_model_id"` // 	embedding模型id
}

type RagQADeleteParams struct {
	UserId string `json:"userId"`
	QABase string `json:"QABase"`
	QAId   string `json:"QAId"`
}

type BatchRagQAMetaParams struct {
	UserId string        `json:"userId"`
	QABase string        `json:"QABase"`
	QAId   string        `json:"QAId"`
	Metas  []*QAMetaInfo `json:"metas"`
}

type QAMetaInfo struct {
	QAPairId     string        `json:"qa_pair_id"`
	MetaDataList []*QAMetaData `json:"metadata_list"`
}

type QAMetaData struct {
	Key       string      `json:"key"`
	Value     interface{} `json:"value"`
	ValueType string      `json:"value_type"`
}

type RagQABatchDeleteMetaParams struct {
	UserId string   `json:"userId"` // 用户id
	QABase string   `json:"QABase"` // 问答库名称
	QAId   string   `json:"QAId"`   // 问答库id
	Keys   []string `json:"keys"`   // 删除的元数据key列表
}

type RagQABatchUpdateMetaKeyParams struct {
	UserId   string            `json:"userId"`   // 用户id
	QABase   string            `json:"QABase"`   // 问答库名称
	QAId     string            `json:"QAId"`     // 问答库id
	Mappings []*RagMetaMapKeys `json:"mappings"` // 元数据key映射列表
}

type RagQAMetaMapKeys struct {
	OldKey string `json:"old_key"`
	NewKey string `json:"new_key"`
}

type RagQAPairItem struct {
	QAPairId string `json:"qa_pair_id"`
	Question string `json:"question,omitempty"`
	Answer   string `json:"answer,omitempty"`
}

type RagAddQAPairParams struct {
	UserId     string           `json:"userId"`
	QAId       string           `json:"QAId"`
	QABaseName string           `json:"QABase"`
	QAPairs    []*RagQAPairItem `json:"QAPairs"`
}

type RagUpdateQAPairParams struct {
	UserId     string         `json:"userId"`
	QAId       string         `json:"QAId"`
	QABaseName string         `json:"QABase"`
	QAPair     *RagQAPairItem `json:"QAPair"`
}

type RagUpdateQAPairStatusParams struct {
	UserId     string `json:"userId"`
	QAId       string `json:"QAId"`
	QABaseName string `json:"QABase"`
	QAPairId   string `json:"QAPairId"`
	Status     bool   `json:"status"`
}

type RagDeleteQAPairParams struct {
	UserId     string   `json:"userId"`
	QAId       string   `json:"QAId"`
	QABaseName string   `json:"QABase"`
	QAPairIds  []string `json:"QAPairIds"`
}

type RagKnowledgeQAHitResp struct {
	Code    int                 `json:"code"`
	Message string              `json:"message"`
	Data    *KnowledgeQAHitData `json:"data"`
}

type KnowledgeQAHitData struct {
	SearchList []*QASearchItem `json:"searchList"`
	Score      []float64       `json:"score"`
}

type QASearchItem struct {
	Title       string      `json:"title"`
	Question    string      `json:"question"`
	Answer      string      `json:"answer"`
	QAPairId    string      `json:"qa_pair_id"`
	QABase      string      `json:"QABase"`
	QAId        string      `json:"QAId"`
	MetaData    interface{} `json:"meta_data"`
	ContentType string      `json:"content_type"` // graph：知识图谱（文本）, text：文档分段（文本）, community_report：社区报告（markdown）
}

type KnowledgeQAHitParams struct {
	UserId                      string                  `json:"userId" validate:"required"`
	KnowledgeIdList             []string                `json:"knowledgeIdList,omitempty" validate:"required"`
	Question                    string                  `json:"question"`
	ReturnMeta                  bool                    `json:"returnMeta"`
	Threshold                   float64                 `json:"threshold"`
	TopK                        int64                   `json:"topK"`
	RetrieveMethod              string                  `json:"retrieveMethod"`
	RerankMod                   string                  `json:"rerankMod"`
	RerankModelId               string                  `json:"rerankModelId"`
	MetadataFiltering           bool                    `json:"metadataFiltering"`
	MetadataFilteringConditions []*QAMetadataFilterItem `json:"metadataFilteringConditions"`
	Weight                      *WeightParams           `json:"weights"`
}

type QAMetadataFilterItem struct {
	FilteringQaBaseName string        `json:"filtering_qa_base_name"`
	LogicalOperator     string        `json:"logical_operator"`
	Conditions          []*QAMetaItem `json:"conditions"`
}

type QAMetaItem struct {
	MetaName           string      `json:"meta_name"`           // 元数据名称
	MetaType           string      `json:"meta_type"`           // 元数据类型
	ComparisonOperator string      `json:"comparison_operator"` // 比较运算符
	Value              interface{} `json:"value,omitempty"`     // 用于过滤的条件值
}

// RagQACreate rag创建问答库
func RagQACreate(ctx context.Context, ragQAParams *RagQACreateParams) error {
	ragServer := config.GetConfig().RagServer
	url := ragServer.Endpoint + ragServer.InitQABaseUri
	paramsByte, err := json.Marshal(ragQAParams)
	if err != nil {
		return err
	}
	result, err := http.GetClient().PostJson(ctx, &http_client.HttpRequestParams{
		Url:        url,
		Body:       paramsByte,
		Timeout:    time.Duration(ragServer.Timeout) * time.Second,
		MonitorKey: "rag_qa_create",
		LogLevel:   http_client.LogAll,
	})
	if err != nil {
		return err
	}
	var resp RagCommonResp
	if err := json.Unmarshal(result, &resp); err != nil {
		log.Errorf(err.Error())
		return err
	}

	if resp.Code != successCode {
		return errors.New(resp.Message)
	}
	return nil
}

// RagQADelete rag删除问答库
func RagQADelete(ctx context.Context, ragDeleteParams *RagQADeleteParams) error {
	ragServer := config.GetConfig().RagServer
	url := ragServer.Endpoint + ragServer.DeleteQABaseUri
	paramsByte, err := json.Marshal(ragDeleteParams)
	if err != nil {
		return err
	}
	result, err := http.GetClient().PostJson(ctx, &http_client.HttpRequestParams{
		Url:        url,
		Body:       paramsByte,
		Timeout:    time.Duration(ragServer.Timeout) * time.Second,
		MonitorKey: "rag_qa_delete",
		LogLevel:   http_client.LogAll,
	})
	if err != nil {
		return err
	}
	var resp RagCommonResp
	if err := json.Unmarshal(result, &resp); err != nil {
		log.Errorf(err.Error())
		return err
	}
	if resp.Code != successCode {
		return errors.New(resp.Message)
	}
	return nil
}

// BatchRagQAMeta 更新问答库元数据
func BatchRagQAMeta(ctx context.Context, batchRagQAMetaParams *BatchRagQAMetaParams) error {
	ragServer := config.GetConfig().RagServer
	url := ragServer.Endpoint + ragServer.UpdateQAMetasUri
	paramsByte, err := json.Marshal(batchRagQAMetaParams)
	if err != nil {
		return err
	}
	result, err := http.GetClient().PostJson(ctx, &http_client.HttpRequestParams{
		Url:        url,
		Body:       paramsByte,
		Timeout:    time.Duration(ragServer.Timeout) * time.Second,
		MonitorKey: "batch_rag_qa_meta",
		LogLevel:   http_client.LogAll,
	})
	if err != nil {
		return err
	}
	var resp RagCommonResp
	if err := json.Unmarshal(result, &resp); err != nil {
		log.Errorf(err.Error())
		return err
	}
	if resp.Code != successCode {
		return errors.New(resp.Message)
	}
	return nil
}

func RagQABatchDeleteMeta(ctx context.Context, ragDeleteParams *RagQABatchDeleteMetaParams) error {
	ragServer := config.GetConfig().RagServer
	url := ragServer.Endpoint + ragServer.DeleteQAMetaKeyUri
	paramsByte, err := json.Marshal(ragDeleteParams)
	if err != nil {
		return err
	}
	result, err := http.GetClient().PostJson(ctx, &http_client.HttpRequestParams{
		Url:        url,
		Body:       paramsByte,
		Timeout:    time.Duration(ragServer.Timeout) * time.Second,
		MonitorKey: "rag_qa_delete_meta_key",
		LogLevel:   http_client.LogAll,
	})
	if err != nil {
		return err
	}
	var resp RagCommonResp
	if err := json.Unmarshal(result, &resp); err != nil {
		log.Errorf(err.Error())
		return err
	}
	if resp.Code != successCode {
		return errors.New(resp.Message)
	}
	return nil
}

func RagQABatchUpdateMeta(ctx context.Context, ragUpdateParams *RagQABatchUpdateMetaKeyParams) error {
	ragServer := config.GetConfig().RagServer
	url := ragServer.Endpoint + ragServer.RenameQAMetakeyUri
	paramsByte, err := json.Marshal(ragUpdateParams)
	if err != nil {
		return err
	}
	result, err := http.GetClient().PostJson(ctx, &http_client.HttpRequestParams{
		Url:        url,
		Body:       paramsByte,
		Timeout:    time.Duration(ragServer.Timeout) * time.Second,
		MonitorKey: "rag_qa_update_meta_key",
		LogLevel:   http_client.LogAll,
	})
	if err != nil {
		return err
	}
	var resp RagCommonResp
	if err := json.Unmarshal(result, &resp); err != nil {
		log.Errorf(err.Error())
		return err
	}
	if resp.Code != successCode {
		return errors.New(resp.Message)
	}
	return nil
}

// RagBatchAddQAPairs rag批量新增问答对
func RagBatchAddQAPairs(ctx context.Context, ragAddQAPairParams *RagAddQAPairParams) error {
	ragServer := config.GetConfig().RagServer
	url := ragServer.Endpoint + ragServer.BatchAddQAPairsUri
	paramsByte, err := json.Marshal(ragAddQAPairParams)
	if err != nil {
		return err
	}
	result, err := http.GetClient().PostJson(ctx, &http_client.HttpRequestParams{
		Url:        url,
		Body:       paramsByte,
		Timeout:    time.Duration(ragServer.Timeout) * time.Second,
		MonitorKey: "rag_add_qa_pair",
		LogLevel:   http_client.LogAll,
	})
	if err != nil {
		return err
	}
	var resp RagCommonResp
	if err := json.Unmarshal(result, &resp); err != nil {
		log.Errorf(err.Error())
		return err
	}
	if resp.Code != successCode {
		return errors.New(resp.Message)
	}
	return nil
}

// RagUpdateQAPair rag更新问答对
func RagUpdateQAPair(ctx context.Context, ragUpdateQAPairParams *RagUpdateQAPairParams) error {
	ragServer := config.GetConfig().RagServer
	url := ragServer.Endpoint + ragServer.UpdateQAPairUri
	paramsByte, err := json.Marshal(ragUpdateQAPairParams)
	if err != nil {
		return err
	}
	result, err := http.GetClient().PostJson(ctx, &http_client.HttpRequestParams{
		Url:        url,
		Body:       paramsByte,
		Timeout:    time.Duration(ragServer.Timeout) * time.Second,
		MonitorKey: "rag_update_qa_pair",
		LogLevel:   http_client.LogAll,
	})
	if err != nil {
		return err
	}
	var resp RagCommonResp
	if err := json.Unmarshal(result, &resp); err != nil {
		log.Errorf(err.Error())
		return err
	}
	if resp.Code != successCode {
		return errors.New(resp.Message)
	}
	return nil
}

// RagUpdateQAPairStatus rag启停问答对
func RagUpdateQAPairStatus(ctx context.Context, ragUpdateQAPairStatusParams *RagUpdateQAPairStatusParams) error {
	ragServer := config.GetConfig().RagServer
	url := ragServer.Endpoint + ragServer.UpdateQAPairStatusUri
	paramsByte, err := json.Marshal(ragUpdateQAPairStatusParams)
	if err != nil {
		return err
	}
	result, err := http.GetClient().PostJson(ctx, &http_client.HttpRequestParams{
		Url:        url,
		Body:       paramsByte,
		Timeout:    time.Duration(ragServer.Timeout) * time.Second,
		MonitorKey: "rag_update_qa_pair_status",
		LogLevel:   http_client.LogAll,
	})
	if err != nil {
		return err
	}
	var resp RagCommonResp
	if err := json.Unmarshal(result, &resp); err != nil {
		log.Errorf(err.Error())
		return err
	}
	if resp.Code != successCode {
		return errors.New(resp.Message)
	}
	return nil
}

// RagDeleteQAPair rag删除问答对
func RagDeleteQAPair(ctx context.Context, ragDeleteQAPairParams *RagDeleteQAPairParams) error {
	ragServer := config.GetConfig().RagServer
	url := ragServer.Endpoint + ragServer.BatchDeleteQAPairsUri
	paramsByte, err := json.Marshal(ragDeleteQAPairParams)
	if err != nil {
		return err
	}
	result, err := http.GetClient().PostJson(ctx, &http_client.HttpRequestParams{
		Url:        url,
		Body:       paramsByte,
		Timeout:    time.Duration(ragServer.Timeout) * time.Second,
		MonitorKey: "rag_delete_qa_pair",
		LogLevel:   http_client.LogAll,
	})
	if err != nil {
		return err
	}
	var resp RagCommonResp
	if err := json.Unmarshal(result, &resp); err != nil {
		log.Errorf(err.Error())
		return err
	}
	if resp.Code != successCode {
		return errors.New(resp.Message)
	}
	return nil
}

// RagKnowledgeQAHit rag问答库命中测试
func RagKnowledgeQAHit(ctx context.Context, knowledgeQAHitParams *KnowledgeQAHitParams) (*RagKnowledgeQAHitResp, error) {
	ragServer := config.GetConfig().RagServer
	url := ragServer.ProxyPoint + ragServer.KnowledgeQAHitUri
	paramsByte, err := json.Marshal(knowledgeQAHitParams)
	if err != nil {
		return nil, err
	}
	result, err := http.GetClient().PostJson(ctx, &http_client.HttpRequestParams{
		Url:        url,
		Body:       paramsByte,
		Timeout:    time.Duration(ragServer.Timeout) * time.Second,
		MonitorKey: "rag_knowledge_qa_hit",
		LogLevel:   http_client.LogAll,
	})
	if err != nil {
		return nil, err
	}
	var resp RagKnowledgeQAHitResp
	if err := json.Unmarshal(result, &resp); err != nil {
		log.Errorf(err.Error())
		return nil, err
	}
	if resp.Code != successCode {
		return nil, errors.New(resp.Message)
	}
	if resp.Data == nil {
		return nil, errors.New("QA Hit data is empty")
	}
	return &resp, nil
}
