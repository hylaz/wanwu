package nodes

import (
	agentPrompt "github.com/UnicomAI/wanwu/internal/agent-service/service/agent-message-flow/prompt"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
)

var (
	ChatPrompt = prompt.FromMessages(schema.Jinja2,
		schema.SystemMessage(agentPrompt.REACT_SYSTEM_PROMPT_JINJA3),
		schema.MessagesPlaceholder(placeholderOfChatHistory, true),
		schema.MessagesPlaceholder(placeholderOfUserInput, false),
	)
)
