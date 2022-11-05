package endpoints

import (
	"fmt"
	"github.com/codemicro/workboat/workboat/views"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
)

func (e *Endpoints) InstallPage(ctx *fiber.Ctx) error {
	_, hasSession, err := e.getSession(ctx)
	if err != nil {
		return errors.WithStack(err)
	}

	if !hasSession {
		return e.loginThenReturn(ctx)
	}

	return views.Render(ctx, views.InstallPage)
}

func (e *Endpoints) InstallPage_SelectRepository(ctx *fiber.Ctx) error {
	session, hasSession, err := e.getSession(ctx)
	if err != nil {
		return errors.WithStack(err)
	}

	if !hasSession {
		return e.loginThenReturn(ctx)
	}

	repos, err := e.giteaClient.ListUserRepositories(session.GiteaToken)
	if err != nil {
		return errors.WithStack(err)
	}

	var res string
	for i, repo := range repos {
		res += fmt.Sprintf("%d %s %s<br>", i, repo.Name, repo.HTMLURL)
	}

	return ctx.SendString(res)
}
