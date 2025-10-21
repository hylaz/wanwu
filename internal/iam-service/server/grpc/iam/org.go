package iam

import (
	"context"
	"strconv"

	errs "github.com/UnicomAI/wanwu/api/proto/err-code"
	iam_service "github.com/UnicomAI/wanwu/api/proto/iam-service"
	"github.com/UnicomAI/wanwu/internal/iam-service/client/model"
	"github.com/UnicomAI/wanwu/internal/iam-service/client/orm"
	"github.com/UnicomAI/wanwu/pkg/util"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Service) GetOrgSelect(ctx context.Context, req *iam_service.GetOrgSelectReq) (*iam_service.Select, error) {
	orgs, err := s.cli.SelectOrgs(ctx, util.MustU32(req.UserId))
	if err != nil {
		return nil, errStatus(errs.Code_IAMOrg, err)
	}
	return &iam_service.Select{Selects: toIDNames(orgs)}, nil
}

func (s *Service) GetOrgList(ctx context.Context, req *iam_service.GetOrgListReq) (*iam_service.GetOrgListResp, error) {
	orgs, count, err := s.cli.GetOrgs(ctx, util.MustU32(req.ParentId), req.Name, toOffset(req), req.PageSize)
	if err != nil {
		return nil, errStatus(errs.Code_IAMOrg, err)
	}
	resp := &iam_service.GetOrgListResp{
		Total:    count,
		PageNo:   req.PageNo,
		PageSize: req.PageSize,
	}
	for _, org := range orgs {
		resp.Orgs = append(resp.Orgs, toOrgInfo(org))
	}
	return resp, nil
}

func (s *Service) GetOrgInfo(ctx context.Context, req *iam_service.GetOrgInfoReq) (*iam_service.OrgInfo, error) {
	org, err := s.cli.GetOrg(ctx, util.MustU32(req.OrgId))
	if err != nil {
		return nil, errStatus(errs.Code_IAMOrg, err)
	}
	return toOrgInfo(org), nil
}

func (s *Service) GetOrgByOrgIDs(ctx context.Context, req *iam_service.GetOrgByOrgIDsReq) (*iam_service.GetOrgByOrgIDsResp, error) {
	var orgIDs []uint32
	for _, userID := range req.OrgIds {
		orgIDs = append(orgIDs, util.MustU32(userID))
	}
	orgs, err := s.cli.GetOrgByOrgIDs(ctx, orgIDs)
	if err != nil {
		return nil, errStatus(errs.Code_IAMOrg, err)
	}
	return &iam_service.GetOrgByOrgIDsResp{
		Orgs: toIDNames(orgs),
	}, nil
}

func (s *Service) GetOrgAndSubOrgSelectByUser(ctx context.Context, req *iam_service.GetOrgAndSubOrgSelectByUserReq) (*iam_service.GetOrgAndSubOrgSelectByUserResp, error) {
	orgs, err := s.cli.GetOrgAndSubOrgSelectByUser(ctx, util.MustU32(req.UserId), util.MustU32(req.OrgId))
	if err != nil {
		return nil, errStatus(errs.Code_IAMOrg, err)
	}
	return &iam_service.GetOrgAndSubOrgSelectByUserResp{
		Orgs: toIDNames(orgs),
	}, nil
}

func (s *Service) CreateOrg(ctx context.Context, req *iam_service.CreateOrgReq) (*iam_service.IDName, error) {
	orgID, err := s.cli.CreateOrg(ctx, &model.Org{
		Status:    true,
		CreatorID: util.MustU32(req.CreatorId),
		ParentID:  util.MustU32(req.ParentId),
		Name:      req.Name,
		Remark:    req.Remark,
	})
	if err != nil {
		return nil, errStatus(errs.Code_IAMOrg, err)
	}
	return &iam_service.IDName{Id: strconv.Itoa(int(orgID)), Name: req.Name}, nil
}

func (s *Service) UpdateOrg(ctx context.Context, req *iam_service.UpdateOrgReq) (*emptypb.Empty, error) {
	if err := s.cli.UpdateOrg(ctx, &model.Org{
		ID:       util.MustU32(req.OrgId),
		ParentID: util.MustU32(req.ParentId),
		Name:     req.Name,
		Remark:   req.Remark,
	}); err != nil {
		return nil, errStatus(errs.Code_IAMOrg, err)
	}
	return &emptypb.Empty{}, nil
}

func (s *Service) DeleteOrg(ctx context.Context, req *iam_service.DeleteOrgReq) (*emptypb.Empty, error) {
	if err := s.cli.DeleteOrg(ctx, util.MustU32(req.OrgId)); err != nil {
		return nil, errStatus(errs.Code_IAMOrg, err)
	}
	return &emptypb.Empty{}, nil
}

func (s *Service) ChangeOrgStatus(ctx context.Context, req *iam_service.ChangeOrgStatusReq) (*emptypb.Empty, error) {
	if err := s.cli.ChangeOrgStatus(ctx, util.MustU32(req.OrgId), req.Status); err != nil {
		return nil, errStatus(errs.Code_IAMOrg, err)
	}
	return &emptypb.Empty{}, nil
}

func (s *Service) AddOrgUser(ctx context.Context, req *iam_service.AddOrgUserReq) (*emptypb.Empty, error) {
	if err := s.cli.AddOrgUser(ctx, util.MustU32(req.OrgId), util.MustU32(req.UserId), util.MustU32(req.RoleId)); err != nil {
		return nil, errStatus(errs.Code_IAMOrg, err)
	}
	return &emptypb.Empty{}, nil
}

func (s *Service) RemoveOrgUser(ctx context.Context, req *iam_service.RemoveOrgUserReq) (*emptypb.Empty, error) {
	if err := s.cli.RemoveOrgUser(ctx, util.MustU32(req.OrgId), util.MustU32(req.UserId)); err != nil {
		return nil, errStatus(errs.Code_IAMOrg, err)
	}
	return &emptypb.Empty{}, nil
}

// --- internal function ---

func toOrgInfo(org *orm.OrgInfo) *iam_service.OrgInfo {
	return &iam_service.OrgInfo{
		OrgId:     strconv.Itoa(int(org.ID)),
		Name:      org.Name,
		Remark:    org.Remark,
		Status:    org.Status,
		CreatedAt: org.CreatedAt,
		Creator:   toIDName(org.Creator),
	}
}
