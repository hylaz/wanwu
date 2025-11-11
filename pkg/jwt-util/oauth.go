package jwt_util

import (
	"context"
	"errors"
	"time"

	oauth2_util "github.com/UnicomAI/wanwu/internal/bff-service/pkg/oauth2-util"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

const (
	// jwt subject
	ACCESS = "access_token"

	AccessTokenTimeout  = int64(60 * 30)          // 30 min
	IDTokenTimeout      = int64(60 * 30)          //30 min
	RefreshTokenTimeout = int64(60 * 60 * 24 * 7) // 7天
)

type AccessTokenClaims struct {
	Scope    []string `json:"scope"`    // access token访问范围
	UserID   string   `json:"userId"`   // 用户ID
	ClientID string   `json:"clientId"` // Client ID
	jwt.StandardClaims
}

func GenerateAccessToken(userID, clientID, issuer string, scopes []string, timeout int64) (string, error) {
	return generateAccessToken(userID, clientID, issuer, scopes, timeout, userSecretKey) //与登录同一个密钥
}

func ParseAccessToken(token string) (*AccessTokenClaims, error) {
	return parseAccessToken(token, userSecretKey)
}

func GenerateRefreshToken(ctx context.Context, userID, clientID string, timeout int64) (string, error) {
	return generateRefreshToken(ctx, userID, clientID, timeout)
}

func generateAccessToken(id, clientID, issuer string, scopes []string, timeout int64, secretKey string) (string, error) {
	if secretKey == "" {
		return "", errors.New("jwt secret key empty")
	}
	nowTime := time.Now().Unix()
	//access token
	claims := &AccessTokenClaims{
		UserID:   id,
		ClientID: clientID,
		Scope:    scopes,
		StandardClaims: jwt.StandardClaims{
			Issuer:    issuer,
			Subject:   ACCESS,            // 用途，目前固定access
			NotBefore: nowTime,           // 生效时间
			ExpiresAt: nowTime + timeout, // 过期时间
		},
	}
	access_token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return access_token, err
}

func parseAccessToken(token, secretKey string) (*AccessTokenClaims, error) {
	if secretKey == "" {
		return nil, errors.New("jwt secret key empty")
	}
	tokenClaims, err := jwt.ParseWithClaims(token, &AccessTokenClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, ErrTokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, ErrTokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, ErrTokenNotValidYet
			} else {
				return nil, ErrTokenInvalid
			}
		}
	}
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*AccessTokenClaims); ok && tokenClaims.Valid {
			return claims, nil
		}
		return nil, ErrTokenInvalid

	} else {
		return nil, ErrTokenInvalid
	}
}

func generateRefreshToken(ctx context.Context, userID, clientID string, timeout int64) (string, error) {
	refreshToken := uuid.NewString()
	//save refresh token to redis
	expied := time.Duration(timeout) * time.Second
	err := oauth2_util.SaveRefreshToken(ctx, refreshToken, expied, oauth2_util.RefreshTokenPayload{
		UserID:   userID,
		ClientID: clientID,
	})
	if err != nil {
		return "", err
	}
	return refreshToken, nil
}
