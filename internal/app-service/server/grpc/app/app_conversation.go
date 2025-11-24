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
