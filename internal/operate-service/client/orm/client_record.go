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
	"github.com/UnicomAI/wanwu/pkg/util"
	"gorm.io/gorm"
)

func (c *Client) AddClientRecord(ctx context.Context, clientId string) *err_code.Status {
	// 检查数据库中是否已存在该clientId的记录
	existingRecord := &model.ClientRecord{}
	nowTs := time.Now().UnixMilli()
	if err := sqlopt.WithClientID(clientId).Apply(c.db).WithContext(ctx).First(existingRecord).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 记录不存在，创建新记录
			if err := sqlopt.WithClientID(clientId).Apply(c.db).WithContext(ctx).Create(&model.ClientRecord{
				ClientId:  clientId,
				UpdatedAt: nowTs,
			}).Error; err != nil {
				return toErrStatus("ope_client_record_create", err.Error())
			}
		} else {
			// 其他数据库错误
			return toErrStatus("ope_client_record_create", err.Error())
		}
	} else {
		// 记录已存在，更新updated_at字段
		if err := sqlopt.WithClientID(clientId).Apply(c.db).WithContext(ctx).Update("updated_at", nowTs).Error; err != nil {
			return toErrStatus("ope_client_record_create", err.Error())
		}
	}
	return nil
}

func (c *Client) GetClientOverview(ctx context.Context, startDate, endDate string) (*ClientOverView, *err_code.Status) {
	if startDate > endDate {
		return nil, toErrStatus("ope_client_overview_get", fmt.Errorf("startDate %v greater than endDate %v", startDate, endDate).Error())

	}
	overview, err := statisticClientOverView(ctx, c.db, startDate, endDate)
	if err != nil {
		return nil, toErrStatus("ope_client_overview_get", err.Error())
	}
	return overview, nil
}

func (c *Client) GetClientTrend(ctx context.Context, startDate, endDate string) (*ClientTrends, *err_code.Status) {
	panic("not implemented")
}

func (c *Client) GetCumulativeClient(ctx context.Context, endAt int64) (int32, *err_code.Status) {
	// 查询累计client（与时间段无关，所有时间的client总数）
	var total int64
	if err := c.db.WithContext(ctx).
		Model(&model.ClientRecord{}).
		Where("created_at < ? ", endAt).
		Count(&total).Error; err != nil {
		return 0, toErrStatus("ope_client_cumulative_get", err.Error())
	}
	return int32(total), nil
}

// --- internal ---

func statisticClientOverView(ctx context.Context, db *gorm.DB, startDate, endDate string) (*ClientOverView, error) {
	newClientOverview, err := statisticNewClientOverView(ctx, db, startDate, endDate)
	if err != nil {
		return nil, err
	}
	activeClientOverview, err := statisticActiveClientOverView(ctx, db, startDate, endDate)
	if err != nil {
		return nil, err
	}
	return &ClientOverView{
		ActiveClient: *activeClientOverview,
		NewClient:    *newClientOverview,
	}, nil
}

// 统计新增client
func statisticNewClientOverView(ctx context.Context, db *gorm.DB, startDate, endDate string) (*ClientOverviewItem, error) {
	prevPeriod, currPeriod, err := util.PreviousDateRange(startDate, endDate)
	if err != nil {
		return nil, err
	}
	currNewCount, err := statisticNewClient(ctx, db, currPeriod[0], currPeriod[len(currPeriod)-1])
	if err != nil {
		return nil, err
	}
	prevNewCount, err := statisticNewClient(ctx, db, prevPeriod[0], prevPeriod[len(prevPeriod)-1])
	if err != nil {
		return nil, err
	}
	return &ClientOverviewItem{
		Value:            float32(currNewCount),
		PeriodOverPeriod: calculatePoP(float32(currNewCount), float32(prevNewCount)),
	}, nil
}

