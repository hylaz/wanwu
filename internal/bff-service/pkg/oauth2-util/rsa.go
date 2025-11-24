package oauth2_util

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
)

func toRsaPublicKey(rsaPublicKey []byte) (*rsa.PublicKey, error) {
	// 1. PEM decode
	block, _ := pem.Decode(rsaPublicKey)
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, errors.New("rsa public key pem decode err")
	}

	// 2. Parse certificate
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("rsa public key parse certificate err: %v", err)
	}

	// 3. Extract public key
	pubKey, ok := pubInterface.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("rsa public key extract error")
	}
	return pubKey, nil
}

func toRsaPrivateKey(rsaPrivateKey []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(rsaPrivateKey)
	if block == nil {
		return nil, errors.New("rsa private key pem decode err")
	}

	privInterface, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("rsa private key parse certificate err: %v", err)
	}

	privKey, ok := privInterface.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("rsa private key extract error")
	}
	return privKey, nil
}
