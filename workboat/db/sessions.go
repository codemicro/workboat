package db

import (
	"database/sql"
	"github.com/codemicro/workboat/workboat/db/models"
	"github.com/pkg/errors"
)

func (db *DB) InsertSession(session *models.Session) error {
	ctx, cancel := db.newContext()
	defer cancel()

	_, err := db.bun.NewInsert().Model(session).Exec(ctx)
	return errors.WithStack(err)
}

func (db *DB) GetSession(key string) (*models.Session, error) {
	ctx, cancel := db.newContext()
	defer cancel()

	res := new(models.Session)
	if err := db.bun.NewSelect().Model(res).Where("token = ?", key).Scan(ctx); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, errors.WithStack(err)
	}
	return res, nil
}
