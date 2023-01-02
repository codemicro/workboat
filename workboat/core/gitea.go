package core

import (
	"code.gitea.io/sdk/gitea"
	"github.com/codemicro/workboat/workboat/config"
	"github.com/pkg/errors"
)

func newAPISystemClient() (*gitea.Client, error) {
	return gitea.NewClient(config.Gitea.BaseURL,
		gitea.SetToken(config.Gitea.AccessToken),
	)
}

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

// GetFileFromRepository fetches a file from a repository. ref may be an empty string to use the default branch.
func GetFileFromRepository(repoOwner, repoName, fpath, ref string) (*gitea.ContentsResponse, error) {
	client, err := newAPISystemClient()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	fcont, _, err := client.GetContents(repoOwner, repoName, ref, fpath)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return fcont, nil
}

func ReportRepoStatus(repoOwner, repoName, commitSha string, opt *gitea.CreateStatusOption) error {
	opt.Context = "Workboat"
	client, err := newAPISystemClient()
	if err != nil {
		return errors.WithStack(err)
	}
	_, _, err = client.CreateStatus(repoOwner, repoName, commitSha, *opt)

	return errors.WithStack(err)
}
