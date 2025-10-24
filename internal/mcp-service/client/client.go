package client

import (
	"context"

	errs "github.com/UnicomAI/wanwu/api/proto/err-code"
	"github.com/UnicomAI/wanwu/internal/mcp-service/client/model"
)

type IClient interface {
	CheckMCPExist(ctx context.Context, orgID, userID, mcpSquareID string) (bool, *errs.Status)
	GetMCP(ctx context.Context, mcpID uint32) (*model.MCPClient, *errs.Status)
	CreateMCP(ctx context.Context, mcp *model.MCPClient) *errs.Status
	UpdateMCP(ctx context.Context, mcp *model.MCPClient) *errs.Status
	DeleteMCP(ctx context.Context, mcpID uint32) *errs.Status
	ListMCPs(ctx context.Context, orgID, userID, name string) ([]*model.MCPClient, *errs.Status)
	ListMCPsByMCPIdList(ctx context.Context, mcpIDList []uint32) ([]*model.MCPClient, *errs.Status)

	CreateCustomTool(ctx context.Context, customTool *model.CustomTool) *errs.Status
	GetCustomTool(ctx context.Context, customTool *model.CustomTool) (*model.CustomTool, *errs.Status)
	ListCustomTools(ctx context.Context, orgID, userID, name string) ([]*model.CustomTool, *errs.Status)
	ListCustomToolsByCustomToolIDs(ctx context.Context, customToolIDs []uint32) ([]*model.CustomTool, *errs.Status)
	UpdateCustomTool(ctx context.Context, customTool *model.CustomTool) *errs.Status
	DeleteCustomTool(ctx context.Context, customToolID uint32) *errs.Status
}
