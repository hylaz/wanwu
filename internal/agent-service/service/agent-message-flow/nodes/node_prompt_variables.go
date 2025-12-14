package nodes

import (
	"context"
	"fmt"
	"time"

	"github.com/UnicomAI/wanwu/internal/agent-service/model/request"
	"github.com/UnicomAI/wanwu/internal/agent-service/service/agent-message-flow/prompt"
	"github.com/cloudwego/eino/schema"
)

const (
	placeholderOfUserInput   = "_user_input"
	placeholderOfChatHistory = "_chat_history"
)

type PromptVariables struct {
	Avs map[string]string
}

func (p *PromptVariables) AssemblePromptVariables(ctx context.Context, reqContext *request.AgentChatContext) (variables map[string]any, err error) {
	req := reqContext.AgentChatReq
	variables = make(map[string]any)

	variables[prompt.PlaceholderOfTime] = time.Now().Format("Monday 2006/01/02 15:04:05 -07")
	variables[prompt.PlaceholderOfAgentName] = req.AgentBaseParams.Name

	var input = req.Input
	if len(req.UploadFile) > 0 {
		input += "\n用户上传的文档连接为:" + req.UploadFile[0]
	}
	variables[placeholderOfUserInput] = []*schema.Message{schema.UserMessage(input)}

	// Handling conversation history
	if len(req.ModelParams.History) > 0 {
		// Add chat history to variable
		variables[placeholderOfChatHistory] = buildHistory(req.ModelParams.History, req.ModelParams.MaxHistory)
	}

	if p.Avs != nil {
		var memoryVariablesList []string
		for k, v := range p.Avs {
			variables[k] = v
			memoryVariablesList = append(memoryVariablesList, fmt.Sprintf("%s: %s\n", k, v))
		}
		variables[prompt.PlaceholderOfVariables] = memoryVariablesList
	}

	return variables, nil
}

func buildHistory(history []request.AssistantConversionHistory, maxHistory int) []*schema.Message {
	var historyList []*schema.Message

	// 处理所有历史记录
	for _, conversionHistory := range history {
		historyList = append(historyList, schema.UserMessage(conversionHistory.Query))
	}
	if maxHistory <= 0 {
		return historyList
	}
	// 只返回最后maxHistory条
	if len(historyList) > maxHistory {
		return historyList[len(historyList)-maxHistory:]
	}
	return historyList
}
