package main

import (
	"net/http"
	"os"
	"strconv"

	"github.com/moonkeat/chainstack/handlers"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/unrolled/render"
)

func main() {
	debug, _ := strconv.ParseBool(os.Getenv("IS_DEBUG"))
	if !debug {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	addr := os.Getenv("SERVER_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	log.Info().Msgf("Server is running and listen on %s", addr)
	err := http.ListenAndServe(addr, handlers.NewHandler(&handlers.Env{
		Render: render.New(),
	}))
	if err != nil && err != http.ErrServerClosed {
		log.Fatal().Err(err).Msgf("Server could not listen on %s", addr)
	}
	log.Info().Msgf("Server stopped")
}
