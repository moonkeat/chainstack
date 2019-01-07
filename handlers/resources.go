package handlers

import (
	"fmt"
	"net/http"
)

func ListResourcesHandler(env *Env, w http.ResponseWriter, r *http.Request) error {
	userID := -1
	userIDFromCtx, ok := r.Context().Value("user_id").(int)
	if ok {
		userID = userIDFromCtx
	}
	if userID == -1 {
		return fmt.Errorf("failed to parse user_id from context")
	}

	resources, err := env.ResourceService.ListResources(userID)
	if err != nil {
		return err
	}

	env.Render.JSON(w, http.StatusOK, resources)
	return nil
}
