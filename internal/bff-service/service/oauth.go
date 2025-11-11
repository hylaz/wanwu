package service

import (
	iam_service "github.com/UnicomAI/wanwu/api/proto/iam-service"
	"github.com/UnicomAI/wanwu/internal/bff-service/model/request"
	"github.com/UnicomAI/wanwu/internal/bff-service/model/response"
	"github.com/gin-gonic/gin"
)

func Authorize(ctx *gin.Context, req *request.AuthRequest) (string, string, error) {
	return "", "", nil
	// //TODO: 调用获取client
	// credential, err := auth.GetCredential(ctx, &auth_service.GetCredentialReq{
	// 	ClientId: req.ClientID,
	// })
	// if err != nil {
	// 	return "", "", fmt.Errorf("%v get auth info err:%v", req.ClientID, err)
	// }
	// if !credential.Status {
	// 	return "", "", fmt.Errorf("client id: %v has been disabled", credential.ClientId)
	// }
	// if credential.ClientId != req.ClientID || (req.RedirectURI != "" && req.RedirectURI != credential.Callback) {
	// 	return "", "", fmt.Errorf("client id: %v or redirecturi: %v missmatch", req.ClientID, req.RedirectURI)
	// }
	// _, token, _ := jwt_util.GenerateCode(req.UserID, config.OAUTH_CODE, credential.ClientId, credential.ClientSecret, jwt_util.CodeTimeout)
	// return credential.Callback, token, nil
}

func Token(ctx *gin.Context, req *request.TokenRequest) (*response.TokenResponse, error) {
	//TODO:
	return nil, nil

	// userId, err := validateCode(ctx, req.Code, req.ClientID, req.ClientSecret, req.RedirectURI)
	// if err != nil {
	// 	return nil, err
	// }
	// user, err := iam.GetUserInfo(ctx, &iam_service.GetUserInfoReq{
	// 	UserId: userId,
	// 	OrgId:  "",
	// })
	// if err != nil {
	// 	return nil, err
	// }
	// _, token, refreshToken, _ := jwt_util.GenerateOauthToken(user.UserId, req.ClientID, jwt_util.UserTokenTimeout, jwt_util.RefreshTokenTimeout)
	// return &response.TokenResponse{
	// 	AccessToken:  token,
	// 	ExpiresIn:    jwt_util.UserTokenTimeout,
	// 	TokenType:    "Bearer",
	// 	RefreshToken: refreshToken,
	// 	Scope:        "",
	// }, nil
}

func Credential(ctx *gin.Context, req *request.CredentialRequest) (*response.TokenResponse, error) {
	//TODO:
	return nil, nil

	// credential, err := auth.GetCredential(ctx, &auth_service.GetCredentialReq{
	// 	ClientId: req.ClientID,
	// })
	// if err != nil {
	// 	return nil, fmt.Errorf("get %v credential info err:%v", req.ClientID, err)
	// }
	// if !credential.Status {
	// 	return nil, fmt.Errorf("client id: %v has been disabled", credential.ClientId)
	// }
	// if credential.ClientId != req.ClientID || req.ClientSecret != credential.ClientSecret {
	// 	return nil, fmt.Errorf("clinet id: %v or client secret missmatch", req.ClientID)
	// }
	// _, token, refreshToken, _ := jwt_util.GenerateOauthToken("", req.ClientID, jwt_util.UserTokenTimeout, jwt_util.RefreshTokenTimeout)
	// return &response.TokenResponse{
	// 	AccessToken:  token,
	// 	ExpiresIn:    jwt_util.UserTokenTimeout,
	// 	TokenType:    "Bearer",
	// 	RefreshToken: refreshToken,
	// 	Scope:        "",
	// }, nil
}

