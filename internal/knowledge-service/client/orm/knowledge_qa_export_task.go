package orm

import (
	"context"

	"github.com/UnicomAI/wanwu/internal/knowledge-service/client/model"
	"github.com/UnicomAI/wanwu/internal/knowledge-service/client/orm/sqlopt"
	async_task "github.com/UnicomAI/wanwu/internal/knowledge-service/pkg/async-task"
	"github.com/UnicomAI/wanwu/internal/knowledge-service/pkg/db"
	"github.com/UnicomAI/wanwu/internal/knowledge-service/service"
	"github.com/UnicomAI/wanwu/pkg/log"
	"gorm.io/gorm"
)

// DeleteQAExportTaskByKnowledgeId 根据问答库id 删除导出任务
func DeleteQAExportTaskByKnowledgeId(tx *gorm.DB, knowledgeId string) error {
	var count int64
	err := sqlopt.SQLOptions(sqlopt.WithKnowledgeID(knowledgeId)).
		Apply(tx, &model.KnowledgeQAPairExportTask{}).
		Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		return tx.Unscoped().Model(&model.KnowledgeQAPairExportTask{}).Where("knowledge_id = ?", knowledgeId).Delete(&model.KnowledgeQAPairExportTask{}).Error
	}
	return nil
}

// CreateKnowledgeQAPairExportTask 导出任务
func CreateKnowledgeQAPairExportTask(ctx context.Context, exportTask *model.KnowledgeQAPairExportTask) error {
	return db.GetHandle(ctx).Transaction(func(tx *gorm.DB) error {
		//1.创建问答库导出任务
		err := createQAPairExportTask(tx, exportTask)
		if err != nil {
			return err
		}
		//2.提交问答库异步任务
		return async_task.SubmitTask(ctx, async_task.KnowledgeQAPairExportTaskType, &async_task.KnowledgeQAPairExportTaskParams{
			TaskId: exportTask.ExportId,
		})
	})
}

// SelectKnowledgeQAPairExportTaskById 根据id查询导出信息
func SelectKnowledgeQAPairExportTaskById(ctx context.Context, exportId string) (*model.KnowledgeQAPairExportTask, error) {
	var exportTask model.KnowledgeQAPairExportTask
	err := sqlopt.SQLOptions(sqlopt.WithExportID(exportId)).
		Apply(db.GetHandle(ctx), &model.KnowledgeQAPairExportTask{}).
		First(&exportTask).Error
	if err != nil {
		log.Errorf("SelectKnowledgeQAPairRunningExportTask exportId %s err: %v", exportId, err)
		return nil, err
	}
	return &exportTask, nil
}

// SelectKnowledgeQAPairExportTaskByQAId 根据knowledge id查询导出信息
func SelectKnowledgeQAPairExportTaskByQAId(ctx context.Context, knowledgeId string) ([]*model.KnowledgeQAPairExportTask, error) {
	var exportTask []*model.KnowledgeQAPairExportTask
	err := sqlopt.SQLOptions(sqlopt.WithKnowledgeID(knowledgeId)).
		Apply(db.GetHandle(ctx), &model.KnowledgeQAPairExportTask{}).
		Find(&exportTask).Error
	if err != nil {
		log.Errorf("SelectKnowledgeQAPairRunningExportTask knowledgeId %s err: %v", knowledgeId, err)
		return nil, err
	}
	return exportTask, nil
}

// UpdateKnowledgeQAPairExportTask 更新导出任务状态
func UpdateKnowledgeQAPairExportTask(ctx context.Context, taskId string, status int, errMsg string, totalCount int64, successCount int64, filePath string, fileSize int64) error {
	return db.GetHandle(ctx).Model(&model.KnowledgeQAPairExportTask{}).
		Where("export_id = ?", taskId).
		Updates(map[string]interface{}{
			"status":           status,
			"error_msg":        errMsg,
			"success_count":    successCount,
			"total_count":      totalCount,
			"export_file_path": filePath,
			"export_file_size": fileSize,
		}).Error
}

// DeleteQAExportTaskById 根据导出任务Id 删除导出任务
func DeleteQAExportTaskById(ctx context.Context, taskId string) error {
	var exportTask model.KnowledgeQAPairExportTask
	err := db.GetHandle(ctx).Transaction(func(tx *gorm.DB) error {
		var count int64
		err := sqlopt.SQLOptions(sqlopt.WithExportID(taskId)).
			Apply(tx, &model.KnowledgeQAPairExportTask{}).
			Find(&exportTask).Count(&count).Error
		if err != nil {
			return err
		}
		if count > 0 {
			err = tx.Unscoped().Model(&model.KnowledgeQAPairExportTask{}).Where("export_id = ?", taskId).Delete(&model.KnowledgeQAPairExportTask{}).Error
			if err != nil {
				return err
			}
			//删除minio中的文件
			err = service.DeleteFile(ctx, exportTask.ExportFilePath)
			if err != nil {
				log.Errorf("minioDelete error %v", err)
				return err
			}
			return err
		}
		return nil
	})
	return err
}

func createQAPairExportTask(tx *gorm.DB, exportTask *model.KnowledgeQAPairExportTask) error {
	return tx.Model(&model.KnowledgeQAPairExportTask{}).Create(exportTask).Error
}
