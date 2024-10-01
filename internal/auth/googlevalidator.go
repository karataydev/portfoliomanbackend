package auth

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2/log"
)

type GoogleValidator struct {
	clientID   string
	httpClient *http.Client
}

func NewGoogleValidator(clientID string) *GoogleValidator {
	return &GoogleValidator{
		clientID:   clientID,
		httpClient: &http.Client{},
	}
}

func (v *GoogleValidator) verifyAccessToken(accessToken string) error {
	url := fmt.Sprintf("https://www.googleapis.com/oauth2/v1/tokeninfo?access_token=%s", accessToken)
	resp, err := v.httpClient.Get(url)
	if err != nil {
		return fmt.Errorf("failed to verify token: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("invalid token")
	}

	var tokenInfo GoogleTokenInfo
	if err := json.NewDecoder(resp.Body).Decode(&tokenInfo); err != nil {
		return fmt.Errorf("failed to decode token info: %v", err)
	}

	log.Info(tokenInfo.Aud)
	log.Info(v.clientID)
	if tokenInfo.Aud != v.clientID {
		return fmt.Errorf("token is not intended for this application")
	}

	return nil
}

func (v *GoogleValidator) ValidateToken(accessToken string) (*GoogleTokenClaims, error) {
	if err := v.verifyAccessToken(accessToken); err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := v.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var userInfo GoogleTokenClaims
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return &userInfo, nil
}
