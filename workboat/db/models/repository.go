package models

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Repository struct {
	bun.BaseModel `bun:"table:sessions"`

	ID                uuid.UUID `bun:"id,pk"`
	GiteaRepositoryID int64     `bun:"gitea_repository_id,notnull"`
}
