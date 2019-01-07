package handlers

import (
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/rs/zerolog/log"
	"github.com/unrolled/render"

	"github.com/moonkeat/chainstack/responses"
	"github.com/moonkeat/chainstack/services"
)

type Env struct {
	Render          *render.Render
	UserService     services.UserService
	TokenService    services.TokenService
	ResourceService services.ResourceService
}

type Handler struct {
	*Env
	H func(e *Env, w http.ResponseWriter, r *http.Request) error
}

type HandlerError struct {
	StatusCode  int
	ActualError error
}

func (e HandlerError) Error() string {
	return e.ActualError.Error()
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h.H(h.Env, w, r)
	if err != nil {
		var body []byte
		if r.Body != nil {
			body, _ = ioutil.ReadAll(r.Body)
		}

		switch err := err.(type) {
		case HandlerError:
			r.ParseForm()
			log.Debug().
				Err(err).
				Int("status_code", err.StatusCode).
				Bytes("reqbody", body).
				Interface("reqForm", r.Form).
				Str("requrl", r.URL.Path).
				Msg("Handler error.")
			h.Render.JSON(w, err.StatusCode, responses.Error{
				Code:    err.StatusCode,
				Message: err.Error(),
			})
		default:
			log.Error().Err(err).Bytes("reqbody", body).Str("requrl", r.URL.Path).Msg("Internal server error.")
			h.Render.JSON(w, http.StatusInternalServerError, responses.Error{
				Code:    http.StatusInternalServerError,
				Message: "internal server error",
			})
		}
	}
}

func NewHandler(env *Env) http.Handler {
	r := mux.NewRouter()

	// authentication
	r.Handle("/token", Handler{Env: env, H: TokenHandler}).Methods("POST")

	chain := alice.New(AuthMiddleware(env, "resources"))
	r.Handle("/resources", chain.Then(Handler{Env: env, H: ListResourcesHandler})).Methods("GET")
	r.Handle("/resources/{key}", chain.Then(Handler{Env: env, H: GetResourceHandler})).Methods("GET")
	r.Handle("/resources", chain.Then(Handler{Env: env, H: CreateResourceHandler})).Methods("POST")

	return r
}
