package model

type RagKnowledgeHitResp struct {
	Code    int               `json:"code"`
	Message string            `json:"message"`
	Data    *KnowledgeHitData `json:"data"`
}

type KnowledgeHitData struct {
	Prompt     string             `json:"prompt"`
	SearchList []*ChunkSearchList `json:"searchList"`
	Score      []float64          `json:"score"`
}

type ChunkSearchList struct {
	Title            string          `json:"title"`
	Snippet          string          `json:"snippet"`
	KbName           string          `json:"kb_name"`
	MetaData         interface{}     `json:"meta_data"`
	ChildContentList []*ChildContent `json:"child_content_list"`
	ChildScore       []float64       `json:"child_score"`
}

type ChildContent struct {
	ChildSnippet string  `json:"child_snippet"`
	Score        float64 `json:"score"`
}
