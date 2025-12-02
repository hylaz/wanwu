package v1

import (
	"encoding/json"
	"net/http"
	"net/url"

	errs "github.com/UnicomAI/wanwu/api/proto/err-code"
	"github.com/UnicomAI/wanwu/internal/bff-service/model/request"
	"github.com/UnicomAI/wanwu/internal/bff-service/model/response"
	"github.com/UnicomAI/wanwu/internal/bff-service/service"
	"github.com/UnicomAI/wanwu/pkg/constant"
	gin_util "github.com/UnicomAI/wanwu/pkg/gin-util"
	grpc_util "github.com/UnicomAI/wanwu/pkg/grpc-util"
	mp "github.com/UnicomAI/wanwu/pkg/model-provider"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
)

// ListLlmModelsByWorkflow
//
//	@Tags		workflow
//	@Summary	llm模型列表（用于workflow）
//	@Description
//	@Security	JWT
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	response.Response{data=response.ListResult{list=response.CozeWorkflowModelInfo}}
//	@Router		/appspace/workflow/model/select/llm [get]
func ListLlmModelsByWorkflow(ctx *gin.Context) {
	resp, err := service.ListLlmModelsByWorkflow(ctx, getUserID(ctx), getOrgID(ctx), mp.ModelTypeLLM)
	gin_util.Response(ctx, resp, err)
}

// CreateWorkflow
//
//	@Tags		workflow
//	@Summary	创建Workflow
//	@Description
//	@Security	JWT
//	@Accept		json
//	@Produce	json
//	@Param		data	body		request.AppBriefConfig	true	"创建Workflow的请求参数"
//	@Success	200		{object}	response.Response{data=response.CozeWorkflowIDData}
//	@Router		/appspace/workflow [post]
func CreateWorkflow(ctx *gin.Context) {
	var req request.AppBriefConfig
	if !gin_util.Bind(ctx, &req) {
		return
	}
	resp, err := service.CreateWorkflow(ctx, getOrgID(ctx), req.Name, req.Desc, req.Avatar.Key)
	gin_util.Response(ctx, resp, err)
}

// CopyWorkflow
//
//	@Tags		workflow
//	@Summary	拷贝Workflow
//	@Description
//	@Security	JWT
//	@Accept		json
//	@Produce	json
//	@Param		data	body		request.WorkflowIDReq	true	"创建Workflow的请求参数"
//	@Success	200		{object}	response.Response{data=response.CozeWorkflowIDData}
//	@Router		/appspace/workflow/copy [post]
func CopyWorkflow(ctx *gin.Context) {
	var req request.WorkflowIDReq
	if !gin_util.Bind(ctx, &req) {
		return
	}
	resp, err := service.CopyWorkflow(ctx, getOrgID(ctx), req.WorkflowID)
	gin_util.Response(ctx, resp, err)
}

// ExportWorkflow
//
//	@Tags			workflow
//	@Summary		导出Workflow
//	@Description	导出工作流的json文件
//	@Security		JWT
//	@Accept			json
//	@Produce		application/octet-stream
//	@Param			workflow_id	query		string	true	"工作流ID"
//	@Success		200			{object}	response.Response{}
//	@Router			/appspace/workflow/export [get]
func ExportWorkflow(ctx *gin.Context) {
	fileName := "workflow_export.json"
	resp, err := service.ExportWorkflow(ctx, getOrgID(ctx), ctx.Query("workflow_id"))
	if err != nil {
		gin_util.Response(ctx, nil, err)
		return
	}
	// 设置响应头
	ctx.Header("Content-Disposition", "attachment; filename*=utf-8''"+url.QueryEscape(fileName))
	ctx.Header("Content-Type", "application/octet-stream")
	ctx.Header("Access-Control-Expose-Headers", "Content-Disposition")
	// 直接写入字节数据
	ctx.Data(http.StatusOK, "application/octet-stream", resp)
}

// ImportWorkflow
//
//	@Tags			workflow
//	@Summary		导入Workflow
//	@Description	通过JSON文件导入工作流
//	@Security		JWT
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			file	formData	file	true	"工作流JSON文件"
//	@Success		200		{object}	response.Response{data=response.CozeWorkflowIDData}
//	@Router			/appspace/workflow/import [post]
func ImportWorkflow(ctx *gin.Context) {
	resp, err := service.ImportWorkflow(ctx, getOrgID(ctx), constant.AppTypeWorkflow)
	gin_util.Response(ctx, resp, err)
}

// GetWorkflowToolSelect
//
//	@Tags		workflow
//	@Summary	工具列表（用于workflow）
//	@Description工具列表（用于workflow）
//	@Security	JWT
//	@Accept		json
//	@Produce	json
//	@Param		toolType	query		string	true	"工具类型"	Enums(builtin,custom)
//	@Param		name		query		string	false	"工具名称"
//	@Success	200			{object}	response.Response{data=response.ListResult{list=[]response.ToolSelect4Workflow}}
//	@Router		/workflow/tool/select [get]
func GetWorkflowToolSelect(ctx *gin.Context) {
	tools, err := service.GetWorkflowToolSelect(ctx, getUserID(ctx), getOrgID(ctx), ctx.Query("toolType"), ctx.Query("name"))
	if err != nil {
		gin_util.Response(ctx, nil, err)
		return
	}
	gin_util.Response(ctx, tools, err)
}

