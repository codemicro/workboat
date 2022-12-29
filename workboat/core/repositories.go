package core

import (
	"code.gitea.io/sdk/gitea"
	"github.com/codemicro/workboat/workboat/config"
	"github.com/pkg/errors"
)

func GetRepository(repositoryID int64) (*gitea.Repository, error) {
	client, err := newAPISystemClient()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	repo, _, err := client.GetRepoByID(repositoryID)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return repo, nil
}

// CreateRepositoryHook creates a webhook on the repository specified by user/repo and returns the new hook object and
// the secret used to create it.
func CreateRepositoryHook(user, repo string) (*gitea.Hook, string, error) {
	client, err := newAPISystemClient()
	if err != nil {
		return nil, "", errors.WithStack(err)
	}

	var secret = "secret" // TODO: Don't hardcode this

	hook, _, err := client.CreateRepoHook(user, repo, gitea.CreateHookOption{
		Type: "gitea",
		Config: map[string]string{
			"url":          config.HTTP.InternalURL + "/webhook/inbound",
			"content_type": "json",
			"http_method":  "POST",
			"secret":       secret,
		},
		Events: []string{"push"},
		Active: true,
	})

	if err != nil {
		return nil, "", errors.WithStack(err)
	}

	return hook, secret, nil
}
