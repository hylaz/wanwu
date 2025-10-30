package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	errs "github.com/UnicomAI/wanwu/api/proto/err-code"
	operate_service "github.com/UnicomAI/wanwu/api/proto/operate-service"
	"github.com/UnicomAI/wanwu/internal/bff-service/config"
	"github.com/UnicomAI/wanwu/internal/bff-service/model/request"
	"github.com/UnicomAI/wanwu/internal/bff-service/model/response"
	gin_util "github.com/UnicomAI/wanwu/pkg/gin-util"
	grpc_util "github.com/UnicomAI/wanwu/pkg/grpc-util"
	"github.com/UnicomAI/wanwu/pkg/log"
	"github.com/UnicomAI/wanwu/pkg/redis"
	"github.com/UnicomAI/wanwu/pkg/util"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

const (
	redisGlobalBrowseKey             = "globalBrowse"
	redisWorkflowTemplateDownloadKey = "workflowTemplateDownloadCount"
)

func GetWorkflowTemplateList(ctx *gin.Context, clientId, category, name string) (*response.GetWorkflowTemplateListResp, error) {
	// 记录工作流模板浏览数据
	if err := recordTemplateBrowse(ctx.Request.Context()); err != nil {
		log.Errorf("record template browse count error: %v", err)
	}
	// 记录client数据
	if _, err := operate.AddClientRecord(ctx, &operate_service.AddClientRecordReq{
		ClientId: clientId,
	}); err != nil {
		log.Errorf("get workflow template list record err:%v", err)
	}

	switch config.Cfg().WorkflowTemplatePath.ServerMode {
	case "remote":
		return getRemoteWorkflowTemplateList(ctx, category, name)
	case "local":
		return getLocalWorkflowTemplateList(ctx.Request.Context(), category, name)
	default:
		// 默认使用本地模式
		return getLocalWorkflowTemplateList(ctx.Request.Context(), category, name)
	}
}

func GetWorkflowTemplateDetail(ctx *gin.Context, clientId, templateId string) (*response.WorkflowTemplateDetail, error) {
	switch config.Cfg().WorkflowTemplatePath.ServerMode {
	case "remote":
		return getRemoteWorkflowTemplateDetail(ctx, templateId)
	case "local":
		return getLocalWorkflowTemplateDetail(ctx.Request.Context(), templateId)
	default:
		// 默认使用本地模式
		return getLocalWorkflowTemplateDetail(ctx.Request.Context(), templateId)
	}
}

func GetWorkflowTemplateRecommend(ctx *gin.Context, clientId, templateId string) (*response.GetWorkflowTemplateListResp, error) {
	switch config.Cfg().WorkflowTemplatePath.ServerMode {
	case "remote":
		res, err := getRemoteWorkflowTemplateList(ctx, "", "")
		if err != nil {
			return nil, err
		}
		return res, nil
	case "local":
		res, err := getLocalWorkflowTemplateList(ctx.Request.Context(), "", "")
		if err != nil {
			return nil, err
		}
		return res, nil
	default:
		// 默认使用本地模式
		res, err := getLocalWorkflowTemplateList(ctx.Request.Context(), "", "")
		if err != nil {
			return nil, err
		}
		return res, nil
	}
}

func DownloadWorkflowTemplate(ctx *gin.Context, clientId, templateId string) ([]byte, error) {
	// 记录工作流模板下载数据
	if err := recordTemplateDownloadCount(ctx.Request.Context(), templateId); err != nil {
		log.Errorf("record template download count error: %v", err)
	}
	switch config.Cfg().WorkflowTemplatePath.ServerMode {
	case "remote":
		res, err := getRemoteDownloadWorkflowTemplate(ctx, templateId)
		if err != nil {
			return nil, err
		}
		return res, nil
	case "local":
		res, err := getLocalDownloadWorkflowTemplate(templateId)
		if err != nil {
			return nil, err
		}
		return res, nil
	default:
		// 默认使用本地模式
		res, err := getLocalDownloadWorkflowTemplate(templateId)
		if err != nil {
			return nil, err
		}
		return res, nil
	}
}

