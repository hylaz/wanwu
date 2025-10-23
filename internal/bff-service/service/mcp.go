package service

import (
	"os"
	"path/filepath"

	assistant_service "github.com/UnicomAI/wanwu/api/proto/assistant-service"
	err_code "github.com/UnicomAI/wanwu/api/proto/err-code"
	mcp_service "github.com/UnicomAI/wanwu/api/proto/mcp-service"
	"github.com/UnicomAI/wanwu/internal/bff-service/model/request"
	"github.com/UnicomAI/wanwu/internal/bff-service/model/response"
	mcp_util "github.com/UnicomAI/wanwu/internal/bff-service/pkg/mcp-util"
	bff_util "github.com/UnicomAI/wanwu/internal/bff-service/pkg/util"
	grpc_util "github.com/UnicomAI/wanwu/pkg/grpc-util"
	"github.com/UnicomAI/wanwu/pkg/log"
	"github.com/gin-gonic/gin"
)

func GetMCPSquareDetail(ctx *gin.Context, userID, orgID, mcpSquareID string) (*response.MCPSquareDetail, error) {
	mcpSquare, err := mcp.GetSquareMCP(ctx.Request.Context(), &mcp_service.GetSquareMCPReq{
		OrgId:       orgID,
		UserId:      userID,
		McpSquareId: mcpSquareID,
	})
	if err != nil {
		return nil, err
	}
	return toMCPSquareDetail(ctx, mcpSquare), nil
}

func GetMCPSquareList(ctx *gin.Context, userID, orgID, category, name string) (*response.ListResult, error) {
	resp, err := mcp.GetSquareMCPList(ctx.Request.Context(), &mcp_service.GetSquareMCPListReq{
		OrgId:    orgID,
		UserId:   userID,
		Category: category,
		Name:     name,
	})
	if err != nil {
		return nil, err
	}
	var list []response.MCPSquareInfo
	for _, mcpSquare := range resp.Infos {
		list = append(list, toMCPSquareInfo(ctx, mcpSquare))
	}
	return &response.ListResult{
		List:  list,
		Total: int64(len(list)),
	}, nil
}

func CreateMCP(ctx *gin.Context, userID, orgID string, req request.MCPCreate) error {
	_, err := mcp.CreateCustomMCP(ctx.Request.Context(), &mcp_service.CreateCustomMCPReq{
		OrgId:       orgID,
		UserId:      userID,
		McpSquareId: req.MCPSquareID,
		Name:        req.Name,
		Desc:        req.Desc,
		From:        req.From,
		SseUrl:      req.SSEURL,
	})
	return err
}

func UpdateMCP(ctx *gin.Context, userID, orgID string, req request.MCPUpdate) error {
	_, err := mcp.UpdateCustomMCP(ctx.Request.Context(), &mcp_service.UpdateCustomMCPReq{
		OrgId:  orgID,
		UserId: userID,
		McpId:  req.MCPID,
		Name:   req.Name,
		Desc:   req.Desc,
		From:   req.From,
		SseUrl: req.SSEURL,
	})
	return err
}

func GetMCP(ctx *gin.Context, mcpID string) (*response.MCPDetail, error) {
	mcpDetail, err := mcp.GetCustomMCP(ctx.Request.Context(), &mcp_service.GetCustomMCPReq{
		McpId: mcpID,
	})
	if err != nil {
		return nil, err
	}
	return toMCPCustomDetail(ctx, mcpDetail), nil
}

func DeleteMCP(ctx *gin.Context, mcpID string) error {
	// 删除智能体表AssistantMCP相关记录
	_, err := assistant.AssistantMCPDeleteByMCPId(ctx.Request.Context(), &assistant_service.AssistantMCPDeleteByMCPIdReq{
		McpId: mcpID,
	})
	if err != nil {
		return err
	}

	_, err = mcp.DeleteCustomMCP(ctx.Request.Context(), &mcp_service.DeleteCustomMCPReq{
		McpId: mcpID,
	})
	return err
}

func GetMCPList(ctx *gin.Context, userID, orgID, name string) (*response.ListResult, error) {
	resp, err := mcp.GetCustomMCPList(ctx.Request.Context(), &mcp_service.GetCustomMCPListReq{
		OrgId:  orgID,
		UserId: userID,
		Name:   name,
	})
	if err != nil {
		return nil, err
	}
	var list []response.MCPInfo
	for _, mcpInfo := range resp.Infos {
		list = append(list, toMCPCustomInfo(ctx, mcpInfo))
	}
	return &response.ListResult{
		List:  list,
		Total: int64(len(list)),
	}, nil
}

func GetMCPSelect(ctx *gin.Context, userID, orgID string, name string) (*response.ListResult, error) {
	resp, err := mcp.GetCustomMCPList(ctx.Request.Context(), &mcp_service.GetCustomMCPListReq{
		OrgId:  orgID,
		UserId: userID,
		Name:   name,
	})
	if err != nil {
		return nil, err
	}
	var list []response.MCPSelect
	for _, mcpInfo := range resp.Infos {
		list = append(list, response.MCPSelect{
			MCPID:       mcpInfo.McpId,
			MCPSquareID: mcpInfo.Info.McpSquareId,
			UniqueId:    bff_util.ConcatAssistantToolUniqueId("mcp", mcpInfo.McpId),
			Name:        mcpInfo.Info.Name,
			Description: mcpInfo.Info.Desc,
			ServerFrom:  mcpInfo.Info.From,
			ServerURL:   mcpInfo.SseUrl,
		})
	}
	return &response.ListResult{
		List:  list,
		Total: int64(len(list)),
	}, nil
}

