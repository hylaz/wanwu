package response

import "github.com/UnicomAI/wanwu/internal/bff-service/model/request"

type GetWorkflowTemplateListResp struct {
	Total        int64                   `json:"total"`
	List         []*WorkflowTemplateInfo `json:"list"`
	DownloadLink WorkflowTemplateURL     `json:"downloadLink"`
}

// WorkflowTemplateDetail 工作流模板详情响应
type WorkflowTemplateDetail struct {
	WorkflowTemplateInfo
	Summary  string `json:"summary"`  // 模板介绍概览
	Feature  string `json:"feature"`  // 模板特性说明
	Scenario string `json:"scenario"` // 模板应用场景
	Note     string `json:"note"`     // 注意事项
}

// WorkflowTemplateListItem 工作流模板列表项
type WorkflowTemplateInfo struct {
	TemplateId    string         `json:"templateId"`    // 模板ID
	Avatar        request.Avatar `json:"avatar"`        // 图标
	Name          string         `json:"name"`          // 模板名称
	Desc          string         `json:"desc"`          // 模板描述
	Category      string         `json:"category"`      // 模板分类
	Author        string         `json:"author"`        // 作者
	DownloadCount int32          `json:"downloadCount"` // 下载次数
}

type WorkflowTemplateURL struct {
	Url string `json:"url"`
}
