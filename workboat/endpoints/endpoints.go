package endpoints

import (
	"github.com/codemicro/workboat/workboat/config"
	"github.com/codemicro/workboat/workboat/core"
	"github.com/codemicro/workboat/workboat/db"
	"github.com/codemicro/workboat/workboat/paths"
	"github.com/codemicro/workboat/workboat/util"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
	"github.com/pkg/errors"
	"sync"
	"time"
)

type Endpoints struct {
	db                *db.DB
	loginStateManager *loginStateManager
	giteaClient       *core.GiteaClient
	sessions          map[string]*session
	sessionLock       *sync.Mutex
}

func New(dbi *db.DB, giteaClient *core.GiteaClient) *Endpoints {
	e := new(Endpoints)

	e.db = dbi
	e.loginStateManager = newLoginStateManager()
	e.giteaClient = giteaClient
	e.sessions = make(map[string]*session)
	e.sessionLock = new(sync.Mutex)

	return e
}

func (e *Endpoints) SetupApp() *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler:          util.JSONErrorHandler,
		DisableStartupMessage: !config.Debug.Enabled,

		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	})

	app.Use("/api", func(ctx *fiber.Ctx) error {
		ctx.Set(fiber.HeaderAccessControlAllowOrigin, "*")
		return ctx.Next()
	})

	app.Get(paths.AuthOauthInbound, e.AuthOauthInbound)
	app.Get(paths.APIAuthNewLogin, e.AuthOauthGetURL)
	app.Get(paths.InstallGetRepository, e.Install_GetRepositories)

	app.Use("/", func(ctx *fiber.Ctx) error {
		return proxy.Do(ctx, config.HTTP.FrontendURL+ctx.Path())
	})

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
