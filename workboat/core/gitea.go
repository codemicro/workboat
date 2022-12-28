package core

import (
	"code.gitea.io/sdk/gitea"
	"context"
	"github.com/codemicro/workboat/workboat/config"
	"github.com/codemicro/workboat/workboat/paths"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"time"
)

type GiteaClient struct {
	oauthConfig *oauth2.Config
}

func NewGiteaClient() *GiteaClient {
	return &GiteaClient{
		oauthConfig: &oauth2.Config{
			ClientID:     config.Gitea.OauthClientID,
			ClientSecret: config.Gitea.OauthClientSecret,
			Endpoint: oauth2.Endpoint{
				AuthURL:  config.Gitea.BaseURL + "/login/oauth/authorize",
				TokenURL: config.Gitea.BaseURL + "/login/oauth/access_token",
			},
			RedirectURL: paths.Make(paths.AuthOauthInbound),
		},
	}
}

func (gc *GiteaClient) newAPIUserClient(token *oauth2.Token) (*gitea.Client, error) {
	return gitea.NewClient(config.Gitea.BaseURL,
		gitea.SetHTTPClient(gc.oauthConfig.Client(context.Background(), token)),
	)
}

func (gc *GiteaClient) newAPISystemClient() (*gitea.Client, error) {
	return gitea.NewClient(config.Gitea.BaseURL,
		gitea.SetToken(config.Gitea.AccessToken),
	)
}

func (gc *GiteaClient) OauthAuthCodeURL(state string) string {
	return gc.oauthConfig.AuthCodeURL(state)
}

func (gc *GiteaClient) OauthExchange(code string) (*oauth2.Token, error) {
	exchangeCtx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	token, err := gc.oauthConfig.Exchange(exchangeCtx, code)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return token, nil
}
