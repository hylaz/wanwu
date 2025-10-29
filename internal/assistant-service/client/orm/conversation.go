package orm

import (
	"context"

	err_code "github.com/UnicomAI/wanwu/api/proto/err-code"
	"github.com/UnicomAI/wanwu/internal/assistant-service/client/model"
	"github.com/UnicomAI/wanwu/internal/assistant-service/client/orm/sqlopt"
	"gorm.io/gorm"
)

func (c *Client) CreateConversation(ctx context.Context, conversation *model.Conversation) *err_code.Status {
	if conversation.ID != 0 {
		return toErrStatus("assistant_conversation_create", "create conversation but id not 0")
	}
	return c.transaction(ctx, func(tx *gorm.DB) *err_code.Status {
		if err := tx.Create(conversation).Error; err != nil {
			return toErrStatus("assistant_conversation_create", err.Error())
		}
		return nil
	})
}

func (c *Client) UpdateConversation(ctx context.Context, conversation *model.Conversation) *err_code.Status {
	if conversation.ID == 0 {
		return toErrStatus("assistant_conversation_update", "update conversation but id 0")
	}
	return c.transaction(ctx, func(tx *gorm.DB) *err_code.Status {
		if err := tx.Model(conversation).Updates(map[string]interface{}{
			"title": conversation.Title,
		}).Error; err != nil {
			return toErrStatus("assistant_conversation_update", err.Error())
		}
		return nil
	})
}

func (c *Client) DeleteConversation(ctx context.Context, conversationID uint32) *err_code.Status {
	return c.transaction(ctx, func(tx *gorm.DB) *err_code.Status {
		if err := sqlopt.WithID(conversationID).Apply(tx).Delete(&model.Conversation{}).Error; err != nil {
			return toErrStatus("assistant_conversation_delete", err.Error())
		}
		return nil
	})
}

func (c *Client) GetConversation(ctx context.Context, conversationID uint32) (*model.Conversation, *err_code.Status) {
	var conversation *model.Conversation
	return conversation, c.transaction(ctx, func(tx *gorm.DB) *err_code.Status {
		conversation = &model.Conversation{}
		if err := sqlopt.WithID(conversationID).Apply(tx).First(conversation).Error; err != nil {
			return toErrStatus("assistant_conversation_get", err.Error())
		}
		return nil
	})
}

func (c *Client) GetConversationList(ctx context.Context, assistantID, userID, orgID string, offset, limit int32) ([]*model.Conversation, int64, *err_code.Status) {
	var conversations []*model.Conversation
	var count int64
	return conversations, count, c.transaction(ctx, func(tx *gorm.DB) *err_code.Status {
		query := sqlopt.DataPerm(userID, orgID).Apply(tx.Model(&model.Conversation{}))

		if assistantID != "" {
			query = query.Where("assistant_id = ?", assistantID)
		}

		if err := query.Count(&count).Error; err != nil {
			return toErrStatus("assistant_conversations_get_list", err.Error())
		}

		if err := query.Offset(int(offset)).Limit(int(limit)).Order("created_at DESC").Find(&conversations).Error; err != nil {
			return toErrStatus("assistant_conversations_get_list", err.Error())
		}

		return nil
	})
}

func (c *Client) DeleteConversationByAssistantID(ctx context.Context, assistantID, userID, orgID string) *err_code.Status {
	return c.transaction(ctx, func(tx *gorm.DB) *err_code.Status {
		if err := tx.Where("assistant_id = ?", assistantID).Delete(&model.Conversation{}).Error; err != nil {
			return toErrStatus("assistant_conversation_delete", err.Error())
		}
		return nil
	})
}
