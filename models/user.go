package models

import (
	"fmt"

	"github.com/asaskevich/govalidator"
)

type UserValidationError struct {
	Field  string
	Reason string
}

func (e UserValidationError) Error() string {
	return fmt.Sprintf("invalid %s: %s", e.Field, e.Reason)
}

type User struct {
	ID       int    `db:"id" json:"id"`
	Email    string `db:"email" json:"email"`
	Password string `db:"password" json:"password,omitempty"`
	Admin    bool   `db:"admin" json:"admin"`
	Quota    *int   `db:"quota" json:"quota,omitempty"`
}

func ValidateUser(email string, password string) error {
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

	return nil
}

// TODO: test admin create user , non admin create user, test update quota < resources
