package service

import (
	"fmt"
	"os"
	"strconv"
	"time"

	iam_service "github.com/UnicomAI/wanwu/api/proto/iam-service"
	"github.com/UnicomAI/wanwu/internal/bff-service/model/request"
	"github.com/UnicomAI/wanwu/internal/bff-service/model/response"
	oauth2_util "github.com/UnicomAI/wanwu/internal/bff-service/pkg/oauth2-util"
	jwt_util "github.com/UnicomAI/wanwu/pkg/jwt-util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var ISSUER = os.Getenv("WANWU_CALLBACK_LLM_BASE_URL") +
	os.Getenv("WANWU_EXTERNAL_SCHEME") +
	"://" + os.Getenv("WANWU_EXTERNAL_ENDPOINT")

func Authorize(ctx *gin.Context, req *request.AuthRequest, userID string) (string, string, error) {
	oauthApp, err := iam.GetOauthApp(ctx, &iam_service.GetOauthAppReq{
		ClientId: req.ClientID,
	})
	if err != nil {
		return "", "", fmt.Errorf("%v get auth info err:%v", req.ClientID, err)
	}
	if !oauthApp.Status {
		return "", "", fmt.Errorf("client id: %v has been disabled", oauthApp.ClientId)
	}
	if oauthApp.ClientId != req.ClientID || (req.RedirectURI != "" && req.RedirectURI != oauthApp.RedirectUri) {
		return "", "", fmt.Errorf("client id: %v or redirecturi: %v missmatch", req.ClientID, req.RedirectURI)
	}
	//code save to redis
	code := uuid.NewString()
	oauth2_util.SaveCode(ctx, code, oauth2_util.CodePayload{
		ClientID: req.ClientID,
		UserID:   userID,
	})
	return oauthApp.RedirectUri, code, nil
}

func Token(ctx *gin.Context, req *request.TokenRequest) (*response.TokenResponse, error) {
	codePayload, err := oauth2_util.ValidateCode(ctx, req.Code, req.ClientID)
	if err != nil {
		return nil, fmt.Errorf("validate code timeout err", err)
	}
	oauthApp, err := iam.GetOauthApp(ctx, &iam_service.GetOauthAppReq{
		ClientId: req.ClientID,
	})
	if err != nil {
		return nil, fmt.Errorf("%v get auth info err:%v", req.ClientID, err)
	}
	err = validateCode(req.ClientID, req.ClientSecret, req.RedirectURI, codePayload, oauthApp)
	if err != nil {
		return nil, err
	}
	user, err := iam.GetUserInfo(ctx, &iam_service.GetUserInfoReq{
		UserId: codePayload.UserID,
		OrgId:  "",
	})
	if err != nil {
		return nil, err
	}
	//access token
	var scopes []string = []string{""} //预留scope处理
	accessToken, err := jwt_util.GenerateAccessToken(user.UserId, req.ClientID, ISSUER, scopes, jwt_util.AccessTokenTimeout)
	if err != nil {
		return nil, err
	}

	//id token
	idToken, err := jwt_util.GenerateIDToken(user.UserId, user.UserName, req.ClientID, ISSUER, jwt_util.IDTokenTimeout)
	if err != nil {
		return nil, err
	}
	//refresh token
	refreshToken, err := jwt_util.GenerateRefreshToken(ctx, user.UserId, req.ClientID, jwt_util.RefreshTokenTimeout)
	if err != nil {
		return nil, err
	}
	return &response.TokenResponse{
		AccessToken:  accessToken,
		ExpiresIn:    jwt_util.AccessTokenTimeout,
		TokenType:    "Bearer",
		IDToken:      idToken,
		RefreshToken: refreshToken,
		Scope:        scopes,
	}, nil
}

func Refresh(ctx *gin.Context, req *request.RefreshRequest) (*response.RefreshTokenResponse, error) {
	refreshPayload, err := oauth2_util.ValidateRefreshToken(ctx, req.RefreshToken, req.ClientID)
	if err != nil {
		return nil, err
	}
	oauthApp, err := iam.GetOauthApp(ctx, &iam_service.GetOauthAppReq{
		ClientId: req.ClientID,
	})
	if err != nil {
		return nil, err
	}
	if !oauthApp.Status {
		return nil, fmt.Errorf("client id: %v has been disabled", oauthApp.ClientId)
	}
	if req.ClientSecret != oauthApp.ClientSecret {
		return nil, fmt.Errorf("clinetId:%v or clientSecret missmatch", req.ClientID)
	}
	scopes := []string{} //scopes处理预留
	//new access token
	accessToken, err := jwt_util.GenerateAccessToken(refreshPayload.UserID, req.ClientID, ISSUER, scopes, jwt_util.AccessTokenTimeout)
	if err != nil {
		return nil, err
	}
	//new refresh token
	refreshToken, err := jwt_util.GenerateRefreshToken(ctx, refreshPayload.UserID, refreshPayload.ClientID, jwt_util.RefreshTokenTimeout)
	if err != nil {
		return nil, err
	}
	return &response.RefreshTokenResponse{
		AccessToken:  accessToken,
		ExpiresAt:    strconv.Itoa(int(time.Now().Add(time.Duration(jwt_util.UserTokenTimeout) * time.Second).UnixMilli())),
		RefreshToken: refreshToken,
	}, nil
}

func OauthConfig(ctx *gin.Context) (*response.OauthConfig, error) {

	return &response.OauthConfig{
		Issuer:           ISSUER,
		AuthEndpoint:     ISSUER + "/user/api/v1" + "/oauth/code/authorize",
		TokenEndpoint:    ISSUER + "/user/api/openapi/v1" + "/oauth/code/token",
		JwksUri:          ISSUER + "/user/api/openapi/v1" + "/oauth/jwks",
		UserInfoEndpoint: ISSUER + "/user/api/openapi/v1" + "/oauth/userinfo",
		ResponseTypes:    []string{"code"},
		IDtokenSignAlg:   []string{"RS256"},
		SubjectTypes:     []string{"public"},
	}, nil
}

func JWKS(ctx *gin.Context) (*response.JWKS, error) {

	return &response.JWKS{Keys: []jwt_util.JWK{*jwt_util.JWKInstance}}, nil
}

func OAuthGetUserInfo(ctx *gin.Context, userID string) (*response.OAuthGetUserInfo, error) {
	user, err := iam.GetUserInfo(ctx, &iam_service.GetUserInfoReq{
		UserId: userID,
		OrgId:  "",
	})
	if err != nil {
		return nil, err
	}
	avataUri := cacheUserAvatar(ctx, user.AvatarPath)
	return &response.OAuthGetUserInfo{
		UserID:    user.UserId,
		Username:  user.UserName,
		Email:     user.Email,
		Nickname:  user.NickName,
		Phone:     user.Phone,
		Gender:    user.Gender,
		AvatarUri: ISSUER + "/user/api" + avataUri.Path,
		Remark:    user.Remark,
		Company:   user.Company,
	}, nil
}

func validateCode(clientID, clientSecret, redirectUri string, codePayload oauth2_util.CodePayload, appInfo *iam_service.OauthApp) error {

	if !appInfo.Status {
		return fmt.Errorf("client id: %v has been disabled", codePayload.ClientID)
	}
	if codePayload.ClientID != clientID { //两次传的不一样
		return fmt.Errorf("client_id mismatch: expected %v, got %v", codePayload.ClientID, clientID)
	}

	if appInfo.ClientSecret != clientSecret {
		return fmt.Errorf("client_secret error for client_id: %v", codePayload.ClientID)
	}

	if redirectUri != "" && redirectUri != appInfo.RedirectUri {
		return fmt.Errorf("redirect_uri err: got %v", redirectUri)
	}
	return nil
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
