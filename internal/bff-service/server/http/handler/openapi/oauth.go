package openapi

import (
	"net/http"

	"github.com/UnicomAI/wanwu/internal/bff-service/model/request"
	"github.com/UnicomAI/wanwu/internal/bff-service/service"
	gin_util "github.com/UnicomAI/wanwu/pkg/gin-util"
	"github.com/gin-gonic/gin"
)

// OAuthAuthorize
//
//	@Tags			openapi.OIDC
//	@Summary		授权码方式
//	@Description	授权码方式-获取授权码
//	@Accept			json
//	@Produce		json
//	@Param			client_id		query		string	true	"备案ID"
//	@Param			redirect_uri	query		string	true	"重定向URI"
//	@Param			response_type	query		string	true	"响应类型"
//	@Param			scope			query		string	false	"权限范围"
//	@Param			state			query		string	true	"状态参数"
//	@Param			jwt_token		query		string	true	"token"
//	@Success		302				{string}	string	"重定向到指定URI"
//	@Router			/oauth/code/authorize [get]
func OAuthAuthorize(ctx *gin.Context) {
	var req request.OAuthRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		gin_util.Response(ctx, nil, err)
		return
	}
	callbackUri, err := service.OAuthAuthorize(ctx, &req)
	if err != nil {
		gin_util.Response(ctx, nil, err)
		return
	}

	ctx.Redirect(http.StatusFound, callbackUri)
}

// OAuthToken
//
//	@Tags			openapi.OIDC
//	@Summary		授权码方式
//	@Description	授权码方式-获取Token
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Param			grant_type		formData	string	true	"授权类型"
//	@Param			code			formData	string	true	"授权码"
//	@Param			client_id		formData	string	true	"Client ID"
//	@Param			redirect_uri	formData	string	true	"回调地址"
//	@Param			client_secret	formData	string	true	"备案密钥"
//	@Success		200				{object}	response.OAuthTokenResponse
//	@Router			/oauth/code/token [post]
func OAuthToken(ctx *gin.Context) {
	var req request.OAuthTokenRequest
	if !gin_util.BindForm(ctx, &req) {
		return
	}
	resp, err := service.OAuthToken(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err})
	}
	ctx.JSON(http.StatusOK, resp)
}

// OAuthRefresh
//
//	@Tags			openapi.OIDC
//	@Summary		刷新令牌
//	@Description	刷新令牌
//	@Accept			json
//	@Produce		json
//	@Param			data	body		request.OAuthRefreshRequest	true	"RefreshToken"
//	@Success		200		{object}	response.OAuthRefreshTokenResponse
//	@Router			/oauth/code/token/refresh [post]
func OAuthRefresh(ctx *gin.Context) {
	var req request.OAuthRefreshRequest
	if !gin_util.Bind(ctx, &req) {
		return
	}
	resp, err := service.OAuthRefresh(ctx, &req)
	gin_util.Response(ctx, resp, err)
}

// OAuthConfig
//
//	@Tags			openapi.OIDC
//	@Summary		动态客户端发现配置
//	@Description	自动获取 OP 的配置信息
//	@Produce		json
//	@Success		200	{object}	response.OAuthConfig
//	@Router			/.well-known/openid-configuration [get]
func OAuthConfig(ctx *gin.Context) {
	resp, err := service.OAuthConfig(ctx)
	if err != nil {
		gin_util.Response(ctx, resp, err)
	}
	ctx.JSON(http.StatusOK, resp)
}

// OAuthJWKS
//
//	@Tags			openapi.OIDC
//	@Summary		公钥获取链接
//	@Description	自动获取OAuthJWKS
//	@Produce		json
//	@Success		200	{object}	response.OAuthJWKS
//	@Router			/oauth/jwks [get]
func OAuthJWKS(ctx *gin.Context) {
	resp, err := service.OAuthJWKS(ctx)
	if err != nil {
		gin_util.Response(ctx, resp, err)
	}
	ctx.JSON(http.StatusOK, resp)
}

// OAuthGetUserInfo
//
//	@Tags			openapi.OIDC
//	@Summary		OAuth获取用户信息
//	@Description	通过access token获取用户信息
//	@Produce		json
//	@Success		200	{object}	response.OAuthGetUserInfo
//	@Router			/oauth/userinfo [get]
func OAuthGetUserInfo(ctx *gin.Context) {
	userID := getUserID(ctx)
	resp, err := service.OAuthGetUserInfo(ctx, userID)
	gin_util.Response(ctx, resp, err)
}

// OAuthLogin
//
//	@Tags			openapi.OIDC
//	@Summary		OAuth登录授权
//	@Description	返回OAuth登录页面
//	@Accept			json
//	@Produce		json
//	@Param			client_id		query		string	true	"备案ID"
//	@Param			redirect_uri	query		string	true	"重定向URI"
//	@Param			response_type	query		string	true	"响应类型"
//	@Param			scope			query		string	false	"权限范围"
//	@Param			state			query		string	true	"状态参数"
//	@Success		302				{string}	string	"重定向到指定URI"
//	@Router			/oauth/login [get]
func OAuthLogin(ctx *gin.Context) {
	var req request.OAuthLoginRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		gin_util.Response(ctx, nil, err)
		return
	}
	loginUri, err := service.OAuthLogin(ctx, &req)
	if err != nil {
		gin_util.Response(ctx, nil, err)
		return
	}

	ctx.Redirect(http.StatusFound, loginUri)
}
