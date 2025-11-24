package service

import (
	"net/url"

	app_service "github.com/UnicomAI/wanwu/api/proto/app-service"
	err_code "github.com/UnicomAI/wanwu/api/proto/err-code"
	"github.com/UnicomAI/wanwu/internal/bff-service/config"
	"github.com/UnicomAI/wanwu/internal/bff-service/model/request"
	"github.com/UnicomAI/wanwu/internal/bff-service/model/response"
	"github.com/UnicomAI/wanwu/pkg/constant"
	grpc_util "github.com/UnicomAI/wanwu/pkg/grpc-util"
	"github.com/UnicomAI/wanwu/pkg/util"
	"github.com/gin-gonic/gin"
)

func GetApiBaseUrl(ctx *gin.Context, req request.GetApiBaseUrlRequest) (string, error) {
	if req.AppType == constant.AppTypeWorkflow {
		apiBaseUrl, err := url.JoinPath(config.Cfg().Server.ApiBaseUrl, "/openapi/v1", req.AppType, "/run")
		if err != nil {
			return "", grpc_util.ErrorStatus(err_code.Code_BFFGeneral, err.Error())
		}
		return apiBaseUrl, nil
	}
	apiBaseUrl, err := url.JoinPath(config.Cfg().Server.ApiBaseUrl, "/openapi/v1", req.AppType, "/chat")
	if err != nil {
		return "", grpc_util.ErrorStatus(err_code.Code_BFFGeneral, err.Error())
	}
	return apiBaseUrl, nil
}

func GetApiKeyByKey(ctx *gin.Context, apiKey string) (*app_service.ApiKeyInfo, error) {
	return app.GetApiKeyByKey(ctx.Request.Context(), &app_service.GetApiKeyByKeyReq{ApiKey: apiKey})
}

func GenApiKey(ctx *gin.Context, userId, orgId string, req request.GenApiKeyRequest) (*response.ApiResponse, error) {
	key, err := app.GenApiKey(ctx.Request.Context(), &app_service.GenApiKeyReq{
		AppId:   req.AppId,
		AppType: req.AppType,
		UserId:  userId,
		OrgId:   orgId,
	})
	if err != nil {
		return nil, err
	}
	return &response.ApiResponse{
		ApiID:     key.ApiId,
		ApiKey:    key.ApiKey,
		CreatedAt: util.Time2Str(key.CreatedAt),
	}, nil
}

func DelApiKey(ctx *gin.Context, req request.DelApiKeyRequest) error {
	_, err := app.DelApiKey(ctx.Request.Context(), &app_service.DelApiKeyReq{
		ApiId: req.ApiId,
	})
	if err != nil {
		return err
	}
	return nil
}

func GetApiKeyList(ctx *gin.Context, userId string, req request.GetApiKeyListRequest) ([]*response.ApiResponse, error) {
	apiKeyList, err := app.GetApiKeyList(ctx.Request.Context(), &app_service.GetApiKeyListReq{
		AppId:   req.AppId,
		AppType: req.AppType,
		UserId:  userId,
	})
	if err != nil {
		return nil, err
	}
	var apiRes []*response.ApiResponse
	for _, apiKeyInfo := range apiKeyList.Info {
		apiRes = append(apiRes, toApiResp(apiKeyInfo))
	}
	return apiRes, nil
}

func toApiResp(apiKeyInfo *app_service.ApiKeyInfo) *response.ApiResponse {
	return &response.ApiResponse{
		ApiID:     apiKeyInfo.ApiId,
		ApiKey:    apiKeyInfo.ApiKey,
		CreatedAt: util.Time2Str(apiKeyInfo.CreatedAt),
	}
}
