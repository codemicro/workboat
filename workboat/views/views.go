package views

import (
	"bytes"
	"context"
	"github.com/gofiber/fiber/v2"
	"io"
)

//go:generate ego -v

type RenderFunc func(ctx context.Context, w io.Writer)

func Render(ctx *fiber.Ctx, f RenderFunc) error {
	ctx.Type("html")
	sb := new(bytes.Buffer)
	f(context.Background(), sb)
	return ctx.Send(sb.Bytes())
}
