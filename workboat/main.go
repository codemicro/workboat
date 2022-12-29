package main

import (
	"fmt"
	"github.com/codemicro/workboat/workboat/config"
	"github.com/rs/zerolog/log"
)

func run() error {
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
