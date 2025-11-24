package request

type CustomTabConfig struct {
	TabLogo  Avatar `json:"tabLogo"`  // 标签页图标
	TabTitle string `json:"tabTitle"` // 标签页标题
	CommonCheck
}

type CustomLoginConfig struct {
	LoginBg          Avatar `json:"loginBg"`          // 登录页背景图
	LoginLogo        Avatar `json:"loginLogo"`        // 登录页图标
	LoginWelcomeText string `json:"loginWelcomeText"` // 登录页欢迎语
	LoginButtonColor string `json:"loginButtonColor"` // 登录按钮颜色
	CommonCheck
}

type CustomHomeConfig struct {
	HomeLogo    Avatar `json:"homeLogo"`    // 平台图标
	HomeName    string `json:"homeName"`    // 平台名称
	HomeBgColor string `json:"homeBgColor"` // 平台背景颜色
	CommonCheck
}
