package request

type AppUrlIdRequest struct {
	UrlId string `json:"urlId" form:"urlId" validate:"required"` // UrlID
}

func (a *AppUrlIdRequest) Check() error { return nil }

type AppUrlConfig struct {
	Name                string `json:"name" validate:"required"` // 名称
	ExpiredAt           string `json:"expiredAt"`                // 过期时间
	Copyright           string `json:"copyright"`                // 版权
	CopyrightEnable     bool   `json:"copyrightEnable"`          // 版权开关
	PrivacyPolicy       string `json:"privacyPolicy"`            // 隐私协议
	PrivacyPolicyEnable bool   `json:"privacyPolicyEnable"`      // 隐私协议开关
	Disclaimer          string `json:"disclaimer"`               // 免责声明
	DisclaimerEnable    bool   `json:"disclaimerEnable"`         // 免责声明开关
	Description         string `json:"description"`
}

func (cfg AppUrlConfig) Check() error {
	return nil
}

type AppUrlCreateRequest struct {
	AppId   string `json:"appId" validate:"required"`   // 应用id
	AppType string `json:"appType" validate:"required"` // 应用类型
	AppUrlConfig
}

type AppUrlUpdateRequest struct {
	UrlId string `json:"urlId" validate:"required"` // UrlID
	AppUrlConfig
}

type AppUrlListRequest struct {
	AppId   string `json:"appId" form:"appId" validate:"required"`     // 应用id
	AppType string `json:"appType" form:"appType" validate:"required"` // 应用类型
}

func (a *AppUrlListRequest) Check() error { return nil }

type AppUrlStatusRequest struct {
	UrlId  string `json:"urlId" validate:"required"` // UrlID
	Status bool   `json:"status"`                    // 启停状态
}

func (a *AppUrlStatusRequest) Check() error { return nil }
