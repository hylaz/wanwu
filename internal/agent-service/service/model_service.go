package service

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/UnicomAI/wanwu/internal/agent-service/pkg/config"
	"github.com/UnicomAI/wanwu/internal/agent-service/pkg/http"
	"github.com/UnicomAI/wanwu/internal/bff-service/model/request"
	http_client "github.com/UnicomAI/wanwu/pkg/http-client"
	"github.com/UnicomAI/wanwu/pkg/log"
	mp "github.com/UnicomAI/wanwu/pkg/model-provider"
	mp_common "github.com/UnicomAI/wanwu/pkg/model-provider/mp-common"
)

const (
	successCode = 0
)

// BffResponse bff返回结果
type BffResponse struct {
	Code int64      `json:"code"`
	Data *ModelInfo `json:"data"`
	Msg  string     `json:"msg"`
}

type ModelInfo struct {
	ModelId     string                  `json:"modelId"`
	Provider    string                  `json:"provider" validate:"required" enums:"OpenAI-API-compatible,YuanJing"` // 模型供应商
	Model       string                  `json:"model" validate:"required"`                                           // 模型名称
	ModelType   string                  `json:"modelType" validate:"required" enums:"llm,embedding,rerank"`
	DisplayName string                  `json:"displayName"` // 模型显示名称
	Avatar      request.Avatar          `json:"avatar" `     // 模型图标路径
	PublishDate string                  `json:"publishDate"` // 模型发布时间
	IsActive    bool                    `json:"isActive"`    // 启用状态（true: 启用，false: 禁用）
	UserId      string                  `json:"userId"`
	OrgId       string                  `json:"orgId"`
	CreatedAt   string                  `json:"createdAt"`
	UpdatedAt   string                  `json:"updatedAt"`
	ModelDesc   string                  `json:"modelDesc"`
	Tags        []mp_common.Tag         `json:"tags"`
	Config      *LLMModelConfig         `json:"config"`
	Examples    *mp.ProviderModelConfig `json:"examples,omitempty"` // 仅用于swagger展示；模型对应供应商中的对应llm、embedding或rerank结构是config实际的参数
}

type LLMModelConfig struct {
	ApiKey          string `json:"apiKey"`
	EndpointUrl     string `json:"endpointUrl"`     // 模型名称
	FunctionCalling string `json:"functionCalling"` // 是否支持functionCall
	VisionSupport   string `json:"visionSupport"`   // 是否支持多模态
}

// SearchModel 查询model信息
func SearchModel(ctx context.Context, modelId string) (*ModelInfo, error) {
	bffServer := config.GetConfig().BffServer
	url := bffServer.Endpoint + bffServer.SearchModelUri + modelId
	result, err := http.GetClient().Get(ctx, &http_client.HttpRequestParams{
		Url:        url,
		Timeout:    time.Duration(bffServer.Timeout) * time.Second,
		MonitorKey: "search_model",
		LogLevel:   http_client.LogAll,
	})
	if err != nil {
		return nil, err
	}
	var resp BffResponse
	if err := json.Unmarshal(result, &resp); err != nil {
		log.Errorf(err.Error())
		return nil, err
	}
	if resp.Code != successCode {
		return nil, errors.New(resp.Msg)
	}
	return resp.Data, nil
}
