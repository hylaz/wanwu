package v1

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/UnicomAI/wanwu/internal/bff-service/model/request"
	"github.com/UnicomAI/wanwu/internal/bff-service/service"
	gin_util "github.com/UnicomAI/wanwu/pkg/gin-util"
	"github.com/gin-gonic/gin"
)

// Authorize
//
//	@Summary		授权码方式
//	@Description	授权码方式-获取授权码
//	@Tags			oauth
//	@Accept			json
//	@Produce		json
//	@Param			client_id		query		string	true	"备案ID"
//	@Param			redirect_uri	query		string	true	"重定向URI"
//	@Param			response_type	query		string	true	"响应类型"
//	@Param			scope			query		string	false	"权限范围"
//	@Param			state			query		string	true	"状态参数"
//	@Success		302				{string}	string	"重定向到指定URI"
//	@Router			/oauth/code/authorize [get]
func Authorize(ctx *gin.Context) {
	var req request.AuthRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		gin_util.Response(ctx, nil, err)
		return
	}
	callback, authCode, err := service.Authorize(ctx, &req)
	if err != nil {
		gin_util.Response(ctx, nil, err)
		return
	}
	redirectURI := fmt.Sprintf(
		"%s?code=%s&state=%s",
		callback,
		url.QueryEscape(authCode),
		url.QueryEscape(req.State), // 对state也进行编码
	)
	ctx.Redirect(http.StatusFound, redirectURI)
}

// CreateOauthApp
//
//	@Tags			oauth
//	@Summary		创建OAuth应用
//	@Description	创建新的OAuth应用
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Param			data	body		request.CreateOauthAppReq	true	"OAuth应用创建请求参数"
//	@Success		200		{object}	response.Response
//	@Router			/oauth/app [post]
func CreateOauthApp(ctx *gin.Context) {
	var req request.CreateOauthAppReq
	if !gin_util.Bind(ctx, &req) {
		return
	}
	gin_util.Response(ctx, nil, service.CreateOauthApp(ctx, getUserID(ctx), &req))
}

// DeleteOauthApp
//
//	@Tags			oauth
//	@Summary		删除OAuth应用
//	@Description	删除指定的OAuth应用
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Param			data	body		request.DeleteOauthAppReq	true	"OAuth应用ID"
//	@Success		200		{object}	response.Response
//	@Router			/oauth/app [delete]
func DeleteOauthApp(ctx *gin.Context) {
	var req request.DeleteOauthAppReq
	if !gin_util.Bind(ctx, &req) {
		return
	}
	gin_util.Response(ctx, nil, service.DeleteOauthApp(ctx, &req))
}

// UpdateOauthApp
//
//	@Tags			oauth
//	@Summary		更新OAuth应用
//	@Description	更新OAuth应用信息
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Param			data	body		request.UpdateOauthAppReq	true	"OAuth应用更新请求参数"
//	@Success		200		{object}	response.Response
//	@Router			/oauth/app [put]
func UpdateOauthApp(ctx *gin.Context) {
	var req request.UpdateOauthAppReq
	if !gin_util.Bind(ctx, &req) {
		return
	}
	gin_util.Response(ctx, nil, service.UpdateOauthApp(ctx, &req))
}

// GetOauthAppList
//
//	@Tags			oauth
//	@Summary		获取OAuth应用列表
//	@Description	获取OAuth应用分页列表
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	response.Response{data=response.ListResult{list=[]response.OauthAppInfo}}
//	@Router			/oauth/app/list [get]
func GetOauthAppList(ctx *gin.Context) {
	resp, err := service.GetOauthAppList(ctx, getUserID(ctx))
	gin_util.Response(ctx, resp, err)

}

// UpdateOauthAppStatus
//
//	@Tags			oauth
//	@Summary		更新OAuth应用状态
//	@Description	启用或禁用OAuth应用
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Param			data	body		request.UpdateOauthAppStatusReq	true	"OAuth应用状态更新请求参数"
//	@Success		200		{object}	response.Response
//	@Router			/oauth/app/status [put]
func UpdateOauthAppStatus(ctx *gin.Context) {
	var req request.UpdateOauthAppStatusReq
	if !gin_util.Bind(ctx, &req) {
		return
	}
	gin_util.Response(ctx, nil, service.UpdateOauthAppStatus(ctx, &req))
}
