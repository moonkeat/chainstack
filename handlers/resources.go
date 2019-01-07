package handlers

import (
	"net/http"
)

func ListResourcesHandler(env *Env, w http.ResponseWriter, r *http.Request) error {
	env.Render.JSON(w, http.StatusOK, []string{})
	return nil
}
