package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func StaticFSHeader(ctx *gin.Context) {
	// gin静态文件服务，修正response header Content-Type
	if strings.Contains(ctx.Request.URL.Path, "/v1/static") || strings.Contains(ctx.Request.URL.Path, "/v1/cache") {
		if strings.Contains(ctx.Request.URL.Path, ".csv") {
			ctx.Header("Content-Type", "text/csv; charset=utf-8")
		}
	}
	ctx.Next()
}
