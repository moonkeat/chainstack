package handlers_test

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/moonkeat/chainstack/handlers"
	"github.com/moonkeat/chainstack/models"
	"github.com/moonkeat/chainstack/services"
	"github.com/unrolled/render"
)

type fakeHandlerOptions struct {
	userServiceReturnError     bool
	tokenServiceReturnError    bool
	resourceServiceReturnError bool
}

func fakeHandler(opt *fakeHandlerOptions) http.Handler {
	userServiceReturnError := false
	if opt != nil && opt.userServiceReturnError {
		userServiceReturnError = opt.userServiceReturnError
	}

	tokenServiceReturnError := false
	if opt != nil && opt.tokenServiceReturnError {
		tokenServiceReturnError = opt.tokenServiceReturnError
	}

	resourceServiceReturnError := false
	if opt != nil && opt.resourceServiceReturnError {
		resourceServiceReturnError = opt.resourceServiceReturnError
	}

	return handlers.NewHandler(&handlers.Env{
		Render:          render.New(),
		UserService:     &fakeUserService{ReturnError: userServiceReturnError},
		TokenService:    &fakeTokenService{ReturnError: tokenServiceReturnError},
		ResourceService: &fakeResourceService{ReturnError: resourceServiceReturnError},
	})
}

type fakeUserService struct {
	ReturnError bool
}

func (s fakeUserService) CreateUser(email string, password string, isAdmin bool) error {
	return nil
}

func (s fakeUserService) AuthenticateUser(email string, password string) (*models.User, error) {
	if s.ReturnError {
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
		return "", fmt.Errorf("token service error")
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

	if token == "correcttoken" {
		return &models.Token{UserID: 1}, nil
	}

	return nil, services.TokenAuthenticationError{}
}

type fakeResourceService struct {
	ReturnError bool
}

func (s fakeResourceService) CreateResource(userID int) (*models.Resource, error) {
	if s.ReturnError {
		return nil, fmt.Errorf("resource service error")
	}

	return &models.Resource{
		Key:       "resource1",
		CreatedAt: time.Now().Truncate(24 * time.Hour),
	}, nil
}

func (s fakeResourceService) GetResource(userID int, key string) (*models.Resource, error) {
	if s.ReturnError {
		return nil, fmt.Errorf("resource service error")
	}

	if key == "resource1" {
		return &models.Resource{
			Key:       "resource1",
			CreatedAt: time.Now().Truncate(24 * time.Hour),
		}, nil
	}

	return nil, sql.ErrNoRows
}

func (s fakeResourceService) DeleteResource(userID int, key string) error {
	if s.ReturnError {
		return fmt.Errorf("resource service error")
	}

	if key == "resource1" {
		return nil
	}

	return sql.ErrNoRows
}

func (s fakeResourceService) ListResources(userID int) ([]models.Resource, error) {
	if s.ReturnError {
		return nil, fmt.Errorf("resource service error")
	}

	if userID == 1 {
		return []models.Resource{
			{
				Key:       "resource1",
				CreatedAt: time.Now().Truncate(24 * time.Hour),
			},
		}, nil
	}

	return nil, nil
}