func Refresh(ctx *gin.Context, req *request.RefreshRequest) (*response.RefreshTokenResponse, error) {
	//TODO:
	return nil, nil

	// oldClaims, err := jwt_util.ParseToken(req.RefreshToken)
	// if err != nil {
	// 	return nil, err
	// }
	// credential, err := auth.GetCredential(ctx, &auth_service.GetCredentialReq{
	// 	ClientId: req.ClientID,
	// })
	// if err != nil {
	// 	return nil, fmt.Errorf("get %v credential info err:%v", req.ClientID, err)
	// }
	// if !credential.Status {
	// 	return nil, fmt.Errorf("client id: %v has been disabled", credential.ClientId)
	// }
	// if credential.ClientId != req.ClientID || req.ClientSecret != credential.ClientSecret || oldClaims.ID != req.ClientID {
	// 	return nil, fmt.Errorf("clinetId:%v or clientSecret missmatch", req.ClientID)
	// }
	// _, token, refreshToken, err := jwt_util.GenerateOauthToken(
	// 	oldClaims.ID,
	// 	req.ClientID,
	// 	jwt_util.UserTokenTimeout,
	// 	jwt_util.RefreshTokenTimeout,
	// )
	// return &response.RefreshTokenResponse{
	// 	AccessToken:  token,
	// 	ExpiresAt:    strconv.Itoa(int(time.Now().Add(time.Duration(jwt_util.UserTokenTimeout) * time.Second).UnixMilli())),
	// 	RefreshToken: refreshToken,
	// }, nil
}

func OauthConfig(ctx *gin.Context) (*response.OauthConfig, error) {
	return nil, nil
	// icon.URL, _ = url.JoinPath(os.Getenv("WANWU_EXTERNAL_SCHEME")+"://"+os.Getenv("WANWU_EXTERNAL_ENDPOINT"),
	// 		os.Getenv("WANWU_WORKFLOW_DEFAULT_ICON"))
}

func validateCode(ctx *gin.Context, code, clientID, clientSecret, redirectURI string) (string, error) {

	//TODO:
	return "", nil
	// credential, err := auth.GetCredential(ctx, &auth_service.GetCredentialReq{
	// 	ClientId: clientID,
	// })
	// if err != nil {
	// 	return "", fmt.Errorf("validate code get credential info err:%v", err)
	// }
	// if !credential.Status {
	// 	return "", fmt.Errorf("client id: %v has been disabled", credential.ClientId)
	// }
	// if credential.ClientId != clientID || credential.ClientSecret != clientSecret || (redirectURI != "" && redirectURI != credential.Callback) {
	// 	return "", fmt.Errorf("client id:%v client secret or redirecturi mismatch", credential.ClientId)
	// }
	// claims, err := jwt_util.ParseCode(code, clientSecret)
	// if err != nil {
	// 	return "", err
	// }
	// if claims.Issuer != clientID || claims.Subject != config.OAUTH_CODE {
	// 	return "", fmt.Errorf("issuer or subject error")
	// }
	// return claims.ID, nil
}

func CreateOauthApp(ctx *gin.Context, userId string, req *request.CreateOauthAppReq) error {
	_, err := iam.CreateOauthApp(ctx, &iam_service.CreateOauthAppReq{
		UserId:      userId,
		Name:        req.Name,
		Desc:        req.Desc,
		RedirectUri: req.RedirectURI,
	})
	if err != nil {
		return err
	}
	return nil
}

func DeleteOauthApp(ctx *gin.Context, req *request.DeleteOauthAppReq) error {
	_, err := iam.DeleteOauthApp(ctx, &iam_service.DeleteOauthAppReq{
		ClientId: req.ClientID,
	})
	if err != nil {
		return err
	}
	return nil
}

func UpdateOauthApp(ctx *gin.Context, req *request.UpdateOauthAppReq) error {
	_, err := iam.UpdateOauthApp(ctx, &iam_service.UpdateOauthAppReq{
		ClientId:    req.ClientID,
		Name:        req.Name,
		Desc:        req.Desc,
		RedirectUri: req.RedirectURI,
	})
	if err != nil {
		return err
	}
	return nil
}

func GetOauthAppList(ctx *gin.Context, userId string) ([]*response.OauthAppInfo, error) {
	resp, err := iam.GetOauthAppList(ctx, &iam_service.GetOauthAppListReq{
		UserId: userId,
	})
	if err != nil {
		return nil, err
	}
	var retList []*response.OauthAppInfo
	for _, app := range resp.Apps {
		retList = append(retList, &response.OauthAppInfo{
			ClientID:     app.ClientId,
			Name:         app.Name,
			Desc:         app.Desc,
			ClientSecret: app.ClientSecret,
			RedirectURI:  app.RedirectUri,
			Status:       app.Status,
		})
	}
	return retList, nil
}

func UpdateOauthAppStatus(ctx *gin.Context, req *request.UpdateOauthAppStatusReq) error {
	_, err := iam.UpdateOauthAppStatus(ctx, &iam_service.UpdateOauthAppStatusReq{
		ClientId: req.ClientID,
		Status:   req.Status,
	})
	if err != nil {
		return err
	}
	return nil
}
