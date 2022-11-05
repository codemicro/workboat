package views

import (
	"context"
	"io"
	"strings"
)

//go:generate ego -v

type RenderFunc func(ctx context.Context, w io.Writer)

func Render(f RenderFunc) string {
	sb := new(strings.Builder)
	f(context.Background(), sb)
	return sb.String()
}
