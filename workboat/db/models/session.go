package models

import (
	"github.com/codemicro/workboat/workboat/util"
	"github.com/pkg/errors"
	"github.com/uptrace/bun"
	"golang.org/x/oauth2"
	"time"
)

type Session struct {
	bun.BaseModel `bun:"table:sessions"`

	Token     string    `bun:"token,pk"`
	ExpiresAt time.Time `bun:"expires_at,notnull"`

	GiteaToken *oauth2.Token `bun:"gitea_token,notnull"`
}

func NewSession(giteaToken *oauth2.Token) (*Session, error) {
	sess := new(Session)
	sess.GiteaToken = giteaToken
	sess.ExpiresAt = time.Now().UTC().Add(time.Hour * 24 * 7)

	var err error
	sess.Token, err = util.GenerateRandomDataString(60)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return sess, nil
}
