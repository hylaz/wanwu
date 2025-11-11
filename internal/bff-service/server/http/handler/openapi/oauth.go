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
	gin_util.Response(ctx, resp, err)
}
