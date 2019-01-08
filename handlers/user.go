package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/moonkeat/chainstack/models"
)

func CreateUserHandler(env *Env, w http.ResponseWriter, r *http.Request) error {
	if r.Body == nil {
		return HandlerError{
			StatusCode:  400,
			ActualError: fmt.Errorf("request body is nil"),
		}
	}

	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		return HandlerError{
			StatusCode:  400,
			ActualError: fmt.Errorf("failed to parse request body as json, err: %s", err),
		}
	}
	defer r.Body.Close()

	userData, err := env.UserService.CreateUser(user.Email, user.Password, user.Admin, user.Quota)
	if err != nil {
		switch err.(type) {
		case models.UserValidationError:
			return HandlerError{
				StatusCode:  http.StatusBadRequest,
				ActualError: err,
			}
		default:
			return err
		}
	}

	env.Render.JSON(w, http.StatusCreated, userData)
	return nil
}

func ListUsersHandler(env *Env, w http.ResponseWriter, r *http.Request) error {
	users, err := env.UserService.ListUsers()
	if err != nil {
		return err
	}

	env.Render.JSON(w, http.StatusOK, users)
	return nil
}

func GetUserHandler(env *Env, w http.ResponseWriter, r *http.Request) error {
	userID := getUserIDFromRequest(r)
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

func UpdateUserQuotaHandler(env *Env, w http.ResponseWriter, r *http.Request) error {
	if r.Body == nil {
		return HandlerError{
			StatusCode:  400,
			ActualError: fmt.Errorf("request body is nil"),
		}
	}

	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		return HandlerError{
			StatusCode:  400,
			ActualError: fmt.Errorf("failed to parse request body as json, err: %s", err),
		}
	}
	defer r.Body.Close()

	userID := getUserIDFromRequest(r)
	userData, err := env.UserService.UpdateUserQuota(userID, user.Quota)
	if err != nil {
		return err
	}

	env.Render.JSON(w, http.StatusOK, userData)
	return nil
}

func DeleteUserHandler(env *Env, w http.ResponseWriter, r *http.Request) error {
	userID := getUserIDFromRequest(r)
	err := env.UserService.DeleteUser(userID)
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

func getUserIDFromRequest(r *http.Request) int {
	vars := mux.Vars(r)
	userIDStr := vars["user_id"]
	userID := -1
	parsedUserID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err == nil {
		userID = int(parsedUserID)
	}

	return userID
}
