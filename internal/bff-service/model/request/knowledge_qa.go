package request

type KnowledgeQAPairListReq struct {
	KnowledgeId string `json:"knowledgeId" form:"knowledgeId" validate:"required"`
	Name        string `json:"name" form:"name"`
	Status      int    `json:"status" form:"status"`
	MetaValue   string `json:"metaValue" form:"metaValue"`
	PageSearch
	CommonCheck
}

type CreateKnowledgeQAPairReq struct {
	KnowledgeId string `json:"knowledgeId" validate:"required"`
	Question    string `json:"question"` //问题
	Answer      string `json:"answer"`   //答案
	CommonCheck
}

type UpdateKnowledgeQAPairReq struct {
	QAPairId string `json:"qaPairId" validate:"required"`
	Question string `json:"question"` //问题
	Answer   string `json:"answer"`   //答案
	CommonCheck
}

type UpdateKnowledgeQAPairSwitchReq struct {
	QAPairId string `json:"qaPairId" validate:"required"`
	Switch   bool   `json:"switch"`
	CommonCheck
}

type DeleteKnowledgeQAPairReq struct {
	KnowledgeId  string   `json:"knowledgeId" validate:"required"`  //问答库id
	QAPairIdList []string `json:"QAPairIdList" validate:"required"` //问答对列表
	CommonCheck
}

type KnowledgeQAPairImportReq struct {
	KnowledgeId string     `json:"knowledgeId" validate:"required"` //问答库id
	DocInfo     []*DocInfo `json:"docInfoList" validate:"required"` //上传文档列表
	CommonCheck
}

type KnowledgeQAPairImportTipReq struct {
	KnowledgeId string `json:"knowledgeId" form:"knowledgeId" validate:"required"` //问答库id
	CommonCheck
}

type KnowledgeQAPairExportReq struct {
	KnowledgeId string `json:"knowledgeId" form:"knowledgeId" validate:"required"` //问答库id
	CommonCheck
}

type KnowledgeExportRecordListReq struct {
	KnowledgeId string `json:"knowledgeId" form:"knowledgeId" validate:"required"` //知识库id
	PageSearch
	CommonCheck
}

type DeleteKnowledgeExportRecordReq struct {
	KnowledgeId    string `json:"knowledgeId" validate:"required"`    //问答库id
	ExportRecordId string `json:"exportRecordId" validate:"required"` //问答库导出记录id
	CommonCheck
}

type KnowledgeQAHitReq struct {
	KnowledgeList        []*AppKnowledgeBase   `json:"knowledgeList"`
	Question             string                `json:"question"   validate:"required"`
	KnowledgeMatchParams *KnowledgeMatchParams `json:"knowledgeMatchParams"   validate:"required"`
	CommonCheck
}
