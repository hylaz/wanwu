package v1

import (
	"net/http"

	"github.com/UnicomAI/wanwu/internal/bff-service/server/http/handler/callback"
	mid "github.com/UnicomAI/wanwu/pkg/gin-util/mid-wrap"
	"github.com/gin-gonic/gin"
)

func registerV1Callback(callbackV1 *gin.RouterGroup) {
	mid.Sub("callback").Reg(callbackV1, "/api/docstatus", http.MethodPost, callback.UpdateDocStatus, "算法更新知识库文档状态（模型扩展调用）")
	mid.Sub("callback").Reg(callbackV1, "/api/deploy/info", http.MethodGet, callback.GetDeployInfo, "获取Maas平台部署信息（模型扩展调用）")
	mid.Sub("callback").Reg(callbackV1, "/api/category/info", http.MethodGet, callback.SelectKnowledgeInfoByName, "获取Maas平台知识库信息（模型扩展调用）")
	mid.Sub("callback").Reg(callbackV1, "/api/doc_status_init", http.MethodGet, callback.DocStatusInit, "将正在解析的文档设置为解析失败")
	mid.Sub("callback").Reg(callbackV1, "/api/knowledge/status", http.MethodPost, callback.UpdateKnowledgeStatus, "更新知识库相关状态")
}
