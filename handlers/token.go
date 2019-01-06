package handlers

import (
	"fmt"
	"net/http"
	"strings"

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

	email := strings.TrimSpace(r.Form.Get("client_id"))
	if email == "" {
		return HandlerError{
			StatusCode:  http.StatusBadRequest,
			ActualError: fmt.Errorf("client_id is required"),
		}
	}

	password := strings.TrimSpace(r.Form.Get("client_secret"))
	if password == "" {
		return HandlerError{
			StatusCode:  http.StatusBadRequest,
			ActualError: fmt.Errorf("client_secret is required"),
		}
	}

	authenticated, err := env.UserService.AuthenticateUser(email, password)
	if err != nil {
		return err
	}

	if !authenticated {
		return HandlerError{
			StatusCode:  http.StatusUnauthorized,
			ActualError: fmt.Errorf("invalid credentials"),
		}
	}

	env.Render.JSON(w, http.StatusOK, &responses.Token{})
	return nil
}
