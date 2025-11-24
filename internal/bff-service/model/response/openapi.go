package response

type OpenAPIAgentCreateConversationResponse struct {
	ConversationID string `json:"conversation_id"`
}

type OpenAPIAgentChatResponse struct {
	Code           int                    `json:"code"`
	Message        string                 `json:"message"`
	Response       string                 `json:"response"`
	GenFileUrlList []OpenAPIAgentChatFile `json:"gen_file_url_list"`
	SearchList     []OpenAIChatSearch     `json:"search_list"`
	History        []OpenAIChatHistory    `json:"history"`
	Usage          OpenAPIAgentChatUsage  `json:"usage"`
	Finish         int                    `json:"finish"`
}

type OpenAPIAgentChatFile struct {
	OutputFileUrl string `json:"output_file_url"`
}

type OpenAPIAgentChatUsage struct {
	CompletionTokens int `json:"completion_tokens"`
	PromptTokens     int `json:"prompt_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type OpenAPIRagChatResponse struct {
	Code    int                 `json:"code"`
	Message string              `json:"message"`
	MsgID   string              `json:"msg_id"`
	Data    OpenAPIRagChatData  `json:"data"`
	History []OpenAIChatHistory `json:"history"`
	Finish  int                 `json:"finish"`
}

type OpenAPIRagChatData struct {
	Output     string             `json:"output"`
	SearchList []OpenAIChatSearch `json:"searchList"`
}

type OpenAIChatSearch struct {
	KBName  string `json:"kb_name"`
	Title   string `json:"title"`
	Snippet string `json:"snippet"`
}

type OpenAIChatHistory struct {
	Query    string `json:"query"`
	Response string `json:"response"`
}

type OpenAPIChatflowCreateConversationResponse struct {
	ConversationId string `json:"conversation_id"`
}
