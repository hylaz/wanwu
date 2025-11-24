package service

import (
	"fmt"
	"strconv"
	"strings"

	assistant_service "github.com/UnicomAI/wanwu/api/proto/assistant-service"
	err_code "github.com/UnicomAI/wanwu/api/proto/err-code"
	"github.com/UnicomAI/wanwu/internal/bff-service/model/request"
	"github.com/UnicomAI/wanwu/internal/bff-service/model/response"
	template "github.com/UnicomAI/wanwu/internal/bff-service/pkg/assistant-template"
	gin_util "github.com/UnicomAI/wanwu/pkg/gin-util"
	grpc_util "github.com/UnicomAI/wanwu/pkg/grpc-util"
	"github.com/gin-gonic/gin"
)

func GetAssistantTemplateList(ctx *gin.Context, category, name string) (response.ListResult, error) {
	var ret []*response.AssistantTemplateInfo
	for _, cfg := range template.Cfg() {
		if category != "" && cfg.Category != category {
			continue
		}
		if name != "" && !strings.Contains(cfg.Name, name) {
			continue
		}
		ret = append(ret, toAssistantTemplate(ctx, cfg))
	}
	return response.ListResult{
		List:  ret,
		Total: int64(len(ret)),
	}, nil
}

func GetAssistantTemplate(ctx *gin.Context, assistantTemplateId string) (*response.AssistantTemplateInfo, error) {
	var assistantTemplate *template.Assistant
	for _, cfg := range template.Cfg() {
		if cfg.TemplateId == assistantTemplateId {
			assistantTemplate = cfg
			break
		}
	}
	if assistantTemplate == nil {
		return nil, grpc_util.ErrorStatus(err_code.Code_BFFGeneral, fmt.Sprintf("template %v not found", assistantTemplateId))
	}
	return toAssistantTemplate(ctx, assistantTemplate), nil
}

func AssistantTemplateCreate(ctx *gin.Context, userId, orgId string, req request.AssistantTemplateRequest) (*response.AssistantCreateResp, error) {
	// check
	var assistantTemplate *template.Assistant
	for _, cfg := range template.Cfg() {
		if cfg.TemplateId == req.AssistantTemplateId {
			assistantTemplate = cfg
			break
		}
	}
	if assistantTemplate == nil {
		return nil, grpc_util.ErrorStatus(err_code.Code_BFFGeneral, fmt.Sprintf("template %v not found", req.AssistantTemplateId))
	}
	// check name
	listResp, err := assistant.GetAssistantListMyAll(ctx.Request.Context(), &assistant_service.GetAssistantListMyAllReq{
		Identity: &assistant_service.Identity{
			UserId: userId,
			OrgId:  orgId,
		},
	})
	if err != nil {
		return nil, err
	}
	var i int
	var exist bool
	var name string
	for {
		i++
		exist = false
		name = gin_util.I18nKey(ctx, "bff_assistant_template_name", assistantTemplate.Name, strconv.Itoa(i))
		for _, assistantInfo := range listResp.AssistantInfos {
			if assistantInfo.Name == name {
				exist = true
				break
			}
		}
		if !exist {
			break
		}
	}
	// create
	createResp, err := AssistantCreate(ctx, userId, orgId, request.AppBriefConfig{
		Avatar: request.Avatar{Key: assistantTemplate.AvatarKey},
		Name:   name,
		Desc:   assistantTemplate.Desc,
	})
	if err != nil {
		return nil, err
	}
	// update
	if _, err = AssistantConfigUpdate(ctx, userId, orgId, request.AssistantConfig{
		AssistantId:       createResp.AssistantId,
		Prologue:          assistantTemplate.Prologue,
		Instructions:      assistantTemplate.Instructions,
		RecommendQuestion: assistantTemplate.RecommendQuestion,
	}); err != nil {
		return nil, err
	}
	return createResp, err
}

func toAssistantTemplate(ctx *gin.Context, cfg *template.Assistant) *response.AssistantTemplateInfo {
	return &response.AssistantTemplateInfo{
		AssistantTemplateId: cfg.TemplateId,
		AppType:             "agentTemplate",
		Category:            cfg.Category,
		AppBriefConfig: request.AppBriefConfig{
			Avatar: CacheAvatar(ctx, cfg.AvatarKey, false),
			Name:   cfg.Name,
			Desc:   cfg.Desc,
		},
		Prologue:                  cfg.Prologue,
		Instructions:              cfg.Instructions,
		RecommendQuestion:         cfg.RecommendQuestion,
		Summary:                   cfg.Summary,
		Feature:                   cfg.Feature,
		Scenario:                  cfg.Scenario,
		WorkFlowConfigInstruction: cfg.WorkflowDesc,
	}
}
