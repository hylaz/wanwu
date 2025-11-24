package app

import (
	"context"

	app_service "github.com/UnicomAI/wanwu/api/proto/app-service"
	errs "github.com/UnicomAI/wanwu/api/proto/err-code"
	"github.com/UnicomAI/wanwu/internal/app-service/client/model"
	"github.com/UnicomAI/wanwu/pkg/util"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Service) GenApiKey(ctx context.Context, req *app_service.GenApiKeyReq) (*app_service.ApiKeyInfo, error) {
	apiKey, err := s.cli.GenApiKey(ctx, req.UserId, req.OrgId, req.AppId, req.AppType, util.GenApiUUID())
	if err != nil {
		return nil, errStatus(errs.Code_AppApikey, err)
	}
	return toProtoApiKey(apiKey), nil
}

func (s *Service) GetApiKeyList(ctx context.Context, req *app_service.GetApiKeyListReq) (*app_service.ApiKeyInfoList, error) {
	apiKeyList, err := s.cli.GetApiKeyList(ctx, req.UserId, req.OrgId, req.AppId, req.AppType)
	if err != nil {
		return nil, errStatus(errs.Code_AppApikey, err)
	}
	ret := &app_service.ApiKeyInfoList{
		Total: int64(len(apiKeyList)),
	}
	for _, apiKey := range apiKeyList {
		ret.Info = append(ret.Info, toProtoApiKey(apiKey))
	}
	return ret, nil
}

func (s *Service) DelApiKey(ctx context.Context, req *app_service.DelApiKeyReq) (*emptypb.Empty, error) {
	err := s.cli.DelApiKey(ctx, util.MustU32(req.ApiId))
	if err != nil {
		return nil, errStatus(errs.Code_AppApikey, err)
	}
	return &emptypb.Empty{}, nil
}

func (s *Service) GetApiKeyByKey(ctx context.Context, req *app_service.GetApiKeyByKeyReq) (*app_service.ApiKeyInfo, error) {
	apiKey, err := s.cli.GetApiKeyByKey(ctx, req.ApiKey)
	if err != nil {
		return nil, errStatus(errs.Code_AppApikey, err)
	}
	return toProtoApiKey(apiKey), nil
}

// --- internal ---

func toProtoApiKey(apiKey *model.ApiKey) *app_service.ApiKeyInfo {
	return &app_service.ApiKeyInfo{
		ApiId:     util.Int2Str(apiKey.ID),
		ApiKey:    apiKey.ApiKey,
		UserId:    apiKey.UserID,
		OrgId:     apiKey.OrgID,
		AppId:     apiKey.AppID,
		AppType:   apiKey.AppType,
		CreatedAt: apiKey.CreatedAt,
	}
}
