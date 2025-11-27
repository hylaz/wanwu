package response

type KnowledgeQAPairPageResult struct {
	List            []*ListKnowledgeQAPairResp `json:"list"`
	QAKnowledgeInfo *QAKnowledgeInfo           `json:"qaKnowledgeInfo"`
	Total           int64                      `json:"total"`
	PageNo          int                        `json:"pageNo"`
	PageSize        int                        `json:"pageSize"`
}

type QAKnowledgeInfo struct {
	KnowledgeId   string `json:"knowledgeId"`
	KnowledgeName string `json:"knowledgeName"`
}

type ListKnowledgeQAPairResp struct {
	QAPairId     string         `json:"qaPairId"`     //问答对id
	KnowledgeId  string         `json:"knowledgeId"`  //问答库id
	Question     string         `json:"question"`     //问题
	Answer       string         `json:"answer"`       //答案
	MetaDataList []*DocMetaData `json:"metaDataList"` //元数据
	Author       string         `json:"author"`       //作者
	UploadTime   string         `json:"uploadTime"`   //上传时间
	Status       int            `json:"status"`       //处理状态
	Switch       bool           `json:"switch"`       //启停开关
	ErrorMsg     string         `json:"errorMsg"`     //处理错误信息
}

type CreateKnowledgeQAPairResp struct {
	QAPairId string `json:"qaPairId"`
}

type KnowledgeQAExportRecordPageResult struct {
	List     []*ListKnowledgeQAExportRecordResp `json:"list"`
	Total    int64                              `json:"total"`
	PageNo   int                                `json:"pageNo"`
	PageSize int                                `json:"pageSize"`
}

type ListKnowledgeQAExportRecordResp struct {
	QAExportRecordId string `json:"qaExportRecordId"` //问答库导出记录id
	Author           string `json:"author"`           //导出人
	ExportTime       string `json:"exportTime"`       //导出时间
	FilePath         string `json:"filePath"`         //导出文件路径
	Status           int    `json:"status"`           //状态
	ErrorMsg         string `json:"errorMsg"`         //导出状态错误信息
}

type KnowledgeQAHitResp struct {
	SearchList []*QAHitSearchList `json:"searchList"` //种种结果
	Score      []float64          `json:"score"`      //打分信息
}

type QAHitSearchList struct {
	Title       string `json:"title"`
	Question    string `json:"question"`
	Answer      string `json:"answer"`
	QAPairId    string `json:"qaPairId"`
	QABase      string `json:"qaBase"`
	QAId        string `json:"qaId"`
	ContentType string `json:"contentType"` // graph：知识图谱（文本）, text：文档分段（文本）, community_report：社区报告（markdown），qa：问答库（文本）
}
