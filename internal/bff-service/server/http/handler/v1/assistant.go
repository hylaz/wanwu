package v1

import (
	"github.com/UnicomAI/wanwu/internal/bff-service/model/request"
	"github.com/UnicomAI/wanwu/internal/bff-service/service"
	gin_util "github.com/UnicomAI/wanwu/pkg/gin-util"
	"github.com/gin-gonic/gin"
)

// AssistantCreate
//
//	@Tags			agent
//	@Summary		创建智能体
//	@Description	创建智能体，填写基本信息，创建完成为草稿状态
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Param			data	body		request.AppBriefConfig	true	"智能体基本信息"
//	@Success		200		{object}	response.Response{data=response.AssistantCreateResp}
//	@Router			/assistant [post]
func AssistantCreate(ctx *gin.Context) {
	userId, orgId := getUserID(ctx), getOrgID(ctx)
	var req request.AppBriefConfig
	if !gin_util.Bind(ctx, &req) {
		return
	}
	resp, err := service.AssistantCreate(ctx, userId, orgId, req)
	gin_util.Response(ctx, resp, err)
}

// AssistantUpdate
//
//	@Tags			agent
//	@Summary		修改智能体基本信息
//	@Description	修改智能体基本信息，名称，头像，简介
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Param			data	body		request.AssistantBrief	true	"智能体基本信息参数"
//	@Success		200		{object}	response.Response
//	@Router			/assistant [put]
func AssistantUpdate(ctx *gin.Context) {
	userId, orgId := getUserID(ctx), getOrgID(ctx)
	var req request.AssistantBrief
	if !gin_util.Bind(ctx, &req) {
		return
	}
	resp, err := service.AssistantUpdate(ctx, userId, orgId, req)
	gin_util.Response(ctx, resp, err)
}

// AssistantConfigUpdate
//
//	@Tags			agent
//	@Summary		修改智能体配置信息
//	@Description	修改智能体配置信息，模型配置，知识库配置等等
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Param			data	body		request.AssistantConfig	true	"智能体配置信息参数"
//	@Success		200		{object}	response.Response
//	@Router			/assistant/config [put]
func AssistantConfigUpdate(ctx *gin.Context) {
	userId, orgId := getUserID(ctx), getOrgID(ctx)
	var req request.AssistantConfig
	if !gin_util.Bind(ctx, &req) {
		return
	}
	resp, err := service.AssistantConfigUpdate(ctx, userId, orgId, req)
	gin_util.Response(ctx, resp, err)
}

// GetAssistantInfo
//
//	@Tags			agent
//	@Summary		查看智能体详情
//	@Description	查看智能体详情
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Param			assistantId	query		string	true	"智能体id"
//	@Success		200			{object}	response.Response{data=response.Assistant}
//	@Router			/assistant [get]
func GetAssistantInfo(ctx *gin.Context) {
	userId, orgId := getUserID(ctx), getOrgID(ctx)
	var req request.AssistantIdRequest
	if !gin_util.BindQuery(ctx, &req) {
		return
	}
	resp, err := service.GetAssistantInfo(ctx, userId, orgId, req)
	gin_util.Response(ctx, resp, err)
}

// AssistantCopy
//
//	@Tags			agent
//	@Summary		复制智能体
//	@Description	复制智能体，创建一个新的智能体，基本信息和配置都和原智能体一致
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Param			data		body		request.AssistantIdRequest	true	"智能体id"
//	@Success		200			{object}	response.Response{data=response.AssistantCreateResp}
//	@Router			/assistant/copy [post]
func AssistantCopy(ctx *gin.Context) {
	userId, orgId := getUserID(ctx), getOrgID(ctx)
	var req request.AssistantIdRequest
	if !gin_util.Bind(ctx, &req) {
		return
	}
	resp, err := service.AssistantCopy(ctx, userId, orgId, req)
	gin_util.Response(ctx, resp, err)
}

// AssistantWorkFlowCreate
//
//	@Tags			agent
//	@Summary		添加工作流
//	@Description	为智能体绑定已发布的工作流
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Param			data	body		request.AssistantWorkFlowAddRequest	true	"工作流新增参数"
//	@Success		200		{object}	response.Response
//	@Router			/assistant/tool/workflow [post]
func AssistantWorkFlowCreate(ctx *gin.Context) {
	userId, orgId := getUserID(ctx), getOrgID(ctx)
	var req request.AssistantWorkFlowAddRequest
	if !gin_util.Bind(ctx, &req) {
		return
	}
	err := service.AssistantWorkFlowCreate(ctx, userId, orgId, req)
	gin_util.Response(ctx, nil, err)
}

