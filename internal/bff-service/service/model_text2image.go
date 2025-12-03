package service

import (
	"fmt"
	"net/http"

	err_code "github.com/UnicomAI/wanwu/api/proto/err-code"
	model_service "github.com/UnicomAI/wanwu/api/proto/model-service"
	gin_util "github.com/UnicomAI/wanwu/pkg/gin-util"
	grpc_util "github.com/UnicomAI/wanwu/pkg/grpc-util"
	mp "github.com/UnicomAI/wanwu/pkg/model-provider"
	mp_common "github.com/UnicomAI/wanwu/pkg/model-provider/mp-common"
	"github.com/gin-gonic/gin"
)

func ModelText2Image(ctx *gin.Context, modelID string, req *mp_common.Text2ImageReq) {
	// modelInfo by modelID
	modelInfo, err := model.GetModelById(ctx.Request.Context(), &model_service.GetModelByIdReq{ModelId: modelID})
	if err != nil {
		gin_util.Response(ctx, nil, err)
		return
	}

	// text2Image config
	text2Image, err := mp.ToModelConfig(modelInfo.Provider, modelInfo.ModelType, modelInfo.ProviderConfig)
	if err != nil {
		gin_util.Response(ctx, nil, grpc_util.ErrorStatus(err_code.Code_BFFGeneral, fmt.Sprintf("model %v text2Image err: %v", modelInfo.ModelId, err)))
		return
	}
	iText2Image, ok := text2Image.(mp.IText2Image)
	if !ok {
		gin_util.Response(ctx, nil, grpc_util.ErrorStatus(err_code.Code_BFFGeneral, fmt.Sprintf("model %v text2Image err: invalid provider", modelInfo.ModelId)))
		return
	}
	// text2Image
	text2ImageReq, err := iText2Image.NewReq(req)
	if err != nil {
		gin_util.Response(ctx, nil, grpc_util.ErrorStatus(err_code.Code_BFFGeneral, fmt.Sprintf("model %v text2Image NewReq err: %v", modelInfo.ModelId, err)))
		return
	}
	resp, err := iText2Image.Text2Image(ctx, text2ImageReq)
	if err != nil {
		gin_util.Response(ctx, nil, grpc_util.ErrorStatus(err_code.Code_BFFGeneral, fmt.Sprintf("model %v text2Image err: %v", modelInfo.ModelId, err)))
		return
	}
	if data, ok := resp.ConvertResp(); ok {
		status := http.StatusOK
		ctx.Set(gin_util.STATUS, status)
		//ctx.Set(config.RESULT, resp.String())
		ctx.JSON(status, data)
		return
	}
	gin_util.Response(ctx, nil, grpc_util.ErrorStatus(err_code.Code_BFFGeneral, fmt.Sprintf("model %v text2Image err: invalid resp", modelInfo.ModelId)))
}
