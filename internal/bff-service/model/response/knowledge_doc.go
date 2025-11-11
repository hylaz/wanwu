package response

type DocPageResult struct {
	List             []*ListDocResp    `json:"list"`
	Total            int64             `json:"total"`
	PageNo           int               `json:"pageNo"`
	PageSize         int               `json:"pageSize"`
	DocKnowledgeInfo *DocKnowledgeInfo `json:"docKnowledgeInfo"`
}

type DocKnowledgeInfo struct {
	KnowledgeId   string `json:"knowledgeId"`
	KnowledgeName string `json:"knowledgeName"`
	GraphSwitch   int32  `json:"graphSwitch"`
}

type ListDocResp struct {
	DocId         string `json:"docId"`
	DocName       string `json:"docName"`       //文档名称
	DocType       string `json:"docType"`       //文档类型
	KnowledgeId   string `json:"knowledgeId"`   //知识库id
	UploadTime    string `json:"uploadTime"`    //上传时间
	Status        int    `json:"status"`        //处理状态
	ErrorMsg      string `json:"errorMsg"`      //解析错误信息，预留
	FileSize      string `json:"fileSize"`      //文件大小，预留
	SegmentMethod string `json:"segmentMethod"` //分段模式 0:通用分段，1：父子分段
	Author        string `json:"author"`        // 上传文档 作者
	GraphStatus   int32  `json:"graphStatus"`   //图谱状态 0:待处理，1.解析中，2.解析成功，3.解析失败
}

type DocImportTipResp struct {
	Message       string `json:"msg"`
	UploadStatus  int32  `json:"uploadstatus"`  //上传状态
	KnowledgeId   string `json:"knowledgeId"`   //知识库id
	KnowledgeName string `json:"knowledgeName"` //知识库名称
}

type DocSegmentResp struct {
	FileName            string            `json:"fileName"`            //名称
	PageTotal           int               `json:"pageTotal"`           //总页数
	SegmentTotalNum     int               `json:"segmentTotalNum"`     //分段数量
	MaxSegmentSize      int               `json:"maxSegmentSize"`      //设置最大长度
	SegmentType         string            `json:"segmentType"`         //分段方式 0自动分段 1自定义分段
	UploadTime          string            `json:"uploadTime"`          //上传时间
	Splitter            string            `json:"splitter"`            // 分隔符（只有自定义分段必填）
	MetaDataList        []*DocMetaData    `json:"metaDataList"`        //文档元数据
	SegmentContentList  []*SegmentContent `json:"contentList"`         //内容
	SegmentImportStatus string            `json:"segmentImportStatus"` //分段导入状态描述
	SegmentMethod       string            `json:"segmentMethod"`       // 分段方式 父子分段/通用分段
}

type DocMetaData struct {
	MetaKey       string `json:"metaKey"`       // key
	MetaValue     string `json:"metaValue"`     // 确定值
	MetaValueType string `json:"metaValueType"` // number，time, string
	MetaRule      string `json:"metaRule"`      // 正则表达式
	MetaId        string `json:"metaId"`        // 元数据id
}

type SegmentContent struct {
	Content    string   `json:"content"`
	Available  bool     `json:"available"`
	ContentId  string   `json:"contentId"`
	ContentNum int      `json:"contentNum"`
	Labels     []string `json:"labels"`
	IsParent   bool     `json:"isParent"` // 父子分段/通用分段 true是父分段，false是通用分段
	ChildNum   int      `json:"childNum"` // 子分段数量
}

type ChildSegmentInfo struct {
	Content  string `json:"content"`  // 内容
	ChildId  string `json:"childId"`  // 子分段id
	ChildNum int    `json:"childNum"` // 子分段序号
	ParentId string `json:"parentId"` // 父分段id
}

type AnalysisDocUrlResp struct {
	UrlList []*Url `json:"urlList"`
}

type Url struct {
	Url      string `json:"url"`
	FileName string `json:"fileName"`
	FileSize int    `json:"fileSize"`
}

type DocChildSegmentResp struct {
	SegmentContentList []*ChildSegmentInfo `json:"contentList"` //内容
}
