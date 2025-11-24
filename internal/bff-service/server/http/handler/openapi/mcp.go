package openapi

import (
	"github.com/UnicomAI/wanwu/internal/bff-service/service"
	gin_util "github.com/UnicomAI/wanwu/pkg/gin-util"
	"github.com/gin-gonic/gin"
)

// GetMCPServerSSE
//
//	@Tags			openapi
//	@Summary		获取MCPServer SSE
//	@Description	获取MCPServer SSE
//	@Accept			json
//	@Produce		json
//	@Param			key	query		string	true	"key"
//	@Success		200	{object}	response.Response{}
//	@Router			/mcp/server/sse [get]
func GetMCPServerSSE(ctx *gin.Context) {
	err := service.GetMCPServerSSE(ctx, getAppID(ctx), ctx.Query("key"))
	if err != nil {
		gin_util.Response(ctx, nil, err)
		return
	}
}

// GetMCPServerMessage
//
//	@Tags			openapi
//	@Summary		获取MCPServer Message
//	@Description	获取MCPServer Message
//	@Accept			json
//	@Produce		json
//	@Param			id	query		string	true	"mcpServerId"
//	@Success		200	{object}	response.Response{}
//	@Router			/mcp/server/message [post]
func GetMCPServerMessage(ctx *gin.Context) {
	err := service.GetMCPServerMessage(ctx, getAppID(ctx))
	if err != nil {
		gin_util.Response(ctx, nil, err)
		return
	}
}

// GetMCPServerStreamable
//
//	@Tags			openapi
//	@Summary		获取MCPServer streamable 类型消息
//	@Description	获取MCPServer streamable 类型消息
//	@Accept			json
//	@Produce		json
//	@Param			key	query		string	true	"key"
//	@Success		200	{object}	response.Response{}
//	@Router			/mcp/server/streamable [post]
func GetMCPServerStreamable(ctx *gin.Context) {
	err := service.GetMCPServerStreamable(ctx, getAppID(ctx))
	if err != nil {
		gin_util.Response(ctx, nil, err)
		return
	}
}
