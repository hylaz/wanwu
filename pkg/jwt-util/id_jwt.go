package jwt_util

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/UnicomAI/wanwu/pkg/log"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

var (
	rsaPublicKey  []byte
	rsaPrivateKey []byte
	Kid           string
	JWKInstance   *JWK
)

func InitIDTokenJWT(publicKeyPath, privateKeyPath string) {
	if rsaPrivateKey != nil && rsaPublicKey != nil {
		log.Panicf("jwt already init")
	}
	if publicKeyPath == "" {
		log.Panicf("oauth rsa public key empty")
	}
	if privateKeyPath == "" {
		log.Panicf("oauth rsa private key empty")
	}
	var err error
	rsaPrivateKey, err = os.ReadFile("/app/tmp/private.pem")
	if err != nil {
		log.Panicf("oauth private key read error")
	}
	rsaPublicKey, err = os.ReadFile("/app/tmp/public.pem")
	if err != nil {
		log.Panicf("oauth public key read error")
	}
	JWKInstance, Kid = GetJWKSAndKid()
}

type IDTokenClaims struct {
	UserID   string `json:"userId"`   // 用户ID
	UserName string `json:"userName"` // 用户名称
	jwt.StandardClaims
}

type JWK struct {
	Kty string `json:"kty"`
	Use string `json:"use"`
	Kid string `json:"kid"`
	Alg string `json:"alg"`
	N   string `json:"n"`
	E   string `json:"e"`
}

func GetJWKSAndKid() (*JWK, string) {
	rsaPub := toRsaPublicKey()
	return getJWKSAndKid(rsaPub)
}

func GenerateIDToken(userID, userName, clientID, issuer string, timeout int64) (string, error) {
	return generateIDToken(userID, userName, clientID, issuer, timeout, string(rsaPrivateKey))
}

func generateIDToken(id, userName, clientID, issuer string, timeout int64, privateKey string) (string, error) {
	if privateKey == "" {
		return "", errors.New("jwt secret key empty")
	}
	rsaPrivateKey := toRsaPrivateKey()
	nowTime := time.Now().Unix()
	claims := &IDTokenClaims{
		UserID:   id,
		UserName: userName,
		StandardClaims: jwt.StandardClaims{
			Issuer:    issuer, //oidc root path
			Subject:   id,     // 用途，目前固定user
			Audience:  clientID,
			NotBefore: nowTime,           // 生效时间
			ExpiresAt: nowTime + timeout, // 过期时间
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = Kid
	tokenString, err := token.SignedString(rsaPrivateKey)
	if err != nil {
		return "", err
	}
	return tokenString, err
}

func getJWKSAndKid(pub *rsa.PublicKey) (*JWK, string) {

	kid := uuid.New().String()

	// 2️⃣ 导出公钥参数 n 和 e
	n := base64.RawURLEncoding.EncodeToString(pub.N.Bytes())
	eBytes := big.NewInt(int64(pub.E)).Bytes()
	e := base64.RawURLEncoding.EncodeToString(eBytes)

	jwk := JWK{
		Kty: "RSA",
		Use: "sig",
		Kid: kid,
		Alg: "RS256",
		N:   n,
		E:   e,
	}

	return &jwk, kid

}

func toRsaPublicKey() *rsa.PublicKey {
	// 1. PEM decode
	block, _ := pem.Decode(rsaPublicKey)
	fmt.Println("pem解码", block.Type)
	if block == nil || block.Type != "PUBLIC KEY" {
		log.Panicf("failed to decode PEM block containing certificate")
	}

	// 2. Parse certificate
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		log.Panicf(err.Error())
	}

	// 3. Extract public key
	pubKey, ok := pubInterface.(*rsa.PublicKey)
	if !ok {
		log.Panicf("not RSA public key")
	}
	return pubKey
}

func toRsaPrivateKey() *rsa.PrivateKey {
	block, _ := pem.Decode(rsaPrivateKey)
	if block == nil {
		log.Panicf("failed to decode PEM")
	}

	privKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		log.Panicf(err.Error())
	}
	return privKey
}
