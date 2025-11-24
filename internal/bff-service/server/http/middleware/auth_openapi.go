package middleware

import (
	"fmt"
	"net/http"
	"strings"

	err_code "github.com/UnicomAI/wanwu/api/proto/err-code"
	"github.com/UnicomAI/wanwu/internal/bff-service/service"
	gin_util "github.com/UnicomAI/wanwu/pkg/gin-util"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
)

func AuthOpenAPI(appType string) func(*gin.Context) {
	return func(ctx *gin.Context) {
		token, err := getApiKey(ctx)
		if err != nil {
			gin_util.ResponseDetail(ctx, http.StatusUnauthorized, codes.Code(err_code.Code_BFFAuth), nil, err.Error())
			ctx.Abort()
			return
		}
		apiKey, err := service.GetApiKeyByKey(ctx, token)
		if err != nil {
			gin_util.ResponseDetail(ctx, http.StatusUnauthorized, codes.Code(err_code.Code_BFFAuth), nil, err.Error())
			ctx.Abort()
			return
		}
		if apiKey.AppType != appType {
			gin_util.ResponseDetail(ctx, http.StatusUnauthorized, codes.Code(err_code.Code_BFFAuth), nil, "invalid appType")
			ctx.Abort()
			return
		}
		ctx.Set(gin_util.USER_ID, apiKey.UserId)
		ctx.Set(gin_util.X_ORG_ID, apiKey.OrgId)
		ctx.Set(gin_util.APP_ID, apiKey.AppId)
	}

}

func AuthOpenAPIByQuery(appType string) func(*gin.Context) {
	return func(ctx *gin.Context) {
		token, err := getApiKeyByQuery(ctx)
		if err != nil {
			gin_util.ResponseDetail(ctx, http.StatusUnauthorized, codes.Code(err_code.Code_BFFAuth), nil, err.Error())
			ctx.Abort()
			return
		}
		apiKey, err := service.GetApiKeyByKey(ctx, token)
		if err != nil {
			gin_util.ResponseDetail(ctx, http.StatusUnauthorized, codes.Code(err_code.Code_BFFAuth), nil, err.Error())
			ctx.Abort()
			return
		}
		if apiKey.AppType != appType {
			gin_util.ResponseDetail(ctx, http.StatusUnauthorized, codes.Code(err_code.Code_BFFAuth), nil, "invalid appType")
			ctx.Abort()
			return
		}
		ctx.Set(gin_util.USER_ID, apiKey.UserId)
		ctx.Set(gin_util.X_ORG_ID, apiKey.OrgId)
		ctx.Set(gin_util.APP_ID, apiKey.AppId)
	}

}

func getApiKey(ctx *gin.Context) (string, error) {
	authorization := ctx.Request.Header.Get("Authorization")
	if authorization != "" {
		tks := strings.Split(authorization, " ")
		if len(tks) > 1 && tks[0] == "Bearer" {
			return tks[1], nil
		} else {
			return "", fmt.Errorf("not Bearer token format")
		}
	} else {
		return "", fmt.Errorf("token is nil")
	}
}

func getApiKeyByQuery(ctx *gin.Context) (string, error) {
	key := ctx.Query("key")
	if key != "" {
		return key, nil
	} else {
		return "", fmt.Errorf("token is nil")
	}
}
