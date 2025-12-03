package v1

import (
	"github.com/UnicomAI/wanwu/internal/bff-service/model/request"
	"github.com/UnicomAI/wanwu/internal/bff-service/service"
	gin_util "github.com/UnicomAI/wanwu/pkg/gin-util"
	"github.com/gin-gonic/gin"
)

// GreateKnowledgeQAPair
//
//	@Tags			knowledge.qa
//	@Summary		新增问答对
//	@Description	新增问答对
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Param			data	body		request.CreateKnowledgeQAPairReq	true	"新增问答对请求参数"
//	@Success		200		{object}	response.Response{data=response.CreateKnowledgeQAPairResp}
//	@Router			/knowledge/qa/pair [post]
func GreateKnowledgeQAPair(ctx *gin.Context) {
	userId, orgId := getUserID(ctx), getOrgID(ctx)
	var req request.CreateKnowledgeQAPairReq
	if !gin_util.Bind(ctx, &req) {
		return
	}
	resp, err := service.CreateKnowledgeQAPair(ctx, userId, orgId, &req)
	gin_util.Response(ctx, resp, err)
}

// UpdateKnowledgeQAPair
//
//	@Tags			knowledge.qa
//	@Summary		编辑问答对
//	@Description	编辑问答对
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Param			data	body		request.UpdateKnowledgeQAPairReq	true	"编辑问答对请求参数"
//	@Success		200		{object}	response.Response
//	@Router			/knowledge/qa/pair [put]
func UpdateKnowledgeQAPair(ctx *gin.Context) {
	userId, orgId := getUserID(ctx), getOrgID(ctx)
	var req request.UpdateKnowledgeQAPairReq
	if !gin_util.Bind(ctx, &req) {
		return
	}
	err := service.UpdateKnowledgeQAPair(ctx, userId, orgId, &req)
	gin_util.Response(ctx, nil, err)
}

// UpdateKnowledgeQAPairSwitch
//
//	@Tags			knowledge.qa
//	@Summary		启停问答对
//	@Description	启停问答对
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Param			data	body		request.UpdateKnowledgeQAPairSwitchReq	true	"启停问答对请求参数"
//	@Success		200		{object}	response.Response
//	@Router			/knowledge/qa/pair/switch [put]
func UpdateKnowledgeQAPairSwitch(ctx *gin.Context) {
	userId, orgId := getUserID(ctx), getOrgID(ctx)
	var req request.UpdateKnowledgeQAPairSwitchReq
	if !gin_util.Bind(ctx, &req) {
		return
	}
	err := service.UpdateKnowledgeQAPairSwitch(ctx, userId, orgId, &req)
	gin_util.Response(ctx, nil, err)
}

// DeleteKnowledgeQAPair
//
//	@Tags			knowledge.qa
//	@Summary		刪除问答对
//	@Description	刪除问答对
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Param			data	body		request.DeleteKnowledgeQAPairReq	true	"刪除问答对请求参数"
//	@Success		200		{object}	response.Response
//	@Router			/knowledge/qa/pair [delete]
func DeleteKnowledgeQAPair(ctx *gin.Context) {
	userId, orgId := getUserID(ctx), getOrgID(ctx)
	var req request.DeleteKnowledgeQAPairReq
	if !gin_util.Bind(ctx, &req) {
		return
	}
	err := service.DeleteKnowledgeQAPair(ctx, userId, orgId, &req)
	gin_util.Response(ctx, nil, err)
}

// GetKnowledgeQAPairList
//
//	@Tags			knowledge.qa
//	@Summary		获取问答对列表
//	@Description	获取问答对列表
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Param			data	query		request.KnowledgeQAPairListReq	true	"问答对列表查询请求参数"
//	@Success		200		{object}	response.Response{data=response.KnowledgeQAPairPageResult}
//	@Router			/knowledge/qa/pair/list [get]
func GetKnowledgeQAPairList(ctx *gin.Context) {
	userId, orgId := getUserID(ctx), getOrgID(ctx)
	var req request.KnowledgeQAPairListReq
	if !gin_util.BindQuery(ctx, &req) {
		return
	}
	resp, err := service.GetKnowledgeQAPairList(ctx, userId, orgId, &req)
	gin_util.Response(ctx, resp, err)
}

