package request

type KnowledgeOrgSelectReq struct {
	KnowledgeId string `json:"knowledgeId" form:"knowledgeId" validate:"required"` //知识库id
	CommonCheck
}

type KnowledgeUserSelectReq struct {
	KnowledgeId string `json:"knowledgeId" form:"knowledgeId" validate:"required"` //知识库id
	CommonCheck
}

type KnowledgeUserNoPermitSelectReq struct {
	KnowledgeId string `json:"knowledgeId" form:"knowledgeId" validate:"required"` //知识库id
	OrgId       string `json:"orgId" form:"orgId" validate:"required"`             //选择组织id
	Transfer    bool   `json:"transfer" form:"transfer"`                           //是否是转让列表
	CommonCheck
}

type KnowledgeUserAddReq struct {
	KnowledgeId       string               `json:"knowledgeId" validate:"required"`       //知识库id
	PermissionType    int                  `json:"permissionType"`                        // 权限类型:0: 查看权限; 10: 编辑权限; 20: 授权权限,数值不连续的原因防止后续有中间权限，目前逻辑 授权权限>编辑权限>查看权限
	KnowledgeUserList []*KnowledgeUserInfo `json:"knowledgeUserList" validate:"required"` // 知识库用户信息
	CommonCheck
}

type KnowledgeUserEditReq struct {
	KnowledgeId   string             `json:"knowledgeId" validate:"required"`   //知识库id
	KnowledgeUser *KnowledgeUserInfo `json:"knowledgeUser" validate:"required"` // 知识库用户信息
	CommonCheck
}

type KnowledgeUserDeleteReq struct {
	KnowledgeId  string `json:"knowledgeId" validate:"required"`  //知识库id
	PermissionId string `json:"permissionId" validate:"required"` // 知识库用户信息权限信息
	CommonCheck
}

type KnowledgeTransferUserAdminReq struct {
	KnowledgeId   string             `json:"knowledgeId" validate:"required"`   //知识库id
	PermissionId  string             `json:"permissionId" form:"permissionId"`  //权限id编辑时传入
	KnowledgeUser *KnowledgeUserInfo `json:"knowledgeUser" validate:"required"` // 知识库用户信息
	CommonCheck
}

type KnowledgeUserInfo struct {
	UserId         string `json:"userId" validate:"required"`
	OrgId          string `json:"orgId" validate:"required"`
	PermissionId   string `json:"permissionId" form:"permissionId"` //权限id编辑时传入
	PermissionType int    `json:"permissionType"`                   // 权限类型: -1 删除此用户权限；0: 查看权限; 10: 编辑权限; 20: 授权权限,数值不连续的原因防止后续有中间权限，目前逻辑 授权权限>编辑权限>查看权限
}
