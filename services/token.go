package services

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/moonkeat/chainstack/models"
	uuid "github.com/satori/go.uuid"
)

type TokenService interface {
	CreateToken(expiresIn time.Duration, scope []string, userID int) (string, error)
	CleanExpiredTokens() error
	AuthenticateToken(token string, path string) (*models.Token, error)
}

type TokenAuthenticationError struct{}

func (e TokenAuthenticationError) Error() string {
	return fmt.Sprint("token invalid")
}

type tokenService struct {
	DB *sqlx.DB
}

func (s tokenService) CreateToken(expiresIn time.Duration, scope []string, userID int) (string, error) {
	token := uuid.NewV4()
	_, err := s.DB.Exec("INSERT INTO access_tokens (token, expires, scope, user_id) VALUES ($1, $2, $3, $4)", token.String(), time.Now().Add(expiresIn), strings.Join(scope, " "), userID)
	if err != nil {
		return "", err
	}

	return token.String(), nil
}

func (s tokenService) AuthenticateToken(tokenString string, path string) (*models.Token, error) {
	token := models.Token{}
	err := s.DB.Get(&token, "SELECT token, expires, scope, user_id FROM access_tokens WHERE token = $1 AND expires > NOW()", tokenString)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if !strings.Contains(token.Scope, path) {
		return nil, TokenAuthenticationError{}
	}

	return &token, nil
}

func (s tokenService) CleanExpiredTokens() error {
	_, err := s.DB.Exec("DELETE FROM access_tokens WHERE expires < NOW()")
	if err != nil {
		return err
	}

	return nil
}

func NewTokenService(db *sqlx.DB) TokenService {
	return &tokenService{
		DB: db,
	}
}
