package endpoints

import (
	"github.com/codemicro/workboat/workboat/util"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"time"
)

type session struct {
	Token      string
	ExpiresAt  time.Time
	GiteaToken *oauth2.Token
}

func newSession(giteaToken *oauth2.Token) (*session, error) {
	sess := new(session)
	sess.GiteaToken = giteaToken
	sess.ExpiresAt = time.Now().UTC().Add(time.Hour * 24)

	var err error
	sess.Token, err = util.GenerateRandomDataString(60)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return sess, nil
}

func (e *Endpoints) createSessionFromOauthExchange(code string) (*session, error) {
	token, err := e.giteaClient.OauthExchange(code)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	sess, err := newSession(token)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	e.sessionLock.Lock()
	defer e.sessionLock.Unlock()

	e.sessions[sess.Token] = sess

	return sess, nil
}

func (e *Endpoints) getSession(ctx *fiber.Ctx) (*session, bool, error) {
	if cookieValue := ctx.Cookies(sessionCookieKey); cookieValue == "" {
		return nil, false, nil
	} else {
		e.sessionLock.Lock()
		defer e.sessionLock.Unlock()

		sess, found := e.sessions[cookieValue]
		if !found {
			return nil, false, nil
		}

		if !sess.ExpiresAt.After(time.Now().UTC()) {
			delete(e.sessions, cookieValue)
			return nil, false, nil
		}

		return sess, true, nil
	}
}
