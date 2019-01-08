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
	userServiceReturnError                   bool
	userServiceQuota                         *int
	tokenServiceReturnError                  bool
	resourceServiceCreateReturnError         bool
	resourceServiceGetResourceError          bool
	resourceServiceDeleteResourceReturnError bool
	resourceServiceListResourcesReturnError  bool
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

	resourceServiceCreateReturnError := false
	if opt != nil && opt.resourceServiceCreateReturnError {
		resourceServiceCreateReturnError = opt.resourceServiceCreateReturnError
	}

	resourceServiceGetResourceError := false
	if opt != nil && opt.resourceServiceGetResourceError {
		resourceServiceGetResourceError = opt.resourceServiceGetResourceError
	}

	resourceServiceDeleteResourceReturnError := false
	if opt != nil && opt.resourceServiceDeleteResourceReturnError {
		resourceServiceDeleteResourceReturnError = opt.resourceServiceDeleteResourceReturnError
	}

	resourceServiceListResourcesReturnError := false
	if opt != nil && opt.resourceServiceListResourcesReturnError {
		resourceServiceListResourcesReturnError = opt.resourceServiceListResourcesReturnError
	}

	userServiceQuota := -1
	if opt != nil && opt.userServiceQuota != nil {
		userServiceQuota = *opt.userServiceQuota
	}

	return handlers.NewHandler(&handlers.Env{
		Render: render.New(),
		UserService: &fakeUserService{
			ReturnError: userServiceReturnError,
			UserQuota:   userServiceQuota,
		},
		TokenService: &fakeTokenService{
			ReturnError: tokenServiceReturnError,
		},
		ResourceService: &fakeResourceService{
			CreateReturnError:         resourceServiceCreateReturnError,
			GetResourceError:          resourceServiceGetResourceError,
			DeleteResourceReturnError: resourceServiceDeleteResourceReturnError,
			ListResourcesReturnError:  resourceServiceListResourcesReturnError,
		},
	})
}

type fakeUserService struct {
	ReturnError bool
	UserQuota   int
}

func (s fakeUserService) CreateUser(email string, password string, isAdmin bool, quota *int) (*models.User, error) {
	return nil, nil
}

func (s fakeUserService) GetUser(userID int) (*models.User, error) {
	if s.ReturnError {
		return nil, fmt.Errorf("user service error")
	}

	if userID != 1 {
		return nil, sql.ErrNoRows
	}

	return &models.User{
		ID:    1,
		Email: "test@test.com",
		Admin: false,
		Quota: &s.UserQuota,
	}, nil
}

func (s fakeUserService) DeleteUser(userID int) error {
	if s.ReturnError {
		return fmt.Errorf("user service error")
	}

	if userID == 1 {
		return nil
	}

	return sql.ErrNoRows
}

func (s fakeUserService) ListUsers() ([]models.User, error) {
	if s.ReturnError {
		return nil, fmt.Errorf("user service error")
	}

	return []models.User{
		{
			ID:    1,
			Email: "test@test.com",
			Admin: false,
			Quota: &s.UserQuota,
		},
	}, nil
}

func (s fakeUserService) AuthenticateUser(email string, password string) (*models.User, error) {
	if s.ReturnError {
		return nil, fmt.Errorf("user service error")
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
	CreateReturnError         bool
	GetResourceError          bool
	DeleteResourceReturnError bool
	ListResourcesReturnError  bool
}

func (s fakeResourceService) CreateResource(userID int) (*models.Resource, error) {
	if s.CreateReturnError {
		return nil, fmt.Errorf("resource service error")
	}

	return &models.Resource{
		Key:       "resource1",
		CreatedAt: time.Now().Truncate(24 * time.Hour),
	}, nil
}

func (s fakeResourceService) GetResource(userID int, key string) (*models.Resource, error) {
	if s.GetResourceError {
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
	if s.DeleteResourceReturnError {
		return fmt.Errorf("resource service error")
	}

	if key == "resource1" {
		return nil
	}

	return sql.ErrNoRows
}

func (s fakeResourceService) ListResources(userID int) ([]models.Resource, error) {
	if s.ListResourcesReturnError {
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
