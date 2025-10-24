package v1

import (
	"net/http"

	v1 "github.com/UnicomAI/wanwu/internal/bff-service/server/http/handler/v1"
	"github.com/UnicomAI/wanwu/internal/bff-service/server/http/middleware"
	"github.com/UnicomAI/wanwu/pkg/constant"
	mid "github.com/UnicomAI/wanwu/pkg/gin-util/mid-wrap"
	"github.com/gin-gonic/gin"
)

func registerAssistant(apiV1 *gin.RouterGroup) {
	mid.Sub("agent").Reg(apiV1, "/assistant/template/list", http.MethodGet, v1.GetAssistantTemplateList, "智能体模板列表")
	mid.Sub("agent").Reg(apiV1, "/assistant/template", http.MethodGet, v1.GetAssistantTemplate, "获取智能体模板")
	mid.Sub("agent").Reg(apiV1, "/assistant/template", http.MethodPost, v1.AssistantTemplateCreate, "复制智能体模板")

	mid.Sub("agent").Reg(apiV1, "/assistant", http.MethodPost, v1.AssistantCreate, "创建智能体")
	mid.Sub("agent").Reg(apiV1, "/assistant", http.MethodPut, v1.AssistantUpdate, "修改智能体基本信息")
	mid.Sub("agent").Reg(apiV1, "/assistant/config", http.MethodPut, v1.AssistantConfigUpdate, "修改智能体配置信息")
	mid.Sub("agent").Reg(apiV1, "/assistant", http.MethodGet, v1.GetAssistantInfo, "查看智能体详情")
	mid.Sub("agent").Reg(apiV1, "/assistant/copy", http.MethodPost, v1.AssistantCopy, "智能体复制")

	mid.Sub("agent").Reg(apiV1, "/assistant/tool/workflow", http.MethodPost, v1.AssistantWorkFlowCreate, "添加工作流")
	mid.Sub("agent").Reg(apiV1, "/assistant/tool/workflow", http.MethodDelete, v1.AssistantWorkFlowDelete, "删除工作流")
	mid.Sub("agent").Reg(apiV1, "/assistant/tool/workflow/switch", http.MethodPut, v1.AssistantWorkFlowEnableSwitch, "启用/停用工作流")

	mid.Sub("agent").Reg(apiV1, "/assistant/tool/mcp", http.MethodPost, v1.AssistantMCPCreate, "添加mcp工具")
	mid.Sub("agent").Reg(apiV1, "/assistant/tool/mcp", http.MethodDelete, v1.AssistantMCPDelete, "删除mcp工具")
	mid.Sub("agent").Reg(apiV1, "/assistant/tool/mcp/switch", http.MethodPut, v1.AssistantMCPEnableSwitch, "启用/停用mcp")

	mid.Sub("agent").Reg(apiV1, "/assistant/tool/custom", http.MethodPost, v1.AssistantCustomToolCreate, "智能体添加自定义工具")
	mid.Sub("agent").Reg(apiV1, "/assistant/tool/custom", http.MethodDelete, v1.AssistantCustomToolDelete, "智能体删除自定义工具")
	mid.Sub("agent").Reg(apiV1, "/assistant/tool/custom/switch", http.MethodPut, v1.AssistantCustomToolEnableSwitch, "智能体启用/停用自定义工具")

	mid.Sub("agent").Reg(apiV1, "/assistant/conversation", http.MethodPost, v1.ConversationCreate, "创建智能体对话")
	mid.Sub("agent").Reg(apiV1, "/assistant/conversation", http.MethodDelete, v1.ConversationDelete, "删除智能体对话")
	mid.Sub("agent").Reg(apiV1, "/assistant/conversation/list", http.MethodGet, v1.GetConversationList, "智能体对话列表")
	mid.Sub("agent").Reg(apiV1, "/assistant/conversation/detail", http.MethodGet, v1.GetConversationDetailList, "智能体对话详情历史列表")

	mid.Sub("agent").Reg(apiV1, "/assistant/stream", http.MethodPost, v1.AssistantConversionStream, "智能体流式问答", middleware.AppHistoryRecord("assistantId", constant.AppTypeAgent))
}
