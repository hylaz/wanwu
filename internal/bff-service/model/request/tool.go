package request

type CustomToolCreate struct {
	Name          string                      `json:"name" validate:"required"`        // 名称
	Description   string                      `json:"description" validate:"required"` // 描述
	ApiAuth       CustomToolApiAuthWebRequest `json:"apiAuth" validate:"required"`     // api身份认证
	Schema        string                      `json:"schema"  validate:"required"`     // schema
	PrivacyPolicy string                      `json:"privacyPolicy"`                   // 隐私政策
}

func (req *CustomToolCreate) Check() error { return nil }

type CustomToolApiAuthWebRequest struct {
	Type             string `json:"type" validate:"required,oneof='None' 'API Key'"` // 认证类型 None 或 APIKey
	APIKey           string `json:"apiKey"`                                          // apiKey 仅当认证类型为API Key时必填
	CustomHeaderName string `json:"customHeaderName"`                                // Custom Header Name 仅当认证类型为API Key时必填
	AuthType         string `json:"authType" validate:"omitempty,oneof=Custom"`      // Auth类型 仅当认证类型为API Key时必填，也可以为空
}

func (req *CustomToolApiAuthWebRequest) Check() error { return nil }

type CustomToolUpdateReq struct {
	CustomToolID  string                      `json:"customToolId" validate:"required"` // 自定义工具ID
	Name          string                      `json:"name" validate:"required"`         // 名称
	Description   string                      `json:"description" validate:"required"`  // 描述
	ApiAuth       CustomToolApiAuthWebRequest `json:"apiAuth" validate:"required"`      // api身份认证
	Schema        string                      `json:"schema"  validate:"required"`      // schema
	PrivacyPolicy string                      `json:"privacyPolicy"`                    // 隐私政策
}

func (req *CustomToolUpdateReq) Check() error { return nil }

type CustomToolIDReq struct {
	CustomToolID string `json:"customToolId" validate:"required"` // 自定义工具id
}

func (req *CustomToolIDReq) Check() error { return nil }

type CustomToolSchemaReq struct {
	Schema string `json:"schema" validate:"required"` // schema
}

func (req *CustomToolSchemaReq) Check() error { return nil }

type ToolSquareAPIKeyReq struct {
	ToolSquareID string `json:"toolSquareId" validate:"required"` // 广场toolId
	APIKey       string `json:"apiKey"`                           // apiKey
}

func (req *ToolSquareAPIKeyReq) Check() error { return nil }

type ToolActionListReq struct {
	ToolId   string `form:"toolId" json:"toolId" validate:"required"`                          // 工具id
	ToolType string `form:"toolType" json:"toolType" validate:"required,oneof=builtin custom"` // 工具类型
}

func (req *ToolActionListReq) Check() error { return nil }

type ToolActionReq struct {
	ToolId     string `form:"toolId" json:"toolId" validate:"required"`                          // 工具id
	ToolType   string `form:"toolType" json:"toolType" validate:"required,oneof=builtin custom"` // 工具类型
	ActionName string `form:"actionName" json:"actionName" validate:"required"`                  // action名称
}

func (req *ToolActionReq) Check() error { return nil }

type CreatePromptByTemplateReq struct {
	TemplateId string `json:"templateId" validate:"required"`
	AppBriefConfig
}

func (req *CreatePromptByTemplateReq) Check() error { return nil }
