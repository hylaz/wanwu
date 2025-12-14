package callback

import (
	"github.com/UnicomAI/wanwu/internal/bff-service/model/request"
	"github.com/UnicomAI/wanwu/internal/bff-service/model/response"
	"github.com/UnicomAI/wanwu/internal/bff-service/pkg/util"
	"github.com/UnicomAI/wanwu/internal/bff-service/service"
	gin_util "github.com/UnicomAI/wanwu/pkg/gin-util"
	mp_common "github.com/UnicomAI/wanwu/pkg/model-provider/mp-common"
	"github.com/gin-gonic/gin"
)

//	@title		AI Agent Productivity Platform - Callback
//	@version	v0.0.1

//	@BasePath	/callback/v1

// FileUrlConvertBase64
//
//	@Tags		callback
//	@Summary	文件Url转换为base64
//	@Accept		json
//	@Produce	json
//	@Param		data	body		request.FileUrlConvertBase64Req	true	"文件Url转换base64请求参数"
//	@Success	200		{object}	response.Response{data=string}
//	@Router		/file/url/base64 [post]
func FileUrlConvertBase64(ctx *gin.Context) {
	var req request.FileUrlConvertBase64Req
	if !gin_util.Bind(ctx, &req) {
		return
	}
	resp, err := service.FileUrlConvertBase64(ctx, &req)
	gin_util.Response(ctx, resp, err)
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

// SearchQABase
//
//	@Tags			callback
//	@Summary		查询问答列表（命中测试）
//	@Description	查询问答列表（命中测试）
//	@Accept			json
//	@Produce		json
//	@Param			data	body		request.RagSearchQABaseReq	true	"查询知识库列表请求参数"
//	@Success		200		{object}	response.Response
//	@Router			/rag/search-qa-base [post]
func SearchQABase(ctx *gin.Context) {
	var req request.RagSearchQABaseReq
	if !gin_util.Bind(ctx, &req) {
		return
	}
	resp, httpStatus := service.RagSearchQABase(ctx, &req)
	gin_util.ResponseRawByte(ctx, httpStatus, resp)
}

// AudioBase64ConvertText
//
//	@Tags		callback
//	@Summary	语音文件（base64格式）转文本内置工具服务
//	@Accept		json
//	@Produce	json
//	@Param		file	formData	file	true	"语音文件 base64"
//	@Param		config	formData	string	true	"请求参数"
//	@Param		apiKey	formData	string	true	"api token"
//	@Success	200		{object}	response.Response{data=string}
//	@Router		/tool/builtin/asr [post]
func AudioBase64ConvertText(ctx *gin.Context) {
	var data request.AudioBase64ConvertTextReq
	if !gin_util.BindForm(ctx, &data) {
		return
	}
	file, err := util.Base64ToFileHeader(data.File, "audio")
	if err != nil {
		gin_util.Response(ctx, nil, err)
		return
	}
	asrServiceReq := &mp_common.AsrReq{
		Config: data.Config,
		File:   file,
		ApiKey: data.ApiKey,
	}
	service.AudioBase64ConvertText(ctx, asrServiceReq)
}
