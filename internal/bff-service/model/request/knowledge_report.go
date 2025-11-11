package request

type GetReportReq struct {
	KnowledgeId string `json:"knowledgeId"  form:"knowledgeId"  validate:"required"`
	PageSearch
	CommonCheck
}

type GenerateReportReq struct {
	KnowledgeId string `json:"knowledgeId"  form:"knowledgeId"  validate:"required"`
	CommonCheck
}

type DeleteReportReq struct {
	KnowledgeId string `json:"knowledgeId"  form:"knowledgeId"  validate:"required"`
	ContentId   string `json:"contentId"  form:"contentId"  validate:"required"`
	CommonCheck
}

type UpdateReportReq struct {
	KnowledgeId string `json:"knowledgeId"  form:"knowledgeId"  validate:"required"`
	ContentId   string `json:"contentId"  form:"contentId"  validate:"required"`
	Content     string `json:"content"  form:"content"  validate:"required"`
	Title       string `json:"title"  form:"title"  validate:"required"`
	CommonCheck
}

type AddReportReq struct {
	KnowledgeId string `json:"knowledgeId"  form:"knowledgeId"  validate:"required"`
	Content     string `json:"content"  form:"content"  validate:"required"`
	Title       string `json:"title"  form:"title"  validate:"required"`
	CommonCheck
}

type BatchAddReportReq struct {
	KnowledgeId  string `json:"knowledgeId"  form:"knowledgeId"  validate:"required"`
	FileUploadId string `json:"fileUploadId"  validate:"required"`
	CommonCheck
}
