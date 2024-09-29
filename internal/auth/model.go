package auth

import (
	"crypto/rsa"

	"github.com/golang-jwt/jwt"
)

type TokenClaims struct {
	UserId int64
	Email  string
}

type GoogleTokenClaims struct {
	Id         string
	Email      string
	FamilyName string
	GivenName  string
	Picture    string
}

type RSAKeys struct {
	PublicKey  *rsa.PublicKey
	PrivateKey *rsa.PrivateKey
}

func NewRSAKeysFromByte(prvKeyByte []byte, pubKeyByte []byte) (*RSAKeys, error) {
	prvKey, err := jwt.ParseRSAPrivateKeyFromPEM(prvKeyByte)
	if err != nil {
		return nil, err
	}
	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(pubKeyByte)
	if err != nil {
		return nil, err
	}
	return &RSAKeys{
		PublicKey:  pubKey,
		PrivateKey: prvKey,
	}, nil
}
