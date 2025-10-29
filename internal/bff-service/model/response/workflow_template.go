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

type ClientTrends struct {
	Client StatisticChart `json:"client"` // 客户端活跃数据
}

type WorkflowStatistic struct {
	Overview WorkflowTemplateOverView `json:"overview"` // 工作流模板统计面板
	Trend    WorkflowTemplateTrends   `json:"trend"`    // 工作流模板统计趋势
}

type WorkflowTemplateOverView struct {
	Browse WorkflowTemplateOverviewItem `json:"browse"` // 工作流模板浏览数据总览
}

type WorkflowTemplateOverviewItem struct {
	Value            float32 `json:"value"`            // 数量
	PeriodOverPeriod float32 `json:"periodOverPeriod"` // 环比上周期百分比
}

type WorkflowTemplateTrends struct {
	Browse StatisticChart `json:"browse"` // 工作流模板浏览数据趋势
}

type StatisticChart struct {
	TableName string               `json:"tableName"` // 统计表名字
	Lines     []StatisticChartLine `json:"lines"`     // 统计表中线段集合
}

type StatisticChartLine struct {
	LineName string                   `json:"lineName"` // 线段名字
	Items    []StatisticChartLineItem `json:"items"`    // 线段横纵坐标值
}

type StatisticChartLineItem struct {
	Key   string  `json:"key"`
	Value float32 `json:"value"`
}
