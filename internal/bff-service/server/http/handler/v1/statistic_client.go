package v1

import (
	"github.com/UnicomAI/wanwu/internal/bff-service/service"
	gin_util "github.com/UnicomAI/wanwu/pkg/gin-util"
	"github.com/gin-gonic/gin"
)

// GetClientStatistic
//
//	@Tags			statistic_client
//	@Summary		获取客户端统计数据
//	@Description	获取客户端统计数据
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Param			startDate	query		string	true	"开始时间（格式yyyy-mm-dd）"
//	@Param			endDate		query		string	true	"结束时间（格式yyyy-mm-dd）"
//	@Success		200			{object}	response.Response{data=response.ClientStatistic}
//	@Router			/statistic/client [get]
func GetClientStatistic(ctx *gin.Context) {
	resp, err := service.GetClientStatistic(ctx, ctx.Query("startDate"), ctx.Query("endDate"))
	gin_util.Response(ctx, resp, err)
}
