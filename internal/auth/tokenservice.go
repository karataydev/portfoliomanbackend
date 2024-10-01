package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

type TokenService struct {
	rsaKeys         *RSAKeys
	Duration        time.Duration
	googleValidator *GoogleValidator
}

func NewTokenService(rsaKeys *RSAKeys, duration time.Duration, googleValidator *GoogleValidator) *TokenService {
	return &TokenService{
		rsaKeys:         rsaKeys,
		Duration:        duration,
		googleValidator: googleValidator,
	}
}

var InvalidTokenErr error = errors.New("Invalid token!")
var CreateTokenErr error = errors.New("Create token error!")

func (v *TokenService) ValidateGoogleToken(idToken string) (*GoogleTokenClaims, error) {
	payload, err := v.googleValidator.ValidateToken(idToken)
	if err != nil {
		return nil, err
	}

	if payload.Id == "" || payload.Email == "" || payload.GivenName == "" {
		return nil, InvalidTokenErr
	}
	return payload, nil
}

func (v *TokenService) ValidateToken(token string) (*TokenClaims, error) {
	tok, err := jwt.Parse(token, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, InvalidTokenErr
		}

		return v.rsaKeys.PublicKey, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok || !tok.Valid {
		return nil, InvalidTokenErr
	}
	floatId, ok := claims["uid"].(float64)
	if !ok {
		return nil, InvalidTokenErr
	}
	id := int64(floatId)
	if id == 0 {
		return nil, InvalidTokenErr
	}
	email, ok := claims["email"].(string)
	if !ok {
		return nil, InvalidTokenErr
	}

	return &TokenClaims{UserId: id, Email: email}, nil
}

func (v *TokenService) CreateToken(id int64, email string) (string, error) {

	now := time.Now().UTC()
	// create claims
	claims := make(jwt.MapClaims)
	claims["uid"] = id                         // User id
	claims["email"] = email                    // User email
	claims["exp"] = now.Add(v.Duration).Unix() // The expiration time after which the token must be disregarded.
	claims["iat"] = now.Unix()                 // The time at which the token was issued.
	claims["nbf"] = now.Unix()                 // The time before which the token must be disregarded.

	// generate token with private key
	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(v.rsaKeys.PrivateKey)
	if err != nil {
		return "", CreateTokenErr
	}

	return token, nil
}