// GetWorkflowToolDetail
//
//	@Tags		workflow
//	@Summary	工具具体action（用于workflow）
//	@Description工具具体action（用于workflow）
//	@Security	JWT
//	@Accept		json
//	@Produce	json
//	@Param		toolType	query		string	true	"工具类型"	Enums(builtin,custom)
//	@Param		actionName	query		string	true	"工具具体action名称"
//	@Param		toolId		query		string	true	"工具ID"
//	@Success	200			{object}	response.Response{data=response.ToolDetail4Workflow}
//	@Router		/workflow/tool/action [get]
func GetWorkflowToolDetail(ctx *gin.Context) {
	data, err := service.GetWorkflowToolDetail(ctx, getUserID(ctx), getOrgID(ctx), ctx.Query("toolId"), ctx.Query("toolType"), ctx.Query("actionName"))
	if err != nil {
		gin_util.Response(ctx, nil, err)
		return
	}
	gin_util.Response(ctx, data, err)
}

// GetWorkflowSelect
//
//	@Tags		workflow
//	@Summary	智能体工作流列表
//	@Description智能体工作流列表
//	@Security	JWT
//	@Accept		json
//	@Produce	json
//	@Param		name	query		string	false	"workflow名称"
//	@Success	200		{object}	response.Response{data=response.ListResult{list=[]response.ExplorationAppInfo}}
//	@Router		/workflow/select [get]
func GetWorkflowSelect(ctx *gin.Context) {
	req := request.GetExplorationAppListRequest{
		Name: ctx.Query("name"),
	}
	resp, err := service.GetWorkflowSelect(ctx, getUserID(ctx), getOrgID(ctx), req)
	gin_util.Response(ctx, resp, err)
}

// CreateWorkflowByTemplate
//
//	@Tags		workflow
//	@Summary	复制工作流模板
//	@Description复制工作流模板
//	@Security	JWT
//	@Accept		json
//	@Produce	json
//	@Param		data	body		request.CreateWorkflowByTemplateReq	true	"通过模板创建Workflow的请求参数"
//	@Success	200		{object}	response.Response{data=response.CozeWorkflowIDData}
//	@Router		/workflow/template [post]
func CreateWorkflowByTemplate(ctx *gin.Context) {
	var req request.CreateWorkflowByTemplateReq
	if !gin_util.Bind(ctx, &req) {
		return
	}
	resp, err := service.CreateWorkflowByTemplate(ctx, getOrgID(ctx), getClientID(ctx), req)
	gin_util.Response(ctx, resp, err)
}

// WorkflowConvert
//
//	@Tags			workflow
//	@Summary		workflow转为chatflow
//	@Description	workflow转为chatflow
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Param			data	body		request.WorkflowConvertReq	true	"对话流工作流转换参数"
//	@Success		200		{object}	response.Response{}
//	@Router			/appspace/workflow/convert [post]
func WorkflowConvert(ctx *gin.Context) {
	var req request.WorkflowConvertReq
	if !gin_util.Bind(ctx, &req) {
		return
	}
	gin_util.Response(ctx, nil, service.WorkflowConvert(ctx, getOrgID(ctx), req.WorkflowID, constant.AppTypeChatflow))
}

// ExplorationWorkflowRun
//
//	@Tags			workflow
//	@Summary		工作流运行接口
//	@Description	工作流运行接口
//	@Security		JWT
//	@Accept			json
//	@Produce		json
//	@Param			data	body		request.WorkflowRunReq	true	"工作流运行参数"
//	@Success		200		{object}	response.ExplorationWorkflowRunResp
//	@Router			/workflow/run [post]
func ExplorationWorkflowRun(ctx *gin.Context) {
	var req request.WorkflowRunReq
	if !gin_util.Bind(ctx, &req) {
		return
	}

	jsonData, err := json.Marshal(req.Input)
	if err != nil {
		gin_util.Response(ctx, nil, grpc_util.ErrorStatus(errs.Code_BFFInvalidArg, "invalid input"))
		return
	}

	resp, err := service.OpenAPIWorkflowRun(ctx, req.WorkflowID, jsonData)
	if err != nil {
		gin_util.Response(ctx, nil, err)
		return
	}

	ewrResp := &response.ExplorationWorkflowRunResp{
		Code: int64(codes.OK),
		Data: string(resp),
		Msg:  "",
	}
	// 判断 resp 是否为 JSON 对象
	if err := json.Unmarshal(resp, &map[string]any{}); err != nil {
		ewrResp.TerminatePlan = "useAnswerContent"
	} else {
		ewrResp.TerminatePlan = "returnVariables"
	}

	b, err := json.Marshal(ewrResp)
	if err != nil {
		gin_util.Response(ctx, nil, grpc_util.ErrorStatus(errs.Code_BFFGeneral, "failed to marshal response"))
		return
	}
	gin_util.ResponseRawByte(ctx, http.StatusOK, b)
}
