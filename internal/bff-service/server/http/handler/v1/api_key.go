package v1

import (
	"github.com/UnicomAI/wanwu/internal/bff-service/model/request"
	"github.com/UnicomAI/wanwu/internal/bff-service/service"
	gin_util "github.com/UnicomAI/wanwu/pkg/gin-util"
	"github.com/gin-gonic/gin"
)

// GetApiBaseUrl
//
//	@Tags			app.key
//	@Summary		获取Api根地址
//	@Description	获取Api根地址
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Param			data	query		request.GetApiBaseUrlRequest	true	"获取Api根地址参数"
//	@Success		200		{object}	response.Response{data=string}
//	@Router			/appspace/app/url [get]
func GetApiBaseUrl(ctx *gin.Context) {
	var req request.GetApiBaseUrlRequest
	if !gin_util.BindQuery(ctx, &req) {
		return
	}
	resp, err := service.GetApiBaseUrl(ctx, req)
	gin_util.Response(ctx, resp, err)
}

// GenApiKey
//
//	@Tags			app.key
//	@Summary		生成ApiKey
//	@Description	生成ApiKey
//	@Accept			json
//	@Produce		json
//	@Param			data	body		request.GenApiKeyRequest	true	"生成ApiKey参数"
//	@Success		200		{object}	response.Response{data=response.ApiResponse}
//	@Router			/appspace/app/key [post]
func GenApiKey(ctx *gin.Context) {
	var req request.GenApiKeyRequest
	if !gin_util.Bind(ctx, &req) {
		return
	}
	resp, err := service.GenApiKey(ctx, getUserID(ctx), getOrgID(ctx), req)
	gin_util.Response(ctx, resp, err)
}

// DelApiKey
//
//	@Tags			app.key
//	@Summary		删除ApiKey
//	@Description	删除ApiKey
//	@Accept			json
//	@Produce		json
//	@Param			data	body		request.DelApiKeyRequest	true	"删除Apikey参数"
//	@Success		200		{object}	response.Response
//	@Router			/appspace/app/key [delete]
func DelApiKey(ctx *gin.Context) {
	var req request.DelApiKeyRequest
	if !gin_util.Bind(ctx, &req) {
		return
	}
	err := service.DelApiKey(ctx, req)
	gin_util.Response(ctx, nil, err)
}

// GetApiKeyList
//
//	@Tags			app.key
//	@Summary		获取ApiKey
//	@Description	获取ApiKey
//	@Accept			json
//	@Produce		json
//	@Param			data	query		request.GetApiKeyListRequest	true	"获取ApiKey参数"
//	@Success		200		{object}	response.Response{data=[]response.ApiResponse}
//	@Router			/appspace/app/key/list [get]
func GetApiKeyList(ctx *gin.Context) {
	var req request.GetApiKeyListRequest
	if !gin_util.BindQuery(ctx, &req) {
		return
	}
	resp, err := service.GetApiKeyList(ctx, getUserID(ctx), req)
	gin_util.Response(ctx, resp, err)
}
