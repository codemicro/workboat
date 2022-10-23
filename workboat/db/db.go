package db

import (
	"context"
	"database/sql"
	"github.com/codemicro/workboat/workboat/config"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/extra/bundebug"
	"time"
)

type DB struct {
	pool           *sql.DB
	bun            *bun.DB
	ContextTimeout time.Duration
}

var ErrNotFound = errors.New("not found")

func New() (*DB, error) {
	dsn := config.Database.Filename
	log.Info().Msg("connecting to database")
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, errors.Wrap(err, "could not open SQL connection")
	}

	db.SetMaxOpenConns(1) // https://github.com/mattn/go-sqlite3/issues/274#issuecomment-191597862

	bundb := bun.NewDB(db, sqlitedialect.New())
	bundb.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithEnabled(config.Debug.Enabled),
	))

	rtn := &DB{
		pool:           db,
		bun:            bundb,
		ContextTimeout: time.Second,
	}

	return rtn, nil
}

func (db *DB) newContext() (context.Context, func()) {
	return context.WithTimeout(context.Background(), db.ContextTimeout)
}