// ImportKnowledgeQAPair
//
//	@Tags			knowledge.qa
//	@Summary		问答库文档导入
//	@Description	问答库文档导入
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Param			data	body		request.KnowledgeQAPairImportReq	true	"问答库文档导入请求参数"
//	@Success		200		{object}	response.Response
//	@Router			/knowledge/qa/pair/import [post]
func ImportKnowledgeQAPair(ctx *gin.Context) {
	userId, orgId := getUserID(ctx), getOrgID(ctx)
	var req request.KnowledgeQAPairImportReq
	if !gin_util.Bind(ctx, &req) {
		return
	}
	err := service.ImportKnowledgeQAPair(ctx, userId, orgId, &req)
	gin_util.Response(ctx, nil, err)
}

// KnowledgeQAHit
//
//	@Tags			knowledge
//	@Summary		问答库命中测试
//	@Description	问答库命中测试
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Param			data	body		request.KnowledgeHitReq	true	"问答库命中测试请求参数"
//	@Success		200		{object}	response.Response{data=response.KnowledgeQAHitResp}
//	@Router			/knowledge/qa/hit [post]
func KnowledgeQAHit(ctx *gin.Context) {
	userId, orgId := getUserID(ctx), getOrgID(ctx)
	var req request.KnowledgeQAHitReq
	if !gin_util.Bind(ctx, &req) {
		return
	}
	resp, err := service.KnowledgeQAHit(ctx, userId, orgId, &req)
	gin_util.Response(ctx, resp, err)
}

// ExportKnowledgeQAPair
//
//	@Tags			knowledge.qa
//	@Summary		问答库文档导出
//	@Description	问答库文档导出
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Param			data	query		request.KnowledgeQAPairExportReq	true	"问答库文档导出请求参数"
//	@Success		200		{object}	response.Response
//	@Router			/knowledge/qa/export [get]
func ExportKnowledgeQAPair(ctx *gin.Context) {
	userId, orgId := getUserID(ctx), getOrgID(ctx)
	var req request.KnowledgeQAPairExportReq
	if !gin_util.BindQuery(ctx, &req) {
		return
	}
	err := service.ExportKnowledgeQAPair(ctx, userId, orgId, &req)
	gin_util.Response(ctx, nil, err)
}

// GetKnowledgeQAExportRecordList
//
//	@Tags			knowledge.qa
//	@Summary		获取问答库导出记录列表
//	@Description	获取问答库导出记录列表
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Param			data	query		request.KnowledgeQAExportRecordListReq	true	"获取问答库导出记录列表请求参数"
//	@Success		200		{object}	response.Response{data=response.KnowledgeQAExportRecordPageResult}
//	@Router			/knowledge/qa/export/record/list [get]
func GetKnowledgeQAExportRecordList(ctx *gin.Context) {
	userId, orgId := getUserID(ctx), getOrgID(ctx)
	var req request.KnowledgeQAExportRecordListReq
	if !gin_util.BindQuery(ctx, &req) {
		return
	}
	resp, err := service.GetKnowledgeQAExportRecordList(ctx, userId, orgId, &req)
	gin_util.Response(ctx, resp, err)
}

// DeleteKnowledgeQAExportRecord
//
//	@Tags			knowledge.qa
//	@Summary		删除问答库导出记录
//	@Description	删除问答库导出记录
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Param			data	body		request.DeleteKnowledgeQAExportRecordReq	true	"删除问答库导出记录请求参数"
//	@Success		200		{object}	response.Response
//	@Router			/knowledge/qa/export/record [delete]
func DeleteKnowledgeQAExportRecord(ctx *gin.Context) {
	userId, orgId := getUserID(ctx), getOrgID(ctx)
	var req request.DeleteKnowledgeQAExportRecordReq
	if !gin_util.Bind(ctx, &req) {
		return
	}
	err := service.DeleteKnowledgeExportRecord(ctx, userId, orgId, &req)
	gin_util.Response(ctx, nil, err)
}