func GetWorkflowTemplateStatistic(ctx *gin.Context, startDate, endDate string) (*response.WorkflowStatistic, error) {
	// 获取当前周期和上一个周期的日期列表
	prevDates, currentDates, err := util.PreviousDateRange(startDate, endDate)
	if err != nil {
		return nil, grpc_util.ErrorStatusWithKey(errs.Code_BFFGeneral, "bff_workflow_template_stats", fmt.Sprintf("get date range error: %v", err))
	}

	// 获取浏览数据
	currentBrowseData, err := getBrowseDataFromRedis(ctx.Request.Context(), currentDates)
	if err != nil {
		return nil, err
	}
	prevBrowseData, err := getBrowseDataFromRedis(ctx.Request.Context(), prevDates)
	if err != nil {
		return nil, err
	}

	// 计算总览数据
	overview := calculateGlobalBrowseOverview(currentBrowseData, prevBrowseData)

	// 计算趋势数据
	trend := calculateGlobalBrowseTrend(ctx, currentBrowseData, currentDates)

	return &response.WorkflowStatistic{
		Overview: overview,
		Trend:    trend,
	}, nil
}

// 从Redis获取多个日期的浏览数据
func getBrowseDataFromRedis(ctx context.Context, dates []string) (map[string]int64, error) {
	items, err := redis.OP().HGetAll(ctx, redisGlobalBrowseKey)
	if err != nil {
		return nil, grpc_util.ErrorStatusWithKey(errs.Code_BFFGeneral, "bff_workflow_template_stats", fmt.Sprintf("redis HGetAll key %v fields %v err: %v", redisGlobalBrowseKey, dates, err))
	}

	data := make(map[string]int64)
	if len(items) == 0 {
		return data, nil
	}
	for _, date := range dates {
		for _, item := range items {
			if item.K == date {
				data[date] = util.MustI64(item.V)
				break
			}
		}
		// 如果某个日期没有数据，默认值为0
		if _, exist := data[date]; !exist {
			data[date] = 0
		}
	}

	return data, nil
}

// 计算总览数据
func calculateGlobalBrowseOverview(currentData, prevData map[string]int64) response.WorkflowTemplateOverView {
	// 计算当前周期总浏览量
	var currentTotal int64
	for _, count := range currentData {
		currentTotal += count
	}

	// 计算上一个周期总浏览量
	var prevTotal int64
	for _, count := range prevData {
		prevTotal += count
	}

	// 计算环比
	var pop float32
	if prevTotal > 0 {
		pop = (float32(currentTotal) - float32(prevTotal)) / float32(prevTotal) * 100
	} else if currentTotal > 0 {
		// 如果上期为0，本期有数据，增长率为100%
		pop = 100
	}

	return response.WorkflowTemplateOverView{
		Browse: response.WorkflowTemplateOverviewItem{
			Value:            float32(currentTotal),
			PeriodOverPeriod: pop,
		},
	}
}

// 计算趋势数据
func calculateGlobalBrowseTrend(ctx *gin.Context, browseData map[string]int64, dates []string) response.WorkflowTemplateTrends {
	var items []response.StatisticChartLineItem
	for _, date := range dates {
		count := browseData[date]
		items = append(items, response.StatisticChartLineItem{
			Key:   date,
			Value: float32(count),
		})
	}
	return response.WorkflowTemplateTrends{
		Browse: response.StatisticChart{
			TableName: gin_util.I18nKey(ctx, "ope_workflow_template_browse_table"),
			Lines: []response.StatisticChartLine{
				{
					LineName: gin_util.I18nKey(ctx, "ope_workflow_template_browse_line"),
					Items:    items,
				},
			},
		},
	}
}

// --- internal ---

func buildWorkflowTempInfo(ctx context.Context, wtfCfg config.WorkflowTempConfig) *response.WorkflowTemplateInfo {
	iconUrl, _ := url.JoinPath(config.Cfg().Server.ApiBaseUrl, config.Cfg().DefaultIcon.WorkflowIcon)
	return &response.WorkflowTemplateInfo{
		TemplateId: wtfCfg.TemplateId,
		Avatar: request.Avatar{
			Path: iconUrl,
		},
		Name:          wtfCfg.Name,
		Author:        wtfCfg.Author,
		Desc:          wtfCfg.Desc,
		Category:      wtfCfg.Category,
		DownloadCount: getTemplateDownloadCount(ctx, wtfCfg.TemplateId),
	}
}

