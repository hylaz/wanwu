package app

import (
	"context"

	app_service "github.com/UnicomAI/wanwu/api/proto/app-service"
	errs "github.com/UnicomAI/wanwu/api/proto/err-code"
	"github.com/UnicomAI/wanwu/internal/app-service/client/model"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Service) GetConversationByID(ctx context.Context, req *app_service.GetConversationByIDReq) (*app_service.ConversationInfo, error) {
	conversation, err := s.cli.GetConversationByID(ctx, req.ConversionId)
	if err != nil {
		return nil, errStatus(errs.Code_AppConversation, err)
	}
	return toProtoConversation(conversation), nil
}

func (s *Service) CreateConversation(ctx context.Context, req *app_service.CreateConversationReq) (*emptypb.Empty, error) {
	err := s.cli.CreateConversation(ctx, req.UserId, req.OrgId, req.AppId, req.AppType, req.ConversationId, req.ConversationName)
	if err != nil {
		return nil, errStatus(errs.Code_AppConversation, err)
	}
	return &emptypb.Empty{}, nil
}

func (s *Service) GetChatflowApplication(ctx context.Context, req *app_service.GetChatflowApplicationReq) (*app_service.ChatflowApplicationInfo, error) {
	ret, err := s.cli.GetChatflowApplication(ctx, req.OrgId, req.UserId, req.WorkflowId)
	if err != nil {
		return nil, errStatus(errs.Code_AppConversation, err)
	}
	return &app_service.ChatflowApplicationInfo{
		ApplicationId: ret.ApplicationID,
		UserId:        ret.UserID,
		OrgId:         ret.OrgID,
		CreatedAt:     ret.CreatedAt,
		UpdatedAt:     ret.UpdatedAt,
	}, nil
}

func (s *Service) GetChatflowByApplicationID(ctx context.Context, req *app_service.GetChatflowByApplicationIDReq) (*app_service.ChatflowApplicationInfo, error) {
	ret, err := s.cli.GetChatflowApplicationByApplicationID(ctx, req.OrgId, req.UserId, req.ApplicationId)
	if err != nil {
		return nil, errStatus(errs.Code_AppConversation, err)
	}
	return &app_service.ChatflowApplicationInfo{
		WorkflowId: ret.WorkflowID,
		UserId:     ret.UserID,
		OrgId:      ret.OrgID,
		CreatedAt:  ret.CreatedAt,
		UpdatedAt:  ret.UpdatedAt,
	}, nil
}

func (s *Service) CreateChatflowApplication(ctx context.Context, req *app_service.CreateChatflowApplicationReq) (*emptypb.Empty, error) {
	err := s.cli.CreateChatflowApplication(ctx, req.OrgId, req.UserId, req.WorkflowId, req.ApplicationId)
	if err != nil {
		return nil, errStatus(errs.Code_AppConversation, err)
	}
	return &emptypb.Empty{}, nil
}

// --- internal ---
func toProtoConversation(record *model.AppConversation) *app_service.ConversationInfo {
	return &app_service.ConversationInfo{
		ConversationId:   record.ConversationID,
		ConversationName: record.ConversationName,
		AppId:            record.AppID,
		AppType:          record.AppType,
		UserId:           record.UserID,
		OrgId:            record.OrgID,
		CreatedAt:        record.CreatedAt,
		UpdatedAt:        record.UpdatedAt,
	}
}
