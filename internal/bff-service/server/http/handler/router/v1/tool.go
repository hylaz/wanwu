package v1

import (
	"net/http"

	v1 "github.com/UnicomAI/wanwu/internal/bff-service/server/http/handler/v1"
	mid "github.com/UnicomAI/wanwu/pkg/gin-util/mid-wrap"
	"github.com/gin-gonic/gin"
)

func registerTool(apiV1 *gin.RouterGroup) {
	// MCP
	mid.Sub("tool").Reg(apiV1, "/mcp", http.MethodPost, v1.CreateMCP, "创建自定义MCP")
	mid.Sub("tool").Reg(apiV1, "/mcp", http.MethodPut, v1.UpdateMCP, "修改自定义MCP")
	mid.Sub("tool").Reg(apiV1, "/mcp", http.MethodGet, v1.GetMCP, "获取自定义MCP详情")
	mid.Sub("tool").Reg(apiV1, "/mcp", http.MethodDelete, v1.DeleteMCP, "删除自定义MCP")
	mid.Sub("tool").Reg(apiV1, "/mcp/list", http.MethodGet, v1.GetMCPList, "获取MCP自定义列表")
	mid.Sub("tool").Reg(apiV1, "/mcp/tool/list", http.MethodGet, v1.GetMCPTools, "获取MCP Tool列表")

	// 自定义工具
	mid.Sub("tool").Reg(apiV1, "/tool/custom", http.MethodPost, v1.CreateCustomTool, "创建自定义工具")
	mid.Sub("tool").Reg(apiV1, "/tool/custom", http.MethodGet, v1.GetCustomTool, "获取自定义工具详情")
	mid.Sub("tool").Reg(apiV1, "/tool/custom", http.MethodDelete, v1.DeleteCustomTool, "删除自定义工具")
	mid.Sub("tool").Reg(apiV1, "/tool/custom", http.MethodPut, v1.UpdateCustomTool, "修改自定义工具")
	mid.Sub("tool").Reg(apiV1, "/tool/custom/list", http.MethodGet, v1.GetCustomToolList, "获取自定义工具列表")
	mid.Sub("tool").Reg(apiV1, "/tool/custom/select", http.MethodGet, v1.GetCustomToolSelect, "获取自定义工具列表（用于下拉选择）")
	mid.Sub("tool").Reg(apiV1, "/tool/custom/schema", http.MethodPost, v1.GetCustomToolActions, "获取可用API列表（根据Schema）")

	// 内置工具
	mid.Sub("tool").Reg(apiV1, "/tool/square", http.MethodGet, v1.GetToolSquareDetail, "获取内置工具详情")
	mid.Sub("tool").Reg(apiV1, "/tool/square/list", http.MethodGet, v1.GetToolSquareList, "获取内置工具列表")
	mid.Sub("tool").Reg(apiV1, "/tool/builtin", http.MethodPost, v1.UpdateToolSquareAPIKey, "修改内置工具")
}
