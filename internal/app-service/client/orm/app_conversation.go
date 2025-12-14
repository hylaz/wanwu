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

func (c *Client) GetChatflowApplication(ctx context.Context, orgId, userId, workflowId string) (*model.ChatflowApplcation, *errs.Status) {
	//如果记录不存在就返回空字符串
	var chatflowApp model.ChatflowApplcation
	err := sqlopt.SQLOptions(
		sqlopt.WithUserID(userId),
		sqlopt.WithWorkflowID(workflowId),
	).Apply(c.db.WithContext(ctx)).First(&chatflowApp).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &model.ChatflowApplcation{}, nil
		}
		return nil, toErrStatus("chatflow_application_get", workflowId, err.Error())
	}
	return &chatflowApp, nil
}

func (c *Client) GetChatflowApplicationByApplicationID(ctx context.Context, orgId, userId, applicationId string) (*model.ChatflowApplcation, *errs.Status) {
	//如果记录不存在就返回空字符串
	var chatflowApp model.ChatflowApplcation
	err := sqlopt.SQLOptions(
		sqlopt.WithUserID(userId),
		sqlopt.WithApplicationID(applicationId),
	).Apply(c.db.WithContext(ctx)).First(&chatflowApp).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &model.ChatflowApplcation{}, nil
		}
		return nil, toErrStatus("chatflow_application_get_by_application_id", applicationId, err.Error())
	}
	return &chatflowApp, nil
}

func (c *Client) CreateChatflowApplication(ctx context.Context, orgId, userId, workflowId, applicationId string) *errs.Status {
	//存储关联关系
	chatflowApp := &model.ChatflowApplcation{
		OrgID:         orgId,
		UserID:        userId,
		ApplicationID: applicationId,
		WorkflowID:    workflowId,
	}
	if err := c.db.WithContext(ctx).Create(chatflowApp).Error; err != nil {
		return toErrStatus("chatflow_application_create", workflowId, err.Error())
	}
	return nil
}
