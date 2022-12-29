package main

import (
	"fmt"
	"github.com/codemicro/workboat/workboat/config"
	"github.com/codemicro/workboat/workboat/core"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

func run() error {
	httpApp := fiber.New()

	httpApp.Post("/webhook/inbound", core.WebhookHandler)

	if err := httpApp.Listen(fmt.Sprintf("%s:%d", config.HTTP.Host, config.HTTP.Port)); err != nil {
		return errors.WithStack(err)
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
