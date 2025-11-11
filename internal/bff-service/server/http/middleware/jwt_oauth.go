package middleware

import (
	"net/http"

	err_code "github.com/UnicomAI/wanwu/api/proto/err-code"
	gin_util "github.com/UnicomAI/wanwu/pkg/gin-util"
	jwt_util "github.com/UnicomAI/wanwu/pkg/jwt-util"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
)

func JWTAccess(ctx *gin.Context) {
	token, err := getJWTToken(ctx)
	if err != nil {
		gin_util.ResponseDetail(ctx, http.StatusUnauthorized, codes.Code(err_code.Code_BFFJWT), nil, err.Error())
		ctx.Abort()
		return
	}
	jwtAccessAuth(ctx, token)
}

func jwtAccessAuth(ctx *gin.Context, token string) {
	httpStatus := http.StatusUnauthorized
	claims, err := jwt_util.ParseAccessToken(token)
	if err != nil {
		gin_util.ResponseDetail(ctx, httpStatus, codes.Code(err_code.Code_BFFJWT), nil, err.Error())
		ctx.Abort()
		return
	}
	//验证sub，是否是access token
	if claims.Subject != jwt_util.ACCESS {
		gin_util.ResponseDetail(ctx, httpStatus, codes.Code(err_code.Code_BFFJWT), nil, "token subject错误")
		ctx.Abort()
		return
	}

	ctx.Set(gin_util.USER_ID, claims.UserID)
	ctx.Set(gin_util.OAUTH_CLIENT_ID, claims.ClientID)
	ctx.Set(gin_util.OAUTH_SCOPE, claims.Scope)
	ctx.Set(gin_util.CLAIMS, claims)
	ctx.Next()
}
