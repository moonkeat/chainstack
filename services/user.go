package services

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"

	"github.com/moonkeat/chainstack/models"
)

const UserQuotaUndefined = -1

type UserService interface {
	CreateUser(email string, password string, isAdmin bool, quota *int) (*models.User, error)
	GetUser(userID int) (*models.User, error)
	UpdateUserQuota(userID int, quota *int) (*models.User, error)
	DeleteUser(userID int) error
	ListUsers() ([]models.User, error)
	AuthenticateUser(email string, password string) (*models.User, error)
}

type userService struct {
	DB *sqlx.DB
}

func (s userService) CreateUser(email string, password string, isAdmin bool, quota *int) (*models.User, error) {
	email = strings.TrimSpace(email)
	err := models.ValidateUser(email, password)
	if err != nil {
		return nil, err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	_, err = s.DB.Exec("INSERT INTO users (email, password, admin, quota) VALUES (lower($1), $2, $3, $4)", email, passwordHash, isAdmin, quota)
	if err != nil {
		if strings.Contains(err.Error(), "users_unique_lower_email_idx") {
			return nil, models.UserValidationError{
				Field:  "email",
				Reason: fmt.Sprintf("user with email already exists"),
			}
		}
		return nil, err
	}

	user, err := s.AuthenticateUser(email, password)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s userService) GetUser(userID int) (*models.User, error) {
	user := models.User{}
	err := s.DB.Get(&user, fmt.Sprintf("SELECT id, email, admin, COALESCE(quota, %d) as quota FROM users WHERE id = $1", UserQuotaUndefined), userID)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s userService) UpdateUserQuota(userID int, quota *int) (*models.User, error) {
	_, err := s.DB.Exec("UPDATE users SET quota = $1 WHERE id = $2", quota, userID)
	if err != nil {
		return nil, err
	}

	user, err := s.GetUser(userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s userService) DeleteUser(userID int) error {
	user := models.User{}
	err := s.DB.Get(&user, "SELECT id FROM users WHERE id = $1", userID)
	if err != nil {
		return err
	}

	_, err = s.DB.Query("DELETE FROM users WHERE id = $1", user.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s userService) ListUsers() ([]models.User, error) {
	users := []models.User{}
	err := s.DB.Select(&users, fmt.Sprintf("SELECT id, email, admin, COALESCE(quota, %d) as quota FROM users", UserQuotaUndefined))
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return users, nil
}

func (s userService) AuthenticateUser(email string, password string) (*models.User, error) {
	user := models.User{}
	err := s.DB.Get(&user, fmt.Sprintf("SELECT id, email, password, admin, COALESCE(quota, %d) as quota FROM users WHERE lower(email) = lower($1)", UserQuotaUndefined), email)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, nil
	}

	user.Password = ""

	return &user, nil
}

func NewUserService(db *sqlx.DB) UserService {
	return &userService{
		DB: db,
	}
}
