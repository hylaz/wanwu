package callback

import (
	"encoding/json"
	"fmt"

	err_code "github.com/UnicomAI/wanwu/api/proto/err-code"
	"github.com/UnicomAI/wanwu/internal/bff-service/model/request"
	"github.com/UnicomAI/wanwu/internal/bff-service/service"
	gin_util "github.com/UnicomAI/wanwu/pkg/gin-util"
	grpc_util "github.com/UnicomAI/wanwu/pkg/grpc-util"
	mp "github.com/UnicomAI/wanwu/pkg/model-provider"
	mp_common "github.com/UnicomAI/wanwu/pkg/model-provider/mp-common"
	"github.com/gin-gonic/gin"
)

// GetModelById
//
//	@Tags		callback
//	@Summary	根据ModelId获取模型
//	@Accept		json
//	@Produce	json
//	@Param		modelId	path		string	true	"模型ID"
//	@Success	200		{object}	response.Response{data=response.ModelInfo}
//	@Router		/model/{modelId} [get]
func GetModelById(ctx *gin.Context) {
	modelId := ctx.Param("modelId")
	resp, err := service.GetModelById(ctx, &request.GetModelByIdRequest{
		BaseModelRequest: request.BaseModelRequest{ModelId: modelId}})
	// 替换callback返回的模型中的apiKey/endpointUrl信息
	if resp != nil && resp.Config != nil {
		cfg := make(map[string]interface{})
		b, err := json.Marshal(resp.Config)
		if err != nil {
			gin_util.Response(ctx, nil, grpc_util.ErrorStatus(err_code.Code_BFFGeneral, fmt.Sprintf("model %v marshal config err: %v", modelId, err)))
			return
		}
		if err = json.Unmarshal(b, &cfg); err != nil {
			gin_util.Response(ctx, nil, grpc_util.ErrorStatus(err_code.Code_BFFGeneral, fmt.Sprintf("model %v unmarshal config err: %v", modelId, err)))
			return
		}
		// 替换apiKey, endpointUrl
		cfg["apiKey"] = "useless-api-key"
		endpoint := mp.ToModelEndpoint(resp.ModelId, resp.Model)
		for k, v := range endpoint {
			if k == "model_url" {
				cfg["endpointUrl"] = v
				break
			}
		}
		// 替换Config
		resp.Config = cfg
	}
	gin_util.Response(ctx, resp, err)
}

// ModelChatCompletions
//
//	@Tags		callback
//	@Summary	Model Chat Completions
//	@Accept		json
//	@Produce	json
//	@Param		modelId	path		string				true	"模型ID"
//	@Param		data	body		mp_common.LLMReq{}	true	"请求参数"
//	@Success	200		{object}	mp_common.LLMResp{}
//	@Router		/model/{modelId}/chat/completions [post]
func ModelChatCompletions(ctx *gin.Context) {
	var data mp_common.LLMReq
	if !gin_util.Bind(ctx, &data) {
		return
	}
	service.ModelChatCompletions(ctx, ctx.Param("modelId"), &data)
}

// ModelEmbeddings
//
//	@Tags		callback
//	@Summary	Model Embeddings
//	@Accept		json
//	@Produce	json
//	@Param		modelId	path		string						true	"模型ID"
//	@Param		data	body		mp_common.EmbeddingReq{}	true	"请求参数"
//	@Success	200		{object}	mp_common.EmbeddingResp{}
//	@Router		/model/{modelId}/embeddings [post]
func ModelEmbeddings(ctx *gin.Context) {
	var data mp_common.EmbeddingReq
	if !gin_util.Bind(ctx, &data) {
		return
	}
	service.ModelEmbeddings(ctx, ctx.Param("modelId"), &data)
}

// ModelRerank
//
//	@Tags		callback
//	@Summary	Model Rerank
//	@Accept		json
//	@Produce	json
//	@Param		modelId	path		string					true	"模型ID"
//	@Param		data	body		mp_common.RerankReq{}	true	"请求参数"
//	@Success	200		{object}	mp_common.RerankResp{}
//	@Router		/model/{modelId}/rerank [post]
func ModelRerank(ctx *gin.Context) {
	var data mp_common.RerankReq
	if !gin_util.Bind(ctx, &data) {
		return
	}
	service.ModelRerank(ctx, ctx.Param("modelId"), &data)
}

// ModelOcr
//
//	@Tags		callback
//	@Summary	Model Ocr
//	@Accept		multipart/form-data
//	@Produce	json
//	@Param		modelId	path		string	true	"模型ID"
//	@Param		file	formData	file	true	"文件"
//	@Success	200		{object}	mp_common.OcrResp{}
//	@Router		/model/{modelId}/ocr [post]
func ModelOcr(ctx *gin.Context) {
	var data mp_common.OcrReq
	if !gin_util.BindForm(ctx, &data) {
		return
	}
	service.ModelOcr(ctx, ctx.Param("modelId"), &data)
}

// ModelPdfParser
//
//	@Tags		callback
//	@Summary	Model PdfParser
//	@Accept		multipart/form-data
//	@Produce	json
//	@Param		modelId		path		string	true	"模型ID"
//	@Param		file		formData	file	true	"文件"
//	@Param		file_name	formData	string	true	"文件名"
//	@Success	200			{object}	mp_common.PdfParserResp{}
//	@Router		/model/{modelId}/pdf-parser [post]
func ModelPdfParser(ctx *gin.Context) {
	var data mp_common.PdfParserReq
	if !gin_util.BindForm(ctx, &data) {
		return
	}
	service.ModelPdfParser(ctx, ctx.Param("modelId"), &data)
}

// ModelGui
//
//	@Tags		callback
//	@Summary	Model Gui
//	@Accept		json
//	@Produce	json
//	@Param		modelId	path		string				true	"模型ID"
//	@Param		data	body		mp_common.GuiReq{}	true	"请求参数"
//	@Success	200		{object}	mp_common.GuiResp{}
//	@Router		/model/{modelId}/gui [post]
func ModelGui(ctx *gin.Context) {
	var data mp_common.GuiReq
	if !gin_util.Bind(ctx, &data) {
		return
	}
	service.ModelGui(ctx, ctx.Param("modelId"), &data)
}

// ModelAsr
//
//	@Tags		callback
//	@Summary	Model Asr
//	@Accept		multipart/form-data
//	@Produce	json
//	@Param		modelId	path		string	true	"模型ID"
//	@Param		file	formData	file	true	"语音文件"
//	@Param		config	formData	string	true	"请求参数"
//	@Success	200		{object}	mp_common.AsrResp{}
//	@Router		/model/{modelId}/asr [post]
func ModelAsr(ctx *gin.Context) {
	var data mp_common.AsrReq
	if !gin_util.BindForm(ctx, &data) {
		return
	}
	service.ModelAsr(ctx, ctx.Param("modelId"), &data)
}

// ModelText2Image
//
//	@Tags		callback
//	@Summary	Model Text-to-Image
//	@Accept		multipart/form-data
//	@Accept		json
//	@Produce	json
//	@Param		modelId	path		string						true	"模型ID"
//	@Param		data	body		mp_common.Text2ImageReq{}	true	"请求参数"
//	@Success	200		{object}	mp_common.Text2ImageResp{}
//	@Router		/model/{modelId}/text2image [post]
func ModelText2Image(ctx *gin.Context) {
	var data mp_common.Text2ImageReq
	if !gin_util.BindForm(ctx, &data) {
		return
	}
	service.ModelText2Image(ctx, ctx.Param("modelId"), &data)
}
