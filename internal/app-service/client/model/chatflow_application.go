package model

// 存储Chatflow关联的智能体应用信息（应用广场）
type ChatflowApplcation struct {
	ID        uint32 `gorm:"primary_key"`
	CreatedAt int64  `gorm:"autoCreateTime:milli;index:idx_chatflow_application_created_at"`
	UpdatedAt int64  `gorm:"autoUpdateTime:milli"`
	// 组织ID
	OrgID string `gorm:"index:idx_chatflow_application_org_id"`
	// 用户ID
	UserID string `gorm:"index:idx_chatflow_application_user_id"`
	// 应用ID
	ApplicationID string `gorm:"index:idx_chatflow_application_application_id"`
	// 对话流ID
	WorkflowID string `gorm:"index:idx_chatflow_application_workflow_id"`
}
