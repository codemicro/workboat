package core

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"strings"
)

func WebhookHandler(ctx *fiber.Ctx) error {
	// TODO: Signature validation?

	if !strings.EqualFold(ctx.Get(fiber.HeaderContentType), "application/json") {
		ctx.Status(fiber.StatusBadRequest)
		return nil
	}

	event := ctx.Get("X-Gitea-Event")

	if event == "" {
		ctx.Status(fiber.StatusBadRequest)
		return ctx.SendString("missing X-Gitea-Event")
	}

	o := make(map[string]any)
	if err := json.Unmarshal(ctx.Body(), &o); err != nil {
		ctx.Status(fiber.StatusBadRequest)
		return ctx.SendString(err.Error())
	}

	go dispatch(strings.Clone(event), o)

	ctx.Status(fiber.StatusNoContent)
	return nil
}
