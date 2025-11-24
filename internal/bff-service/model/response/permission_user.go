package response

import "github.com/UnicomAI/wanwu/internal/bff-service/model/request"

type UserID struct {
	UserID string `json:"userId"`
}

type UserInfo struct {
	UserID    string         `json:"userId"`
	Username  string         `json:"username"`
	Nickname  string         `json:"nickname"`
	Phone     string         `json:"phone"`
	Email     string         `json:"email"`
	Gender    string         `json:"gender"`
	Remark    string         `json:"remark"`
	Company   string         `json:"company"`
	CreatedAt string         `json:"createdAt"`
	Creator   IDName         `json:"creator"` // 创建人
	Status    bool           `json:"status"`
	Language  Language       `json:"language"`
	Orgs      []OrgRole      `json:"orgs"` // 用户的组织角色列表
	Avatar    request.Avatar `json:"avatar"`
}

type OrgRole struct {
	Org   IDName   `json:"org"`   // 组织
	Roles []IDName `json:"roles"` // 角色列表
}
