package v1

import (
	"net/http"

	v1 "github.com/UnicomAI/wanwu/internal/bff-service/server/http/handler/v1"
	mid "github.com/UnicomAI/wanwu/pkg/gin-util/mid-wrap"
	"github.com/gin-gonic/gin"
)

func registerStatisticClient(apiV1 *gin.RouterGroup) {
	mid.Sub("statistic_client").Reg(apiV1, "/statistic/client", http.MethodGet, v1.GetClientStatistic, "获取使用工作流模板客户端统计")
}
