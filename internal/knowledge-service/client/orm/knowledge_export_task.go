package orm

import (
	"context"

	"github.com/UnicomAI/wanwu/internal/knowledge-service/client/model"
	"github.com/UnicomAI/wanwu/internal/knowledge-service/client/orm/sqlopt"
	async_task "github.com/UnicomAI/wanwu/internal/knowledge-service/pkg/async-task"
	"github.com/UnicomAI/wanwu/internal/knowledge-service/pkg/config"
	"github.com/UnicomAI/wanwu/internal/knowledge-service/pkg/db"
	"github.com/UnicomAI/wanwu/internal/knowledge-service/service"
	"github.com/UnicomAI/wanwu/pkg/log"
	"gorm.io/gorm"
)

// DeleteExportTaskByKnowledgeId 根据知识库id 删除导出任务
func DeleteExportTaskByKnowledgeId(tx *gorm.DB, knowledgeId string) error {
	var count int64
	err := sqlopt.SQLOptions(sqlopt.WithKnowledgeID(knowledgeId)).
		Apply(tx, &model.KnowledgeExportTask{}).
		Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		return tx.Unscoped().Model(&model.KnowledgeExportTask{}).Where("knowledge_id = ?", knowledgeId).Delete(&model.KnowledgeExportTask{}).Error
	}
	return nil
}

// CreateKnowledgeQAPairExportTask 问答库导出任务
func CreateKnowledgeQAPairExportTask(ctx context.Context, exportTask *model.KnowledgeExportTask) error {
	return db.GetHandle(ctx).Transaction(func(tx *gorm.DB) error {
		//1.创建问答库导出任务
		err := createExportTask(tx, exportTask)
		if err != nil {
			return err
		}
		//2.提交问答库异步任务
		return async_task.SubmitTask(ctx, async_task.KnowledgeQAPairExportTaskType, &async_task.KnowledgeQAPairExportTaskParams{
			TaskId: exportTask.ExportId,
		})
	})
}

// CreateKnowledgeDocExportTask 知识库导出任务
func CreateKnowledgeDocExportTask(ctx context.Context, exportTask *model.KnowledgeExportTask) error {
	return db.GetHandle(ctx).Transaction(func(tx *gorm.DB) error {
		//1.创建知识库导出任务
		err := createExportTask(tx, exportTask)
		if err != nil {
			return err
		}
		//2.提交知识库异步任务
		return async_task.SubmitTask(ctx, async_task.KnowledgeDocExportTaskType, &async_task.KnowledgeDocExportTaskParams{
			TaskId: exportTask.ExportId,
		})
	})
}

// SelectKnowledgeExportTaskById 根据id查询导出信息
func SelectKnowledgeExportTaskById(ctx context.Context, exportId string) (*model.KnowledgeExportTask, error) {
	var exportTask model.KnowledgeExportTask
	err := sqlopt.SQLOptions(sqlopt.WithExportID(exportId)).
		Apply(db.GetHandle(ctx), &model.KnowledgeExportTask{}).
		First(&exportTask).Error
	if err != nil {
		log.Errorf("SelectKnowledgeRunningExportTask exportId %s err: %v", exportId, err)
		return nil, err
	}
	return &exportTask, nil
}

// SelectKnowledgeExportTaskByKnowledgeId 根据knowledge id查询导出信息
func SelectKnowledgeExportTaskByKnowledgeId(ctx context.Context, knowledgeId string, userId string, orgId string, pageSize int32, pageNum int32) ([]*model.KnowledgeExportTask, error) {
	limit := pageSize
	offset := pageSize * (pageNum - 1)
	var exportTask []*model.KnowledgeExportTask
	err := sqlopt.SQLOptions(sqlopt.WithKnowledgeID(knowledgeId),
		sqlopt.WithPermit(orgId, userId)).Apply(db.GetHandle(ctx), &model.KnowledgeExportTask{}).
		Order("create_at desc").Limit(int(limit)).Offset(int(offset)).
		Find(&exportTask).Error
	if err != nil {
		log.Errorf("SelectKnowledgeRunningExportTask knowledgeId %s err: %v", knowledgeId, err)
		return nil, err
	}
	return exportTask, nil
}

// SelectExportTaskByKnowledgeIdNoDeleteCheck 根据knowledge id查询导出信息
func SelectExportTaskByKnowledgeIdNoDeleteCheck(ctx context.Context, knowledgeId string, userId string, orgId string) ([]*model.KnowledgeExportTask, error) {
	var exportTask []*model.KnowledgeExportTask
	err := sqlopt.SQLOptions(sqlopt.WithKnowledgeID(knowledgeId),
		sqlopt.WithPermit(orgId, userId)).Apply(db.GetHandle(ctx), &model.KnowledgeExportTask{}).
		Order("create_at desc").
		Find(&exportTask).Error
	if err != nil {
		log.Errorf("SelectKnowledgeRunningExportTask knowledgeId %s err: %v", knowledgeId, err)
		return nil, err
	}
	return exportTask, nil
}

// UpdateKnowledgeExportTask 更新导出任务状态
func UpdateKnowledgeExportTask(ctx context.Context, taskId string, status int, errMsg string, totalCount int64, successCount int64, filePath string, fileSize int64) error {
	return db.GetHandle(ctx).Model(&model.KnowledgeExportTask{}).
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

// DeleteExportTaskById 根据导出任务Id 删除导出任务
func DeleteExportTaskById(ctx context.Context, taskId string) error {
	var exportTask model.KnowledgeExportTask
	err := db.GetHandle(ctx).Transaction(func(tx *gorm.DB) error {
		var count int64
		err := sqlopt.SQLOptions(sqlopt.WithExportID(taskId)).
			Apply(tx, &model.KnowledgeExportTask{}).
			Find(&exportTask).Count(&count).Error
		if err != nil {
			return err
		}
		if count > 0 {
			err = tx.Unscoped().Model(&model.KnowledgeExportTask{}).Where("export_id = ?", taskId).Delete(&model.KnowledgeExportTask{}).Error
			if err != nil {
				return err
			}
			//删除minio中的文件
			filePath := "http://" + config.GetConfig().Minio.EndPoint + "/" + exportTask.ExportFilePath
			err = service.DeleteFile(ctx, filePath)
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

func createExportTask(tx *gorm.DB, exportTask *model.KnowledgeExportTask) error {
	return tx.Model(&model.KnowledgeExportTask{}).Create(exportTask).Error
}
