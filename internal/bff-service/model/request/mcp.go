package request

type MCPIDReq struct {
	MCPID string `json:"mcpId" validate:"required"`
}

func (req *MCPIDReq) Check() error {
	return nil
}

type MCPCreate struct {
	Avatar      Avatar `json:"avatar"`                     // 图标
	MCPSquareID string `json:"mcpSquareId"`                // 广场mcpId(非空表示来源于广场)
	Name        string `json:"name" validate:"required"`   // 名称
	Desc        string `json:"desc" validate:"required"`   // 描述
	From        string `json:"from" validate:"required"`   // 来源
	SSEURL      string `json:"sseUrl" validate:"required"` // SSE URL
}

func (req *MCPCreate) Check() error {
	return nil
}

type MCPUpdate struct {
	Avatar Avatar `json:"avatar"` // 图标
	MCPID  string `json:"mcpId" validate:"required"`
	Name   string `json:"name" validate:"required"`   // 名称
	Desc   string `json:"desc" validate:"required"`   // 描述
	From   string `json:"from" validate:"required"`   // 来源
	SSEURL string `json:"sseUrl" validate:"required"` // SSE URL
}

func (req *MCPUpdate) Check() error {
	return nil
}

type MCPActionListReq struct {
	ToolId   string `form:"toolId" json:"toolId" validate:"required"`
	ToolType string `form:"toolType" json:"toolType" validate:"required,oneof=mcp mcpserver"`
}

func (req *MCPActionListReq) Check() error {
	return nil
}

type MCPActionReq struct {
	ToolId     string `json:"toolId" validate:"required"`
	ToolType   string `json:"toolType" validate:"required,oneof=mcp mcpserver"`
	ActionName string `json:"actionName" validate:"required"`
}

func (req *MCPActionReq) Check() error {
	return nil
}
