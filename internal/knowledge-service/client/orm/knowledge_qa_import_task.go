package orm

import (
	"context"

	"github.com/UnicomAI/wanwu/internal/knowledge-service/client/model"
	"github.com/UnicomAI/wanwu/internal/knowledge-service/client/orm/sqlopt"
	async_task "github.com/UnicomAI/wanwu/internal/knowledge-service/pkg/async-task"
	"github.com/UnicomAI/wanwu/internal/knowledge-service/pkg/db"
	"github.com/UnicomAI/wanwu/pkg/log"
	"gorm.io/gorm"
)

// DeleteQAImportTaskByKnowledgeId 根据问答库id 删除导入任务
func DeleteQAImportTaskByKnowledgeId(tx *gorm.DB, knowledgeId string) error {
	var count int64
	err := sqlopt.SQLOptions(sqlopt.WithKnowledgeID(knowledgeId)).
		Apply(tx, &model.KnowledgeQAPairImportTask{}).
		Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		return tx.Unscoped().Model(&model.KnowledgeQAPairImportTask{}).Where("knowledge_id = ?", knowledgeId).Delete(&model.KnowledgeQAPairImportTask{}).Error
	}
	return nil
}

// CreateKnowledgeQAPairImportTask 导入任务
func CreateKnowledgeQAPairImportTask(ctx context.Context, importTask *model.KnowledgeQAPairImportTask) error {
	return db.GetHandle(ctx).Transaction(func(tx *gorm.DB) error {
		//1.创建知识库导入任务
		err := createQAPairImportTask(tx, importTask)
		if err != nil {
			return err
		}
		//2.通知rag更新知识库
		return async_task.SubmitTask(ctx, async_task.KnowledgeQAPairImportTaskType, &async_task.DocImportTaskParams{
			TaskId: importTask.ImportId,
		})
	})
}

// SelectKnowledgeQAPairImportTaskById 根据id查询导入信息
func SelectKnowledgeQAPairImportTaskById(ctx context.Context, importId string) (*model.KnowledgeQAPairImportTask, error) {
	var importTask model.KnowledgeQAPairImportTask
	err := sqlopt.SQLOptions(sqlopt.WithImportID(importId)).
		Apply(db.GetHandle(ctx), &model.KnowledgeQAPairImportTask{}).
		First(&importTask).Error
	if err != nil {
		log.Errorf("SelectKnowledgeQAPairRunningImportTask importId %s err: %v", importId, err)
		return nil, err
	}
	return &importTask, nil
}

// UpdateKnowledgeQAPairImportTaskStatus 更新导入任务状态
func UpdateKnowledgeQAPairImportTaskStatus(ctx context.Context, taskId string, status int, errMsg string, totalCount int64, successCount int64) error {
	return db.GetHandle(ctx).Model(&model.KnowledgeQAPairImportTask{}).
		Where("import_id = ?", taskId).
		Updates(map[string]interface{}{
			"status":        status,
			"error_msg":     errMsg,
			"success_count": successCount,
			"total_count":   totalCount,
		}).Error
}

func createQAPairImportTask(tx *gorm.DB, importTask *model.KnowledgeQAPairImportTask) error {
	return tx.Model(&model.KnowledgeQAPairImportTask{}).Create(importTask).Error
}
