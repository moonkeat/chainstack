package handlers

import (
	"fmt"
	"net/http"

	"github.com/moonkeat/chainstack/responses"
)

const GrantTypeClientCredentials = "client_credentials"

func TokenHandler(env *Env, w http.ResponseWriter, r *http.Request) error {
	r.ParseForm()

	grantType := r.Form.Get("grant_type")
	if grantType != GrantTypeClientCredentials {
		return HandlerError{
			StatusCode:  http.StatusBadRequest,
			ActualError: fmt.Errorf("invalid grant type: '%s'", grantType),
		}
	}

	env.Render.JSON(w, http.StatusOK, &responses.Token{})
	return nil
}
