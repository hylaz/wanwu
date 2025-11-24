package v1

import (
	err_code "github.com/UnicomAI/wanwu/api/proto/err-code"
	"github.com/UnicomAI/wanwu/internal/bff-service/model/request"
	"github.com/UnicomAI/wanwu/internal/bff-service/service"
	gin_util "github.com/UnicomAI/wanwu/pkg/gin-util"
	grpc_util "github.com/UnicomAI/wanwu/pkg/grpc-util"
	"github.com/gin-gonic/gin"
)

// GetRoleTemplate
//
//	@Tags			permission.role
//	@Summary		获取角色模板（用于创建角色）
//	@Description	获取当前用户在X-Org-Id组织的角色模板
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	response.Response{data=response.RoleTemplate}
//	@Router			/role/template [get]
func GetRoleTemplate(ctx *gin.Context) {
	resp, err := service.GetRoleTemplate(ctx, getUserID(ctx), getOrgID(ctx))
	gin_util.Response(ctx, resp, err)
}

// CreateRole
//
//	@Tags			permission.role
//	@Summary		创建角色
//	@Description	创建X-Org-Id组织的角色
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Param			data	body		request.RoleCreate	true	"角色信息"
//	@Success		200		{object}	response.Response{data=response.RoleID}
//	@Router			/role [post]
func CreateRole(ctx *gin.Context) {
	var req request.RoleCreate
	if !gin_util.Bind(ctx, &req) {
		return
	}
	if !isAdmin(ctx) {
		gin_util.Response(ctx, nil, grpc_util.ErrorStatusWithKey(err_code.Code_BFFGeneral, "bff_role_cannot_create"))
		return
	}
	resp, err := service.CreateRole(ctx, getUserID(ctx), getOrgID(ctx), &req)
	gin_util.Response(ctx, resp, err)
}

// ChangeRole
//
//	@Tags			permission.role
//	@Summary		编辑角色
//	@Description	编辑X-Org-Id组织的角色
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Param			data	body		request.RoleUpdate	true	"角色信息"
//	@Success		200		{object}	response.Response
//	@Router			/role [put]
func ChangeRole(ctx *gin.Context) {
	var req request.RoleUpdate
	if !gin_util.Bind(ctx, &req) {
		return
	}
	if !isAdmin(ctx) {
		gin_util.Response(ctx, nil, grpc_util.ErrorStatusWithKey(err_code.Code_BFFGeneral, "bff_role_cannot_change"))
		return
	}
	err := service.ChangeRole(ctx, getUserID(ctx), getOrgID(ctx), &req)
	gin_util.Response(ctx, nil, err)
}

// DeleteRole
//
//	@Tags			permission.role
//	@Summary		删除角色
//	@Description	删除X-Org-Id组织的角色
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Param			data	body		request.RoleID	true	"角色ID"
//	@Success		200		{object}	response.Response
//	@Router			/role [delete]
func DeleteRole(ctx *gin.Context) {
	var req request.RoleID
	if !gin_util.Bind(ctx, &req) {
		return
	}
	if !isAdmin(ctx) {
		gin_util.Response(ctx, nil, grpc_util.ErrorStatusWithKey(err_code.Code_BFFGeneral, "bff_role_cannot_delete"))
		return
	}
	err := service.DeleteRole(ctx, getOrgID(ctx), req.RoleID)
	gin_util.Response(ctx, nil, err)
}

// GetRoleInfo
//
//	@Tags			permission.role
//	@Summary		获取角色信息
//	@Description	获取X-Org-Id组织的角色信息
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Param			roleId	query		string	true	"角色ID"
//	@Success		200		{object}	response.Response{data=response.RoleInfo}
//	@Router			/role/info [get]
func GetRoleInfo(ctx *gin.Context) {
	resp, err := service.GetRoleInfo(ctx, getUserID(ctx), getOrgID(ctx), ctx.Query("roleId"))
	gin_util.Response(ctx, resp, err)
}

// GetRoleList
//
//	@Tags			permission.role
//	@Summary		获取角色列表
//	@Description	获取X-Org-Id组织的角色列表
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Param			name		query		string	false	"角色名(模糊查询)"
//	@Param			pageNo		query		int		true	"页面编号，从1开始"
//	@Param			pageSize	query		int		true	"单页数量，从1开始"
//	@Success		200			{object}	response.Response{data=response.PageResult{list=[]response.RoleInfo}}
//	@Router			/role/list [get]
func GetRoleList(ctx *gin.Context) {
	resp, err := service.GetRoleList(ctx, getUserID(ctx), getOrgID(ctx), ctx.Query("name"), getPageNo(ctx), getPageSize(ctx))
	gin_util.Response(ctx, resp, err)
}

// ChangeRoleStatus
//
//	@Tags			permission.role
//	@Summary		修改角色状态
//	@Description	修改X-Org-Id组织的角色状态
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Param			data	body		request.RoleStatus	true	"角色信息"
//	@Success		200		{object}	response.Response
//	@Router			/role/status [put]
func ChangeRoleStatus(ctx *gin.Context) {
	var req request.RoleStatus
	if !gin_util.Bind(ctx, &req) {
		return
	}
	if !isAdmin(ctx) {
		gin_util.Response(ctx, nil, grpc_util.ErrorStatusWithKey(err_code.Code_BFFGeneral, "bff_role_cannot_change_status"))
		return
	}
	err := service.ChangeRoleStatus(ctx, getOrgID(ctx), req.RoleID.RoleID, req.Status)
	gin_util.Response(ctx, nil, err)
}
