package model

type AppConversation struct {
	ID        uint32 `gorm:"primary_key"`
	CreatedAt int64  `gorm:"autoCreateTime:milli;index:idx_app_conversation_created_at"`
	UpdatedAt int64  `gorm:"autoUpdateTime:milli"`
	// 组织ID
	OrgID string `gorm:"index:idx_app_conversation_org_id"`
	// 用户ID
	UserID string `gorm:"index:idx_app_conversation_user_id"`
	// 应用ID
	AppID string `gorm:"index:idx_app_conversation_app_id"`
	// 应用类型
	AppType string `gorm:"index:idx_app_conversation_app_type"`
	// 会话ID
	ConversationID string `gorm:"index:idx_app_conversation_conversation_id"`
	// 会话名称
	ConversationName string `gorm:"index:idx_app_conversation_conversation_name"`
}
