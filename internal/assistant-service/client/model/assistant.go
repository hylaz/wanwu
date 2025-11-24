package model

type Assistant struct {
	ID                  uint32 `gorm:"primarykey;column:id;comment:智能体Id"`
	AvatarPath          string `gorm:"column:avatar_path;comment:智能体头像"`
	Name                string `gorm:"column:name;comment:智能体名称"`
	Desc                string `gorm:"column:desc;comment:智能体介绍"`
	Instructions        string `gorm:"column:instructions;comment:系统提示词"`
	Prologue            string `gorm:"column:prologue;comment:开场白"`
	RecommendQuestion   string `gorm:"column:recommend_question;comment:推荐问题列表"`
	ModelConfig         string `gorm:"column:model_config;type:longtext;comment:模型配置"`
	RerankConfig        string `gorm:"column:rerank_config;type:longtext;comment:rerank模型配置"`
	KnowledgebaseConfig string `gorm:"column:knowledgebase_config;type:longtext;comment:知识库配置"`
	SafetyConfig        string `gorm:"column:safety_config;type:longtext;comment:安全配置"`
	VisionConfig        string `gorm:"column:vision_config;type:longtext;comment:视觉配置"`
	Scope               int    `gorm:"column:scope;type:tinyint;comment:智能体可见范围"`
	UserId              string `gorm:"column:user_id;index:idx_assistant_user_id;comment:用户id"`
	OrgId               string `gorm:"column:org_id;index:idx_assistant_org_id;comment:组织id"`
	CreatedAt           int64  `gorm:"autoCreateTime:milli;comment:创建时间"`
	UpdatedAt           int64  `gorm:"autoUpdateTime:milli;comment:更新时间"`
}
