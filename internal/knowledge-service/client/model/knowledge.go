package model

type KnowledgeBase struct {
	Id             uint32 `gorm:"column:id;primary_key;type:bigint(20) auto_increment;not null;comment:'id';" json:"id"`       // Primary Key
	KnowledgeId    string `gorm:"uniqueIndex:idx_unique_knowledge_id;column:knowledge_id;type:varchar(64)" json:"knowledgeId"` // Business Primary Key
	Name           string `gorm:"column:name;index:idx_user_id_name,priority:2;type:varchar(256);not null;default:''" json:"name"`
	Description    string `gorm:"column:description;type:text;comment:'知识库描述';" json:"description"`
	DocCount       int    `gorm:"column:doc_count;type:int(11);not null;default:0;comment:'文档数量';" json:"docCount"`
	ShareCount     int    `gorm:"column:share_count;type:int(11);not null;default:0;comment:'文档共享数量';" json:"shareCount"`
	DocSize        int64  `gorm:"column:doc_size;type:bigint(20);not null;default:0;comment:'文档大小单位：字节';" json:"docSize"`
	EmbeddingModel string `gorm:"column:embedding_model;type:longtext;not null;comment:'embedding模型信息';" json:"embeddingModel"`
	CreatedAt      int64  `gorm:"column:create_at;type:bigint(20);not null;" json:"createAt"` // Create Time
	UpdatedAt      int64  `gorm:"column:update_at;type:bigint(20);not null;" json:"updateAt"` // Update Time
	UserId         string `gorm:"column:user_id;index:idx_user_id_name,priority:1;type:varchar(64);not null;default:'';" json:"userId"`
	OrgId          string `gorm:"column:org_id;type:varchar(64);not null;default:'';" json:"orgId"`
	Deleted        int    `gorm:"column:deleted;type:tinyint(1);not null;default:0;comment:'是否逻辑删除';" json:"deleted"`
}

func (KnowledgeBase) TableName() string {
	return "knowledge_base"
}
