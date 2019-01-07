package services

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"

	"github.com/moonkeat/chainstack/models"
)

type UserService interface {
	GetUser(userID int) (*models.User, error)
	CreateUser(email string, password string, isAdmin bool) error
	AuthenticateUser(email string, password string) (*models.User, error)
}

type UserValidationError struct {
	Field  string
	Reason string
}

func (e UserValidationError) Error() string {
	return fmt.Sprintf("invalid %s: %s", e.Field, e.Reason)
}

type userService struct {
	DB *sqlx.DB
}

func (s userService) CreateUser(email string, password string, isAdmin bool) error {
	email = strings.TrimSpace(email)
	if !govalidator.IsEmail(email) {
		return UserValidationError{
			Field:  "email",
			Reason: fmt.Sprintf("'%s' is not a valid email", email),
		}
	}

	if len(password) < 8 {
		return UserValidationError{
			Field:  "password",
			Reason: fmt.Sprintf("password should be at least 8 characters"),
		}
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = s.DB.Query("INSERT INTO users (email, password, admin) VALUES (lower($1), $2, $3)", email, passwordHash, isAdmin)
	if err != nil {
		return err
	}

	return nil
}

func (s userService) GetUser(userID int) (*models.User, error) {
	user := models.User{}
	err := s.DB.Get(&user, "SELECT id, email, password, admin, COALESCE(quota, -1) as quota FROM users WHERE id = $1", userID)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return &user, nil
}

func (s userService) AuthenticateUser(email string, password string) (*models.User, error) {
	user := models.User{}
	err := s.DB.Get(&user, "SELECT id, email, password, admin FROM users WHERE lower(email) = lower($1)", email)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, nil
	}

	return &user, nil
}

func NewUserService(db *sqlx.DB) UserService {
	return &userService{
		DB: db,
	}
}
