package handlers

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func CreateResourceHandler(env *Env, w http.ResponseWriter, r *http.Request) error {
	userID, err := getUserIDFromContext(r)
	if err != nil {
		return err
	}

	resource, err := env.ResourceService.CreateResource(userID)
	if err != nil {
		return err
	}

	env.Render.JSON(w, http.StatusCreated, resource)
	return nil
}

func GetResourceHandler(env *Env, w http.ResponseWriter, r *http.Request) error {
	userID, err := getUserIDFromContext(r)
	if err != nil {
		return err
	}

	vars := mux.Vars(r)
	key := vars["key"]

	resource, err := env.ResourceService.GetResource(userID, key)
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

func ListResourcesHandler(env *Env, w http.ResponseWriter, r *http.Request) error {
	userID, err := getUserIDFromContext(r)
	if err != nil {
		return err
	}

	resources, err := env.ResourceService.ListResources(userID)
	if err != nil {
		return err
	}

	env.Render.JSON(w, http.StatusOK, resources)
	return nil
}

func getUserIDFromContext(r *http.Request) (int, error) {
	userID := -1
	userIDFromCtx, ok := r.Context().Value("user_id").(int)
	if ok {
		userID = userIDFromCtx
	}
	if userID == -1 {
		return userID, fmt.Errorf("failed to parse user_id from context")
	}

	return userID, nil
}
