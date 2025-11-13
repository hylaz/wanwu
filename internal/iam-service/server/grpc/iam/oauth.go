package iam

import (
	"context"

	errs "github.com/UnicomAI/wanwu/api/proto/err-code"
	iam_service "github.com/UnicomAI/wanwu/api/proto/iam-service"
	"github.com/UnicomAI/wanwu/internal/iam-service/client/model"
	"github.com/UnicomAI/wanwu/pkg/util"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Service) CreateOauthApp(ctx context.Context, req *iam_service.CreateOauthAppReq) (*emptypb.Empty, error) {
	err := s.cli.CreateOauthApp(ctx, &model.OauthApp{
		UserID:      util.MustU32(req.UserId),
		Name:        req.Name,
		Description: req.Desc,
		RedirectURI: req.RedirectUri,
	})
	if err != nil {
		return nil, errStatus(errs.Code_IAMOauth, err)
	}
	return &emptypb.Empty{}, nil
}

func (s *Service) DeleteOauthApp(ctx context.Context, req *iam_service.DeleteOauthAppReq) (*emptypb.Empty, error) {
	err := s.cli.DeleteOauthApp(ctx, req.ClientId)
	if err != nil {
		return nil, errStatus(errs.Code_IAMOauth, err)
	}
	return &emptypb.Empty{}, nil
}

func (s *Service) UpdateOauthApp(ctx context.Context, req *iam_service.UpdateOauthAppReq) (*emptypb.Empty, error) {
	err := s.cli.UpdateOauthApp(ctx, &model.OauthApp{
		ClientID:    req.ClientId,
		Name:        req.Name,
		Description: req.Desc,
		RedirectURI: req.RedirectUri,
	})
	if err != nil {
		return nil, errStatus(errs.Code_IAMOauth, err)
	}
	return &emptypb.Empty{}, nil
}

func (s *Service) GetOauthAppList(ctx context.Context, req *iam_service.GetOauthAppListReq) (*iam_service.OauthAppListResp, error) {
	apps, total, err := s.cli.GetOauthAppList(ctx, util.MustU32(req.UserId), req.Name, toOffset(req), req.PageSize)
	if err != nil {
		return nil, errStatus(errs.Code_IAMOauth, err)
	}
	resp := &iam_service.OauthAppListResp{
		Total:    total,
		PageNo:   req.PageNo,
		PageSize: req.PageSize,
	}
	for _, app := range apps {
		resp.Apps = append(resp.Apps, toOauthAppProto(app))
	}
	return resp, nil
}

func (s *Service) UpdateOauthAppStatus(ctx context.Context, req *iam_service.UpdateOauthAppStatusReq) (*emptypb.Empty, error) {
	err := s.cli.UpdateOauthAppStatus(ctx, req.ClientId, req.Status)
	if err != nil {
		return nil, errStatus(errs.Code_IAMOauth, err)
	}
	return &emptypb.Empty{}, nil
}

func (s *Service) GetOauthApp(ctx context.Context, req *iam_service.GetOauthAppReq) (*iam_service.OauthApp, error) {
	app, err := s.cli.GetOauthApp(ctx, req.ClientId)
	if err != nil {
		return nil, errStatus(errs.Code_IAMOauth, err)
	}
	return toOauthAppProto(app), nil
}

func toOauthAppProto(app *model.OauthApp) *iam_service.OauthApp {
	if app == nil {
		return nil
	}
	return &iam_service.OauthApp{
		ClientId:     app.ClientID,
		Name:         app.Name,
		Desc:         app.Description,
		RedirectUri:  app.RedirectURI,
		ClientSecret: app.ClientSecret,
		Status:       app.Status,
		CreatedAt:    app.CreatedAt,
		UpdatedAt:    app.UpdatedAt,
	}
}
