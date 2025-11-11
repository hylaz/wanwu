package openapi

import (
	"net/http"

	"github.com/UnicomAI/wanwu/internal/bff-service/model/request"
	"github.com/UnicomAI/wanwu/internal/bff-service/service"
	gin_util "github.com/UnicomAI/wanwu/pkg/gin-util"
	"github.com/gin-gonic/gin"
)

// Token
//
//	@Summary		授权码方式
//	@Description	授权码方式-获取Token
//	@Tags			OIDC
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Param			grant_type		formData	string	true	"授权类型"
//	@Param			code			formData	string	true	"授权码"
//	@Param			client_id		formData	string	true	"Client ID"
//	@Param			redirect_uri	formData	string	true	"回调地址"
//	@Param			client_secret	formData	string	true	"备案密钥"
//	@Success		200				{object}	response.TokenResponse
//	@Router			/oauth/code/token [post]
func Token(ctx *gin.Context) {
	var req request.TokenRequest
	if !gin_util.BindForm(ctx, &req) {
		return
	}
	resp, err := service.Token(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err})
	}
	ctx.JSON(http.StatusOK, resp)
}

// Refresh
//
//	@Summary		刷新令牌
//	@Description	刷新令牌
//	@Tags			OIDC
//	@Accept			json
//	@Produce		json
//	@Param			data	body		request.RefreshRequest	true	"RefreshToken"
//	@Success		200		{object}	response.RefreshTokenResponse
//	@Router			/oauth/code/token/refresh [post]
func Refresh(ctx *gin.Context) {
	var req request.RefreshRequest
	if !gin_util.Bind(ctx, &req) {
		return
	}
	resp, err := service.Refresh(ctx, &req)
	gin_util.Response(ctx, resp, err)
}

// OauthConfig
//
//	@Summary		动态客户端发现配置
//	@Description	自动获取 OP 的配置信息
//	@Tags			OIDC
//	@Produce		json
//	@Success		200	{object}	response.OauthConfig
//	@Router			/.well-known/openid-configuration [get]
func OauthConfig(ctx *gin.Context) {
	resp, err := service.OauthConfig(ctx)
	if err != nil {
		gin_util.Response(ctx, resp, err)
	}
	ctx.JSON(http.StatusOK, resp)
}

// JWKS
//
//	@Summary		公钥获取链接
//	@Description	自动获取JWKS
//	@Tags			OIDC
//	@Produce		json
//	@Success		200	{object}	response.JWKS
//	@Router			/oauth/jwks [get]
func JWKS(ctx *gin.Context) {
	resp, err := service.JWKS(ctx)
	if err != nil {
		gin_util.Response(ctx, resp, err)
	}
	ctx.JSON(http.StatusOK, resp)
}

// OauthGetUserInfo
//
//	@Summary		OAuth获取用户信息
//	@Description	通过access token获取用户信息
//	@Tags			OIDC
//	@Produce		json
//	@Success		200	{object}	response.OAuthGetUserInfo
//	@Router			/oauth/userinfo [get]
func OAuthGetUserInfo(ctx *gin.Context) {
	userID := getUserID(ctx)
	resp, err := service.OAuthGetUserInfo(ctx, userID)
	gin_util.Response(ctx, resp, err)
}
