package response

import "github.com/UnicomAI/wanwu/internal/bff-service/model/request"

type PromptTemplateDetail struct {
	TemplateId string `json:"templateId" validate:"required"`
	request.AppBriefConfig
	Category string `json:"category"` // 模板分类
	Author   string `json:"author"`   // 作者
	Prompt   string `json:"prompt"`   // 提示词
}

type PromptIDData struct {
	PromptId string `json:"promptId"`
}
