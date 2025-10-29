package service

import (
	"context"
	"fmt"
	"strconv"

	operate_service "github.com/UnicomAI/wanwu/api/proto/operate-service"
	"github.com/UnicomAI/wanwu/internal/bff-service/model/response"
	gin_util "github.com/UnicomAI/wanwu/pkg/gin-util"
	"github.com/UnicomAI/wanwu/pkg/log"
	"github.com/UnicomAI/wanwu/pkg/redis"
	"github.com/UnicomAI/wanwu/pkg/util"
	"github.com/gin-gonic/gin"
)

func GetWorkflowTemplateStatistic(ctx *gin.Context, startDate, endDate string) (*response.WorkflowStatistic, error) {
	// 获取当前周期和上一个周期的日期列表
	prevDates, currentDates, err := util.PreviousDateRange(startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("get date range error: %v", err)
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
	overview := calculateOverview(currentBrowseData, prevBrowseData)

	// 计算趋势数据
	trend := calculateTrend(currentBrowseData, currentDates)

	return &response.WorkflowStatistic{
		Overview: overview,
		Trend:    trend,
	}, nil
}

// 从Redis获取多个日期的浏览数据
func getBrowseDataFromRedis(ctx context.Context, dates []string) (map[string]int64, error) {
	data := make(map[string]int64)

	for _, date := range dates {
		redisKey := getRedisWorkflowTemplateBrowseKey(date)

		countStr, err := redis.OP().Cli().Get(ctx, redisKey).Result()
		if err != nil {
			log.Infof("redis get key %s error: %v", redisKey, err)
			continue
		}

		count, err := strconv.ParseInt(countStr, 10, 64)
		if err != nil {
			// 解析失败，当作0处理
			data[date] = 0
		} else {
			data[date] = count
		}
	}

	return data, nil
}

// 计算总览数据
func calculateOverview(currentData, prevData map[string]int64) response.WorkflowTemplateOverView {
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
func calculateTrend(browseData map[string]int64, dates []string) response.WorkflowTemplateTrends {
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
			TableName: "工作流模板浏览趋势",
			Lines: []response.StatisticChartLine{
				{
					LineName: "浏览量",
					Items:    items,
				},
			},
		},
	}
}

func GetCumulativeClientStatistic(ctx *gin.Context, endAt string) (*response.ClientStatistic, error) {
	overview, err := getCumulativeClientStatisticOverview(ctx, endAt)
	if err != nil {
		return nil, err
	}
	return &response.ClientStatistic{
		Overview: *overview,
		Trend:    response.ClientTrends{},
	}, nil
}

func GetClientStatistic(ctx *gin.Context, startDate, endDate string) (*response.ClientStatistic, error) {
	overview, err := getClientStatisticOverview(ctx, startDate, endDate)
	if err != nil {
		return nil, err
	}
	trend, err := getClientStatisticTrend(ctx, startDate, endDate)
	if err != nil {
		return nil, err
	}
	//i18n替换表名
	trend.Client.TableName = gin_util.I18nKey(ctx, trend.Client.TableName)
	for i, line := range trend.Client.Lines {
		trend.Client.Lines[i].LineName = gin_util.I18nKey(ctx, line.LineName)
	}

	return &response.ClientStatistic{
		Overview: *overview,
		Trend:    *trend,
	}, nil
}

func getCumulativeClientStatisticOverview(ctx *gin.Context, date string) (*response.ClientOverView, error) {
	resp, err := operate.GetCumulativeClientOverview(ctx, &operate_service.GetCumulativeClientOverviewReq{
		EndAt: date,
	})
	if err != nil {
		return nil, err
	}
	return &response.ClientOverView{
		TotalClient: response.ClientOverviewItem{
			Value: resp.TotalClient.Value,
		},
	}, nil
}

func getClientStatisticOverview(ctx *gin.Context, startDate, endDate string) (*response.ClientOverView, error) {
	resp, err := operate.GetClientOverview(ctx, &operate_service.GetClientOverviewReq{
		StartDate: startDate,
		EndDate:   endDate,
	})
	if err != nil {
		return nil, err
	}
	return &response.ClientOverView{
		AdditionClient: clientOverviewPb2resp(resp.AdditionClient),
		ActiveClient:   clientOverviewPb2resp(resp.ActiveClient),
	}, nil
}

func getClientStatisticTrend(ctx *gin.Context, startDate, endDate string) (*response.ClientTrends, error) {
	resp, err := operate.GetClientTrend(ctx, &operate_service.GetClientTrendReq{
		StartDate: startDate,
		EndDate:   endDate,
	})
	if err != nil {
		return nil, err
	}
	return &response.ClientTrends{
		Client: convertStatisticChart(resp.Client),
	}, nil
}

func clientOverviewPb2resp(item *operate_service.ClientOverviewItem) response.ClientOverviewItem {
	valueStr := fmt.Sprintf("%.2f", item.Value)
	value, _ := strconv.ParseFloat(valueStr, 64)
	return response.ClientOverviewItem{
		Value:            float32(value),
		PeriodOverPeriod: item.PeriodOverperiod,
	}
}

func convertStatisticChart(pbChart *operate_service.StatisticChart) response.StatisticChart {
	if pbChart == nil {
		return response.StatisticChart{}
	}
	respChart := response.StatisticChart{
		TableName: pbChart.TableName,
		Lines:     make([]response.StatisticChartLine, 0, len(pbChart.ChartLines)),
	}
	for _, pbLine := range pbChart.ChartLines {
		goLine := response.StatisticChartLine{
			LineName: pbLine.LineName,
			Items:    make([]response.StatisticChartLineItem, 0, len(pbLine.Items)),
		}

		for _, pbItem := range pbLine.Items {
			goLine.Items = append(goLine.Items, response.StatisticChartLineItem{
				Key:   pbItem.Key,
				Value: pbItem.Value,
			})
		}
		respChart.Lines = append(respChart.Lines, goLine)
	}
	return respChart
}
