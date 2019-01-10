package handlers

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/moonkeat/chainstack/services"
)

func CreateResourceHandler(env *Env, w http.ResponseWriter, r *http.Request) error {
	userID, err := getUserIDFromRequest(r)
	if err != nil {
		return err
	}

	user, err := env.UserService.GetUser(*userID)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if err == sql.ErrNoRows {
		return HandlerError{
			StatusCode:  http.StatusForbidden,
			ActualError: fmt.Errorf("access denied"),
		}
	}

	resources, err := env.ResourceService.ListResources(*userID)
	if err != nil {
		return err
	}

	if user.Quota != nil && *user.Quota < len(resources)+1 && *user.Quota != services.UserQuotaUndefined {
		return HandlerError{
			StatusCode:  http.StatusForbidden,
			ActualError: fmt.Errorf("resource quota exceeded"),
		}
	}

	resource, err := env.ResourceService.CreateResource(*userID)
	if err != nil {
		return err
	}

	env.Render.JSON(w, http.StatusCreated, resource)
	return nil
}

func GetResourceHandler(env *Env, w http.ResponseWriter, r *http.Request) error {
	userID, err := getUserIDFromRequest(r)
	if err != nil {
		return err
	}

	vars := mux.Vars(r)
	key := vars["key"]

	resource, err := env.ResourceService.GetResource(*userID, key)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if err == sql.ErrNoRows {
		return HandlerError{
			StatusCode:  http.StatusForbidden,
			ActualError: fmt.Errorf("access denied"),
		}
	}

	env.Render.JSON(w, http.StatusOK, resource)
	return nil
}
func DeleteResourceHandler(env *Env, w http.ResponseWriter, r *http.Request) error {
	userID, err := getUserIDFromRequest(r)
	if err != nil {
		return err
	}

	vars := mux.Vars(r)
	key := vars["key"]

	err = env.ResourceService.DeleteResource(*userID, key)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if err == sql.ErrNoRows {
		return HandlerError{
			StatusCode:  http.StatusForbidden,
			ActualError: fmt.Errorf("access denied"),
		}
	}

	env.Render.Data(w, http.StatusNoContent, nil)
	return nil
}

func ListResourcesHandler(env *Env, w http.ResponseWriter, r *http.Request) error {
	userID, err := getUserIDFromRequest(r)
	if err != nil {
		return err
	}

	resources, err := env.ResourceService.ListResources(*userID)
	if err != nil {
		return err
	}

	env.Render.JSON(w, http.StatusOK, resources)
	return nil
}
