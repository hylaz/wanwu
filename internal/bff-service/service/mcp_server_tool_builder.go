package service

import (
	"context"

	mcp_service "github.com/UnicomAI/wanwu/api/proto/mcp-service"
	"github.com/UnicomAI/wanwu/pkg/constant"
	openapi3_util "github.com/UnicomAI/wanwu/pkg/openapi3-util"
	"github.com/UnicomAI/wanwu/pkg/util"
)

type mcpServerToolBuilder interface {
	MCPServerToolType() string
	AppID() string
	AppName() string
	GetOpenapiSchema(ctx context.Context) (string, *openapi3_util.Auth, error)
}

// --- mcpServerCustomToolBuilder ---

type mcpServerCustomToolBuilder struct {
	customToolID   string
	customToolName string
}

func (mcpServerCustomToolBuilder) MCPServerToolType() string {
	return constant.MCPServerToolTypeCustomTool
}

func (builder *mcpServerCustomToolBuilder) AppID() string {
	return builder.customToolID
}

func (builder *mcpServerCustomToolBuilder) AppName() string {
	return builder.customToolName
}

func (builder *mcpServerCustomToolBuilder) GetOpenapiSchema(ctx context.Context) (string, *openapi3_util.Auth, error) {
	customToolInfo, err := mcp.GetCustomToolInfo(ctx, &mcp_service.GetCustomToolInfoReq{
		CustomToolId: builder.customToolID,
	})
	if err != nil {
		return "", nil, err
	}
	builder.customToolName = customToolInfo.Name
	auth, err := util.ConvertApiAuthWebRequestProto(customToolInfo.ApiAuth)
	if err != nil {
		return "", nil, err
	}
	return customToolInfo.Schema, auth, nil
}

// --- mcpServerBuiltInToolBuilder ---

type mcpServerBuiltInToolBuilder struct {
	toolSquareId   string
	toolSquareName string
}

func (mcpServerBuiltInToolBuilder) MCPServerToolType() string {
	return constant.MCPServerToolTypeBuiltInTool
}

func (builder *mcpServerBuiltInToolBuilder) AppID() string {
	return builder.toolSquareId
}

func (builder *mcpServerBuiltInToolBuilder) AppName() string {
	return builder.toolSquareName
}

func (builder *mcpServerBuiltInToolBuilder) GetOpenapiSchema(ctx context.Context) (string, *openapi3_util.Auth, error) {
	toolSquareInfo, err := mcp.GetSquareTool(ctx, &mcp_service.GetSquareToolReq{
		ToolSquareId: builder.toolSquareId,
		Identity:     &mcp_service.Identity{},
	})
	if err != nil {
		return "", nil, err
	}
	builder.toolSquareName = toolSquareInfo.Info.Name
	auth, err := util.ConvertApiAuthWebRequestProto(toolSquareInfo.BuiltInTools.ApiAuth)
	if err != nil {
		return "", nil, err
	}
	return toolSquareInfo.Schema, auth, nil
}

// --- mcpServerOpenapiSchemaBuilder ---

type mcpServerOpenapiSchemaBuilder struct {
	schema string
	name   string
	auth   util.ApiAuthWebRequest
}

func (mcpServerOpenapiSchemaBuilder) MCPServerToolType() string {
	return constant.MCPServerToolTypeOpenAPI
}

func (builder *mcpServerOpenapiSchemaBuilder) AppID() string {
	return ""
}

func (builder *mcpServerOpenapiSchemaBuilder) AppName() string {
	return builder.name
}

func (builder *mcpServerOpenapiSchemaBuilder) GetOpenapiSchema(ctx context.Context) (string, *openapi3_util.Auth, error) {
	auth, err := builder.auth.ToOpenapiAuth()
	if err != nil {
		return "", nil, err
	}
	return builder.schema, auth, nil
}
