package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func ListUsersHandler(env *Env, w http.ResponseWriter, r *http.Request) error {
	users, err := env.UserService.ListUsers()
	if err != nil {
		return err
	}

	env.Render.JSON(w, http.StatusOK, users)
	return nil
}

func GetUserHandler(env *Env, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	userIDStr := vars["user_id"]
	userID := -1
	parsedUserID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err == nil {
		userID = int(parsedUserID)
	}

	user, err := env.UserService.GetUser(userID)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if err == sql.ErrNoRows {
		return HandlerError{
			StatusCode:  http.StatusNotFound,
			ActualError: fmt.Errorf("user not found"),
		}
	}

	env.Render.JSON(w, http.StatusOK, user)
	return nil
}

func DeleteUserHandler(env *Env, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	userIDStr := vars["user_id"]
	userID := -1
	parsedUserID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err == nil {
		userID = int(parsedUserID)
	}

	err = env.UserService.DeleteUser(userID)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if err == sql.ErrNoRows {
		return HandlerError{
			StatusCode:  http.StatusNotFound,
			ActualError: fmt.Errorf("user not found"),
		}
	}

	env.Render.Data(w, http.StatusNoContent, nil)
	return nil
}
