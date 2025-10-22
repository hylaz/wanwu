package orm

import (
	"context"
	"strconv"

	err_code "github.com/UnicomAI/wanwu/api/proto/err-code"
	"github.com/UnicomAI/wanwu/internal/assistant-service/client/model"
	"github.com/UnicomAI/wanwu/internal/assistant-service/client/orm/sqlopt"
	"gorm.io/gorm"
)

func (c *Client) CreateAssistant(ctx context.Context, assistant *model.Assistant) *err_code.Status {
	if assistant.ID != 0 {
		return toErrStatus("assistant_create", "create assistant but id not 0")
	}
	return c.transaction(ctx, func(tx *gorm.DB) *err_code.Status {
		if err := tx.Create(assistant).Error; err != nil {
			return toErrStatus("assistant_create", err.Error())
		}
		return nil
	})
}

func (c *Client) UpdateAssistant(ctx context.Context, assistant *model.Assistant) *err_code.Status {
	if assistant.ID == 0 {
		return toErrStatus("assistant_update", "update assistant but id 0")
	}
	return c.transaction(ctx, func(tx *gorm.DB) *err_code.Status {
		if err := tx.Model(assistant).Updates(map[string]interface{}{
			"avatar_path":          assistant.AvatarPath,
			"name":                 assistant.Name,
			"desc":                 assistant.Desc,
			"instructions":         assistant.Instructions,
			"prologue":             assistant.Prologue,
			"recommend_question":   assistant.RecommendQuestion,
			"model_config":         assistant.ModelConfig,
			"knowledgebase_config": assistant.KnowledgebaseConfig,
			"scope":                assistant.Scope,
			"rerank_config":        assistant.RerankConfig,
			"online_search_config": assistant.OnlineSearchConfig,
			"safety_config":        assistant.SafetyConfig,
			"vision_config":        assistant.VisionConfig,
		}).Error; err != nil {
			return toErrStatus("assistant_update", err.Error())
		}
		return nil
	})
}

func (c *Client) DeleteAssistant(ctx context.Context, assistantID uint32) *err_code.Status {
	return c.transaction(ctx, func(tx *gorm.DB) *err_code.Status {
		if err := sqlopt.WithID(assistantID).Apply(tx).Delete(&model.Assistant{}).Error; err != nil {
			return toErrStatus("assistant_delete", err.Error())
		}
		if err := sqlopt.WithAssistantID(assistantID).Apply(tx).Delete(&model.AssistantWorkflow{}).Error; err != nil {
			return toErrStatus("assistant_delete", err.Error())
		}
		if err := sqlopt.WithAssistantID(assistantID).Apply(tx).Delete(&model.AssistantMCP{}).Error; err != nil {
			return toErrStatus("assistant_delete", err.Error())
		}
		if err := sqlopt.WithAssistantID(assistantID).Apply(tx).Delete(&model.AssistantCustom{}).Error; err != nil {
			return toErrStatus("assistant_delete", err.Error())
		}
		return nil
	})
}

func (c *Client) GetAssistant(ctx context.Context, assistantID uint32) (*model.Assistant, *err_code.Status) {
	var assistant *model.Assistant
	return assistant, c.transaction(ctx, func(tx *gorm.DB) *err_code.Status {
		assistant = &model.Assistant{}
		if err := sqlopt.WithID(assistantID).Apply(tx).First(assistant).Error; err != nil {
			return toErrStatus("assistant_get", err.Error())
		}
		return nil
	})
}

func (c *Client) GetAssistantsByIDs(ctx context.Context, assistantIDs []uint32) ([]*model.Assistant, *err_code.Status) {
	var assistants []*model.Assistant
	return assistants, c.transaction(ctx, func(tx *gorm.DB) *err_code.Status {
		if err := sqlopt.WithIDs(assistantIDs).Apply(tx).Find(&assistants).Error; err != nil {
			return toErrStatus("assistants_get_by_ids", err.Error())
		}
		return nil
	})
}

func (c *Client) GetAssistantList(ctx context.Context, userID, orgID string, name string) ([]*model.Assistant, int64, *err_code.Status) {
	var assistants []*model.Assistant
	var count int64
	return assistants, count, c.transaction(ctx, func(tx *gorm.DB) *err_code.Status {
		query := sqlopt.DataPerm(tx.Model(&model.Assistant{}), userID, orgID)

		if name != "" {
			query = query.Where("name LIKE ?", "%"+name+"%")
		}

		if err := query.Count(&count).Error; err != nil {
			return toErrStatus("assistants_get_list", err.Error())
		}

		if err := query.Order("created_at DESC").Find(&assistants).Error; err != nil {
			return toErrStatus("assistants_get_list", err.Error())
		}

		return nil
	})
}

func (c *Client) CheckSameAssistantName(ctx context.Context, userID, orgID, name, assistantID string) *err_code.Status {
	// 同一组织下不允许重名
	return c.transaction(ctx, func(tx *gorm.DB) *err_code.Status {
		query := sqlopt.SQLOptions(
			sqlopt.WithOrgID(orgID),
		).Apply(tx.Model(&model.Assistant{}))

		if assistantID != "" {
			id, _ := strconv.ParseUint(assistantID, 10, 32)
			query = query.Where("id != ?", uint32(id))
		}

		if name != "" {
			query = query.Where("name = ?", name)
		}
		var count int64
		if err := query.Count(&count).Error; err != nil {
			return toErrStatus("assistant_get_by_name", err.Error())
		}

		// 存在同名智能体
		if count > 0 {
			return toErrStatus("assistant_same_name", name)
		}
		return nil
	})
}
