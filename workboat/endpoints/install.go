package endpoints

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
)

func (e *Endpoints) Install_GetRepositories(ctx *fiber.Ctx) error {
	session, hasSession, err := e.getSession(ctx)
	if err != nil {
		return errors.WithStack(err)
	}

	if !hasSession {
		return fiber.ErrUnauthorized
	}

	repos, err := e.giteaClient.ListUserRepositories(session.GiteaToken)
	if err != nil {
		return errors.WithStack(err)
	}

	alreadyInstalledIDs := make(map[int64]struct{})
	{
		var ids []int64
		for _, repo := range repos {
			ids = append(ids, repo.ID)
		}
		foundRepos, err := e.db.GetRepositoriesWhereExist(ids)
		if err != nil {
			return errors.WithStack(err)
		}
		for _, fr := range foundRepos {
			alreadyInstalledIDs[fr.GiteaRepositoryID] = struct{}{}
		}
	}

	type resStruct struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	}

	var res []*resStruct
	for _, repo := range repos {
		if _, found := alreadyInstalledIDs[repo.ID]; found {
			continue
		}
		res = append(res, &resStruct{
			ID:   repo.ID,
			Name: repo.FullName,
		})
	}

	return ctx.JSON(res)
}

func (e *Endpoints) InstallPage_DoInstall(ctx *fiber.Ctx) error {
	//giteaRepositoryID := ctx.FormValue("id")

	// Make webhook
	// Store webhook token
	// Discover existing configuration files
	return nil
}
