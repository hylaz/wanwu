package response

type AgentScopeWorkFlowInfo struct {
	Id           string `json:"id"`           // 应用id
	ConfigDesc   string `json:"configDesc"`   // 应用简介
	ConfigENName string `json:"configENName"` // 应用英文名称
	ConfigName   string `json:"configName"`   // 应用名称
	ExampleFlag  int    `json:"example_flag"` // 示例标识
	IsStream     int    `json:"is_stream"`    // 流式标识
	OrgID        string `json:"orgID"`        // 组织ID
	Status       string `json:"status"`       // 应用状态
	UpdatedTime  string `json:"updatedTime"`  // 应用更新时间
	UserID       string `json:"userID"`       // 用户ID
}

type AgentScopeWorkFlowListResp struct {
	Code    int                           `json:"code"`
	Message string                        `json:"msg"`
	Data    *AgentScopeWorkFlowPageResult `json:"data"`
}

type AgentScopeWorkFlowPageResult struct {
	List     []AgentScopeWorkFlowInfo `json:"list"`
	Total    int64                    `json:"total"`
	PageNo   int                      `json:"pageNo"`
	PageSize int                      `json:"pageSize"`
}

type AgentScopeDeleteWorkFlowResp struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
}
