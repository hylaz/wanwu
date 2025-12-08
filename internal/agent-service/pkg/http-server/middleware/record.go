package middleware

import (
	"encoding/json"
	"io"
	"strings"
	"time"

	err_code "github.com/UnicomAI/wanwu/api/proto/err-code"
	agent_log "github.com/UnicomAI/wanwu/internal/agent-service/pkg/agent-log"
	gin_util "github.com/UnicomAI/wanwu/pkg/gin-util"
	"github.com/gin-gonic/gin"
)

func Record(ctx *gin.Context) {
	start := time.Now().UnixMilli()
	var req string
	var err error
	if ctx.ContentType() == gin.MIMEJSON {
		if req, err = requestBody(ctx); err != nil {
			agent_log.LogAccessPB(ctx, requestFullPath(ctx), ctx.Request.Method, req, nil, err, start)
			gin_util.ResponseErrCodeKey(ctx, err_code.Code_BFFInvalidArg, "", err.Error())
			ctx.Abort()
			return
		}
	}
	ctx.Next()

	resp := ctx.GetString(gin_util.RESULT)
	agent_log.LogAccessPB(ctx, requestFullPath(ctx), ctx.Request.Method, req, resp, nil, start)
}

func requestFullPath(ctx *gin.Context) string {
	if ctx.Request.URL.RawQuery != "" {
		return ctx.Request.URL.Path + "?" + ctx.Request.URL.RawQuery
	}
	return ctx.Request.URL.Path
}

func requestBody(ctx *gin.Context) (string, error) {
	var body []byte
	var err error
	if cb, ok := ctx.Get(gin.BodyBytesKey); ok {
		if cbb, ok := cb.([]byte); ok {
			body = cbb
		}
	}
	if body == nil {
		body, err = io.ReadAll(ctx.Request.Body)
		if err != nil {
			return "", err
		}
		ctx.Set(gin.BodyBytesKey, body)
	}

	// avoid err: unexpected end of JSON input
	if strings.TrimSpace(string(body)) == "" {
		return "", nil
	}

	kv := make(map[string]interface{})
	if err = json.Unmarshal(body, &kv); err != nil {
		return "", err
	}
	if b, err := json.Marshal(kv); err != nil {
		return "", err
	} else {
		return string(b), nil
	}
}
