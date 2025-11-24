package service

import (
	"fmt"

	"github.com/UnicomAI/wanwu/api/proto/common"
	err_code "github.com/UnicomAI/wanwu/api/proto/err-code"
	"github.com/UnicomAI/wanwu/internal/bff-service/model/request"
	"github.com/UnicomAI/wanwu/internal/bff-service/model/response"
	bff_util "github.com/UnicomAI/wanwu/internal/bff-service/pkg/util"
	grpc_util "github.com/UnicomAI/wanwu/pkg/grpc-util"
	mp "github.com/UnicomAI/wanwu/pkg/model-provider"
	"github.com/UnicomAI/wanwu/pkg/util"
	"github.com/gin-gonic/gin"
)

// --- app breif ---

func appBriefProto2Model(ctx *gin.Context, appBrief *common.AppBrief) response.AppBriefInfo {
	return response.AppBriefInfo{
		AppId:     appBrief.AppId,
		AppType:   appBrief.AppType,
		UniqueId:  bff_util.ConcatAssistantToolUniqueId(appBrief.AppType, appBrief.AppId),
		Avatar:    cacheAppAvatar(ctx, appBrief.AvatarPath, appBrief.AppType),
		Name:      appBrief.Name,
		Desc:      appBrief.Desc,
		CreatedAt: util.Time2Str(appBrief.CreatedAt),
		UpdatedAt: util.Time2Str(appBrief.UpdatedAt),
	}
}

// --- app brief config ---

func appBriefConfigProto2Model(ctx *gin.Context, appBrief *common.AppBriefConfig, appType string) request.AppBriefConfig {
	return request.AppBriefConfig{
		Avatar: cacheAppAvatar(ctx, appBrief.AvatarPath, appType),
		Name:   appBrief.Name,
		Desc:   appBrief.Desc,
	}
}

func appBriefConfigModel2Proto(appBrief request.AppBriefConfig) *common.AppBriefConfig {
	return &common.AppBriefConfig{
		Name:       appBrief.Name,
		Desc:       appBrief.Desc,
		AvatarPath: appBrief.Avatar.Key,
	}
}

// --- app model config ---

func appModelConfigProto2Model(appModel *common.AppModelConfig, displayName string) (request.AppModelConfig, error) {
	ret := request.AppModelConfig{
		Provider:    appModel.Provider,
		Model:       appModel.Model,
		ModelId:     appModel.ModelId,
		ModelType:   appModel.ModelType,
		DisplayName: displayName,
	}
	modelParams, _, err := mp.ToModelParams(appModel.Provider, appModel.ModelType, appModel.Config)
	if err != nil {
		return ret, grpc_util.ErrorStatusWithKey(err_code.Code_BFFGeneral, "bff_model_params", fmt.Sprintf("model %v get app model config err: %v", appModel.ModelId, err))
	}
	ret.Config = modelParams
	return ret, nil
}

func appModelConfigModel2Proto(appModel request.AppModelConfig) (*common.AppModelConfig, error) {
	configStr, err := appModel.ConfigString()
	if err != nil {
		return nil, grpc_util.ErrorStatusWithKey(err_code.Code_BFFGeneral, "bff_model_config_string", fmt.Sprintf("model %v get app model config err: %v", appModel.ModelId, err))
	}
	return &common.AppModelConfig{
		Provider:  appModel.Provider,
		Model:     appModel.Model,
		ModelId:   appModel.ModelId,
		ModelType: appModel.ModelType,
		Config:    configStr,
	}, nil
}
