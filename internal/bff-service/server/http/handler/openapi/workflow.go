package openapi

import (
	"net/http"

	"github.com/UnicomAI/wanwu/internal/bff-service/model/request"
	"github.com/UnicomAI/wanwu/internal/bff-service/service"
	gin_util "github.com/UnicomAI/wanwu/pkg/gin-util"
	"github.com/gin-gonic/gin"
)

// WorkflowRun
//
//	@Tags			openapi
//	@Summary		工作流OpenAPI
//	@Description	工作流OpenAPI
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	response.Response
//	@Router			/workflow/run [post]
func WorkflowRun(ctx *gin.Context) {
	var body []byte
	if cb, ok := ctx.Get(gin.BodyBytesKey); ok {
		if cbb, ok := cb.([]byte); ok {
			body = cbb
		}
	}
	resp, err := service.OpenAPIWorkflowRun(ctx, getAppID(ctx), body)
	if err != nil {
		gin_util.Response(ctx, nil, err)
		return
	}
	_, err = ctx.Writer.Write(resp)
	if err != nil {
		gin_util.Response(ctx, nil, err)
		return
	}
	ctx.Set(gin_util.STATUS, http.StatusOK)
	ctx.Writer.Flush()
}

// CreateChatflowConversation
//
//	@Tags			openapi
//	@Summary		对话流创建对话OpenAPI
//	@Description	对话流创建对话OpenAPI
//	@Accept			json
//	@Produce		json
//	@Param			data	body		request.OpenAPIChatflowCreateConversationRequest	true	"请求参数"
//	@Success		200		{object}	response.Response{data=response.OpenAPIChatflowCreateConversationResponse}
//	@Router			/chatflow/conversation [post]
func CreateChatflowConversation(ctx *gin.Context) {
	var req request.OpenAPIChatflowCreateConversationRequest
	if !gin_util.Bind(ctx, &req) {
		return
	}
	userID := getUserID(ctx)
	orgID := getOrgID(ctx)
	appID := getAppID(ctx)

	resp, err := service.CreateChatflowConversation(ctx, userID, orgID, appID, req.ConversationName)
	if err != nil {
		gin_util.Response(ctx, nil, err)
		return
	}
	gin_util.Response(ctx, resp, nil)
}

// ChatflowChat
//
//	@Tags			openapi
//	@Summary		对话流OpenAPI
//	@Description	对话流OpenAPI
//	@Accept			json
//	@Produce		json
//	@Param			data	body		request.OpenAPIChatflowChatRequest	true	"请求参数"
//	@Success		200		{object}	response.Response
//	@Router			/chatflow/chat [post]
func ChatflowChat(ctx *gin.Context) {
	var req request.OpenAPIChatflowChatRequest
	if !gin_util.Bind(ctx, &req) {
		return
	}
	userID := getUserID(ctx)
	orgID := getOrgID(ctx)
	appID := getAppID(ctx)

	// 流式处理 - 直接操作响应流
	err := service.ChatflowChat(ctx, userID, orgID, appID, req.ConversationId, req.Query)
	if err != nil {
		gin_util.Response(ctx, nil, err)
		return
	}
}

// WorkflowFileUpload
//
//	@Tags			openapi
//	@Summary		工作流OpenAPI文件上传
//	@Description	工作流OpenAPI文件上传
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			file	formData	file	true	"文件"
//	@Success		200		{object}	string
//	@Success		400		{object}	response.Response
//	@Router			/workflow/file/upload [post]
func WorkflowFileUpload(ctx *gin.Context) {
	resp, err := service.OpenAPIWorkflowFileUpload(ctx)
	if err != nil {
		gin_util.Response(ctx, nil, err)
		return
	}
	ctx.String(http.StatusOK, resp)
}
