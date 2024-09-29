package auth

import (
	"context"

	"google.golang.org/api/idtoken"
)

type GoogleValidator struct {
	validator *idtoken.Validator
	clientID  string
}

func NewGoogleValidator(clientID string) (*GoogleValidator, error) {
	ctx := context.Background()
	validator, err := idtoken.NewValidator(ctx)
	if err != nil {
		return nil, err
	}
	return &GoogleValidator{
		validator: validator,
		clientID:  clientID,
	}, nil
}

func (v *GoogleValidator) ValidateToken(idToken string) (*idtoken.Payload, error) {
	return v.validator.Validate(context.Background(), idToken, v.clientID)
}
