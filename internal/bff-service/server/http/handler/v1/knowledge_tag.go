package v1

import (
	"github.com/UnicomAI/wanwu/internal/bff-service/model/request"
	"github.com/UnicomAI/wanwu/internal/bff-service/service"
	gin_util "github.com/UnicomAI/wanwu/pkg/gin-util"
	"github.com/gin-gonic/gin"
)

// GetKnowledgeTagSelect
//
//	@Tags			knowledge.tag
//	@Summary		查询知识库标签列表
//	@Description	查询知识库标签列表
//	@Security		JWT
//	@Accept			json
//	@Param			data	query	request.KnowledgeTagSelectReq	true	"查询知识库请求参数"
//	@Produce		json
//	@Success		200	{object}	response.Response{data=response.KnowledgeTagListResp}
//	@Router			/knowledge/tag [get]
func GetKnowledgeTagSelect(ctx *gin.Context) {
	userId, orgId := getUserID(ctx), getOrgID(ctx)
	var req request.KnowledgeTagSelectReq
	if !gin_util.BindQuery(ctx, &req) {
		return
	}
	resp, err := service.SelectKnowledgeTagList(ctx, userId, orgId, &req)
	gin_util.Response(ctx, resp, err)
}

// SelectTagBindCount
//
//	@Tags			knowledge.tag
//	@Summary		查询标签绑定知识库数量
//	@Description	查询标签绑定知识库数量
//	@Security		JWT
//	@Accept			json
//	@Param			data	query	request.TagBindCountReq	true	"查询tag绑定数量参数请求参数"
//	@Produce		json
//	@Success		200	{object}	response.Response{data=response.KnowledgeTagListResp}
//	@Router			/knowledge/tag/bind/count [get]
func SelectTagBindCount(ctx *gin.Context) {
	userId, orgId := getUserID(ctx), getOrgID(ctx)
	var req request.TagBindCountReq
	if !gin_util.BindQuery(ctx, &req) {
		return
	}
	resp, err := service.SelectTagBindCount(ctx, userId, orgId, &req)
	gin_util.Response(ctx, resp, err)
}

// CreateKnowledgeTag
//
//	@Tags			knowledge.tag
//	@Summary		创建知识库标签
//	@Description	创建知识库标签
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Param			data	body		request.CreateKnowledgeTagReq	true	"创建知识库标签请求参数"
//	@Success		200		{object}	response.Response{data=response.CreateKnowledgeTagResp}
//	@Router			/knowledge/tag [post]
func CreateKnowledgeTag(ctx *gin.Context) {
	userId, orgId := getUserID(ctx), getOrgID(ctx)
	var req request.CreateKnowledgeTagReq
	if !gin_util.Bind(ctx, &req) {
		return
	}
	resp, err := service.CreateKnowledgeTag(ctx, userId, orgId, &req)
	gin_util.Response(ctx, resp, err)
}

// UpdateKnowledgeTag
//
//	@Tags			knowledge.tag
//	@Summary		修改知识库标签
//	@Description	修改知识库标签
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Param			data	body		request.UpdateKnowledgeTagReq	true	"修改知识库标签请求参数"
//	@Success		200		{object}	response.Response
//	@Router			/knowledge/tag [put]
func UpdateKnowledgeTag(ctx *gin.Context) {
	userId, orgId := getUserID(ctx), getOrgID(ctx)
	var req request.UpdateKnowledgeTagReq
	if !gin_util.Bind(ctx, &req) {
		return
	}
	err := service.UpdateKnowledgeTag(ctx, userId, orgId, &req)
	gin_util.Response(ctx, nil, err)
}

// DeleteKnowledgeTag
//
//	@Tags			knowledge.tag
//	@Summary		删除知识库标签
//	@Description	删除知识库标签
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Param			data	body		request.DeleteKnowledgeTagReq	true	"删除知识库标签请求参数"
//	@Success		200		{object}	response.Response
//	@Router			/knowledge/tag [delete]
func DeleteKnowledgeTag(ctx *gin.Context) {
	userId, orgId := getUserID(ctx), getOrgID(ctx)
	var req request.DeleteKnowledgeTagReq
	if !gin_util.Bind(ctx, &req) {
		return
	}
	err := service.DeleteKnowledgeTag(ctx, userId, orgId, &req)
	gin_util.Response(ctx, nil, err)
}

// BindKnowledgeTag
//
//	@Tags			knowledge.tag
//	@Summary		绑定知识库标签
//	@Description	绑定知识库标签
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Param			data	body		request.BindKnowledgeTagReq	true	"绑定知识库标签请求参数"
//	@Success		200		{object}	response.Response
//	@Router			/knowledge/tag/bind [post]
func BindKnowledgeTag(ctx *gin.Context) {
	userId, orgId := getUserID(ctx), getOrgID(ctx)
	var req request.BindKnowledgeTagReq
	if !gin_util.Bind(ctx, &req) {
		return
	}
	err := service.BindKnowledgeTag(ctx, userId, orgId, &req)
	gin_util.Response(ctx, nil, err)
}
