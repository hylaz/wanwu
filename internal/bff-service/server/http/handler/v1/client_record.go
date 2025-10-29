package v1

import (
	"github.com/UnicomAI/wanwu/internal/bff-service/service"
	gin_util "github.com/UnicomAI/wanwu/pkg/gin-util"
	"github.com/gin-gonic/gin"
)

// GetWorkflowTemplateStatistic
//
//	@Tags			common
//	@Summary		获取工作流模板统计
//	@Description	获取工作流模板统计
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Param			startDate	query		string	true	"开始时间（格式yyyy-mm-dd）"
//	@Param			endDate		query		string	true	"结束时间（格式yyyy-mm-dd）"
//	@Success		200			{object}	response.Response{data=response.WorkflowStatistic}
//	@Router			/workflow/template/statistic [get]
func GetWorkflowTemplateStatistic(ctx *gin.Context) {
	resp, err := service.GetWorkflowTemplateStatistic(ctx, ctx.Query("startDate"), ctx.Query("endDate"))
	gin_util.Response(ctx, resp, err)
}

// GetClientStatistic
//
//	@Tags			common
//	@Summary		获取使用工作流模板用户统计
//	@Description	获取使用工作流模板用户统计
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Param			startDate	query		string	true	"开始时间（格式yyyy-mm-dd）"
//	@Param			endDate		query		string	true	"结束时间（格式yyyy-mm-dd）"
//	@Success		200			{object}	response.Response{data=response.ClientStatistic}
//	@Router			/client/statistic [get]
func GetClientStatistic(ctx *gin.Context) {
	resp, err := service.GetClientStatistic(ctx, ctx.Query("startDate"), ctx.Query("endDate"))
	gin_util.Response(ctx, resp, err)
}

// GetCumulativeClientStatistic
//
//	@Tags			common
//	@Summary		获取累计使用工作流模板用户统计
//	@Description	获取累计使用工作流模板用户统计
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Param			endAt	query		string	true	"时间戳"
//	@Success		200		{object}	response.Response{data=response.ClientStatistic}
//	@Router			/client/statistic/cumulative [get]
func GetCumulativeClientStatistic(ctx *gin.Context) {
	resp, err := service.GetCumulativeClientStatistic(ctx, ctx.Query("endAt"))
	gin_util.Response(ctx, resp, err)
}