func buildWorkflowTempDetail(ctx context.Context, wtfCfg config.WorkflowTempConfig) *response.WorkflowTemplateDetail {
	iconUrl, _ := url.JoinPath(config.Cfg().Server.ApiBaseUrl, config.Cfg().DefaultIcon.WorkflowIcon)
	return &response.WorkflowTemplateDetail{
		WorkflowTemplateInfo: response.WorkflowTemplateInfo{
			TemplateId: wtfCfg.TemplateId,
			Avatar: request.Avatar{
				Path: iconUrl,
			},
			Name:          wtfCfg.Name,
			Desc:          wtfCfg.Desc,
			Category:      wtfCfg.Category,
			Author:        wtfCfg.Author,
			DownloadCount: getTemplateDownloadCount(ctx, wtfCfg.TemplateId),
		},
		Summary:  wtfCfg.Summary,
		Feature:  wtfCfg.Feature,
		Scenario: wtfCfg.Scenario,
		Note:     wtfCfg.Note,
	}
}

// --- 获取工作流模板列表 ---

func getRemoteWorkflowTemplateList(ctx *gin.Context, category, name string) (*response.GetWorkflowTemplateListResp, error) {
	client := resty.New()
	client.SetTimeout(30 * time.Second)
	var res response.Response
	var ret response.GetWorkflowTemplateListResp
	resp, err := client.R().
		SetContext(ctx).
		SetQueryParams(map[string]string{
			"category": category,
			"name":     name,
		}).
		SetHeader("Accept", "application/json").
		SetResult(&res).
		Get(config.Cfg().WorkflowTemplatePath.ListUrl)
	if err != nil {
		// 远程调用失败，返回默认下载链接
		return &response.GetWorkflowTemplateListResp{
			Total: 0,
			List:  make([]*response.WorkflowTemplateInfo, 0),
			DownloadLink: response.WorkflowTemplateURL{
				Url: config.Cfg().WorkflowTemplatePath.GlobalWebListUrl,
			},
		}, nil
	}

	if resp.StatusCode() != http.StatusOK {
		// status not ok,  返回默认下载链接
		return &response.GetWorkflowTemplateListResp{
			Total: 0,
			List:  make([]*response.WorkflowTemplateInfo, 0),
			DownloadLink: response.WorkflowTemplateURL{
				Url: config.Cfg().WorkflowTemplatePath.GlobalWebListUrl,
			},
		}, nil
	}
	marshal, err := json.Marshal(res.Data)
	if err != nil {
		return nil, grpc_util.ErrorStatusWithKey(errs.Code_BFFGeneral, "bff_workflow_template_list", fmt.Sprintf("request  marshal response body: %v", err))
	}
	if err = json.Unmarshal(marshal, &ret); err != nil {
		return nil, grpc_util.ErrorStatusWithKey(errs.Code_BFFGeneral, "bff_workflow_template_list", fmt.Sprintf("request  unmarshal response body: %v", err))
	}
	// 远程调用成功，返回远程结果
	return &ret, nil
}

func getLocalWorkflowTemplateList(ctx context.Context, category, name string) (*response.GetWorkflowTemplateListResp, error) {
	var resWorkflowTemp []*response.WorkflowTemplateInfo
	for _, wtfCfg := range config.Cfg().WorkflowTemplateConfig {
		if name != "" && !strings.Contains(wtfCfg.Name, name) {
			continue
		}
		if !(category == "" || category == "all") && !strings.Contains(wtfCfg.Category, category) {
			continue
		}
		resWorkflowTemp = append(resWorkflowTemp, buildWorkflowTempInfo(ctx, *wtfCfg))
	}
	return &response.GetWorkflowTemplateListResp{
		Total:        int64(len(resWorkflowTemp)),
		List:         resWorkflowTemp,
		DownloadLink: response.WorkflowTemplateURL{},
	}, nil
}

// --- 获取工作流模板详情 ---

func getRemoteWorkflowTemplateDetail(ctx *gin.Context, templateId string) (*response.WorkflowTemplateDetail, error) {
	client := resty.New()
	client.SetTimeout(30 * time.Second)
	var res response.Response
	var ret response.WorkflowTemplateDetail
	resp, err := client.R().
		SetContext(ctx).
		SetQueryParams(map[string]string{
			"templateId": templateId,
		}).
		SetHeader("Accept", "application/json").
		SetResult(&res).
		Get(config.Cfg().WorkflowTemplatePath.DetailUrl)
	if err != nil {
		return nil, grpc_util.ErrorStatus(errs.Code_BFFGeneral, "bff_workflow_template_detail", fmt.Sprintf("failed to call remote workflow template API: %v", err))
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, grpc_util.ErrorStatus(errs.Code_BFFGeneral, "bff_workflow_template_detail", fmt.Sprintf("request remote workflow template http code: %v", resp.StatusCode()))
	}
	marshal, err := json.Marshal(res.Data)
	if err != nil {
		return nil, grpc_util.ErrorStatusWithKey(errs.Code_BFFGeneral, "bff_workflow_template_detail", fmt.Sprintf("request marshal response body: %v", err))
	}
	if err = json.Unmarshal(marshal, &ret); err != nil {
		return nil, grpc_util.ErrorStatusWithKey(errs.Code_BFFGeneral, "bff_workflow_template_detail", fmt.Sprintf("request unmarshal response body: %v", err))
	}
	// 远程调用成功，返回远程结果
	return &ret, nil
}

