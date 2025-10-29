package model

// ClientStats 统计汇总表
type ClientStats struct {
	ID        int64  `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	CreatedAt int64  `gorm:"autoCreateTime:milli"`
	UpdatedAt int64  `gorm:"column:updated_at;not null;default:0" json:"updated_at" `
	ClientId  string `gorm:"column:client_id;not null;default:0" json:"client_id"`
}
