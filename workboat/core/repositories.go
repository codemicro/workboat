package core

import (
	"code.gitea.io/sdk/gitea"
	"github.com/codemicro/workboat/workboat/config"
	"github.com/codemicro/workboat/workboat/paths"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

func (gc *GiteaClient) ListUserRepositories(token *oauth2.Token) ([]*gitea.Repository, error) {
	client, err := gc.newAPIUserClient(token)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	repos, _, err := client.ListMyRepos(gitea.ListReposOptions{})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return repos, nil
}

func (gc *GiteaClient) GetRepository(repositoryID int64) (*gitea.Repository, error) {
	client, err := gc.newAPISystemClient()
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
func (gc *GiteaClient) CreateRepositoryHook(user, repo string) (*gitea.Hook, string, error) {
	client, err := gc.newAPISystemClient()
	if err != nil {
		return nil, "", errors.WithStack(err)
	}

	var secret = "secret" // TODO: Don't hardcode this

	hook, _, err := client.CreateRepoHook(user, repo, gitea.CreateHookOption{
		Type: "gitea",
		Config: map[string]string{
			"url":          config.HTTP.InternalURL + paths.WebhookInbound,
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
