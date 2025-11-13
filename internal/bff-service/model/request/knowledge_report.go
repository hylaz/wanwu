package request

type KnowledgeReportSelectReq struct {
	KnowledgeId string `json:"knowledgeId"  form:"knowledgeId"  validate:"required"`
	PageSearch
	CommonCheck
}

type KnowledgeReportGenerateReq struct {
	KnowledgeId string `json:"knowledgeId"  form:"knowledgeId"  validate:"required"`
	CommonCheck
}

type KnowledgeReportDeleteReq struct {
	KnowledgeId string `json:"knowledgeId"  form:"knowledgeId"  validate:"required"`
	ContentId   string `json:"contentId"  form:"contentId"  validate:"required"`
	CommonCheck
}

type KnowledgeReportUpdateReq struct {
	KnowledgeId string `json:"knowledgeId"  form:"knowledgeId"  validate:"required"`
	ContentId   string `json:"contentId"  form:"contentId"  validate:"required"`
	Content     string `json:"content"  form:"content"  validate:"required"`
	Title       string `json:"title"  form:"title"  validate:"required"`
	CommonCheck
}

type KnowledgeReportAddReq struct {
	KnowledgeId string `json:"knowledgeId"  form:"knowledgeId"  validate:"required"`
	Content     string `json:"content"  form:"content"  validate:"required"`
	Title       string `json:"title"  form:"title"  validate:"required"`
	CommonCheck
}

type KnowledgeReportBatchAddReq struct {
	KnowledgeId  string `json:"knowledgeId"  form:"knowledgeId"  validate:"required"`
	FileUploadId string `json:"fileUploadId"  validate:"required"`
	CommonCheck
}
