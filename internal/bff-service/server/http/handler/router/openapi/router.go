package openapi

import (
	"net/http"

	"github.com/UnicomAI/wanwu/internal/bff-service/server/http/handler/openapi"
	"github.com/UnicomAI/wanwu/internal/bff-service/server/http/middleware"
	"github.com/UnicomAI/wanwu/pkg/constant"
	mid "github.com/UnicomAI/wanwu/pkg/gin-util/mid-wrap"
	"github.com/gin-gonic/gin"
)

func Register(openAPI *gin.RouterGroup) {
	// openapi
	mid.Sub("openapi").Reg(openAPI, "/agent/conversation", http.MethodPost, openapi.CreateAgentConversation, "智能体创建对话OpenAPI", middleware.AuthOpenAPI(constant.AppTypeAgent))
	mid.Sub("openapi").Reg(openAPI, "/agent/chat", http.MethodPost, openapi.ChatAgent, "智能体问答OpenAPI", middleware.AuthOpenAPI(constant.AppTypeAgent))
	mid.Sub("openapi").Reg(openAPI, "/rag/chat", http.MethodPost, openapi.ChatRag, "文本问答OpenAPI", middleware.AuthOpenAPI(constant.AppTypeRag))
	mid.Sub("openapi").Reg(openAPI, "/workflow/run", http.MethodPost, openapi.WorkflowRun, "工作流OpenAPI", middleware.AuthOpenAPI(constant.AppTypeWorkflow))
	mid.Sub("openapi").Reg(openAPI, "/chatflow/conversation", http.MethodPost, openapi.CreateChatflowConversation, "对话流创建对话OpenAPI", middleware.AuthOpenAPI(constant.AppTypeChatflow))
	mid.Sub("openapi").Reg(openAPI, "/chatflow/chat", http.MethodPost, openapi.ChatflowChat, "对话流OpenAPI", middleware.AuthOpenAPI(constant.AppTypeChatflow))
	mid.Sub("openapi").Reg(openAPI, "/workflow/file/upload", http.MethodPost, openapi.WorkflowFileUpload, "工作流OpenAPI文件上传", middleware.AuthOpenAPI(constant.AppTypeWorkflow))
	mid.Sub("openapi").Reg(openAPI, "/mcp/server/sse", http.MethodGet, openapi.GetMCPServerSSE, "新建MCP服务sse连接", middleware.AuthOpenAPIByQuery(constant.AppTypeMCPServer))
	mid.Sub("openapi").Reg(openAPI, "/mcp/server/message", http.MethodPost, openapi.GetMCPServerMessage, "获取MCP服务sse消息", middleware.AuthOpenAPIByQuery(constant.AppTypeMCPServer))
	mid.Sub("openapi").Reg(openAPI, "/mcp/server/streamable", http.MethodGet, openapi.GetMCPServerStreamable, "获取MCP服务streamable消息(GET)", middleware.AuthOpenAPIByQuery(constant.AppTypeMCPServer))
	mid.Sub("openapi").Reg(openAPI, "/mcp/server/streamable", http.MethodPost, openapi.GetMCPServerStreamable, "获取MCP服务streamable消息(POST)", middleware.AuthOpenAPIByQuery(constant.AppTypeMCPServer))

	// oauth
	mid.Sub("openapi").Reg(openAPI, "/oauth/jwks", http.MethodGet, openapi.OAuthJWKS, "JWT公钥")
	mid.Sub("openapi").Reg(openAPI, "/oauth/login", http.MethodGet, openapi.OAuthLogin, "OAuth登录授权")
	mid.Sub("openapi").Reg(openAPI, "/oauth/code/token", http.MethodPost, openapi.OAuthToken, "授权码获取token")
	mid.Sub("openapi").Reg(openAPI, "/oauth/code/token/refresh", http.MethodPost, openapi.OAuthRefresh, "刷新Access Token")
	mid.Sub("openapi").Reg(openAPI, "/.well-known/openid-configuration", http.MethodGet, openapi.OAuthConfig, "返回Endpoint配置")
	// oauth user
	mid.Sub("openapi").Reg(openAPI, "/oauth/userinfo", http.MethodGet, openapi.OAuthGetUserInfo, "OAuth获取用户信息", middleware.JWTOAuthAccess)
}
