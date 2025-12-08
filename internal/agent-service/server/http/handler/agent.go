package handler

import (
	"github.com/UnicomAI/wanwu/internal/agent-service/model/request"
	"github.com/UnicomAI/wanwu/internal/agent-service/service"
	gin_util "github.com/UnicomAI/wanwu/pkg/gin-util"
	"github.com/gin-gonic/gin"
)

func AgentChat(ctx *gin.Context) {
	var req request.AgentChatReq
	if !gin_util.Bind(ctx, &req) {
		return
	}
	err := service.AgentChat(ctx, &req)
	gin_util.Response(ctx, nil, err)
}
