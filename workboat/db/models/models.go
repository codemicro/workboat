package models

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users"`

	ID           uuid.UUID `bun:"id,pk,type:varchar"`
	EmailAddress string    `bun:"email_address,notnull"`
}
