package request

type MCPIDReq struct {
	MCPID string `json:"mcpId" validate:"required"`
}

func (req *MCPIDReq) Check() error {
	return nil
}

type MCPCreate struct {
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
	MCPID  string `json:"mcpId" validate:"required"`
	Name   string `json:"name"`   // 名称
	Desc   string `json:"desc"`   // 描述
	From   string `json:"from"`   // 来源
	SSEURL string `json:"sseUrl"` // SSE URL
}

func (req *MCPUpdate) Check() error {
	return nil
}
