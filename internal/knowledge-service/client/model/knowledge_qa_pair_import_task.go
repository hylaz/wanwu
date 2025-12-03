package model

const (
	KnowledgeQAPairImportInit      = 0 //任务待处理
	KnowledgeQAPairImportImporting = 1 //导入中
	KnowledgeQAPairImportSuccess   = 2 //导入成功
	KnowledgeQAPairImportFail      = 3 //导入失败
)

type KnowledgeQAPairImportTask struct {
	Id           uint32 `gorm:"column:id;primary_key;type:bigint(20) auto_increment;not null;comment:'id';" json:"id"`
	ImportId     string `gorm:"uniqueIndex:idx_unique_import_id;column:import_id;type:varchar(64)" json:"importId"` // Business Primary Key
	KnowledgeId  string `gorm:"column:knowledge_id;type:varchar(64);not null;index:idx_knowledge_id" json:"knowledgeId"`
	DocInfo      string `gorm:"column:doc_info;type:longtext;not null;comment:'文件信息'" json:"docInfo"`
	Status       int    `gorm:"column:status;type:tinyint(1);not null;comment:'0-任务待处理；1-任务导入中 ；2-任务完成；3-任务失败'" json:"status"`
	SuccessCount int    `gorm:"column:success_count;type:bigint(20);default:0;comment:'成功数量'" json:"successCount"`
	TotalCount   int    `gorm:"column:total_count;type:bigint(20);default:0;comment:'导入数量，当在导入过程中出现重启，则total为0'" json:"totalCount"`
	ErrorMsg     string `gorm:"column:error_msg;type:longtext;not null;comment:'导入的错误信息'" json:"errorMsg"`
	CreatedAt    int64  `gorm:"column:create_at;type:bigint(20);not null;autoCreateTime:milli" json:"createAt"` // Create Time
	UpdatedAt    int64  `gorm:"column:update_at;type:bigint(20);not null;autoUpdateTime:milli" json:"updateAt"` // Update Time
	UserId       string `gorm:"column:user_id;type:varchar(64);not null;default:'';" json:"userId"`
	OrgId        string `gorm:"column:org_id;type:varchar(64);not null;default:''" json:"orgId"`
}

func (KnowledgeQAPairImportTask) TableName() string {
	return "knowledge_qa_pair_import_task"
}
