package response

type KnowledgeSplitterListResp struct {
	KnowledgeSplitterList []*KnowledgeSplitter `json:"knowledgeSplitterList"`
}

type KnowledgeSplitter struct {
	SplitterId    string `json:"splitterId"`    //知识库分隔符id
	SplitterName  string `json:"splitterName"`  //知识库分隔符名称
	SplitterValue string `json:"splitterValue"` //知识库分隔符值
	Type          string `json:"type"`          //分隔符类型： preset / custom
}
