package jwt_util

import (
	"errors"
	"time"

	"github.com/UnicomAI/wanwu/pkg/log"
	"github.com/dgrijalva/jwt-go"
)

const (
	// jwt subject
	OAUTH = "oauth"

	AccessTokenTimeout  = int64(60 * 30)     // 30 min
	IDTokenTimeout      = int64(60 * 30)     //30 min
	RefreshTokenTimeout = int64(60 * 60 * 7) // 7天
)

type AuthClaims struct {
	UserID     string `json:"userId"` // 用户ID
	BufferTime int64  `json:"bufferTime"`
	jwt.StandardClaims
}

func InitAuthJWT(key string) {
	if userSecretKey != "" {
		log.Panicf("jwt already init")
	}
	if key == "" {
		log.Panicf("jwt secret key empty")
	}
	userSecretKey = key
}

func GenerateOauthToken(userID string, timeout int64) (*CustomClaims, string, error) {
	return generateOauthToken(userID, timeout, userSecretKey)
}

func generateOauthToken(id string, timeout int64, secretKey string) (*CustomClaims, string, error) {
	if secretKey == "" {
		return nil, "", errors.New("jwt secret key empty")
	}
	nowTime := time.Now().Unix()
	claims := &CustomClaims{
		UserID:     id,
		BufferTime: nowTime + BufferTime, // 缓冲时间，当nowTime大于等于BufferTime and nowTime小于ExpiresAt是获得新的token
		StandardClaims: jwt.StandardClaims{
			Issuer:    "wanwu",
			Subject:   USER,              // 用途，目前固定user
			NotBefore: nowTime,           // 生效时间
			ExpiresAt: nowTime + timeout, // 过期时间
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secretKey))
	if err != nil {
		return nil, "", err
	}
	return claims, token, err
}
