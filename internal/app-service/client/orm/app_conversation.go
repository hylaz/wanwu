package orm

import (
	"context"
	"errors"

	errs "github.com/UnicomAI/wanwu/api/proto/err-code"
	"github.com/UnicomAI/wanwu/internal/app-service/client/model"
	"github.com/UnicomAI/wanwu/internal/app-service/client/orm/sqlopt"
	"gorm.io/gorm"
)

func (c *Client) GetConversationByID(ctx context.Context, conversationId string) (*model.AppConversation, *errs.Status) {
	var conversation model.AppConversation
	if err := sqlopt.SQLOptions(
		sqlopt.WithConversationID(conversationId),
	).Apply(c.db.WithContext(ctx)).First(&conversation).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, toErrStatus("app_conversation_not_found", conversationId)
		}
		return nil, toErrStatus("app_conversation_get", conversationId, err.Error())
	}
	return &conversation, nil
}

func (c *Client) CreateConversation(ctx context.Context, userId, orgId, appId, appType, conversationId, conversationName string) *errs.Status {
	err := sqlopt.SQLOptions(
		sqlopt.WithUserID(userId),
		sqlopt.WithOrgID(orgId),
		sqlopt.WithAppID(appId),
		sqlopt.WithAppType(appType),
		sqlopt.WithConversationID(conversationId),
	).Apply(c.db.WithContext(ctx)).First(&model.AppConversation{}).Error
	if err == nil {
		return toErrStatus("app_conversation_exist")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return toErrStatus("app_conversation_get", conversationId, err.Error())
	}
	conversation := &model.AppConversation{
		UserID:           userId,
		OrgID:            orgId,
		AppID:            appId,
		AppType:          appType,
		ConversationID:   conversationId,
		ConversationName: conversationName,
	}
	if err := c.db.WithContext(ctx).Create(conversation).Error; err != nil {
		return toErrStatus("app_conversation_create", conversationId, err.Error())
	}
	return nil
}
