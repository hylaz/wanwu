package middleware

import (
	"github.com/UnicomAI/wanwu/internal/bff-service/service"
	"github.com/UnicomAI/wanwu/pkg/log"
	"github.com/gin-gonic/gin"
)

func AppHistoryRecord(filedId, appType string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		appID := getFieldValue(ctx, filedId)
		userID, _ := getUserID(ctx)
		ctx.Next()
		if appID == "" || userID == "" || appType == "" {
			log.Errorf("record user %v app %v type %v history err", userID, appID, appType)
			return
		}
		if err := service.AddAppHistoryRecord(ctx, userID, appID, appType); err != nil {
			log.Errorf("record user %v app %v type %v history err: %v", userID, appID, appType, err)
		}
	}
}
