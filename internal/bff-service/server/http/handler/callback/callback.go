package callback

import (
	"encoding/json"
	"fmt"

	"github.com/UnicomAI/wanwu/internal/bff-service/model/response"

	err_code "github.com/UnicomAI/wanwu/api/proto/err-code"
	"github.com/UnicomAI/wanwu/internal/bff-service/model/request"
	"github.com/UnicomAI/wanwu/internal/bff-service/service"
	"github.com/UnicomAI/wanwu/pkg/constant"
	gin_util "github.com/UnicomAI/wanwu/pkg/gin-util"
	grpc_util "github.com/UnicomAI/wanwu/pkg/grpc-util"
	mp "github.com/UnicomAI/wanwu/pkg/model-provider"
	mp_common "github.com/UnicomAI/wanwu/pkg/model-provider/mp-common"
	"github.com/gin-gonic/gin"
)

//	@title		AI Agent Productivity Platform - Callback
//	@version	v0.0.1

//	@BasePath	/callback/v1

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

// UpdateDocStatus
//
//	@Tags			callback
//	@Summary		更新文档状态（模型扩展调用）
//	@Description	更新文档状态（模型扩展调用）
//	@Accept			json
//	@Produce		json
//	@Param			data	body		request.CallbackUpdateDocStatusReq	true	"更新文档状态请求参数"
//	@Success		200		{object}	response.Response
//	@Router			/api/docstatus [post]
func UpdateDocStatus(ctx *gin.Context) {
	var req request.CallbackUpdateDocStatusReq
	if !gin_util.Bind(ctx, &req) {
		return
	}
	err := service.UpdateDocStatus(ctx, &req)
	gin_util.Response(ctx, nil, err)
}

// UpdateKnowledgeStatus
//
//	@Tags			callback
//	@Summary		更新知识库状态
//	@Description	更新知识库状态
//	@Accept			json
//	@Produce		json
//	@Param			data	body		request.CallbackUpdateDocStatusReq	true	"更新知识库状态请求参数"
//	@Success		200		{object}	response.Response
//	@Router			/api/knowledge/status [post]
func UpdateKnowledgeStatus(ctx *gin.Context) {
	var req request.CallbackUpdateKnowledgeStatusReq
	if !gin_util.Bind(ctx, &req) {
		return
	}
	err := service.UpdateKnowledgeStatus(ctx, &req)
	gin_util.Response(ctx, nil, err)
}

// DocStatusInit
//
//	@Tags			callback
//	@Summary		将正在解析的文档设置为解析失败
//	@Description	将正在解析的文档设置为解析失败
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	response.Response{}
//	@Router			/api/doc_status_init [get]
func DocStatusInit(ctx *gin.Context) {
	resp, err := service.DocStatusInit(ctx, "", "")
	gin_util.Response(ctx, resp, err)
}

// GetDeployInfo
//
//	@Tags			callback
//	@Summary		获取Maas平台部署信息（模型扩展调用）
//	@Description	获取Maas平台部署信息（模型扩展调用）
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	response.Response{}
//	@Router			/api/deploy/info [get]
func GetDeployInfo(ctx *gin.Context) {
	resp, err := service.GetDeployInfo(ctx)
	gin_util.Response(ctx, resp, err)
}

// SelectKnowledgeInfoByName
//
//	@Tags			callback
//	@Summary		获取Maas平台知识库信息（模型扩展调用）
//	@Description	获取Maas平台知识库信息（模型扩展调用）
//	@Accept			json
//	@Produce		json
//	@Param			data	body		request.SearchKnowledgeInfoReq	true	"根据知识库名称请求参数"
//	@Success		200		{object}	response.Response{}
//	@Router			/api/category/info [get]
func SelectKnowledgeInfoByName(ctx *gin.Context) {
	var req request.SearchKnowledgeInfoReq
	if !gin_util.BindQuery(ctx, &req) {
		return
	}
	resp, err := service.SelectKnowledgeInfoByName(ctx, req.UserId, req.OrgId, &req)
	gin_util.Response(ctx, resp, err)
}

