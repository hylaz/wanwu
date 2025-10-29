package util

import (
	"fmt"
	"time"
)

var UTC8 *time.Location

const timeMsFormat = "2006-01-02 15:04:05.000"

const timeFormat = "2006-01-02 15:04:05"
const dateFormat = "2006-01-02"

func InitTimeLocal() error {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return err
	}
	UTC8 = loc
	return nil
}

// WeekStartTime 某周的开始时间(周一0点)，-1上周，0本周，1下周
func WeekStartTime(t time.Time, week int) time.Time {
	offset := time.Monday - t.Weekday()
	if offset > 0 { // 周日特殊处理
		offset = -6
	}
	y, m, d := t.Date()
	today := time.Date(y, m, d, 0, 0, 0, 0, t.Location())
	return today.AddDate(0, 0, int(offset)+week*7)
}

func Time2Str(millSec int64) string {
	return time.UnixMilli(millSec).In(UTC8).Format(timeFormat)
}

func Time2MsStr(millSec int64) string {
	return time.UnixMilli(millSec).In(UTC8).Format(timeMsFormat)
}

func Str2Time(timeStr string) (int64, error) {
	t, err := time.ParseInLocation(timeFormat, timeStr, UTC8)
	if err != nil {
		return 0, err
	}
	return t.UnixMilli(), nil
}

func Str2Date(timeStr string) (int64, error) {
	t, err := time.ParseInLocation(dateFormat, timeStr, UTC8)
	if err != nil {
		return 0, err
	}
	return t.UnixMilli(), nil
}

func Time2Date(ts int64) string {
	return time.UnixMilli(ts).In(UTC8).Format(dateFormat) // 输出示例: 2025-05-09
}

func Date2Time(date string) (int64, error) {
	t, err := time.ParseInLocation(dateFormat, date, UTC8)
	if err != nil {
		return 0, err
	}
	return t.UnixMilli(), nil
}

// DateRange 返回[startTs, endTs]闭区间日期列表
func DateRange(startTs, endTs int64) []string {
	if startTs > endTs {
		return nil
	}
	endDate := Time2Date(endTs)
	var ret []string
	for {
		date := Time2Date(startTs)
		ret = append(ret, date)
		if date == endDate {
			break
		}
		startTs = startTs + time.Hour.Milliseconds()*24
	}
	return ret
}

// 返回上一个周期和当前周期闭区间日期列表
func PreviousDateRange(startDate, endDate string) ([]string, []string, error) {
	// 1. 解析输入日期
	startAt, err := Date2Time(startDate)
	if err != nil {
		return nil, nil, err
	}
	endAt, err := Date2Time(endDate)
	if err != nil {
		return nil, nil, err
	}
	if startAt > endAt {
		return nil, nil, fmt.Errorf("startDate %v greater than endDate %v", startDate, endDate)
	}
	// 2. 计算上一周期时间戳区间（前后区间日期无重叠）
	deltaDura := endAt - startAt + 24*time.Hour.Milliseconds()
	pervStartTs := startAt - deltaDura
	pervEndTs := endAt - deltaDura
	// 3. 计算上一个周期，当前周期
	return DateRange(pervStartTs, pervEndTs), DateRange(startAt, endAt), nil
}
