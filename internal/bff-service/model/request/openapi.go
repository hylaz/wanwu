package request

type OpenAPIAgentCreateConversationRequest struct {
	Title string `json:"title"`
}

func (req *OpenAPIAgentCreateConversationRequest) Check() error {
	return nil
}

type OpenAPIAgentChatRequest struct {
	ConversationID string `json:"conversation_id" validate:"required"`
	Query          string `json:"query" validate:"required"`
	Stream         bool   `json:"stream"`
}

func (req *OpenAPIAgentChatRequest) Check() error {
	return nil
}

type OpenAPIRagChatRequest struct {
	Query   string     `json:"query" validate:"required"`
	Stream  bool       `json:"stream"`
	History []*History `json:"history"`
}

func (req *OpenAPIRagChatRequest) Check() error {
	return nil
}

type OpenAPIChatflowCreateConversationRequest struct {
	ConversationName string `json:"conversation_name"`
}

func (req *OpenAPIChatflowCreateConversationRequest) Check() error {
	return nil
}

type OpenAPIChatflowChatRequest struct {
	ConversationId string `json:"conversation_id"`
	Query          string `json:"query"`
}

func (req *OpenAPIChatflowChatRequest) Check() error {
	return nil
}
