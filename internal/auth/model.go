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
	Id         string `json:"id"`
	Email      string `json:"email"`
	FamilyName string `json:"family_name"`
	GivenName  string `json:"given_name"`
	Picture    string `json:"picture"`
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

type GoogleTokenInfo struct {
	Aud       string `json:"audience"`
	UserId    string `json:"user_id"`
	Scope     string `json:"scope"`
	ExpiresIn int    `json:"expires_in"`
}
