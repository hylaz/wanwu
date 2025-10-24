// @Author wangxm 10/24/星期五 14:46:00
package service

import (
	"github.com/ThinkInAIXYZ/go-mcp/protocol"
	mcp_service "github.com/UnicomAI/wanwu/api/proto/mcp-service"
	"github.com/UnicomAI/wanwu/internal/bff-service/model/request"
	"github.com/UnicomAI/wanwu/internal/bff-service/model/response"
	"github.com/gin-gonic/gin"
)

func GetToolSelect(ctx *gin.Context, userID, orgID string, req request.ToolSelectReq) (*response.ListResult, error) {
	resp, err := mcp.GetToolSelect(ctx.Request.Context(), &mcp_service.GetToolSelectReq{
		Name: req.ToolName,
		Identity: &mcp_service.Identity{
			UserId: userID,
			OrgId:  orgID,
		},
	})
	if err != nil {
		return nil, err
	}

	var list []response.ToolSelect
	for _, item := range resp.List {
		list = append(list, response.ToolSelect{
			UniqueId: "tool-" + item.ToolId,
			ToolInfo: response.ToolInfo{
				ToolId:          item.ToolId,
				ToolName:        item.ToolName,
				ToolType:        item.ToolType,
				Desc:            item.Desc,
				APIKey:          item.ApiKey,
				NeedApiKeyInput: item.NeedApiKeyInput,
			},
		})
	}
	return &response.ListResult{
		List:  list,
		Total: int64(len(list)),
	}, nil
}

func GetToolActionList(ctx *gin.Context, userID, orgID string, req request.ToolActionListReq) (*response.ToolActionList, error) {
	resp, err := mcp.GetToolActionList(ctx.Request.Context(), &mcp_service.GetToolActionListReq{
		ToolId:   req.ToolId,
		ToolType: req.ToolType,
		Identity: &mcp_service.Identity{
			UserId: userID,
			OrgId:  orgID,
		},
	})
	if err != nil {
		return nil, err
	}

	var list []*protocol.Tool
	for _, item := range resp.Actions {
		list = append(list, toTool(item))
	}

	return &response.ToolActionList{
		Actions: list,
	}, nil
}

func GetToolActionDetail(ctx *gin.Context, userID, orgID string, req request.ToolActionReq) (*response.ToolActionDetail, error) {
	resp, err := mcp.GetToolAction(ctx.Request.Context(), &mcp_service.GetToolActionReq{
		ToolId:     req.ToolId,
		ToolType:   req.ToolType,
		ActionName: req.ActionName,
		Identity: &mcp_service.Identity{
			UserId: userID,
			OrgId:  orgID,
		},
	})
	if err != nil {
		return nil, err
	}

	actionDetail := &response.ToolActionDetail{
		NeedApiKeyInput: resp.NeedApiKeyInput,
		APIKey:          resp.ApiKey,
		Action:          toTool(resp.Action),
	}

	return actionDetail, nil
}

// --- internal ---
func toTool(tool *mcp_service.ToolActionInfo) *protocol.Tool {
	ret := &protocol.Tool{
		Name:        tool.Name,
		Description: tool.Desc,
		InputSchema: protocol.InputSchema{
			Type:       protocol.InputSchemaType(tool.InputSchema.GetType()),
			Required:   tool.InputSchema.GetRequired(),
			Properties: make(map[string]*protocol.Property),
		},
	}
	for k, v := range tool.InputSchema.Properties {
		ret.InputSchema.Properties[k] = &protocol.Property{
			Type:        protocol.DataType(v.Type),
			Description: v.Description,
		}
	}
	return ret
}
