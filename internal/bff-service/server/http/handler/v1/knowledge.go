package v1

import (
	"github.com/UnicomAI/wanwu/internal/bff-service/model/request"
	"github.com/UnicomAI/wanwu/internal/bff-service/service"
	gin_util "github.com/UnicomAI/wanwu/pkg/gin-util"
	"github.com/gin-gonic/gin"
)

// GetKnowledgeSelect
//
//	@Tags			knowledge
//	@Summary		查询知识库列表
//	@Description	查询知识库列表
//	@Security		JWT
//	@Accept			json
//	@Param			data	body	request.KnowledgeSelectReq	true	"查询知识库列表"
//	@Produce		json
//	@Success		200	{object}	response.Response{data=response.KnowledgeListResp}
//	@Router			/knowledge/select [post]
func GetKnowledgeSelect(ctx *gin.Context) {
	userId, orgId := getUserID(ctx), getOrgID(ctx)
	var req request.KnowledgeSelectReq
	if !gin_util.Bind(ctx, &req) {
		return
	}
	resp, err := service.SelectKnowledgeList(ctx, userId, orgId, &req)
	gin_util.Response(ctx, resp, err)
}

// CreateKnowledge
//
//	@Tags			knowledge
//	@Summary		创建知识库
//	@Description	创建知识库
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Param			data	body		request.CreateKnowledgeReq	true	"创建知识库请求参数"
//	@Success		200		{object}	response.Response{data=[]response.CreateKnowledgeResp}
//	@Router			/knowledge [post]
func CreateKnowledge(ctx *gin.Context) {
	userId, orgId := getUserID(ctx), getOrgID(ctx)
	var req request.CreateKnowledgeReq
	if !gin_util.Bind(ctx, &req) {
		return
	}
	resp, err := service.CreateKnowledge(ctx, userId, orgId, &req)
	gin_util.Response(ctx, resp, err)
}

// UpdateKnowledge
//
//	@Tags			knowledge
//	@Summary		修改知识库（文档分类）
//	@Description	修改知识库（文档分类）
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Param			data	body		request.UpdateKnowledgeReq	true	"修改知识库请求参数"
//	@Success		200		{object}	response.Response
//	@Router			/knowledge [put]
func UpdateKnowledge(ctx *gin.Context) {
	userId, orgId := getUserID(ctx), getOrgID(ctx)
	var req request.UpdateKnowledgeReq
	if !gin_util.Bind(ctx, &req) {
		return
	}
	err := service.UpdateKnowledge(ctx, userId, orgId, &req)
	gin_util.Response(ctx, nil, err)
}

// DeleteKnowledge
//
//	@Tags			knowledge
//	@Summary		删除知识库
//	@Description	删除知识库
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Param			data	body		request.DeleteKnowledge	true	"删除知识库请求参数"
//	@Success		200		{object}	response.Response
//	@Router			/knowledge [delete]
func DeleteKnowledge(ctx *gin.Context) {
	userId, orgId := getUserID(ctx), getOrgID(ctx)
	var req request.DeleteKnowledge
	if !gin_util.Bind(ctx, &req) {
		return
	}
	resp, err := service.DeleteKnowledge(ctx, userId, orgId, &req)
	gin_util.Response(ctx, resp, err)
}

// KnowledgeHit
//
//	@Tags			knowledge
//	@Summary		知识库命中测试
//	@Description	知识库命中测试
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Param			data	body		request.KnowledgeHitReq	true	"知识库命中测试请求参数"
//	@Success		200		{object}	response.Response{data=response.KnowledgeHitResp}
//	@Router			/knowledge/hit [post]
func KnowledgeHit(ctx *gin.Context) {
	userId, orgId := getUserID(ctx), getOrgID(ctx)
	var req request.KnowledgeHitReq
	if !gin_util.Bind(ctx, &req) {
		return
	}
	resp, err := service.KnowledgeHit(ctx, userId, orgId, &req)
	gin_util.Response(ctx, resp, err)
}

// GetKnowledgeMetaKeySelect
//
//	@Tags			knowledge
//	@Summary		获取知识库元数据
//	@Description	获取知识库元数据
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Param			data	query		request.GetKnowledgeMetaSelectReq	true	"获取知识库元数据请求参数"
//	@Success		200		{object}	response.Response{data=response.GetKnowledgeMetaSelectResp}
//	@Router			/knowledge/meta/select [get]
func GetKnowledgeMetaKeySelect(ctx *gin.Context) {
	userId, orgId := getUserID(ctx), getOrgID(ctx)
	var req request.GetKnowledgeMetaSelectReq
	if !gin_util.BindQuery(ctx, &req) {
		return
	}
	resp, err := service.GetKnowledgeMetaSelect(ctx, userId, orgId, &req)
	gin_util.Response(ctx, resp, err)
}

// GetKnowledgeMetaValueList
//
//	@Tags			knowledge
//	@Summary		获取文档元数据列表
//	@Description	获取文档元数据列表
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Param			data	body		request.KnowledgeMetaValueListReq	true	"文档元数据列表请求参数"
//	@Success		200		{object}	response.Response{data=response.KnowledgeMetaValueListResp}
//	@Router			/knowledge/meta/value/list [post]
func GetKnowledgeMetaValueList(ctx *gin.Context) {
	userId, orgId := getUserID(ctx), getOrgID(ctx)
	var req request.KnowledgeMetaValueListReq
	if !gin_util.Bind(ctx, &req) {
		return
	}
	resp, err := service.GetKnowledgeMetaValueList(ctx, userId, orgId, &req)
	gin_util.Response(ctx, resp, err)
}

// UpdateKnowledgeMetaValue
//
//	@Tags			knowledge
//	@Summary		更新知识库元数据值
//	@Description	更新知识库元数据值
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Param			data	body		request.UpdateMetaValueReq	true	"更新知识库元数据值请求参数"
//	@Success		200		{object}	response.Response
//	@Router			/knowledge/meta/value/update [post]
func UpdateKnowledgeMetaValue(ctx *gin.Context) {
	userId, orgId := getUserID(ctx), getOrgID(ctx)
	var req request.UpdateMetaValueReq
	if !gin_util.Bind(ctx, &req) {
		return
	}
	err := service.UpdateKnowledgeMetaValue(ctx, userId, orgId, &req)
	gin_util.Response(ctx, nil, err)
}
