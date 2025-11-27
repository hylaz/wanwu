package model

const (
	KnowledgeQAPairExportInit      = 0 //任务待处理
	KnowledgeQAPairExportExporting = 1 //导出中
	KnowledgeQAPairExportSuccess   = 2 //导出成功
	KnowledgeQAPairExportFail      = 3 //导出失败
)

type KnowledgeQAPairExportTaskParams struct {
	KnowledgeId string `json:"knowledge_id"` //问答库id
}

type KnowledgeQAPairExportTask struct {
	Id             uint32 `gorm:"column:id;primary_key;type:bigint(20) auto_increment;not null;comment:'id';" json:"id"`
	ExportId       string `gorm:"uniqueIndex:idx_unique_export_id;column:export_id;type:varchar(64)" json:"exportId"` // Business Primary Key
	KnowledgeId    string `gorm:"column:knowledge_id;type:varchar(64);not null;index:idx_knowledge_id" json:"knowledgeId"`
	ExportFilePath string `gorm:"column:export_file_path;type:text;not null;comment:'导出文件地址'" json:"exportFilePath"`
	ExportFileSize int64  `gorm:"column:export_file_size;type:bigint(20);not null;comment:'导出文件大小'" json:"exportFileSize"`
	Status         int    `gorm:"column:status;type:tinyint(1);not null;comment:'0-任务待处理；1-任务导出中 ；2-任务完成；3-任务失败'" json:"status"`
	SuccessCount   int    `gorm:"column:success_count;type:bigint(20);default:0;comment:'成功数量'" json:"successCount"`
	TotalCount     int    `gorm:"column:total_count;type:bigint(20);default:0;comment:'导出数量，当在导出过程中出现重启，则total为0'" json:"totalCount"`
	ErrorMsg       string `gorm:"column:error_msg;type:longtext;not null;comment:'导出的错误信息'" json:"errorMsg"`
	CreatedAt      int64  `gorm:"column:create_at;type:bigint(20);not null;autoCreateTime:milli" json:"createAt"` // Create Time
	UpdatedAt      int64  `gorm:"column:update_at;type:bigint(20);not null;autoUpdateTime:milli" json:"updateAt"` // Update Time
	UserId         string `gorm:"column:user_id;type:varchar(64);not null;default:'';" json:"userId"`
	OrgId          string `gorm:"column:org_id;type:varchar(64);not null;default:''" json:"orgId"`
}

func (KnowledgeQAPairExportTask) TableName() string {
	return "knowledge_qa_pair_export_task"
}
