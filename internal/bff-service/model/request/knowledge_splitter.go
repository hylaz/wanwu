package request

type CreateKnowledgeSplitterReq struct {
	SplitterName  string `json:"splitterName" form:"splitterName" validate:"required"`
	SplitterValue string `json:"splitterValue" form:"splitterValue" validate:"required"`
	CommonCheck
}

type UpdateKnowledgeSplitterReq struct {
	SplitterId    string `json:"splitterId"  form:"splitterId" validate:"required"`
	SplitterName  string `json:"splitterName" form:"splitterName" validate:"required"`
	SplitterValue string `json:"splitterValue" form:"splitterValue" validate:"required"`
	CommonCheck
}

type DeleteKnowledgeSplitterReq struct {
	SplitterId string `json:"splitterId"  form:"splitterId" validate:"required"`
	CommonCheck
}

type GetKnowledgeSplitterReq struct {
	SplitterName string `json:"splitterName"  form:"splitterName"`
	CommonCheck
}
