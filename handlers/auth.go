package handlers

import (
	"net/http"
	"strings"

	"github.com/moonkeat/chainstack/responses"
)

func AuthMiddleware(env *Env, path string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			accessToken := strings.TrimSpace(strings.Replace(r.Header.Get("Authorization"), "Bearer ", "", -1))
			err := env.TokenService.AuthenticateToken(accessToken, path)
			if err != nil {
				env.Render.JSON(w, http.StatusUnauthorized, responses.Error{
					Code:    http.StatusUnauthorized,
					Message: "access denied",
				})
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
