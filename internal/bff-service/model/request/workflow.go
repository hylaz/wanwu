package request

type WorkflowIDReq struct {
	WorkflowID string `json:"workflow_id" validate:"required"`
}

func (r *WorkflowIDReq) Check() error {
	return nil
}

type GetWorkflowListReq struct {
	UserId string `form:"userId" json:"userId" validate:"required" `
	OrgId  string `form:"orgId" json:"orgId" validate:"required" `
}

func (g *GetWorkflowListReq) Check() error {
	return nil
}

type CreateWorkflowByTemplateReq struct {
	TemplateId string `json:"templateId" validate:"required"`
	AppBriefConfig
}

func (r *CreateWorkflowByTemplateReq) Check() error {
	return nil
}
