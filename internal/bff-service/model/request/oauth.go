package request

type AuthRequest struct {
	ResponseType string   `form:"response_type" validate:"required,eq=code"`
	RedirectURI  string   `form:"redirect_uri"`
	ClientID     string   `form:"client_id" validate:"required"`
	Scopes       []string `form:"scope"`
	State        string   `form:"state" validate:"required"`
}

func (t *AuthRequest) Check() error {
	return nil
}

type TokenRequest struct {
	GrantType    string `form:"grant_type" validate:"required,eq=authorization_code"`
	Code         string `form:"code" validate:"required"`
	RedirectURI  string `form:"redirect_uri"`
	ClientID     string `form:"client_id" validate:"required"`
	ClientSecret string `form:"client_secret" validate:"required"`
}

func (t *TokenRequest) Check() error {
	return nil
}

type RefreshRequest struct {
	GrantType    string `json:"grant_type" validate:"required,eq=refresh_token"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RefreshToken string `json:"refresh_token"`
}

func (t *RefreshRequest) Check() error {
	return nil
}

type CreateOauthAppReq struct {
	Name        string `json:"name" validate:"required"`
	Desc        string `json:"desc"`
	RedirectURI string `json:"redirectUri" validate:"required"`
}

func (c *CreateOauthAppReq) Check() error {
	return nil
}

type DeleteOauthAppReq struct {
	ClientID string `json:"clientId" validate:"required"`
}

func (d *DeleteOauthAppReq) Check() error {
	return nil
}

type UpdateOauthAppReq struct {
	ClientID    string `json:"clientId" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Desc        string `json:"desc"`
	RedirectURI string `json:"redirectUri" validate:"required"`
}

func (u *UpdateOauthAppReq) Check() error {
	return nil
}

type UpdateOauthAppStatusReq struct {
	ClientID string `json:"clientId" validate:"required"`
	Status   bool   `json:"status"` // 启停状态
}

func (u *UpdateOauthAppStatusReq) Check() error {
	return nil
}
