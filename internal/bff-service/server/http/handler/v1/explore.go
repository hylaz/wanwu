package v1

import (
	"github.com/UnicomAI/wanwu/internal/bff-service/model/request"
	"github.com/UnicomAI/wanwu/internal/bff-service/service"
	gin_util "github.com/UnicomAI/wanwu/pkg/gin-util"
	"github.com/gin-gonic/gin"
)

// GetExplorationAppList
//
//	@Tags			exploration
//	@Summary		获取应用广场应用
//	@Description	获取应用广场应用
//	@Accept			json
//	@Produce		json
//	@Param			data	query		request.GetExplorationAppListRequest	true	"获取应用广场应用参数"
//	@Success		200		{object}	response.Response{data=response.ListResult{list=[]response.ExplorationAppInfo}}
//	@Router			/exploration/app/list [get]
func GetExplorationAppList(ctx *gin.Context) {
	var req request.GetExplorationAppListRequest
	if !gin_util.BindQuery(ctx, &req) {
		return
	}
	resp, err := service.GetExplorationAppList(ctx, getUserID(ctx), getOrgID(ctx), req)
	gin_util.Response(ctx, resp, err)
}

// ChangeExplorationAppFavorite
//
//	@Tags			exploration
//	@Summary		更改App收藏状态
//	@Description	更改App收藏状态
//	@Accept			json
//	@Produce		json
//	@Param			data	body		request.ChangeExplorationAppFavoriteRequest	true	"更改App收藏状态参数"
//	@Success		200		{object}	response.Response
//	@Router			/exploration/app/favorite [post]
func ChangeExplorationAppFavorite(ctx *gin.Context) {
	var req request.ChangeExplorationAppFavoriteRequest
	if !gin_util.Bind(ctx, &req) {
		return
	}
	err := service.ChangeExplorationAppFavorite(ctx, getUserID(ctx), getOrgID(ctx), req)
	gin_util.Response(ctx, nil, err)
}
