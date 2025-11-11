package model

type OauthApp struct {
	ID        uint32 `gorm:"primary_key"`
	CreatedAt int64  `gorm:"autoCreateTime:milli"`
	UpdatedAt int64  `gorm:"autoUpdateTime:milli"`
	// 注册该客户端的用户ID
	UserID uint32 `gorm:"index:idx_oauth_app_user_id"`
	// 客户端应用名称
	Name string `gorm:"index:idx_oauth_app_name"`
	// 客户端唯一标识符
	ClientID string `gorm:"index:idx_oauth_app_client_id"`
	// 客户端密钥
	ClientSecret string
	// OAuth回调地址（redirect_uri）
	RedirectURI string
	// 状态（启用/禁用）
	Status bool `gorm:"index:idx_oauth_app_status"`
	// 应用描述
	Description string
}
