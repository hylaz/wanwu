package v1

import (
	"github.com/UnicomAI/wanwu/internal/bff-service/model/request"
	"github.com/UnicomAI/wanwu/internal/bff-service/service"
	gin_util "github.com/UnicomAI/wanwu/pkg/gin-util"
	"github.com/gin-gonic/gin"
)

// GetKnowledgeReport
//
//	@Tags			knowledge
//	@Summary		获取社区报告
//	@Description	获取社区报告
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Param			data	query		request.GetReportReq	true	"获取社区报告请求参数"
//	@Success		200		{object}	response.Response{data=response.ReportPageResult}
//	@Router			/knowledge/report/list [get]
func GetKnowledgeReport(ctx *gin.Context) {
	userId, orgId := getUserID(ctx), getOrgID(ctx)
	var req request.GetReportReq
	if !gin_util.BindQuery(ctx, &req) {
		return
	}
	resp, err := service.GetKnowledgeReport(ctx, userId, orgId, &req)
	gin_util.Response(ctx, resp, err)
}

// GenerateKnowledgeReport
//
//	@Tags			knowledge
//	@Summary		生成社区报告
//	@Description	生成社区报告
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Param			data	body		request.GenerateReportReq	true	"生成社区报告请求参数"
//	@Success		200		{object}	response.Response
//	@Router			/knowledge/report/generate [post]
func GenerateKnowledgeReport(ctx *gin.Context) {
	userId, orgId := getUserID(ctx), getOrgID(ctx)
	var req request.GenerateReportReq
	if !gin_util.Bind(ctx, &req) {
		return
	}
	err := service.GenerateKnowledgeReport(ctx, userId, orgId, &req)
	gin_util.Response(ctx, nil, err)
}

// DeleteKnowledgeReport
//
//	@Tags			knowledge
//	@Summary		删除社区报告
//	@Description	删除社区报告
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Param			data	body		request.DeleteReportReq	true	"删除社区报告请求参数"
//	@Success		200		{object}	response.Response
//	@Router			/knowledge/report/delete [delete]
func DeleteKnowledgeReport(ctx *gin.Context) {
	userId, orgId := getUserID(ctx), getOrgID(ctx)
	var req request.DeleteReportReq
	if !gin_util.Bind(ctx, &req) {
		return
	}
	err := service.DeleteKnowledgeReport(ctx, userId, orgId, &req)
	gin_util.Response(ctx, nil, err)
}

// UpdateKnowledgeReport
//
//	@Tags			knowledge
//	@Summary		编辑社区报告
//	@Description	编辑社区报告
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Param			data	body		request.UpdateReportReq	true	"编辑社区报告请求参数"
//	@Success		200		{object}	response.Response
//	@Router			/knowledge/report/update [post]
func UpdateKnowledgeReport(ctx *gin.Context) {
	userId, orgId := getUserID(ctx), getOrgID(ctx)
	var req request.UpdateReportReq
	if !gin_util.Bind(ctx, &req) {
		return
	}
	err := service.UpdateKnowledgeReport(ctx, userId, orgId, &req)
	gin_util.Response(ctx, nil, err)
}

// AddKnowledgeReport
//
//	@Tags			knowledge
//	@Summary		单条新增社区报告
//	@Description	单条新增社区报告
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Param			data	body		request.AddReportReq	true	"单条新增社区报告请求参数"
//	@Success		200		{object}	response.Response
//	@Router			/knowledge/report/add [post]
func AddKnowledgeReport(ctx *gin.Context) {
	userId, orgId := getUserID(ctx), getOrgID(ctx)
	var req request.AddReportReq
	if !gin_util.Bind(ctx, &req) {
		return
	}
	err := service.AddKnowledgeReport(ctx, userId, orgId, &req)
	gin_util.Response(ctx, nil, err)
}

// BatchAddKnowledgeReport
//
//	@Tags			knowledge
//	@Summary		批量新增社区报告
//	@Description	批量新增社区报告
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Param			data	body		request.BatchAddReportReq	true	"批量新增社区报告请求参数"
//	@Success		200		{object}	response.Response
//	@Router			/knowledge/report/batch/add [post]
func BatchAddKnowledgeReport(ctx *gin.Context) {
	userId, orgId := getUserID(ctx), getOrgID(ctx)
	var req request.BatchAddReportReq
	if !gin_util.Bind(ctx, &req) {
		return
	}
	err := service.BatchAddKnowledgeReport(ctx, userId, orgId, &req)
	gin_util.Response(ctx, nil, err)
}
