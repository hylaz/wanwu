package model

type FileInfo struct {
	FileName string `json:"fileName"`
	FileSize int64  `json:"fileSize"`
	FileUrl  string `json:"fileUrl"`
}

type ConversationDetails struct {
	Id             string     `json:"id"`
	AssistantId    string     `json:"assistantId"`
	ConversationId string     `json:"conversationId"`
	Prompt         string     `json:"prompt"`
	SysPrompt      string     `json:"sysPrompt"`
	Response       string     `json:"response"`
	SearchList     string     `json:"searchList"`
	QaType         int32      `json:"qaType"`
	FileUrl        string     `json:"requestFileUrls"`
	FileSize       int64      `json:"fileSize"`
	FileName       string     `json:"fileName"`
	FileInfo       []FileInfo `json:"fileInfo"`
	UserId         string     `json:"userId"`
	OrgId          string     `json:"orgId"`
	CreatedAt      int64      `json:"createdAt"`
	UpdatedAt      int64      `json:"updatedAt"`
}