func getLocalWorkflowTemplateDetail(ctx context.Context, templateId string) (*response.WorkflowTemplateDetail, error) {
	wtfCfg, exist := config.Cfg().WorkflowTemp(templateId)
	if !exist {
		return nil, grpc_util.ErrorStatus(errs.Code_BFFGeneral, "bff_workflow_template_detail", "get local workflow template detail empty")
	}
	return buildWorkflowTempDetail(ctx, wtfCfg), nil
}

// --- 下载工作流模板 ---

func getRemoteDownloadWorkflowTemplate(ctx *gin.Context, templateId string) ([]byte, error) {
	client := resty.New()
	client.SetTimeout(30 * time.Second)
	var res response.Response
	resp, err := client.R().
		SetContext(ctx).
		SetQueryParams(map[string]string{
			"templateId": templateId,
		}).
		SetHeader("Accept", "application/json").
		SetResult(&res).
		Get(config.Cfg().WorkflowTemplatePath.DetailUrl)
	if err != nil {
		return nil, grpc_util.ErrorStatus(errs.Code_BFFGeneral, "bff_workflow_template_download", fmt.Sprintf("failed to call remote workflow template API: %v", err))
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, grpc_util.ErrorStatus(errs.Code_BFFGeneral, "bff_workflow_template_download", fmt.Sprintf("request remote workflow template http code: %v", resp.StatusCode()))
	}
	// 远程调用成功，返回远程结果
	return convertToBytes(res.Data)
}

func getLocalDownloadWorkflowTemplate(templateId string) ([]byte, error) {
	wtfCfg, exist := config.Cfg().WorkflowTemp(templateId)
	if !exist {
		return nil, grpc_util.ErrorStatus(errs.Code_BFFGeneral, "bff_workflow_template_download", "get local workflow template download empty")
	}
	return []byte(wtfCfg.Schema), nil
}

func convertToBytes(data any) ([]byte, error) {
	if data == nil {
		return nil, nil
	}
	switch v := data.(type) {
	case []byte:
		return v, nil
	case string:
		return []byte(v), nil
	default:
		return nil, grpc_util.ErrorStatus(errs.Code_BFFGeneral, "bff_workflow_template_download", "unsupported data type for conversion to bytes")
	}
}

// 记录模板下载量到单独的Redis Key
func recordTemplateDownloadCount(ctx context.Context, templateID string) error {
	// 使用HINCRBY原子性增加模板下载量
	err := redis.OP().Cli().HIncrBy(ctx, redisWorkflowTemplateDownloadKey, templateID, 1).Err()
	if err != nil {
		return fmt.Errorf("redis HIncrBy key %v field %v err: %v", redisWorkflowTemplateDownloadKey, templateID, err)
	}
	return nil
}

// 根据templateId获取下载量
func getTemplateDownloadCount(ctx context.Context, templateID string) int32 {
	// 使用HGet获取指定模板的下载量
	countStr, err := redis.OP().Cli().HGet(ctx, redisWorkflowTemplateDownloadKey, templateID).Result()
	if err != nil {
		// 键或字段不存在，返回0
		return 0
	}
	return util.MustI32(countStr)
}

// 记录模板下载量到单独的Redis Key
func recordTemplateBrowse(ctx context.Context) error {
	// 使用HINCRBY原子性增加模板下载量
	date := util.Time2Date(time.Now().UnixMilli())
	err := redis.OP().Cli().HIncrBy(ctx, redisGlobalBrowseKey, date, 1).Err()
	if err != nil {
		return fmt.Errorf("redis IncrBy key %v filed %v err: %v", redisGlobalBrowseKey, date, err)
	}
	return nil
}
