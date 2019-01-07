package handlers_test

import (
	"fmt"
	"time"

	"github.com/moonkeat/chainstack/models"
	"github.com/moonkeat/chainstack/services"
)

type fakeUserService struct{}

func (s fakeUserService) CreateUser(email string, password string, isAdmin bool) error {
	return nil
}

func (s fakeUserService) AuthenticateUser(email string, password string) (*models.User, error) {
	if email == "internalerror" {
		return nil, fmt.Errorf("internal server error occurred")
	}

	if email == "correct@email.com" && password == "correctpassword" {
		return &models.User{}, nil
	}

	if email == "admin@email.com" && password == "adminpassword" {
		return &models.User{Admin: true}, nil
	}

	return nil, nil
}

type fakeTokenService struct {
	ReturnError bool
}

func (s fakeTokenService) CreateToken(expiresIn time.Duration, scope []string, userID int) (string, error) {
	if s.ReturnError {
		return "", fmt.Errorf("some error here")
	}
	return "fakeToken", nil
}

func (s fakeTokenService) CleanExpiredTokens() error {
	return nil
}

func (s fakeTokenService) AuthenticateToken(token string, path string) (*models.Token, error) {
	if token == "tokenwithinvaliduserid" {
		return &models.Token{UserID: -1}, nil
	}

	if token == "tokenserviceerror" {
		return &models.Token{UserID: 2}, nil
	}

	if token == "correcttoken" {
		return &models.Token{UserID: 1}, nil
	}

	return nil, services.TokenAuthenticationError{}
}

type fakeResourceService struct{}

func (s fakeResourceService) ListResources(userID int) ([]models.Resource, error) {
	if userID == 1 {
		return []models.Resource{
			{
				Key:       "resource1",
				CreatedAt: time.Now().Truncate(24 * time.Hour),
			},
		}, nil
	}

	if userID == 2 {
		return nil, fmt.Errorf("token service error")
	}

	return nil, nil
}
