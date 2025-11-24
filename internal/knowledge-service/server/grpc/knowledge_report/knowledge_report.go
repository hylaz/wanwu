package knowledge_report

import (
	"context"

	errs "github.com/UnicomAI/wanwu/api/proto/err-code"
	knowledgebase_report_service "github.com/UnicomAI/wanwu/api/proto/knowledgebase-report-service"
	"github.com/UnicomAI/wanwu/internal/knowledge-service/client/model"
	"github.com/UnicomAI/wanwu/internal/knowledge-service/client/orm"
	"github.com/UnicomAI/wanwu/internal/knowledge-service/pkg/util"
	"github.com/UnicomAI/wanwu/internal/knowledge-service/service"
	"github.com/UnicomAI/wanwu/pkg/log"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	Regenerate       = "重新生成"
	ReportProcessing = 1
	ReportSuccess    = 2
	ReportFailed     = 3
)

func (s *Service) GetKnowledgeReport(ctx context.Context, req *knowledgebase_report_service.GetReportReq) (*knowledgebase_report_service.GetReportResp, error) {
	//查询知识库信息
	knowledge, err := orm.SelectKnowledgeById(ctx, req.KnowledgeInfo.KnowledgeId, "", "")
	if err != nil {
		log.Errorf("没有操作该知识库的权限 错误(%v) 参数(%v)", err, req)
		return nil, util.ErrCode(errs.Code_KnowledgeBaseSelectFailed)
	}
	// 查询知识图谱状态
	graphSuccess, err := getGraphStatus(ctx, req)
	if err != nil {
		return nil, err
	}
	// 没有成功生成的知识图谱
	if !graphSuccess {
		return &knowledgebase_report_service.GetReportResp{}, nil
	}
	lastImportTask, _ := orm.SelectReportLatestImportTaskByKnowledgeID(ctx, req.KnowledgeInfo.KnowledgeId)
	status := buildLastImportStatus(lastImportTask)
	// 判断报告状态
	switch knowledge.ReportStatus {
	case model.ReportInit:
		return &knowledgebase_report_service.GetReportResp{CanGenerate: true}, nil

	case model.ReportProcessing:
		if knowledge.ReportCreateCount > 0 {
			return &knowledgebase_report_service.GetReportResp{Status: ReportProcessing, GenerateLabel: Regenerate}, nil
		}
		return &knowledgebase_report_service.GetReportResp{Status: ReportProcessing}, nil

	case model.ReportSuccess:
		return handleSuccessReport(ctx, knowledge, req, status)

	default: // 生成失败状态
		if model.ErrorReportStatus(knowledge.ReportStatus) {
			return &knowledgebase_report_service.GetReportResp{
				Status:        ReportFailed,
				CanGenerate:   true,
				GenerateLabel: Regenerate,
			}, nil
		}
		return nil, nil
	}
}

func (s *Service) GenerateKnowledgeReport(ctx context.Context, req *knowledgebase_report_service.ReportIdentity) (*emptypb.Empty, error) {
	err := orm.CreateKnowledgeReport(ctx, req.KnowledgeId)
	if err != nil {
		return nil, util.ErrCode(errs.Code_KnowledgeGenerateReportFailed)
	}
	return &emptypb.Empty{}, nil
}

func (s *Service) DeleteKnowledgeReport(ctx context.Context, req *knowledgebase_report_service.DeleteReportReq) (*emptypb.Empty, error) {
	//查询知识库信息
	knowledge, err := orm.SelectKnowledgeById(ctx, req.KnowledgeInfo.KnowledgeId, "", "")
	if err != nil {
		log.Errorf("没有操作该知识库的权限 错误(%v) 参数(%v)", err, req)
		return nil, util.ErrCode(errs.Code_KnowledgeBaseSelectFailed)
	}
	err = service.RagDeleteReport(ctx, &service.RagDeleteReportParams{
		UserId:            knowledge.UserId,
		KnowledgeBaseName: knowledge.RagName,
		KnowledgeId:       knowledge.KnowledgeId,
		ReportIds:         []string{req.ContentId},
	})
	if err != nil {
		return nil, util.ErrCode(errs.Code_KnowledgeDeleteReportFailed)
	}
	return &emptypb.Empty{}, nil
}

func (s *Service) UpdateKnowledgeReport(ctx context.Context, req *knowledgebase_report_service.UpdateReportReq) (*emptypb.Empty, error) {
	//查询知识库信息
	knowledge, err := orm.SelectKnowledgeById(ctx, req.KnowledgeInfo.KnowledgeId, "", "")
	if err != nil {
		log.Errorf("没有操作该知识库的权限 错误(%v) 参数(%v)", err, req)
		return nil, util.ErrCode(errs.Code_KnowledgeBaseSelectFailed)
	}
	err = service.RagUpdateReport(ctx, &service.RagUpdateReportParams{
		UserId:            knowledge.UserId,
		KnowledgeBaseName: knowledge.RagName,
		KnowledgeId:       knowledge.KnowledgeId,
		ReportItem: &service.RagUpdateReportItem{
			Content:  req.Data.Content,
			ReportId: req.Data.ContentId,
			Title:    req.Data.Title,
		},
	})
	if err != nil {
		return nil, util.ErrCode(errs.Code_KnowledgeUpdateReportFailed)
	}
	return &emptypb.Empty{}, nil
}

