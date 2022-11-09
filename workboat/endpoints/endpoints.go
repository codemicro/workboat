package endpoints

import (
	"github.com/codemicro/workboat/workboat/config"
	"github.com/codemicro/workboat/workboat/core"
	"github.com/codemicro/workboat/workboat/db"
	"github.com/codemicro/workboat/workboat/paths"
	"github.com/codemicro/workboat/workboat/util"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"time"
)

type Endpoints struct {
	db                *db.DB
	loginStateManager *loginStateManager
	giteaClient       *core.GiteaClient
}

func New(dbi *db.DB, giteaClient *core.GiteaClient) *Endpoints {
	e := new(Endpoints)

	e.db = dbi
	e.loginStateManager = newLoginStateManager()
	e.giteaClient = giteaClient

	return e
}

func (e *Endpoints) SetupApp() *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler:          util.JSONErrorHandler,
		DisableStartupMessage: !config.Debug.Enabled,

		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	})

	app.Get(paths.Index, e.Index)
	app.Get(paths.AuthLogin, e.AuthLogin)

	app.Get(paths.AuthOauthOutbound, e.AuthOauthOutbound)
	app.Get(paths.AuthOauthInbound, e.AuthOauthInbound)

	app.Get(paths.Install, e.InstallPage)
	app.Get(paths.InstallSelectRepository, e.InstallPage_SelectRepository)
	app.Post(paths.InstallDoInstall, e.InstallPage_DoInstall)

	return app
}

func (e *Endpoints) loginThenReturn(ctx *fiber.Ctx) error {
	// This weirdness is to make sure we make a copy of the string as the string isn't valid outside this handler.
	urlBytes := []byte(ctx.OriginalURL())
	url := string(urlBytes)

	state, err := e.loginStateManager.New(url)
	if err != nil {
		return errors.WithStack(err)
	}

	nextPath := paths.Make(paths.AuthOauthOutbound + "?state=" + state)

	return ctx.Redirect(nextPath)
}
