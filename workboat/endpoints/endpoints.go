package endpoints

import (
	"github.com/codemicro/workboat/workboat/config"
	"github.com/codemicro/workboat/workboat/db"
	"github.com/codemicro/workboat/workboat/paths"
	"github.com/codemicro/workboat/workboat/util"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/oauth2"
	"time"
)

type Endpoints struct {
	db    *db.DB
	login struct {
		stateManager *loginStateManager
		oauthConfig  *oauth2.Config
	}
}

func New(dbi *db.DB) *Endpoints {
	e := new(Endpoints)

	e.db = dbi
	e.login.stateManager = newLoginStateManager()
	e.login.oauthConfig = &oauth2.Config{
		ClientID:     config.Gitea.OauthClientID,
		ClientSecret: config.Gitea.OauthClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  config.Gitea.BaseURL + "/login/oauth/authorize",
			TokenURL: config.Gitea.BaseURL + "/login/oauth/access_token",
		},
		RedirectURL: paths.Make(paths.AuthOauthInbound),
	}

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

	return app
}
