package orm

import (
	"fmt"
	"time"

	"github.com/UnicomAI/wanwu/internal/operate-service/client/model"
	"github.com/UnicomAI/wanwu/pkg/log"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

const (
	ClientRecordSync = "ClientRecordSyncTask"
)

type CronTask struct {
	TaskType  string // 任务类型
	Condition string // 执行条件
}

// 定时任务管理器
type CronManager struct {
	cron      *cron.Cron
	tasks     map[string]cron.EntryID
	isRunning bool
	db        *gorm.DB
}

var (
	cronManager *CronManager
)

// 初始化定时任务
func CronInit(db *gorm.DB) error {
	cronManager = &CronManager{
		cron:  cron.New(),
		tasks: make(map[string]cron.EntryID),
		db:    db,
	}

	// 注册工作流模板同步任务
	task := getClientRecordSyncTask()
	log.Infof("register cron task: %s with condition: %s", task.TaskType, task.Condition)

	entryID, err := cronManager.cron.AddFunc(task.Condition, executeClientRecordSync)
	if err != nil {
		log.Errorf("register cron task %s error: %v", task.TaskType, err)
		return err
	}

	cronManager.tasks[task.TaskType] = entryID
	log.Infof("cron task %s registered successfully with entry ID: %d", task.TaskType, entryID)

	cronManager.cron.Start()
	cronManager.isRunning = true
	log.Infof("workflow template sync cron task started successfully")
	return nil
}

// 停止定时任务
func CronStop() {
	if cronManager != nil && cronManager.isRunning {
		cronManager.cron.Stop()
		cronManager.isRunning = false
		log.Infof("cron tasks stopped")
	}
}

// 获取工作流模板同步任务配置
func getClientRecordSyncTask() CronTask {
	return CronTask{
		TaskType:  ClientRecordSync,
		Condition: "10 0 * * *", // 每天凌晨0点10分执行
	}
}

// 执行工作流模板记录同步任务
func executeClientRecordSync() {
	startTime := time.Now()
	defer func() {
		if r := recover(); r != nil {
			log.Errorf("%s panic: %+v", ClientRecordSync, r)
		}
	}()

	log.Infof("%s start at: %s", ClientRecordSync, startTime.Format("2006-01-02 15:04:05"))

	// 统计前一天的活跃客户端
	yesterday := time.Now().AddDate(0, 0, -1)
	yesterdayStr := yesterday.Format("2006-01-02")

	// 计算活跃客户端数量并存储到新表
	activeClientCount, err := calculateAndStoreActiveClients(yesterdayStr)
	if err != nil {
		log.Errorf("%s calculate active clients error: %v", ClientRecordSync, err)
		return
	}

	spendTime := time.Since(startTime).Milliseconds()
	log.Infof("%s completed, active clients: %d, spend %d ms",
		ClientRecordSync, activeClientCount, spendTime)
}

// 计算活跃客户端数量并存储到WorkflowTemplateClientStats表
func calculateAndStoreActiveClients(date string) (int32, error) {
	// 计算时间范围：前一天的00:00:00到23:59:59
	startTime, err := time.Parse("2006-01-02", date)
	if err != nil {
		return 0, fmt.Errorf("parse date %s error: %v", date, err)
	}
	endTime := startTime.Add(24 * time.Hour).Add(-1 * time.Millisecond)

	// 从WorkflowTemplateStats表统计在指定日期范围内有活动的客户端数量
	var activeClientCount int64
	err = cronManager.db.Model(&model.ClientStats{}).
		Where("created_at BETWEEN ? AND ?", startTime.UnixMilli(), endTime.UnixMilli()).
		Or("updated_at BETWEEN ? AND ?", startTime.UnixMilli(), endTime.UnixMilli()).
		Distinct("client_id").
		Count(&activeClientCount).Error

	if err != nil {
		return 0, fmt.Errorf("count active clients error: %v", err)
	}

	// 检查该日期的记录是否已存在
	var existingStat model.ActiveClientStats
	err = cronManager.db.Where("date = ?", date).First(&existingStat).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// 创建新记录
			newStat := model.ActiveClientStats{
				Date:         date,
				ActiveClient: int32(activeClientCount),
			}
			if err := cronManager.db.Create(&newStat).Error; err != nil {
				return 0, fmt.Errorf("create active client stats error: %v", err)
			}
			log.Infof("created active client stats for date %s: %d clients", date, activeClientCount)
		} else {
			return 0, fmt.Errorf("query active client stats error: %v", err)
		}
	} else {
		// 更新现有记录
		if err := cronManager.db.Model(&existingStat).Updates(map[string]any{
			"active_client": int32(activeClientCount),
			"update_at":     time.Now().UnixMilli(),
		}).Error; err != nil {
			return 0, fmt.Errorf("update active client stats error: %v", err)
		}
		log.Infof("updated active client stats for date %s: %d clients", date, activeClientCount)
	}

	return int32(activeClientCount), nil
}
