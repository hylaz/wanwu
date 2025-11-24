package response

import (
	"github.com/UnicomAI/wanwu/internal/bff-service/config"
)

type CozeWorkflowModelInfo struct {
	ModelInfo
	ModelAbility CozeWorkflowModelInfoAbility `json:"model_ability"`
	ModelParams  []config.WorkflowModelParam  `json:"model_params"`
}

type CozeWorkflowModelInfoAbility struct {
	CotDisplay         bool `json:"cot_display"`
	FunctionCall       bool `json:"function_call"`
	ImageUnderstanding bool `json:"image_understanding"`
	AudioUnderstanding bool `json:"audio_understanding"`
	VideoUnderstanding bool `json:"video_understanding"`
}

type CozeWorkflowListResp struct {
	Code int                   `json:"code"`
	Msg  string                `json:"msg"`
	Data *CozeWorkflowListData `json:"data,omitempty"`
}

type CozeWorkflowListData struct {
	Workflows []*CozeWorkflowListDataWorkflow `json:"workflow_list"`
}

type CozeWorkflowListDataWorkflow struct {
	WorkflowId string `json:"workflow_id"`
	Name       string `json:"name"`
	Desc       string `json:"desc"`
	URL        string `json:"url"`
	CreateTime int64  `json:"create_time"`
	UpdateTime int64  `json:"update_time"`
}

type CozeWorkflowIDResp struct {
	Code int                 `json:"code"`
	Msg  string              `json:"msg"`
	Data *CozeWorkflowIDData `json:"data,omitempty"`
}

type CozeWorkflowIDData struct {
	WorkflowID string `json:"workflow_id"`
}

type CozeWorkflowDeleteResp struct {
	Code int                     `json:"code"`
	Msg  string                  `json:"msg"`
	Data *CozeWorkflowDeleteData `json:"data,omitempty"`
}

type CozeWorkflowDeleteData struct {
	Status int64 `json:"status"`
}

func (d *CozeWorkflowDeleteData) GetStatus() int64 {
	if d == nil {
		return 0
	}
	return d.Status
}

type CozeWorkflowExportResp struct {
	Code int                     `json:"code"`
	Msg  string                  `json:"msg"`
	Data *CozeWorkflowExportData `json:"data,omitempty"`
}

type CozeWorkflowExportData struct {
	WorkflowName string `json:"name"`
	WorkflowDesc string `json:"desc"`
	Schema       string `json:"schema"`
}

type ToolDetail4Workflow struct {
	Inputs     []interface{} `json:"inputs"`
	Outputs    []interface{} `json:"outputs"`
	ActionName string        `json:"actionName"`
	ActionID   string        `json:"actionId"`
	IconUrl    string        `json:"iconUrl"`
}

// ToolActionParamWithoutTypeList4Workflow type非list的定义
type ToolActionParamWithoutTypeList4Workflow struct {
	Input       struct{}      `json:"input"`
	Description string        `json:"description"`
	Name        string        `json:"name"`
	Type        string        `json:"type"` // 非list
	Required    bool          `json:"required"`
	Children    []interface{} `json:"schema"`
}

// ToolActionParamWithTypeList4Workflow type是list的定义
type ToolActionParamWithTypeList4Workflow struct {
	Input       struct{}                           `json:"input"`
	Description string                             `json:"description"`
	Name        string                             `json:"name"`
	Type        string                             `json:"type"` // list
	Required    bool                               `json:"required"`
	Schema      ToolActionParamInTypeList4Workflow `json:"schema"`
}

type ToolActionParamInTypeList4Workflow struct {
	Type     string        `json:"type"`
	Children []interface{} `json:"schema"`
}

type CozeCreateConversationResponse struct {
	Code             int64                 `thrift:"code,1" form:"code" json:"code" query:"code"`
	Msg              string                `thrift:"msg,2" form:"msg" json:"msg" query:"msg"`
	ConversationData *CozeConversationData `thrift:"ConversationData,3,optional" form:"data" json:"data,omitempty"`
}

type CozeConversationData struct {
	Id            int64             `thrift:"Id,1" form:"id" json:"id,string"`
	CreatedAt     int64             `thrift:"CreatedAt,2" form:"created_at" json:"created_at"`
	MetaData      map[string]string `thrift:"MetaData,3" form:"meta_data" json:"meta_data"`
	CreatorID     *int64            `thrift:"CreatorID,4,optional" form:"creator_d" json:"creator_d,string,omitempty"`
	ConnectorID   *int64            `thrift:"ConnectorID,5,optional" form:"connector_id" json:"connector_id,string,omitempty"`
	LastSectionID *int64            `thrift:"LastSectionID,6,optional" form:"last_section_id" json:"last_section_id,string,omitempty"`
	AccountID     *int64            `thrift:"AccountID,7,optional" form:"account_id" json:"account_id,omitempty"`
}

type UploadFileByWorkflowResp struct {
	Url string `json:"url"`
	Uri string `json:"uri"`
}
