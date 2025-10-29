package model

// ActiveClientStats 统计汇总表
type ActiveClientStats struct {
	ID           int64  `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	CreatedAt    int64  `gorm:"autoCreateTime:milli"`
	UpdateAt     int64  `gorm:"autoUpdateTime:milli"`
	Date         string `gorm:"column:date"`
	ActiveClient int32  `gorm:"column:active_client;not null;default:0" json:"active_client"`
}