func GetMCPToolList(ctx *gin.Context, mcpID, sseUrl string) (*response.MCPToolList, error) {
	if mcpID != "" {
		mcpDetail, err := mcp.GetCustomMCP(ctx.Request.Context(), &mcp_service.GetCustomMCPReq{
			McpId: mcpID,
		})
		if err != nil {
			return nil, err
		}
		if mcpDetail.SseUrl != "" {
			sseUrl = mcpDetail.SseUrl
		}
	}
	if sseUrl == "" {
		return nil, grpc_util.ErrorStatus(err_code.Code_BFFInvalidArg, "sseUrl empty")
	}

	tools, err := mcp_util.ListTools(ctx.Request.Context(), sseUrl)
	if err != nil {
		return nil, grpc_util.ErrorStatus(err_code.Code_BFFGeneral, err.Error())
	}
	return &response.MCPToolList{Tools: tools}, nil
}

// --- internal ---

func cacheMCPAvatar(ctx *gin.Context, avatarPath string) request.Avatar {
	avatar := request.Avatar{}
	if avatarPath == "" {
		return avatar
	}
	avatarCacheMu.Lock()
	defer avatarCacheMu.Unlock()

	filePath := filepath.Join(mcpAvatarCacheLocalDir, avatarPath)

	_, err := os.Stat(filePath)
	// 1 文件存在
	if err == nil {
		avatar.Path = filepath.Join("/v1", filePath)
		return avatar
	}
	// 2 系统错误
	if !os.IsNotExist(err) {
		log.Errorf("cache mcp avatar %v check %v exist err: %v", avatarPath, filePath, err)
		return avatar
	}
	// 3. 文件不存在
	// 3.1 创建目录
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		log.Errorf("cache mcp avatar %v mkdir %v err: %v", avatarPath, filepath.Dir(filePath))
		return avatar
	}
	// 3.2 下载文件
	resp, err := mcp.GetMCPAvatar(ctx.Request.Context(), &mcp_service.GetMCPAvatarReq{AvatarPath: avatarPath})
	if err != nil {
		log.Errorf("cache mcp avatar %v download err: %v", avatarPath, err)
		return avatar
	}
	// 3.3 写入文件
	if err := os.WriteFile(filePath, resp.Data, 0644); err != nil {
		log.Errorf("cache mcp avatar %v write file %v err: %v", avatarPath, filePath, err)
		return avatar
	}
	avatar.Path = filepath.Join("/v1", filePath)
	return avatar
}

func toMCPCustomDetail(ctx *gin.Context, mcpDetail *mcp_service.CustomMCPDetail) *response.MCPDetail {
	return &response.MCPDetail{
		MCPInfo: response.MCPInfo{
			MCPID:         mcpDetail.McpId,
			SSEURL:        mcpDetail.SseUrl,
			MCPSquareInfo: toMCPSquareInfo(ctx, mcpDetail.Info),
		},
		MCPSquareIntro: toMCPSquareIntro(mcpDetail.Intro),
	}
}

func toMCPCustomInfo(ctx *gin.Context, mcpInfo *mcp_service.CustomMCPInfo) response.MCPInfo {
	return response.MCPInfo{
		MCPID:         mcpInfo.McpId,
		SSEURL:        mcpInfo.SseUrl,
		MCPSquareInfo: toMCPSquareInfo(ctx, mcpInfo.Info),
	}
}

func toMCPSquareDetail(ctx *gin.Context, mcpSquare *mcp_service.SquareMCPDetail) *response.MCPSquareDetail {
	ret := &response.MCPSquareDetail{
		MCPSquareInfo:  toMCPSquareInfo(ctx, mcpSquare.Info),
		MCPSquareIntro: toMCPSquareIntro(mcpSquare.Intro),
		MCPTools: response.MCPTools{
			SSEURL:    mcpSquare.Tool.SseUrl,
			HasCustom: mcpSquare.Tool.HasCustom,
		},
	}
	for _, tool := range mcpSquare.Tool.Tools {
		ret.MCPTools.Tools = append(ret.MCPTools.Tools, toMCPTool(tool))
	}
	return ret
}

func toMCPSquareInfo(ctx *gin.Context, mcpSquareInfo *mcp_service.SquareMCPInfo) response.MCPSquareInfo {
	return response.MCPSquareInfo{
		MCPSquareID: mcpSquareInfo.McpSquareId,
		Avatar:      cacheMCPAvatar(ctx, mcpSquareInfo.AvatarPath),
		Name:        mcpSquareInfo.Name,
		Desc:        mcpSquareInfo.Desc,
		From:        mcpSquareInfo.From,
		Category:    mcpSquareInfo.Category,
	}
}

func toMCPSquareIntro(mcpSquareIntro *mcp_service.SquareMCPIntro) response.MCPSquareIntro {
	if mcpSquareIntro == nil {
		return response.MCPSquareIntro{}
	}
	return response.MCPSquareIntro{
		Summary:  mcpSquareIntro.Summary,
		Feature:  mcpSquareIntro.Feature,
		Scenario: mcpSquareIntro.Scenario,
		Manual:   mcpSquareIntro.Manual,
		Detail:   mcpSquareIntro.Detail,
	}
}

func toMCPTool(tool *mcp_service.MCPTool) response.MCPTool {
	ret := response.MCPTool{
		Name:        tool.Name,
		Description: tool.Description,
		InputSchema: response.MCPToolInputSchema{
			Type:       tool.InputSchema.GetType(),
			Required:   tool.InputSchema.GetRequired(),
			Properties: make(map[string]response.MCPToolInputSchemaValue),
		},
	}
	for k, v := range tool.InputSchema.GetProperties() {
		ret.InputSchema.Properties[k] = response.MCPToolInputSchemaValue{
			Type:        v.Type,
			Description: v.Description,
		}
	}
	return ret
}
