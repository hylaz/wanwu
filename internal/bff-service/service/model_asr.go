package service

import (
	"encoding/json"
	"fmt"
	"net/http"

	err_code "github.com/UnicomAI/wanwu/api/proto/err-code"
	model_service "github.com/UnicomAI/wanwu/api/proto/model-service"
	"github.com/UnicomAI/wanwu/internal/bff-service/config"
	gin_util "github.com/UnicomAI/wanwu/pkg/gin-util"
	grpc_util "github.com/UnicomAI/wanwu/pkg/grpc-util"
	mp "github.com/UnicomAI/wanwu/pkg/model-provider"
	mp_common "github.com/UnicomAI/wanwu/pkg/model-provider/mp-common"
	mp_yuanjing "github.com/UnicomAI/wanwu/pkg/model-provider/mp-yuanjing"
	"github.com/gin-gonic/gin"
)

func ModelAsr(ctx *gin.Context, modelID string, req *mp_common.AsrReq) {
	// modelInfo by modelID
	modelInfo, err := model.GetModelById(ctx.Request.Context(), &model_service.GetModelByIdReq{ModelId: modelID})
	if err != nil {
		gin_util.Response(ctx, nil, err)
		return
	}
	modelAsr(ctx, modelID, modelInfo.Provider, modelInfo.ModelType, modelInfo.ProviderConfig, req)
}

func modelAsr(ctx *gin.Context, modelId, provider, modelType, providerConfig string, req *mp_common.AsrReq) {
	// asr config
	asr, err := mp.ToModelConfig(provider, modelType, providerConfig)
	if err != nil {
		gin_util.Response(ctx, nil, grpc_util.ErrorStatus(err_code.Code_BFFGeneral, fmt.Sprintf("model %v asr err: %v", modelId, err)))
		return
	}
	iAsr, ok := asr.(mp.IAsr)
	if !ok {
		gin_util.Response(ctx, nil, grpc_util.ErrorStatus(err_code.Code_BFFGeneral, fmt.Sprintf("model %v asr err: invalid provider", modelId)))
		return
	}

	asrReq, err := iAsr.NewReq(req)
	if err != nil {
		gin_util.Response(ctx, nil, grpc_util.ErrorStatus(err_code.Code_BFFGeneral, fmt.Sprintf("model %v asr NewReq err: %v", modelId, err)))
		return
	}
	resp, err := iAsr.Asr(ctx, asrReq)
	if err != nil {
		gin_util.Response(ctx, nil, grpc_util.ErrorStatus(err_code.Code_BFFGeneral, fmt.Sprintf("model %v asr err: %v", modelId, err)))
		return
	}
	if data, ok := resp.ConvertResp(); ok {
		status := http.StatusOK
		ctx.Set(gin_util.STATUS, status)
		//ctx.Set(config.RESULT, resp.String())
		ctx.JSON(status, data)
		return
	}
	gin_util.Response(ctx, nil, grpc_util.ErrorStatus(err_code.Code_BFFGeneral, fmt.Sprintf("model %v asr err: invalid resp", modelId)))
}

// AudioBase64ConvertText 语音文件（base64格式）转文本内置工具服务
func AudioBase64ConvertText(ctx *gin.Context, req *mp_common.AsrReq) {
	modelsMap := config.Cfg().GetModelsMap()
	modelInfo, exists := modelsMap[config.AsrModelId]
	if !exists {
		gin_util.Response(ctx, nil, grpc_util.ErrorStatus(err_code.Code_BFFGeneral, fmt.Sprintf("AudioBase64ConvertText err: 模型ID %s 不存在于配置中", config.AsrModelId)))
	}
	// asr config
	asrProviderConfig := &mp_yuanjing.Asr{
		ApiKey:      req.ApiKey,
		EndpointUrl: modelInfo.Endpoint,
	}
	providerConfig, err := json.Marshal(asrProviderConfig)
	if err != nil {
		gin_util.Response(ctx, nil, grpc_util.ErrorStatus(err_code.Code_BFFGeneral, fmt.Sprintf("AudioBase64ConvertText err: %v", err)))
		return
	}
	modelAsr(ctx, modelInfo.ModelId, modelInfo.Provider, modelInfo.ModelType, string(providerConfig), req)
}
