package model

type MCPClient struct {
	ID          uint32 `gorm:"primary_key"`
	CreatedAt   int64  `gorm:"autoCreateTime:milli;index:idx_mcp_created_at"`
	UpdatedAt   int64  `gorm:"autoUpdateTime:milli"`
	OrgID       string `gorm:"index:idx_mcp_org_id"`
	UserID      string `gorm:"index:idx_mcp_user_id"`
	McpSquareId string `gorm:"index:idx_mcp_mcp_square_id"`
	Name        string `gorm:"index:idx_mcp_name"`
	From        string `gorm:"index:idx_mcp_from"`
	Desc        string
	SseUrl      string
	AvatarPath  string `gorm:"column:avatar_path;comment:'自定义工具头像'"`
}