func statisticNewClient(ctx context.Context, db *gorm.DB, startDate, endDate string) (int64, error) {
	startTs, err := util.Date2Time(startDate)
	if err != nil {
		return 0, err
	}
	endTs, err := util.Date2Time(endDate)
	if err != nil {
		return 0, err
	}
	endTs = endTs + 24*time.Hour.Milliseconds()
	// 查询新增client（创建时间在指定时间段内）
	var newCount int64
	if err := db.WithContext(ctx).
		Model(&model.ClientRecord{}).
		Where("created_at BETWEEN ? AND ?", startTs, endTs).
		Count(&newCount).Error; err != nil {
		return 0, fmt.Errorf("new client stat err: %v", err)
	}
	return newCount, nil
}

// 统计活跃client
func statisticActiveClientOverView(ctx context.Context, db *gorm.DB, startDate, endDate string) (*ClientOverviewItem, error) {
	prevPeriod, currPeriod, err := util.PreviousDateRange(startDate, endDate)
	if err != nil {
		return nil, err
	}
	currActiveCount, err := statisticActiveClient(ctx, db, currPeriod[0], currPeriod[len(currPeriod)-1])
	if err != nil {
		return nil, err
	}
	prevActiveCount, err := statisticActiveClient(ctx, db, prevPeriod[0], prevPeriod[len(prevPeriod)-1])
	if err != nil {
		return nil, err
	}
	return &ClientOverviewItem{
		Value:            float32(currActiveCount),
		PeriodOverPeriod: calculatePoP(float32(currActiveCount), float32(prevActiveCount)),
	}, nil
}

func statisticActiveClient(ctx context.Context, db *gorm.DB, startDate, endDate string) (int64, error) {
	// 如果时间范围包含今天，则需要先更新今天的活跃用户统计数据
	startTs, err := util.Date2Time(startDate)
	if err != nil {
		return 0, err
	}
	endTs, err := util.Date2Time(endDate)
	if err != nil {
		return 0, err
	}
	endTs = endTs + 24*time.Hour.Milliseconds()
	nowTs := time.Now().UnixMilli()
	if startTs <= nowTs && nowTs < endTs {

		if err := updateActiveDailyStats(ctx, db, util.Time2Date(nowTs)); err != nil {
			return 0, err
		}
	}
	// 查询活跃client（最后操作时间在指定时间段内）
	var activeClient model.ClientDailyStats
	if err := sqlopt.SQLOptions(
		sqlopt.StartDate(startDate),
		sqlopt.EndDate(endDate),
	).Apply(db).WithContext(ctx).Select("SUM(dau_count) as dau_count").First(&activeClient).Error; err != nil {
		return 0, fmt.Errorf("active client stat err: %v", err)
	}
	return int64(activeClient.DauCount), nil
}

func updateActiveDailyStats(ctx context.Context, db *gorm.DB, date string) error {
	startTs, err := util.Date2Time(date)
	if err != nil {
		return err
	}
	endTs := startTs + 24*time.Hour.Milliseconds()
	// 查询活跃client（更新时间在指定时间段内）
	var activeCount int64
	if err := db.WithContext(ctx).
		Model(&model.ClientRecord{}).
		Where("updated_at BETWEEN ? AND ?", startTs, endTs).
		Count(&activeCount).Error; err != nil {
		return fmt.Errorf("active client stat err: %v", err)
	}
	// 更新或插入每日活跃统计记录
	var existingRecord model.ClientDailyStats
	if err := db.WithContext(ctx).First(&existingRecord).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 记录不存在，创建新记录
			if err := db.WithContext(ctx).Create(&model.ClientDailyStats{
				Date:     date,
				DauCount: int32(activeCount),
			}).Error; err != nil {
				return fmt.Errorf("create client daily stats err: %v", err)
			}
		} else {
			// 其他数据库错误
			return err
		}
	} else {
		// 记录已存在，更新dau_count字段
		if err := db.WithContext(ctx).Model(&existingRecord).Updates(map[string]interface{}{
			"dau_count": int32(activeCount),
		}).Error; err != nil {
			return fmt.Errorf("update client daily stats err: %v", err)
		}
	}
	return nil
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
