package main

import (
	"net/http"
	"os"
	"strconv"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/unrolled/render"

	"github.com/moonkeat/chainstack/handlers"
	"github.com/moonkeat/chainstack/services"
)

func main() {
	debug, _ := strconv.ParseBool(os.Getenv("IS_DEBUG"))
	if !debug {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	dbConnString := os.Getenv("DB_CONNSTRING")
	db, err := sqlx.Connect("postgres", dbConnString)
	if err != nil {
		log.Fatal().Err(err).Msgf("Failed to connect to postgres, connString: '%s'", dbConnString)
	}

	addr := os.Getenv("SERVER_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	log.Info().Msgf("Server is running and listen on %s", addr)
	err = http.ListenAndServe(addr, handlers.NewHandler(&handlers.Env{
		Render:      render.New(),
		UserService: services.NewUserService(db),
	}))
	if err != nil && err != http.ErrServerClosed {
		log.Fatal().Err(err).Msgf("Server could not listen on %s", addr)
	}
	log.Info().Msgf("Server stopped")
}