// GetWorkflowList
//
//	@Tags			callback
//	@Summary		根据userId和spaceId获取Workflow
//	@Description	根据userId和spaceId获取Workflow
//	@Accept			json
//	@Produce		json
//	@Param			userId	query		string	true	"获取工作流参数userId"
//	@Param			orgId	query		string	true	"获取工作流参数orgId"
//	@Success		200		{object}	response.Response
//	@Router			/workflow/list [get]
func GetWorkflowList(ctx *gin.Context) {
	var req request.GetWorkflowListReq
	if !gin_util.BindQuery(ctx, &req) {
		return
	}
	resp, err := service.GetAppList(ctx, req.UserId, req.OrgId, constant.AppTypeWorkflow)
	gin_util.Response(ctx, resp, err)
}

// GetWorkflowCustomTool
//
//	@Tags			callback
//	@Summary		获取自定义工具详情
//	@Description	获取自定义工具详情
//	@Accept			json
//	@Produce		json
//	@Param			customToolId	query		string	true	"customToolId"
//	@Success		200				{object}	response.Response{data=response.CustomToolDetail}
//	@Router			/workflow/tool/custom [get]
func GetWorkflowCustomTool(ctx *gin.Context) {
	resp, err := service.GetCustomTool(ctx, "", "", ctx.Query("customToolId"))
	gin_util.Response(ctx, resp, err)
}

// GetWorkflowSquareTool
//
//	@Tags			callback
//	@Summary		获取内置工具详情
//	@Description	获取内置工具详情
//	@Accept			json
//	@Produce		json
//	@Param			toolSquareId	query		string	true	"toolSquareId"
//	@Param			userID			query		string	true	"用户ID"
//	@Param			orgID			query		string	true	"组织ID"
//	@Success		200				{object}	response.Response{data=response.ToolSquareDetail}
//	@Router			/workflow/tool/square [get]
func GetWorkflowSquareTool(ctx *gin.Context) {
	resp, err := service.GetToolSquareDetail(ctx, "", "", ctx.Query("toolSquareId"))
	gin_util.Response(ctx, resp, err)
}

// SearchKnowledgeBase
//
//	@Tags			callback
//	@Summary		查询知识库列表（命中测试）
//	@Description	查询知识库列表（命中测试）
//	@Accept			json
//	@Produce		json
//	@Param			data	body		request.RagSearchKnowledgeBaseReq	true	"查询知识库列表请求参数"
//	@Success		200		{object}	response.Response
//	@Router			/rag/search-knowledge-base [post]
func SearchKnowledgeBase(ctx *gin.Context) {
	var req request.RagSearchKnowledgeBaseReq
	if !gin_util.Bind(ctx, &req) {
		return
	}
	resp, httpStatus := service.RagSearchKnowledgeBase(ctx, &req)
	gin_util.ResponseRawByte(ctx, httpStatus, resp)
}

// KnowledgeStreamSearch
//
//	@Tags			callback
//	@Summary		知识库流式问答
//	@Description	知识库流式问答
//	@Accept			json
//	@Produce		json
//	@Param			data	body		request.RagKnowledgeChatReq	true	"知识库流式问答请求参数"
//	@Success		200		{object}	response.Response
//	@Router			/rag/knowledge/stream/search [post]
func KnowledgeStreamSearch(ctx *gin.Context) {
	userId := ctx.GetHeader("X-uid")
	var req request.RagKnowledgeChatReq
	if !gin_util.Bind(ctx, &req) {
		return
	}
	req.UserId = userId
	err := service.KnowledgeStreamSearch(ctx, &req)
	if err != nil {
		resp, httpStatus := response.CommonRagKnowledgeError(err)
		gin_util.ResponseRawByte(ctx, httpStatus, resp)
		return
	}
}
