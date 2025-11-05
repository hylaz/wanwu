package service

import (
	"fmt"
	"strings"

	assistant_service "github.com/UnicomAI/wanwu/api/proto/assistant-service"
	errs "github.com/UnicomAI/wanwu/api/proto/err-code"
	"github.com/UnicomAI/wanwu/internal/bff-service/config"
	"github.com/UnicomAI/wanwu/internal/bff-service/model/request"
	"github.com/UnicomAI/wanwu/internal/bff-service/model/response"
	grpc_util "github.com/UnicomAI/wanwu/pkg/grpc-util"
	"github.com/gin-gonic/gin"
)

func CreatePromptByTemplate(ctx *gin.Context, userID, orgID string, req request.CreatePromptByTemplateReq) (*response.PromptIDData, error) {
	promptCfg, exist := config.Cfg().PromptTemp(req.TemplateId)
	if !exist {
		return nil, grpc_util.ErrorStatus(errs.Code_BFFGeneral, "bff_prompt_template_detail", "get prompt template detail empty")
	}
	promptIDResp, err := assistant.CustomPromptCreate(ctx.Request.Context(), &assistant_service.CustomPromptCreateReq{
		AvatarPath: req.Avatar.Key,
		Name:       req.Name,
		Desc:       req.Desc,
		Prompt:     promptCfg.Prompt,
		Identity: &assistant_service.Identity{
			UserId: userID,
			OrgId:  orgID,
		},
	})
	if err != nil {
		return nil, err
	}
	return &response.PromptIDData{
		PromptId: promptIDResp.CustomPromptId,
	}, nil
}

func GetPromptTemplateList(ctx *gin.Context, category, name string) (*response.ListResult, error) {
	var promptTemplateList []*response.PromptTemplateDetail
	for _, promptCfg := range config.Cfg().PromptTemplates {
		if name != "" && !strings.Contains(promptCfg.Name, name) {
			continue
		}
		if !(category == "" || category == "all") && !strings.Contains(promptCfg.Category, category) {
			continue
		}
		promptTemplateList = append(promptTemplateList, buildPromptTempDetail(*promptCfg))
	}
	fmt.Println()
	return &response.ListResult{
		List:  promptTemplateList,
		Total: int64(len(promptTemplateList)),
	}, nil
}

func GetPromptTemplateDetail(ctx *gin.Context, templateId string) (*response.PromptTemplateDetail, error) {
	promptCfg, exist := config.Cfg().PromptTemp(templateId)
	if !exist {
		return nil, grpc_util.ErrorStatus(errs.Code_BFFGeneral, "bff_prompt_template_detail", "get prompt template detail empty")
	}
	return buildPromptTempDetail(promptCfg), nil
}

// --- internal ---
func buildPromptTempDetail(wtfCfg config.PromptTempConfig) *response.PromptTemplateDetail {
	iconUrl := config.Cfg().DefaultIcon.PromptIcon
	return &response.PromptTemplateDetail{
		TemplateId: wtfCfg.TemplateId,
		Category:   wtfCfg.Category,
		Author:     wtfCfg.Author,
		Prompt:     wtfCfg.Prompt,
		AppBriefConfig: request.AppBriefConfig{
			Avatar: request.Avatar{Path: iconUrl},
			Name:   wtfCfg.Name,
			Desc:   wtfCfg.Desc,
		},
	}
}
