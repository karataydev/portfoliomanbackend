package user

import (
	"fmt"

	"github.com/karataydev/portfoliomanbackend/internal/auth"
)

type Service struct {
	repo         *Repository
	tokenService *auth.TokenService
}

func NewService(repo *Repository, tokenService *auth.TokenService) *Service {
	return &Service{repo: repo, tokenService: tokenService}
}

func (s *Service) GetByEmail(email string) (*User, error) {
	return s.repo.GetByEmail(email)
}

func (s *Service) Get(id int64) (*User, error) {
	return s.repo.Get(id)
}

func (s *Service) SignUp(googleToken string) (*SignInUpResponse, error) {
	googleClaims, err := s.tokenService.ValidateGoogleToken(googleToken)
	if err != nil {
		return nil, fmt.Errorf("invalid Google token: %w", err)
	}

	userExists, err := s.userExists(googleClaims.Email)
	if err != nil {
		return nil, err
	}
	if userExists {
		return s.existingUserTokenCreate(googleClaims.Email)
	}

	newUser := &User{
		FirstName:         googleClaims.GivenName,
		LastName:          googleClaims.FamilyName,
		Email:             googleClaims.Email,
		GoogleId:          googleClaims.Id,
		ProfilePictureUrl: googleClaims.Picture,
	}

	createdUser, err := s.repo.Save(newUser)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Generate token
	token, err := s.tokenService.CreateToken(createdUser.Id, createdUser.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to create token: %w", err)
	}

	return &SignInUpResponse{
		User:        createdUser,
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   int(s.tokenService.Duration.Seconds()),
	}, nil
}

func (s *Service) SignIn(googleToken string) (*SignInUpResponse, error) {
	googleClaims, err := s.tokenService.ValidateGoogleToken(googleToken)
	if err != nil {
		return nil, fmt.Errorf("invalid Google token: %w", err)
	}
	return s.existingUserTokenCreate(googleClaims.Email)

}

func (s *Service) existingUserTokenCreate(email string) (*SignInUpResponse, error) {
	user, err := s.GetByEmail(email)
	if err != nil {
		return nil, err
	}

	token, err := s.tokenService.CreateToken(user.Id, user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to create token: %w", err)
	}

	return &SignInUpResponse{
		User:        user,
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   int(s.tokenService.Duration.Seconds()),
		UserExisted: true,
	}, nil
}

func (s *Service) userExists(email string) (bool, error) {
	_, err := s.GetByEmail(email)
	if err != nil {
		if err == UserNotFoundErr {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
