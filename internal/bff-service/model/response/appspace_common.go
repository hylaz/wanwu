package response

import "github.com/UnicomAI/wanwu/internal/bff-service/model/request"

type AppBriefInfo struct {
	UniqueId    string         `json:"uniqueId"`    // 随机unique id(每次动态生成)
	AppId       string         `json:"appId"`       // 应用id
	AppType     string         `json:"appType"`     // 应用类型
	Avatar      request.Avatar `json:"avatar"`      // 应用图标
	Name        string         `json:"name"`        // 应用名称
	Desc        string         `json:"desc"`        // 应用描述
	CreatedAt   string         `json:"createdAt"`   // 应用创建时间
	UpdatedAt   string         `json:"updatedAt"`   // 应用更新时间(用于历史记录排序)
	PublishType string         `json:"publishType"` // 发布类型(public:公开发布,private:私密发布)
}

type AppUrlInfo struct {
	UrlId               string `json:"urlId"`               // UrlID
	AppId               string `json:"appId"`               // 应用ID
	AppType             string `json:"appType"`             // 应用类型
	Name                string `json:"name"`                // Url名称
	CreatedAt           string `json:"createdAt"`           // 创建时间
	ExpiredAt           string `json:"expiredAt"`           // 过期时间
	Copyright           string `json:"copyright"`           // 知识产权
	CopyrightEnable     bool   `json:"copyrightEnable"`     // 知识产权开关
	PrivacyPolicy       string `json:"privacyPolicy"`       // 隐私政策
	PrivacyPolicyEnable bool   `json:"privacyPolicyEnable"` // 隐私政策开关
	Disclaimer          string `json:"disclaimer"`          // 免责声明
	DisclaimerEnable    bool   `json:"disclaimerEnable"`    // 免责声明开关
	Suffix              string `json:"suffix"`              // 生成Url后缀
	Status              bool   `json:"status"`              // 应用Url开关
	UserId              string `json:"userId"`              // 用户ID
	OrgId               string `json:"orgId"`               // 组织ID
	Description         string `json:"description"`         // 应用描述
}

type AppUrlConfig struct {
	Assistant  *Assistant  `json:"assistant"`  // 基本信息
	AppUrlInfo *AppUrlInfo `json:"appUrlInfo"` // 应用Url信息
}

type VisionConfig struct {
	MaxPicNum int32 `json:"maxPicNum"` // 最大图片数量
	PicNum    int32 `json:"picNum"`    // 视觉配置图片数量
}
