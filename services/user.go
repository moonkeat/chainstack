package services

import (
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"

	"github.com/moonkeat/chainstack/models"
)

type UserService interface {
	CreateUser(email string, password string, isAdmin bool) (bool, error)
	AuthenticateUser(email string, password string) (bool, error)
}

type userService struct {
	DB *sqlx.DB
}

func (s userService) CreateUser(email string, password string, isAdmin bool) (bool, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return false, err
	}

	_, err = s.DB.Query("INSERT INTO users (email, password, admin) VALUES ($1, $2, $3)", email, passwordHash, isAdmin)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s userService) AuthenticateUser(email string, password string) (bool, error) {
	user := models.User{}
	err := s.DB.Select(&user, "SELECT password FROM users WHERE email = $1", email)
	if err != nil {
		return false, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return false, nil
	}

	return true, nil
}

func NewUserService(db *sqlx.DB) UserService {
	return &userService{
		DB: db,
	}
}
