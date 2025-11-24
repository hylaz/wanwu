package perm

import (
	"context"

	errs "github.com/UnicomAI/wanwu/api/proto/err-code"
	perm_service "github.com/UnicomAI/wanwu/api/proto/perm-service"
	"github.com/UnicomAI/wanwu/internal/iam-service/client"
	"github.com/UnicomAI/wanwu/internal/iam-service/config"
	grpc_util "github.com/UnicomAI/wanwu/pkg/grpc-util"
	"github.com/UnicomAI/wanwu/pkg/util"
	"github.com/gromitlee/access/pkg/perm"
)

type Service struct {
	perm_service.UnimplementedPermServiceServer
	cli client.IClient
}

func NewService(cli client.IClient) *Service {
	return &Service{
		cli: cli,
	}
}

func (s *Service) CheckUserEnable(ctx context.Context, req *perm_service.CheckUserEnableReq) (*perm_service.CheckUserEnableResp, error) {
	needLogin, language, lastUpdatePasswordAt, err := s.cli.CheckUserOK(ctx, util.MustU32(req.UserId), util.MustI64(req.GenTokenAt))
	if err != nil {
		if needLogin {
			return nil, errStatus(errs.Code_PermRBACReLogin, err)
		}
		return nil, errStatus(errs.Code_PermRBAC, err)
	}
	return &perm_service.CheckUserEnableResp{Language: language, LastUpdatePasswordAt: lastUpdatePasswordAt}, nil
}

func (s *Service) CheckUserPerm(ctx context.Context, req *perm_service.CheckUserPermReq) (*perm_service.CheckUserPermResp, error) {
	var oneOfPerms []perm.Perm
	for _, p := range req.OneOfPerms {
		oneOfPerms = append(oneOfPerms, perm.Perm{Obj: perm.Obj(p)})
	}
	needLogin, isAdmin, language, lastUpdatePasswordAt, err := s.cli.CheckUserPerm(ctx, util.MustU32(req.UserId), util.MustI64(req.GenTokenAt), util.MustU32(req.OrgId), oneOfPerms)
	if err != nil {
		if needLogin {
			return nil, errStatus(errs.Code_PermRBACReLogin, err)
		}
		return nil, errStatus(errs.Code_PermRBAC, err)
	}
	return &perm_service.CheckUserPermResp{IsAdmin: isAdmin,
		IsSystem:             util.MustU32(req.OrgId) == config.TopOrgID(),
		Language:             language,
		LastUpdatePasswordAt: lastUpdatePasswordAt}, nil
}

func errStatus(code errs.Code, status *errs.Status) error {
	return grpc_util.ErrorStatusWithKey(code, status.TextKey, status.Args...)
}