func (s *Service) AddKnowledgeReport(ctx context.Context, req *knowledgebase_report_service.AddReportReq) (*emptypb.Empty, error) {
	//查询知识库信息
	knowledge, err := orm.SelectKnowledgeById(ctx, req.KnowledgeInfo.KnowledgeId, "", "")
	if err != nil {
		log.Errorf("没有操作该知识库的权限 错误(%v) 参数(%v)", err, req)
		return nil, util.ErrCode(errs.Code_KnowledgeBaseSelectFailed)
	}
	err = service.RagAddReport(ctx, &service.RagAddReportParams{
		UserId:            knowledge.UserId,
		KnowledgeBaseName: knowledge.RagName,
		KnowledgeId:       knowledge.KnowledgeId,
		ReportItem: []*service.RagAddReportItem{{
			Content: req.Content,
			Title:   req.Title,
		}},
	})
	if err != nil {
		return nil, util.ErrCode(errs.Code_KnowledgeAddReportFailed)
	}
	return &emptypb.Empty{}, nil
}

func (s *Service) BatchAddKnowledgeReport(ctx context.Context, req *knowledgebase_report_service.BatchAddKnowledgeReportReq) (*emptypb.Empty, error) {
	//查询知识库信息
	knowledge, err := orm.SelectKnowledgeById(ctx, req.KnowledgeInfo.KnowledgeId, "", "")
	if err != nil {
		log.Errorf("没有操作该知识库的权限 错误(%v) 参数(%v)", err, req)
		return nil, util.ErrCode(errs.Code_KnowledgeBaseSelectFailed)
	}
	task, err := orm.BuildReportImportTask(req, knowledge)
	if err != nil {
		return nil, util.ErrCode(errs.Code_KnowledgeAddReportFailed)
	}
	err = orm.BatchCreateKnowledgeReport(ctx, task)
	if err != nil {
		return nil, util.ErrCode(errs.Code_KnowledgeAddReportFailed)
	}
	return &emptypb.Empty{}, nil
}

func handleSuccessReport(ctx context.Context, knowledge *model.KnowledgeBase, req *knowledgebase_report_service.GetReportReq, lastImportStatus int32) (*knowledgebase_report_service.GetReportResp, error) {
	resp, err := service.RagGetReport(ctx, &service.RagGetReportParams{
		UserId:            knowledge.UserId,
		KnowledgeBaseName: knowledge.RagName,
		KnowledgeId:       knowledge.KnowledgeId,
		PageSize:          req.PageSize,
		SearchAfter:       req.PageSize * (req.PageNum - 1),
	})
	if err != nil {
		log.Errorf("获取社区报告失败 错误(%v) 参数(%v)", err, req)
		return nil, util.ErrCode(errs.Code_KnowledgeGetReportFailed)
	}
	return buildReportSuccessResp(resp, lastImportStatus), nil
}

func getGraphStatus(ctx context.Context, req *knowledgebase_report_service.GetReportReq) (bool, error) {
	//获取文档列表
	list, err := orm.GetDocListByKnowledgeId(ctx, "", "", req.KnowledgeInfo.KnowledgeId)
	if err != nil {
		log.Errorf("获取知识库列表失败(%v)  参数(%v)", err, req)
		return false, util.ErrCode(errs.Code_KnowledgeBaseSelectFailed)
	}
	//判断是否有成功的知识图谱
	graphSuccess := false
	if len(list) > 0 {
		for _, doc := range list {
			if doc.GraphStatus == model.GraphSuccess {
				graphSuccess = true
			}
		}
	}
	return graphSuccess, nil
}

func buildReportSuccessResp(resp *service.RagReportListResp, lastImportStatus int32) *knowledgebase_report_service.GetReportResp {
	retList := make([]*knowledgebase_report_service.ReportInfo, 0)
	for _, item := range resp.List {
		retList = append(retList, &knowledgebase_report_service.ReportInfo{
			Content:   item.Content,
			ContentId: item.ContentId,
			Title:     item.ReportTitle,
		})
	}
	return &knowledgebase_report_service.GetReportResp{
		Total:            int32(resp.ChunkTotalNum),
		CreatedAt:        buildCreatedAt(resp),
		Status:           ReportSuccess,
		CanGenerate:      true,
		CanAddReport:     true,
		GenerateLabel:    Regenerate,
		List:             retList,
		LastImportStatus: lastImportStatus,
	}
}

func buildCreatedAt(resp *service.RagReportListResp) string {
	if resp.ChunkTotalNum == 0 {
		return ""
	}
	return resp.List[0].CreateTime
}

func buildLastImportStatus(importTask *model.KnowledgeReportImportTask) int32 {
	if importTask == nil {
		return -1
	}
	return int32(importTask.Status)
}
