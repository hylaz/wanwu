package response

type ToolSelect4Workflow struct {
	ToolID   string                `json:"toolId"`
	ToolName string                `json:"toolName"`
	ToolType string                `json:"toolType"`
	IconUrl  string                `json:"iconUrl"`
	ApiKey   string                `json:"apiKey"`
	Desc     string                `json:"desc"`
	Actions  []ToolAction4Workflow `json:"actions"`
}

type ToolAction4Workflow struct {
	ActionName string `json:"actionName"`
	ActionID   string `json:"actionId"`
	Desc       string `json:"desc"`
}
