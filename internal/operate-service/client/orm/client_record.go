package orm

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	err_code "github.com/UnicomAI/wanwu/api/proto/err-code"
	"github.com/UnicomAI/wanwu/internal/operate-service/client/model"
	"github.com/UnicomAI/wanwu/internal/operate-service/client/orm/sqlopt"
	"github.com/UnicomAI/wanwu/pkg/log"
	"github.com/UnicomAI/wanwu/pkg/util"
	"gorm.io/gorm"
)

func (c *Client) AddClientRecord(ctx context.Context, clientId string) *err_code.Status {
	if err := recordClient(ctx, c.db, clientId); err != nil {
		log.Errorf("record workflow template clientId %v  err: %v", clientId, err)
	}
	return nil
}

func (c *Client) GetClientOverview(ctx context.Context, startDate, endDate string) (*ClientOverView, *err_code.Status) {
	if startDate > endDate {
		return nil, toErrStatus("ope_client_overview_get", fmt.Errorf("startDate %v greater than endDate %v", startDate, endDate).Error())

	}
	userOverView, err := statisticClientOverView(ctx, c.db, startDate, endDate)
	if err != nil {
		return nil, toErrStatus("ope_client_overview_get", err.Error())
	}
	return userOverView, nil
}

func (c *Client) GetClientTrend(ctx context.Context, startDate, endDate string) (*ClientTrends, *err_code.Status) {
	// 参数验证
	if startDate > endDate {
		return nil, toErrStatus("ope_client_trend_get", fmt.Errorf("startDate %v greater than endDate %v", startDate, endDate).Error())
	}

	// 获取活跃客户统计数据
	var stats []*model.ActiveClientStats
	if err := sqlopt.SQLOptions(
		sqlopt.StartDate(startDate),
		sqlopt.EndDate(endDate),
	).Apply(c.db.WithContext(ctx)).Find(&stats).Error; err != nil {
		return nil, toErrStatus("ope_client_trend_get", err.Error())
	}

	// 获取新增客户量统计数据
	additionStats, err := c.getAdditionClientStats(ctx, startDate, endDate)
	if err != nil {
		return nil, toErrStatus("ope_client_trend_get", err.Error())
	}

	// 生成日期范围
	startDateTs, _ := util.Date2Time(startDate)
	endDateTs, _ := util.Date2Time(endDate)
	dates := util.DateRange(startDateTs, endDateTs)

	return &ClientTrends{
		Client: fillTrends(dates, stats, additionStats),
	}, nil
}

func (c *Client) GetCumulativeClientOverview(ctx context.Context, endAt string) (*ClientOverView, *err_code.Status) {
	userOverView, err := statisticCumulativeClientOverView(ctx, c.db, endAt)
	if err != nil {
		return nil, toErrStatus("ope_client_cumulative_overview_get", err.Error())
	}
	return userOverView, nil
}

func (c *Client) getAdditionClientStats(ctx context.Context, startDate, endDate string) ([]*ClientStats, error) {
	startTimestamp, err := util.Date2Time(startDate)
	if err != nil {
		return nil, err
	}
	endTimestamp, err := util.Date2Time(endDate)
	if err != nil {
		return nil, err
	}
	endTimestamp = endTimestamp + 24*time.Hour.Milliseconds()

	// 获取日期范围
	dates := util.DateRange(startTimestamp, endTimestamp)
	var additionClientStats []*ClientStats
	// 为每个日期计算累计客户和新增客户
	for _, date := range dates {
		currentDateTs, _ := util.Date2Time(date)
		// 计算当前日期的新增客户
		var newCount int64
		if err := c.db.WithContext(ctx).
			Model(&model.ClientStats{}).
			Where("created_at BETWEEN ? AND ?", currentDateTs, currentDateTs+24*time.Hour.Milliseconds()).
			Distinct("client_id").
			Count(&newCount).Error; err != nil {
			return nil, fmt.Errorf("addition users stat err: %v", err)
		}
		additionClientStats = append(additionClientStats, &ClientStats{
			Date:  date,
			Value: int32(newCount),
		})
	}

	return additionClientStats, nil
}

func fillTrends(dates []string, activeStats []*model.ActiveClientStats, additionStats []*ClientStats) StatisticChart {
	var activeItems []StatisticChartLineItem
	var additionItems []StatisticChartLineItem
	for _, date := range dates {
		// 活跃客户数据
		activeItem := StatisticChartLineItem{Key: date}
		for _, stat := range activeStats {
			if stat.Date == date {
				activeItem.Value = float32(stat.ActiveClient)
				break
			}
		}
		activeItems = append(activeItems, activeItem)
		// 新增客户数据
		newItem := StatisticChartLineItem{Key: date}
		for _, stat := range additionStats {
			if stat.Date == date {
				newItem.Value = float32(stat.Value)
				break
			}
		}
		additionItems = append(additionItems, newItem)
	}

	return StatisticChart{
		TableName: "operate_workflow_template_client_table",
		Lines: []StatisticChartLine{
			{
				LineName: "operate_workflow_template_active_client_line",
				Items:    activeItems,
			},
			{
				LineName: "operate_workflow_template_addition_client_line",
				Items:    additionItems,
			},
		},
	}
}

