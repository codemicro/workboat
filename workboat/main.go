package main

import (
	"fmt"
	"github.com/codemicro/workboat/workboat/config"
	"github.com/codemicro/workboat/workboat/db"
	"github.com/codemicro/workboat/workboat/endpoints"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func run() error {
	database, err := db.New()
	if err != nil {
		return errors.WithStack(err)
	}

	if err := database.Migrate(); err != nil {
		return errors.Wrap(err, "failed migration")
	}

	e := endpoints.New(database)
	app := e.SetupApp()

	serveAddr := config.HTTP.Host + ":" + strconv.Itoa(config.HTTP.Port)

	go func() {
		shutdownNotifier := make(chan os.Signal, 1)
		signal.Notify(shutdownNotifier, syscall.SIGINT)
		<-shutdownNotifier
		if err := app.Shutdown(); err != nil {
			log.Error().Err(err).Msg("failed to shutdown server on SIGINT")
			log.Fatal().Msg("terminating")
		}
	}()

	log.Info().Msgf("starting server on %s", serveAddr)

	if err := app.Listen(serveAddr); err != nil {
		return errors.Wrap(err, "fiber server run failed")
	}

	log.Info().Msg("shutting down...")
	return nil
}

func main() {
	config.InitLogging()
	if err := run(); err != nil {
		fmt.Printf("%+v\n", err)
		log.Error().Stack().Err(err).Msg("failed to run")
	}
}
