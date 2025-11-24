package v1

import (
	"net/http"

	v1 "github.com/UnicomAI/wanwu/internal/bff-service/server/http/handler/v1"
	mid "github.com/UnicomAI/wanwu/pkg/gin-util/mid-wrap"
	"github.com/gin-gonic/gin"
)

func registerWorkflow(apiV1 *gin.RouterGroup) {
	mid.Sub("workflow").Reg(apiV1, "/appspace/workflow", http.MethodPost, v1.CreateWorkflow, "创建workflow")
	mid.Sub("workflow").Reg(apiV1, "/appspace/workflow/copy", http.MethodPost, v1.CopyWorkflow, "拷贝workflow")
	mid.Sub("workflow").Reg(apiV1, "/appspace/workflow/model/select/llm", http.MethodGet, v1.ListLlmModelsByWorkflow, "lm模型列表（用于workflow）")
	mid.Sub("workflow").Reg(apiV1, "/appspace/workflow/export", http.MethodGet, v1.ExportWorkflow, "导出workflow")
	mid.Sub("workflow").Reg(apiV1, "/appspace/workflow/import", http.MethodPost, v1.ImportWorkflow, "导入workflow")

	mid.Sub("workflow").Reg(apiV1, "/workflow/tool/action", http.MethodGet, v1.GetWorkflowToolDetail, "获取Tool具体action")
	mid.Sub("workflow").Reg(apiV1, "/workflow/tool/select", http.MethodGet, v1.GetWorkflowToolSelect, "获取Tool列表")

	mid.Sub("workflow").Reg(apiV1, "/workflow/select", http.MethodGet, v1.GetWorkflowSelect, "智能体工作流下拉列表接口")

	mid.Sub("workflow").Reg(apiV1, "/workflow/template", http.MethodPost, v1.CreateWorkflowByTemplate, "复制工作流模板")
	// --- chatflow ---
	mid.Sub("workflow").Reg(apiV1, "/appspace/chatflow", http.MethodPost, v1.CreateChatflow, "创建chatflow")
	mid.Sub("workflow").Reg(apiV1, "/appspace/chatflow/copy", http.MethodPost, v1.CopyChatflow, "拷贝chatflow")
	mid.Sub("workflow").Reg(apiV1, "/appspace/chatflow/import", http.MethodPost, v1.ImportChatflow, "导入chatflow")
	mid.Sub("workflow").Reg(apiV1, "/appspace/chatflow/export", http.MethodGet, v1.ExportChatflow, "导出chatflow")
}
