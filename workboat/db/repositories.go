package db

import (
	"database/sql"
	"github.com/codemicro/workboat/workboat/db/models"
	"github.com/pkg/errors"
	"strings"
)

func (db *DB) InsertRepository(repo *models.Repository) error {
	ctx, cancel := db.newContext()
	defer cancel()

	_, err := db.bun.NewInsert().Model(repo).Exec(ctx)
	return errors.WithStack(err)
}

func (db *DB) GetRepositoryByGiteaID(id int64) (*models.Repository, error) {
	ctx, cancel := db.newContext()
	defer cancel()

	res := new(models.Repository)
	if err := db.bun.NewSelect().Model(res).Where("gitea_repository_id = ?", id).Scan(ctx); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, errors.WithStack(err)
	}
	return res, nil
}

func (db *DB) GetRepositoriesWhereExist(ids []int64) ([]*models.Repository, error) {
	ctx, cancel := db.newContext()
	defer cancel()

	var (
		queryParts   []string
		typelessArgs []any
	)
	for _, v := range ids {
		queryParts = append(queryParts, "gitea_repository_id = ?")
		typelessArgs = append(typelessArgs, any(v))
	}

	var res []*models.Repository
	if err := db.bun.NewSelect().Model(&res).Where(strings.Join(queryParts, " OR "), typelessArgs...).Scan(ctx); err != nil {
		return nil, errors.WithStack(err)
	}
	return res, nil
}
