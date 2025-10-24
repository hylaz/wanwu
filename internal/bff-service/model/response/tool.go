package response

import (
	"github.com/ThinkInAIXYZ/go-mcp/protocol"
	"github.com/UnicomAI/wanwu/internal/bff-service/model/request"
)

type CustomToolDetail struct {
	CustomToolInfo
	Schema        string                              `json:"schema"`        // schema
	ApiAuth       request.CustomToolApiAuthWebRequest `json:"apiAuth"`       // apiAuth
	ApiList       []CustomToolActionInfo              `json:"apiList"`       // action列表
	PrivacyPolicy string                              `json:"privacyPolicy"` // 隐私政策
	ToolSquareID  string                              `json:"toolSquareId"`  // 广场mcpId(非空表示来源于广场)
}

type CustomToolInfo struct {
	CustomToolId string `json:"customToolId"` // 自定义工具id
	Name         string `json:"name"`         // 名称
	Description  string `json:"description"`  // 描述
}

type CustomToolActionInfo struct {
	Name   string `json:"name"`
	Desc   string `json:"desc"`
	Method string `json:"method"`
	Path   string `json:"path"`
}

type CustomToolSelect struct {
	UniqueId string `json:"uniqueId"` // 随机unique id(每次动态生成)
	CustomToolInfo
}

type ToolSquareDetail struct {
	ToolSquareInfo
	ToolSquareActions
	Schema string `json:"schema"`
}

type ToolSquareInfo struct {
	ToolSquareID string         `json:"toolSquareId"` // 广场mcpId(非空表示来源于广场)
	Avatar       request.Avatar `json:"avatar"`       // 图标
	Name         string         `json:"name"`         // 名称
	Desc         string         `json:"desc"`         // 描述
	Tags         []string       `json:"tags"`         // 标签
}

type ToolSquareActions struct {
	NeedApiKeyInput bool      `json:"needApiKeyInput"` // 是否需要apiKey输入
	APIKey          string    `json:"apiKey"`          // apiKey
	Tools           []MCPTool `json:"tools"`           // action列表
	Detail          string    `json:"detail"`          // 详细描述
	ActionSum       int64     `json:"actionSum"`       // action总数
}

type ToolSelect struct {
	UniqueId string `json:"uniqueId"` // unique id
	ToolInfo
}

type ToolInfo struct {
	ToolId          string `json:"toolId"`          // 工具id
	ToolName        string `json:"toolName"`        // 工具名称
	ToolType        string `json:"toolType"`        // 工具类型
	Desc            string `json:"desc"`            // 工具描述
	NeedApiKeyInput bool   `json:"needApiKeyInput"` // 是否需要apiKey输入
	APIKey          string `json:"apiKey"`          // apiKey
}

type ToolActionList struct {
	Actions []*protocol.Tool `json:"actions"` // action列表
}

type ToolActionDetail struct {
	NeedApiKeyInput bool           `json:"needApiKeyInput"` // 是否需要apiKey输入
	APIKey          string         `json:"apiKey"`          // apiKey
	Action          *protocol.Tool `json:"action"`          // action列表
}
