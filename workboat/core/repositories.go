package core

import (
	"code.gitea.io/sdk/gitea"
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
