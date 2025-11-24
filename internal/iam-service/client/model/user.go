package model

type User struct {
	ID        uint32 `gorm:"primary_key"`
	CreatedAt int64  `gorm:"autoCreateTime:milli"`
	UpdatedAt int64  `gorm:"autoUpdateTime:milli"`
	// 是否系统内置管理员（目前只表示是否系统唯一内置用户）
	IsAdmin bool `gorm:"index:idx_user_is_admin"`
	// 状态
	Status bool `gorm:"index:idx_user_status"`
	// 创建人
	CreatorID uint32 `gorm:"index:idx_user_creator_id"`
	// 用户名
	Name string `gorm:"index:idx_user_name"`
	// 昵称
	Nick string `gorm:"index:idx_user_nick"`
	// 性别
	Gender string `gorm:"index:idx_user_gender"`
	// 电话
	Phone string `gorm:"index:idx_user_phone"`
	// 邮箱
	Email string `gorm:"index:idx_user_email"`
	// 公司
	Company string `gorm:"index:idx_user_company"`
	// 备注
	Remark string
	// 密码
	Password string
	// 最后更新密码时间
	LastUpdatePasswordAt int64
	// 最后一次登录时间（毫秒时间戳）
	LastLoginAt int64
	// 最新token有效时间（毫秒时间戳，此前生成的token都无效）
	LastTokenAt int64
	// 最后一次操作时间（毫秒时间戳）
	LastExecAt int64
	// 用户语言
	Language string `gorm:"index:idx_user_language"`
	// 用户头像
	AvatarPath string `gorm:"index:idx_user_avatar_path"`
	// isEmailCheck
	IsEmailCheck bool `gorm:"default:false;not null"`
}
