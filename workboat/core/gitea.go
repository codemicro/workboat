package core

import (
	"code.gitea.io/sdk/gitea"
	"github.com/codemicro/workboat/workboat/config"
)

func newAPISystemClient() (*gitea.Client, error) {
	return gitea.NewClient(config.Gitea.BaseURL,
		gitea.SetToken(config.Gitea.AccessToken),
	)
}
