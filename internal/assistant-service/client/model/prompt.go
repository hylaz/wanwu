package model

type CustomPrompt struct {
	ID         uint32 `gorm:"primarykey;column:id"`
	AvatarPath string `gorm:"column:avatar_path;comment:自定义提示头像"`
	Name       string `gorm:"column:name;index:idx_custom_prompt_name;comment:自定义提示名称"`
	Desc       string `gorm:"column:desc;comment:自定义提示描述"`
	Prompt     string `gorm:"column:prompt;comment:自定义提示内容"`
	UserId     string `gorm:"column:user_id;index:idx_custom_prompt_user_id;comment:用户id"`
	OrgId      string `gorm:"column:org_id;index:idx_custom_prompt_org_id;comment:组织id"`
	CreatedAt  int64  `gorm:"autoCreateTime:milli;comment:创建时间"`
	UpdatedAt  int64  `gorm:"autoUpdateTime:milli;comment:更新时间"`
}