// AssistantWorkFlowDelete
//
//	@Tags			agent
//	@Summary		删除工作流
//	@Description	为智能体解绑工作流
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Param			data	body		request.AssistantWorkFlowDelRequest	true	"工作流id,智能体id"
//	@Success		200		{object}	response.Response
//	@Router			/assistant/tool/workflow [delete]
func AssistantWorkFlowDelete(ctx *gin.Context) {
	userId, orgId := getUserID(ctx), getOrgID(ctx)
	var req request.AssistantWorkFlowDelRequest
	if !gin_util.Bind(ctx, &req) {
		return
	}
	err := service.AssistantWorkFlowDelete(ctx, userId, orgId, req)
	gin_util.Response(ctx, nil, err)
}

// AssistantWorkFlowEnableSwitch
//
//	@Tags			agent
//	@Summary		启用/停用工作流
//	@Description	修改智能体绑定的工作流的启用状态
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Param			data	body		request.AssistantWorkFlowToolEnableRequest	true	"工作流id,智能体id,开关"
//	@Success		200		{object}	response.Response
//	@Router			/assistant/tool/workflow/switch [put]
func AssistantWorkFlowEnableSwitch(ctx *gin.Context) {
	userId, orgId := getUserID(ctx), getOrgID(ctx)
	var req request.AssistantWorkFlowToolEnableRequest
	if !gin_util.Bind(ctx, &req) {
		return
	}
	err := service.AssistantWorkFlowEnableSwitch(ctx, userId, orgId, req)
	gin_util.Response(ctx, nil, err)
}

// AssistantMCPCreate
//
//	@Tags			agent
//	@Summary		添加mcp工具
//	@Description	为智能体绑定已发布的mcp工具
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Param			data	body		request.AssistantMCPToolAddRequest	true	"mcp新增参数"
//	@Success		200		{object}	response.Response
//	@Router			/assistant/tool/mcp [post]
func AssistantMCPCreate(ctx *gin.Context) {
	userId, orgId := getUserID(ctx), getOrgID(ctx)
	var req request.AssistantMCPToolAddRequest
	if !gin_util.Bind(ctx, &req) {
		return
	}
	err := service.AssistantMCPCreate(ctx, userId, orgId, req)
	gin_util.Response(ctx, nil, err)
}

// AssistantMCPDelete
//
//	@Tags			agent
//	@Summary		删除mcp
//	@Description	为智能体解绑mcp
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Param			data	body		request.AssistantMCPToolDelRequest	true	"mcp工具id，智能体id"
//	@Success		200		{object}	response.Response
//	@Router			/assistant/tool/mcp [delete]
func AssistantMCPDelete(ctx *gin.Context) {
	userId, orgId := getUserID(ctx), getOrgID(ctx)
	var req request.AssistantMCPToolDelRequest
	if !gin_util.Bind(ctx, &req) {
		return
	}
	err := service.AssistantMCPDelete(ctx, userId, orgId, req)
	gin_util.Response(ctx, nil, err)
}

// AssistantMCPEnableSwitch
//
//	@Tags			agent
//	@Summary		启用/停用 MCP
//	@Description	修改智能体绑定的MCP的启用状态
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Param			data	body		request.AssistantMCPToolEnableRequest	true	"mcp工具id、智能体id、enable"
//	@Success		200		{object}	response.Response
//	@Router			/assistant/tool/mcp/switch [put]
func AssistantMCPEnableSwitch(ctx *gin.Context) {
	userId, orgId := getUserID(ctx), getOrgID(ctx)
	var req request.AssistantMCPToolEnableRequest
	if !gin_util.Bind(ctx, &req) {
		return
	}
	err := service.AssistantMCPEnableSwitch(ctx, userId, orgId, req)
	gin_util.Response(ctx, nil, err)
}

// AssistantCustomToolCreate
//
//	@Tags			agent
//	@Summary		添加自定义工具
//	@Description	为智能体绑定自定义工具
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Param			data	body		request.AssistantCustomToolAddRequest	true	"自定义工具新增参数"
//	@Success		200		{object}	response.Response
//	@Router			/assistant/tool/custom [post]
func AssistantCustomToolCreate(ctx *gin.Context) {
	userId, orgId := getUserID(ctx), getOrgID(ctx)
	var req request.AssistantCustomToolAddRequest
	if !gin_util.Bind(ctx, &req) {
		return
	}
	err := service.AssistantCustomToolCreate(ctx, userId, orgId, req)
	gin_util.Response(ctx, nil, err)
}

// AssistantCustomToolDelete
//
//	@Tags			agent
//	@Summary		删除自定义工具
//	@Description	为智能体解绑自定义工具
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Param			data	body		request.AssistantCustomToolDelRequest	true	"智能体id与自定义工具id"
//	@Success		200		{object}	response.Response
//	@Router			/assistant/tool/custom [delete]
func AssistantCustomToolDelete(ctx *gin.Context) {
	userId, orgId := getUserID(ctx), getOrgID(ctx)
	var req request.AssistantCustomToolDelRequest
	if !gin_util.Bind(ctx, &req) {
		return
	}
	err := service.AssistantCustomToolDelete(ctx, userId, orgId, req)
	gin_util.Response(ctx, nil, err)
}

