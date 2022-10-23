package core

import (
	"code.gitea.io/sdk/gitea"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

func (gc *GiteaClient) ListUserRepositories(token *oauth2.Token) ([]*gitea.Repository, error) {
	client, err := gc.newAPIClient(token)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	repos, _, err := client.ListMyRepos(gitea.ListReposOptions{})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return repos, nil
}
