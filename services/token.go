package services

import (
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"
)

type TokenService interface {
	CreateToken(expiresIn time.Duration, scope []string, userID int) (string, error)
	CleanExpiredTokens() error
}

type tokenService struct {
	DB *sqlx.DB
}

func (s tokenService) CreateToken(expiresIn time.Duration, scope []string, userID int) (string, error) {
	token := uuid.NewV4()

	_, err := s.DB.Query("INSERT INTO access_tokens (token, expires, scope, user_id) VALUES ($1, $2, $3, $4)", token.String(), time.Now().Add(expiresIn), strings.Join(scope, ","), userID)
	if err != nil {
		return "", err
	}

	return token.String(), nil
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