// AssistantCustomToolEnableSwitch
//
//	@Tags			agent
//	@Summary		启用/停用自定义工具
//	@Description	修改智能体绑定的自定义工具的启用状态
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Param			data	body		request.AssistantCustomToolEnableRequest	true	"智能体id与自定义工具id"
//	@Success		200		{object}	response.Response
//	@Router			/assistant/tool/custom/switch [put]
func AssistantCustomToolEnableSwitch(ctx *gin.Context) {
	userId, orgId := getUserID(ctx), getOrgID(ctx)
	var req request.AssistantCustomToolEnableRequest
	if !gin_util.Bind(ctx, &req) {
		return
	}
	err := service.AssistantCustomToolEnableSwitch(ctx, userId, orgId, req)
	gin_util.Response(ctx, nil, err)
}

// ConversationCreate
//
//	@Tags			agent
//	@Summary		创建智能体对话
//	@Description	创建智能体对话
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Param			data	body		request.ConversationCreateRequest	true	"智能体对话创建参数"
//	@Success		200		{object}	response.Response{data=response.ConversationCreateResp}
//	@Router			/assistant/conversation [post]
func ConversationCreate(ctx *gin.Context) {
	userId, orgId := getUserID(ctx), getOrgID(ctx)
	var req request.ConversationCreateRequest
	if !gin_util.Bind(ctx, &req) {
		return
	}
	resp, err := service.ConversationCreate(ctx, userId, orgId, req)
	gin_util.Response(ctx, resp, err)
}

// ConversationDelete
//
//	@Tags			agent
//	@Summary		删除智能体对话
//	@Description	删除智能体对话
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Param			data	body		request.ConversationIdRequest	true	"智能体对话的id"
//	@Success		200		{object}	response.Response
//	@Router			/assistant/conversation [delete]
func ConversationDelete(ctx *gin.Context) {
	userId, orgId := getUserID(ctx), getOrgID(ctx)
	var req request.ConversationIdRequest
	if !gin_util.Bind(ctx, &req) {
		return
	}
	resp, err := service.ConversationDelete(ctx, userId, orgId, req)
	gin_util.Response(ctx, resp, err)
}

// GetConversationList
//
//	@Tags			agent
//	@Summary		智能体对话列表
//	@Description	智能体对话列表
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Param			assistantId	query		string	true	"智能体id"
//	@Param			pageNo		query		int		true	"页面编号，从1开始"
//	@Param			pageSize	query		int		true	"单页数量，从1开始"
//	@Success		200			{object}	response.Response{data=response.PageResult{list=[]response.ConversationInfo}}
//	@Router			/assistant/conversation/list [get]
func GetConversationList(ctx *gin.Context) {
	userId, orgId := getUserID(ctx), getOrgID(ctx)
	var req request.ConversationGetListRequest
	if !gin_util.BindQuery(ctx, &req) {
		return
	}
	resp, err := service.GetConversationList(ctx, userId, orgId, req)
	gin_util.Response(ctx, resp, err)
}

// GetConversationDetailList
//
//	@Tags			agent
//	@Summary		智能体对话详情历史列表
//	@Description	智能体对话详情历史列表
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Param			conversationId	query		string	true	"智能体对话id"
//	@Param			pageNo			query		int		true	"页面编号，从1开始"
//	@Param			pageSize		query		int		true	"单页数量，从1开始"
//	@Success		200				{object}	response.Response{data=response.PageResult{list=[]response.ConversationDetailInfo}}
//	@Router			/assistant/conversation/detail [get]
func GetConversationDetailList(ctx *gin.Context) {
	userId, orgId := getUserID(ctx), getOrgID(ctx)
	var req request.ConversationGetDetailListRequest
	if !gin_util.BindQuery(ctx, &req) {
		return
	}
	resp, err := service.GetConversationDetailList(ctx, userId, orgId, req)
	gin_util.Response(ctx, resp, err)
}

// AssistantConversionStream
//
//	@Tags			agent
//	@Summary		智能体流式问答
//	@Description	智能体流式问答
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Param			data	body		request.ConversionStreamRequest	true	"智能体流式问答参数"
//	@Success		200		{object}	response.Response
//	@Router			/assistant/stream [post]
func AssistantConversionStream(ctx *gin.Context) {
	userId, orgId := getUserID(ctx), getOrgID(ctx)
	var req request.ConversionStreamRequest
	if !gin_util.Bind(ctx, &req) {
		return
	}
	if err := service.AssistantConversionStream(ctx, userId, orgId, req); err != nil {
		gin_util.Response(ctx, nil, err)
	}
}
