package request

import "mime/multipart"

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

type WorkflowUploadFileReq struct {
	File     *multipart.FileHeader `form:"file" json:"file" validate:"required"` // 二进制格式
	FileName string                `form:"fileName" json:"fileName" validate:"required"`
}

func (u *WorkflowUploadFileReq) Check() error {
	return nil
}

type WorkflowUploadFileByBase64Req struct {
	File     string `form:"file" json:"file" validate:"required"` // base64格式
	FileName string `form:"fileName" json:"fileName" validate:"required"`
}

func (u *WorkflowUploadFileByBase64Req) Check() error {
	return nil
}

type WorkflowConvertReq struct {
	WorkflowID string `json:"workflow_id" validate:"required"`
}

func (r *WorkflowConvertReq) Check() error {
	return nil
}

type WorkflowRunReq struct {
	WorkflowID string         `json:"workflow_id" validate:"required"`
	Input      map[string]any `json:"input" `
}

func (r *WorkflowRunReq) Check() error {
	return nil
}
