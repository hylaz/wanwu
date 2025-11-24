package orm

import (
	"context"
	"encoding/json"
	"time"

	knowledgebase_report_service "github.com/UnicomAI/wanwu/api/proto/knowledgebase-report-service"
	"github.com/UnicomAI/wanwu/internal/knowledge-service/client/model"
	"github.com/UnicomAI/wanwu/internal/knowledge-service/client/orm/sqlopt"
	async_task "github.com/UnicomAI/wanwu/internal/knowledge-service/pkg/async-task"
	"github.com/UnicomAI/wanwu/internal/knowledge-service/pkg/db"
	"github.com/UnicomAI/wanwu/internal/knowledge-service/pkg/generator"
	"github.com/UnicomAI/wanwu/internal/knowledge-service/service"
	"github.com/UnicomAI/wanwu/pkg/log"
	"gorm.io/gorm"
)

func SelectKnowledgeReportById(ctx context.Context, importId string) (*model.KnowledgeReportImportTask, error) {
	var ret = model.KnowledgeReportImportTask{}
	err := sqlopt.SQLOptions(sqlopt.WithImportID(importId)).
		Apply(db.GetHandle(ctx), &model.KnowledgeReportImportTask{}).
		First(&ret).Error
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

// SelectReportLatestImportTaskByKnowledgeID 查询最新的导入信息
func SelectReportLatestImportTaskByKnowledgeID(ctx context.Context, knowledgeId string) (*model.KnowledgeReportImportTask, error) {
	var importTask model.KnowledgeReportImportTask
	err := sqlopt.SQLOptions(sqlopt.WithKnowledgeID(knowledgeId)).
		Apply(db.GetHandle(ctx), &model.KnowledgeReportImportTask{}).
		Order("create_at desc").
		First(&importTask).Error
	if err != nil {
		log.Errorf("SelectReportLatestImportTaskByKnowledgeID docId %s err: %v", knowledgeId, err)
		return nil, err
	}
	return &importTask, nil
}

// BatchCreateKnowledgeReport 批量添加知识报告
func BatchCreateKnowledgeReport(ctx context.Context, reportImportTask *model.KnowledgeReportImportTask) error {
	return db.GetHandle(ctx).Transaction(func(tx *gorm.DB) error {
		//1.创建知识报告导入任务
		err := tx.Create(&reportImportTask).Error
		if err != nil {
			return err
		}
		//2.异步执行导入数据
		return async_task.SubmitTask(ctx, async_task.KnowledgeReportTaskType, &async_task.KnowledgeReportImportTaskParams{
			TaskId: reportImportTask.ImportId,
		})
	})
}

// CreateOneKnowledgeReport 创建一个报告
func CreateOneKnowledgeReport(ctx context.Context, importTask *model.KnowledgeReportImportTask, importParams *service.RagAddReportParams) error {
	return db.GetHandle(ctx).Transaction(func(tx *gorm.DB) error {
		//1.创建知识库导入任务
		err := tx.Model(&model.KnowledgeReportImportTask{}).Where("import_id = ?", importTask.ImportId).
			Update("success_count", gorm.Expr("success_count + ?", 1)).Error
		if err != nil {
			return err
		}
		//2.通知rag创建分段
		return service.RagAddReport(ctx, importParams)
	})
}

// UpdateReportImportTaskStatus 更新导入任务状态
func UpdateReportImportTaskStatus(ctx context.Context, taskId string, status int, errMsg string, totalCount int) error {
	return db.GetHandle(ctx).Model(&model.KnowledgeReportImportTask{}).
		Where("import_id = ?", taskId).
		Updates(map[string]interface{}{
			"status":      status,
			"error_msg":   errMsg,
			"total_count": totalCount,
		}).Error
}

func BuildReportImportTask(req *knowledgebase_report_service.BatchAddKnowledgeReportReq, knowledge *model.KnowledgeBase) (*model.KnowledgeReportImportTask, error) {
	params := model.KnowledgeReportImportParams{
		FileUrl: req.FileUrl,
	}
	marshal, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	milli := time.Now().UnixMilli()
	return &model.KnowledgeReportImportTask{
		ImportId:     generator.GetGenerator().NewID(),
		KnowledgeId:  knowledge.KnowledgeId,
		Status:       model.KnowledgeReportImportInit,
		TotalCount:   0,
		SuccessCount: 0,
		ErrorMsg:     "",
		ImportParams: string(marshal),
		UserId:       req.KnowledgeInfo.UserId,
		OrgId:        req.KnowledgeInfo.OrgId,
		CreatedAt:    milli,
		UpdatedAt:    milli,
	}, nil
}
