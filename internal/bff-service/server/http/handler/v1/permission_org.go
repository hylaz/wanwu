package v1

import (
	err_code "github.com/UnicomAI/wanwu/api/proto/err-code"
	"github.com/UnicomAI/wanwu/internal/bff-service/model/request"
	"github.com/UnicomAI/wanwu/internal/bff-service/service"
	gin_util "github.com/UnicomAI/wanwu/pkg/gin-util"
	grpc_util "github.com/UnicomAI/wanwu/pkg/grpc-util"
	"github.com/gin-gonic/gin"
)

// CreateOrg
//
//	@Tags			permission.org
//	@Summary		创建下级组织
//	@Description	创建X-Org-Id组织的下级组织
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Param			data	body		request.OrgCreate	true	"组织信息"
//	@Success		200		{object}	response.Response{data=response.OrgID}
//	@Router			/org [post]
func CreateOrg(ctx *gin.Context) {
	var req request.OrgCreate
	if !gin_util.Bind(ctx, &req) {
		return
	}
	if !isAdmin(ctx) {
		gin_util.Response(ctx, nil, grpc_util.ErrorStatusWithKey(err_code.Code_BFFGeneral, "bff_org_cannot_create"))
		return
	}
	resp, err := service.CreateOrg(ctx, getUserID(ctx), getOrgID(ctx), &req)
	gin_util.Response(ctx, resp, err)
}

// ChangeOrg
//
//	@Tags		permission.org
//	@Summary	编辑下级组织
//	@Security	JWT
//	@Accept		json
//	@Produce	json
//	@Param		data	body		request.OrgUpdate	true	"组织信息"
//	@Success	200		{object}	response.Response
//	@Router		/org [put]
func ChangeOrg(ctx *gin.Context) {
	var req request.OrgUpdate
	if !gin_util.Bind(ctx, &req) {
		return
	}
	if !isAdmin(ctx) {
		gin_util.Response(ctx, nil, grpc_util.ErrorStatusWithKey(err_code.Code_BFFGeneral, "bff_org_cannot_change"))
		return
	}
	err := service.ChangeOrg(ctx, getOrgID(ctx), &req)
	gin_util.Response(ctx, nil, err)
}

// DeleteOrg
//
//	@Tags		permission.org
//	@Summary	删除下级组织
//	@Security	JWT
//	@Accept		json
//	@Produce	json
//	@Param		data	body		request.OrgID	true	"组织ID"
//	@Success	200		{object}	response.Response
//	@Router		/org [delete]
func DeleteOrg(ctx *gin.Context) {
	var req request.OrgID
	if !gin_util.Bind(ctx, &req) {
		return
	}
	if !isAdmin(ctx) {
		gin_util.Response(ctx, nil, grpc_util.ErrorStatusWithKey(err_code.Code_BFFGeneral, "bff_org_cannot_delete"))
		return
	}
	err := service.DeleteOrg(ctx, getOrgID(ctx), req.OrgID)
	gin_util.Response(ctx, nil, err)
}

// GetOrgInfo
//
//	@Tags		permission.org
//	@Summary	获取组织信息
//	@Security	JWT
//	@Accept		json
//	@Produce	json
//	@Param		orgId	query		string	true	"组织ID"
//	@Success	200		{object}	response.Response{data=response.OrgInfo}
//	@Router		/org/info [get]
func GetOrgInfo(ctx *gin.Context) {
	resp, err := service.GetOrgInfo(ctx, ctx.Query("orgId"))
	gin_util.Response(ctx, resp, err)
}

// GetOrgList
//
//	@Tags			permission.org
//	@Summary		获取下级组织列表
//	@Description	获取X-Org-Id组织的下级组织列表
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Param			name		query		string	false	"组织名(模糊查询)"
//	@Param			pageNo		query		int		true	"页面编号，从1开始"
//	@Param			pageSize	query		int		true	"单页数量，从1开始"
//	@Success		200			{object}	response.Response{data=response.PageResult{list=[]response.OrgInfo}}
//	@Router			/org/list [get]
func GetOrgList(ctx *gin.Context) {
	resp, err := service.GetOrgList(ctx, getOrgID(ctx), ctx.Query("name"), getPageNo(ctx), getPageSize(ctx))
	gin_util.Response(ctx, resp, err)
}

// ChangeOrgStatus
//
//	@Tags		permission.org
//	@Summary	修改下级组织状态
//	@Security	JWT
//	@Accept		json
//	@Produce	json
//	@Param		data	body		request.OrgStatus	true	"组织信息"
//	@Success	200		{object}	response.Response
//	@Router		/org/status [put]
func ChangeOrgStatus(ctx *gin.Context) {
	var req request.OrgStatus
	if !gin_util.Bind(ctx, &req) {
		return
	}
	if !isAdmin(ctx) {
		gin_util.Response(ctx, nil, grpc_util.ErrorStatusWithKey(err_code.Code_BFFGeneral, "bff_org_cannot_change_status"))
		return
	}
	err := service.ChangeOrgStatus(ctx, getOrgID(ctx), req.OrgID.OrgID, req.Status)
	gin_util.Response(ctx, nil, err)
}