func recordClient(ctx context.Context, db *gorm.DB, clientId string) error {
	// 检查数据库中是否已存在该clientId的记录
	existingRecord := &model.ClientStats{}
	result := db.WithContext(ctx).Where("client_id = ?", clientId).First(existingRecord)
	now := time.Now().UnixMilli()
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// 记录不存在，创建新记录
			newRecord := &model.ClientStats{
				ClientId:  clientId,
				UpdatedAt: now,
			}
			if err := db.WithContext(ctx).Create(&newRecord).Error; err != nil {
				return fmt.Errorf("create workflow template stats error: %v", err)
			}
		} else {
			// 其他数据库错误
			return fmt.Errorf("query workflow template stats error: %v", result.Error)
		}
	} else {
		// 记录已存在，更新updated_at字段
		if err := db.WithContext(ctx).Model(&existingRecord).Update("updated_at", now).Error; err != nil {
			return fmt.Errorf("update workflow template stats error: %v", err)
		}
	}
	return nil
}
func statisticCumulativeClientOverView(ctx context.Context, db *gorm.DB, endAt string) (*ClientOverView, error) {
	// 查询累计用户（与时间段无关，所有时间的用户总数）
	var totalUsers int64
	if err := db.WithContext(ctx).
		Model(&model.ClientStats{}).
		Where("created_at < ? ", endAt).
		Distinct("client_id").
		Count(&totalUsers).Error; err != nil {
		return nil, fmt.Errorf("total users stat err: %v", err)
	}
	return &ClientOverView{
		TotalClient: ClientOverviewItem{
			Value: float32(totalUsers),
		},
	}, nil
}
func statisticClientOverView(ctx context.Context, db *gorm.DB, startDate, endDate string) (*ClientOverView, error) {
	// 检查日期
	prevPeriod, currPeriod, err := util.PreviousDateRange(startDate, endDate)
	if err != nil {
		return nil, err
	}

	// 当前周期统计
	currNewUsers, currActiveUsers, err := userStatsByDateRange(ctx, db, currPeriod, startDate, endDate)
	if err != nil {
		return nil, err
	}

	// 上一周期统计
	prevNewUsers, prevActiveUsers, err := userStatsByDateRange(ctx, db, prevPeriod, prevPeriod[0], prevPeriod[len(prevPeriod)-1])
	if err != nil {
		return nil, err
	}

	return &ClientOverView{
		ActiveClient: ClientOverviewItem{
			Value:            currActiveUsers,
			PeriodOverPeriod: calculatePoP(currActiveUsers, prevActiveUsers),
		},
		AdditionClient: ClientOverviewItem{
			Value:            currNewUsers,
			PeriodOverPeriod: calculatePoP(currNewUsers, prevNewUsers),
		},
	}, nil
}

func userStatsByDateRange(ctx context.Context, db *gorm.DB, dates []string, startDate, endDate string) (float32, float32, error) {
	if len(dates) == 0 {
		return 0, 0, nil
	}

	startTimestamp, err := util.Date2Time(startDate)
	if err != nil {
		return 0, 0, err
	}
	endTimestamp, err := util.Date2Time(endDate)
	if err != nil {
		return 0, 0, err
	}
	endTimestamp = endTimestamp + 24*time.Hour.Milliseconds()

	// 查询新增用户（创建时间在指定时间段内）
	var newUsers int64
	if err := db.WithContext(ctx).
		Model(&model.ClientStats{}).
		Where("created_at BETWEEN ? AND ?", startTimestamp, endTimestamp).
		Distinct("client_id").
		Count(&newUsers).Error; err != nil {
		return 0, 0, fmt.Errorf("new users stat err: %v", err)
	}

	// 查询活跃用户（最后操作时间在指定时间段内）
	var activeUsers int64
	if err := db.WithContext(ctx).
		Model(&model.ClientStats{}).
		Where("updated_at BETWEEN ? AND ?", startTimestamp, endTimestamp).
		Distinct("client_id").
		Count(&activeUsers).Error; err != nil {
		return 0, 0, fmt.Errorf("active users stat err: %v", err)
	}
	return float32(newUsers), float32(activeUsers), nil
}

// 计算环比
func calculatePoP(current, previous float32) float32 {
	if previous == 0 {
		if current == 0 {
			return 0
		}
		return 100 // 避免除以零的错误
	}
	value, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", ((current-previous)/previous)*100), 32)
	return float32(value)
}
