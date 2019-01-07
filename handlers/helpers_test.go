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

func (s fakeTokenService) AuthenticateToken(token string, path string) error {
	if token == "correcttoken" {
		return nil
	}

	return services.TokenAuthenticationError{}
}
