package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	operate_service "github.com/UnicomAI/wanwu/api/proto/operate-service"
	"github.com/UnicomAI/wanwu/internal/bff-service/config"
	"github.com/UnicomAI/wanwu/internal/bff-service/model/request"
	"github.com/UnicomAI/wanwu/internal/bff-service/model/response"
	"github.com/UnicomAI/wanwu/pkg/log"
	"github.com/UnicomAI/wanwu/pkg/redis"
	"github.com/UnicomAI/wanwu/pkg/util"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

const downloadKey = "workflowtemplateDownloadCount"

func GetWorkflowTemplateList(ctx *gin.Context, clientId, category, name string) (*response.GetWorkflowTemplateListResp, error) {
	// 记录工作流模板浏览数据
	if err := recordTemplateBrowse(ctx.Request.Context()); err != nil {
		log.Errorf("record template browse count error: %v", err)
	}
	_, err := operate.AddClientRecord(ctx, &operate_service.AddClientRecordReq{
		ClientId: clientId,
	})
	if err != nil {
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

func GetWorkflowTemplateRecommend(ctx *gin.Context, clientId, templateId string) ([]*response.WorkflowTemplateInfo, error) {
	switch config.Cfg().WorkflowTemplatePath.ServerMode {
	case "remote":
		res, err := getRemoteWorkflowTemplateList(ctx, "", "")
		if err != nil {
			return nil, err
		}
		return res.List, nil
	case "local":
		res, err := getLocalWorkflowTemplateList(ctx.Request.Context(), "", "")
		if err != nil {
			return nil, err
		}
		return res.List, nil
	default:
		// 默认使用本地模式
		res, err := getLocalWorkflowTemplateList(ctx.Request.Context(), "", "")
		if err != nil {
			return nil, err
		}
		return res.List, nil
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
		return nil, fmt.Errorf("failed to call remote workflow template API: %v", err)
	}

	if resp.StatusCode() != http.StatusOK {
		fmt.Print(resp.StatusCode())
		return &response.GetWorkflowTemplateListResp{
			Total: 0,
			List:  make([]*response.WorkflowTemplateInfo, 0),
			DownloadLink: response.WorkflowTemplateURL{
				Url: "todo",
			},
		}, nil
	}
	marshal, err := json.Marshal(res.Data)
	if err != nil {
		return nil, fmt.Errorf("request  marshal response body: %v", err)
	}
	if err = json.Unmarshal(marshal, &ret); err != nil {
		return nil, fmt.Errorf("request unmarshal response body: %v", err)
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
		return nil, fmt.Errorf("failed to call remote workflow template API: %v", err)
	}

	if resp.StatusCode() != http.StatusOK {
		fmt.Print(resp.StatusCode())
		return nil, fmt.Errorf("todo")
	}
	marshal, err := json.Marshal(res.Data)
	if err != nil {
		return nil, fmt.Errorf("request  marshal response body: %v", err)
	}
	if err = json.Unmarshal(marshal, &ret); err != nil {
		return nil, fmt.Errorf("request unmarshal response body: %v", err)
	}
	// 远程调用成功，返回远程结果
	return &ret, nil
}

func getLocalWorkflowTemplateDetail(ctx context.Context, templateId string) (*response.WorkflowTemplateDetail, error) {
	wtfCfg, exist := config.Cfg().WorkflowTemp(templateId)
	if !exist {
		return nil, fmt.Errorf("todo")
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
		return nil, fmt.Errorf("failed to call remote workflow template API: %v", err)
	}

	if resp.StatusCode() != http.StatusOK {
		fmt.Print(resp.StatusCode())
		return nil, fmt.Errorf("todo")
	}
	// 远程调用成功，返回远程结果
	return convertToBytes(res.Data)
}

func getLocalDownloadWorkflowTemplate(templateId string) ([]byte, error) {
	wtfCfg, exist := config.Cfg().WorkflowTemp(templateId)
	if !exist {
		return nil, fmt.Errorf("template not found: %s", templateId)
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
		return nil, fmt.Errorf("无法转换类型 %T 为 []byte", data)
	}
}

// 记录模板下载量到单独的Redis Key
func recordTemplateDownloadCount(ctx context.Context, templateID string) error {
	// 使用HINCRBY原子性增加模板下载量
	err := redis.OP().Cli().HIncrBy(ctx, downloadKey, templateID, 1).Err()
	if err != nil {
		return fmt.Errorf("redis HIncrBy key %v field %v err: %v", downloadKey, templateID, err)
	}
	return nil
}

// 根据templateId获取下载量
func getTemplateDownloadCount(ctx context.Context, templateID string) int32 {
	// 使用HGet获取指定模板的下载量
	countStr, err := redis.OP().Cli().HGet(ctx, downloadKey, templateID).Result()
	if err != nil {
		// 键或字段不存在，返回0
		return 0
	}
	return util.MustI32(countStr)
}

// 记录模板下载量到单独的Redis Key
func recordTemplateBrowse(ctx context.Context) error {
	key := getRedisWorkflowTemplateBrowseKey(util.Time2Date(time.Now().UnixMilli()))
	// 使用HINCRBY原子性增加模板下载量
	err := redis.OP().Cli().IncrBy(ctx, key, 1).Err()
	if err != nil {
		return fmt.Errorf("redis IncrBy key %v  err: %v", key, err)
	}
	return nil
}

// 获取Redis键
func getRedisWorkflowTemplateBrowseKey(date string) string {
	return fmt.Sprintf("globalBrowse:%s", date)
}
