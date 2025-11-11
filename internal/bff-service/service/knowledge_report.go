package service

import (
	errs "github.com/UnicomAI/wanwu/api/proto/err-code"
	knowledgebase_report_service "github.com/UnicomAI/wanwu/api/proto/knowledgebase-report-service"
	"github.com/UnicomAI/wanwu/internal/bff-service/model/request"
	"github.com/UnicomAI/wanwu/internal/bff-service/model/response"
	grpc_util "github.com/UnicomAI/wanwu/pkg/grpc-util"
	"github.com/UnicomAI/wanwu/pkg/log"
	"github.com/UnicomAI/wanwu/pkg/minio"
	"github.com/gin-gonic/gin"
	"path/filepath"
)

func GetKnowledgeReport(ctx *gin.Context, userId, orgId string, req *request.GetReportReq) (*response.ReportPageResult, error) {
	resp, err := knowledgeBaseReport.GetKnowledgeReport(ctx, &knowledgebase_report_service.GetReportReq{
		KnowledgeInfo: &knowledgebase_report_service.ReportIdentity{
			KnowledgeId: req.KnowledgeId,
			UserId:      userId,
			OrgId:       orgId,
		},
		PageSize: int32(req.PageSize),
		PageNum:  int32(req.PageNo),
	})
	if err != nil {
		return nil, err
	}
	return buildKnowledgeReportList(req, resp), nil
}

func GenerateKnowledgeReport(ctx *gin.Context, userId, orgId string, req *request.GenerateReportReq) error {
	_, err := knowledgeBaseReport.GenerateKnowledgeReport(ctx, &knowledgebase_report_service.ReportIdentity{
		KnowledgeId: req.KnowledgeId,
		UserId:      userId,
		OrgId:       orgId,
	})
	return err
}

func DeleteKnowledgeReport(ctx *gin.Context, userId, orgId string, req *request.DeleteReportReq) error {
	_, err := knowledgeBaseReport.DeleteKnowledgeReport(ctx, &knowledgebase_report_service.DeleteReportReq{
		KnowledgeInfo: &knowledgebase_report_service.ReportIdentity{
			KnowledgeId: req.KnowledgeId,
			UserId:      userId,
			OrgId:       orgId,
		},
		ContentId: req.ContentId,
	})
	return err
}

func UpdateKnowledgeReport(ctx *gin.Context, userId, orgId string, req *request.UpdateReportReq) error {
	_, err := knowledgeBaseReport.UpdateKnowledgeReport(ctx, &knowledgebase_report_service.UpdateReportReq{
		KnowledgeInfo: &knowledgebase_report_service.ReportIdentity{
			KnowledgeId: req.KnowledgeId,
			UserId:      userId,
			OrgId:       orgId,
		},
		Data: &knowledgebase_report_service.ReportInfo{
			ContentId: req.ContentId,
			Title:     req.Title,
			Content:   req.Content,
		},
	})
	return err
}

func AddKnowledgeReport(ctx *gin.Context, userId, orgId string, req *request.AddReportReq) error {
	_, err := knowledgeBaseReport.AddKnowledgeReport(ctx, &knowledgebase_report_service.AddReportReq{
		KnowledgeInfo: &knowledgebase_report_service.ReportIdentity{
			KnowledgeId: req.KnowledgeId,
			UserId:      userId,
			OrgId:       orgId,
		},
		Title:   req.Title,
		Content: req.Content,
	})
	return err
}

func BatchAddKnowledgeReport(ctx *gin.Context, userId, orgId string, req *request.BatchAddReportReq) error {
	docUrl, err := minio.GetUploadFileWithExpire(ctx, req.FileUploadId)
	if err != nil {
		log.Errorf("GetUploadFileWithNotExpire error %v", err)
		return grpc_util.ErrorStatus(errs.Code_KnowledgeDocImportUrlFailed)
	}
	ext := filepath.Ext(docUrl)
	if ext != ".csv" {
		return grpc_util.ErrorStatus(errs.Code_KnowledgeDocSegmentFileCSVTypeFail)
	}
	_, err = knowledgeBaseReport.BatchAddKnowledgeReport(ctx.Request.Context(), &knowledgebase_report_service.BatchAddKnowledgeReportReq{
		KnowledgeInfo: &knowledgebase_report_service.ReportIdentity{
			KnowledgeId: req.KnowledgeId,
			UserId:      userId,
			OrgId:       orgId,
		},
		FileUrl: docUrl,
	})
	return err
}

func buildKnowledgeReportList(req *request.GetReportReq, resp *knowledgebase_report_service.GetReportResp) *response.ReportPageResult {
	retList := make([]*response.ReportInfo, 0)
	for _, v := range resp.List {
		retList = append(retList, &response.ReportInfo{
			Content:   v.Content,
			ContentId: v.ContentId,
			Title:     v.Title,
		})
	}
	return &response.ReportPageResult{
		List:      retList,
		Total:     resp.Total,
		PageNo:    req.PageNo,
		PageSize:  req.PageSize,
		CreatedAt: resp.CreatedAt,
		Status:    resp.Status,
	}
}
