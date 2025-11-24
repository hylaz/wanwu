package response

type ExplorationAppInfo struct {
	AppBriefInfo
	IsFavorite bool `json:"isFavorite"` // 收藏标签
	User       User `json:"user"`       // 作者信息
}

type User struct {
	UserId   string `json:"userId"`   // 用户ID
	UserName string `json:"userName"` // 用户名称
}
