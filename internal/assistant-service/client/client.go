package client

import (
	"context"

	err_code "github.com/UnicomAI/wanwu/api/proto/err-code"
	"github.com/UnicomAI/wanwu/internal/assistant-service/client/model"
)

type IClient interface {
	//================Assistant================
	CreateAssistant(ctx context.Context, assistant *model.Assistant) *err_code.Status
	UpdateAssistant(ctx context.Context, assistant *model.Assistant) *err_code.Status
	DeleteAssistant(ctx context.Context, assistantID uint32) *err_code.Status
	GetAssistant(ctx context.Context, assistantID uint32) (*model.Assistant, *err_code.Status)
	GetAssistantsByIDs(ctx context.Context, assistantIDs []uint32) ([]*model.Assistant, *err_code.Status)
	GetAssistantList(ctx context.Context, userID, orgID string, name string) ([]*model.Assistant, int64, *err_code.Status)
	CheckSameAssistantName(ctx context.Context, userID, orgID, name, assistantID string) *err_code.Status
	CopyAssistant(ctx context.Context, assistant *model.Assistant, workflows []*model.AssistantWorkflow, mcps []*model.AssistantMCP, customTools []*model.AssistantCustom) (uint32, *err_code.Status)

	//================AssistantWorkflow================
	CreateAssistantWorkflow(ctx context.Context, workflow *model.AssistantWorkflow) *err_code.Status
	DeleteAssistantWorkflow(ctx context.Context, assistantId uint32, workflowId string) *err_code.Status
	UpdateAssistantWorkflow(ctx context.Context, workflow *model.AssistantWorkflow) *err_code.Status
	GetAssistantWorkflow(ctx context.Context, assistantId uint32, workflowId string) (*model.AssistantWorkflow, *err_code.Status)
	GetAssistantWorkflowsByAssistantID(ctx context.Context, assistantId uint32) ([]*model.AssistantWorkflow, *err_code.Status)
	DeleteAssistantWorkflowByWorkflowId(ctx context.Context, workflowId string) *err_code.Status

	//================AssistantMCP================
	CreateAssistantMCP(ctx context.Context, assistantId uint32, mcpId string, userId, orgID string) *err_code.Status
	DeleteAssistantMCP(ctx context.Context, assistantId uint32, mcpId string) *err_code.Status
	GetAssistantMCP(ctx context.Context, assistantId uint32, mcpId string) (*model.AssistantMCP, *err_code.Status)
	DeleteAssistantMCPByMCPId(ctx context.Context, mcpId string) *err_code.Status
	GetAssistantMCPList(ctx context.Context, assistantId uint32) ([]*model.AssistantMCP, *err_code.Status)
	UpdateAssistantMCP(ctx context.Context, mcp *model.AssistantMCP) *err_code.Status

	//================AssistantCustom================
	CreateAssistantCustom(ctx context.Context, assistantId uint32, customId string, userId, orgID string) *err_code.Status
	DeleteAssistantCustom(ctx context.Context, assistantId uint32, customId string) *err_code.Status
	DeleteAssistantCustomByCustomToolId(ctx context.Context, customId string) *err_code.Status
	GetAssistantCustom(ctx context.Context, assistantId uint32, customId string) (*model.AssistantCustom, *err_code.Status)
	UpdateAssistantCustom(ctx context.Context, custom *model.AssistantCustom) *err_code.Status
	GetAssistantCustomList(ctx context.Context, assistantId uint32) ([]*model.AssistantCustom, *err_code.Status)

	//================Conversation================
	CreateConversation(ctx context.Context, conversation *model.Conversation) *err_code.Status
	UpdateConversation(ctx context.Context, conversation *model.Conversation) *err_code.Status
	DeleteConversation(ctx context.Context, conversationID uint32) *err_code.Status
	GetConversation(ctx context.Context, conversationID uint32) (*model.Conversation, *err_code.Status)
	GetConversationList(ctx context.Context, assistantID, userID, orgID string, offset, limit int32) ([]*model.Conversation, int64, *err_code.Status)
	DeleteConversationByAssistantID(ctx context.Context, assistantID, userID, orgID string) *err_code.Status
}
