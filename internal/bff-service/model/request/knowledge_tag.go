package request

type KnowledgeTagSelectReq struct {
	KnowledgeId string `json:"knowledgeId" form:"knowledgeId" `
	TagName     string `json:"tagName" form:"tagName" `
	CommonCheck
}

type KnowledgeTagBindCountReq struct {
	TagId string `json:"tagId" form:"tagId" validate:"required"`
	CommonCheck
}

type CreateKnowledgeTagReq struct {
	TagName string `json:"tagName" form:"tagName" validate:"required"`
	CommonCheck
}

type UpdateKnowledgeTagReq struct {
	TagId   string `json:"tagId"  form:"tagId" validate:"required"`
	TagName string `json:"tagName" form:"tagName" validate:"required"`
	CommonCheck
}

type DeleteKnowledgeTagReq struct {
	TagId string `json:"tagId"  form:"tagId" validate:"required"`
	CommonCheck
}

type BindKnowledgeTagReq struct {
	TagIdList   []string `json:"tagIdList"  form:"tagIdList" validate:"required"`
	KnowledgeId string   `json:"knowledgeId"  form:"knowledgeId" validate:"required"`
	Option      int      `json:"option"  form:"option" ` //0:bind,1:unbind
	CommonCheck
}
