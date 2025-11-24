package response

import (
	oauth2_util "github.com/UnicomAI/wanwu/internal/bff-service/pkg/oauth2-util"
)

type OAuthTokenResponse struct {
	AccessToken  string   `json:"access_token"`  // 访问令牌
	ExpiresIn    int64    `json:"expires_in"`    // token过期时间(毫秒时间戳)
	IDToken      string   `json:"id_token"`      // ID令牌
	TokenType    string   `json:"token_type"`    // 令牌类型(bearer)
	RefreshToken string   `json:"refresh_token"` // 刷新令牌
	Scope        []string `json:"scope"`         // 权限范围
}

type OAuthRefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`  // 访问令牌
	RefreshToken string `json:"refresh_token"` // 刷新令牌(可选)
	ExpiresAt    string `json:"expires_at"`    // token过期时间(毫秒时间戳)
}

type OAuthConfig struct {
	Issuer           string   `json:"issuer"`                                // auth的base url
	AuthEndpoint     string   `json:"authorization_endpoint"`                // 获取授权码接口
	TokenEndpoint    string   `json:"token_endpoint"`                        // 获取token接口
	JwksUri          string   `json:"jwks_uri"`                              // 获取jwt公钥
	UserInfoEndpoint string   `json:"userinfo_endpoint"`                     //获取用户信息接口
	ResponseTypes    []string `json:"response_types_supported"`              // 授权模式，默认code
	IDtokenSignAlg   []string `json:"id_token_signing_alg_values_supported"` //jwt签名算法
	SubjectTypes     []string `json:"subject_types_supported"`               // 用户标识类型，即 ID Token中的sub是如何生成的。
}

type OAuthAppInfo struct {
	ClientID     string `json:"clientId"`     // 客户端ID
	Name         string `json:"name"`         // 应用名称
	Desc         string `json:"desc"`         // 应用描述
	ClientSecret string `json:"clientSecret"` // 客户端密钥
	RedirectURI  string `json:"redirectUri"`  // oauth重定向地址
	Status       bool   `json:"status"`       // oauth应用开关
}

type OAuthJWKS struct {
	Keys []oauth2_util.JWK `json:"keys"`
}

type OAuthGetUserInfo struct {
	UserID    string `json:"userId"`
	Username  string `json:"username"`
	Nickname  string `json:"nickname"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	Gender    string `json:"gender"`
	Remark    string `json:"remark"`
	Company   string `json:"company"`
	AvatarUri string `json:"avatar"`
}
