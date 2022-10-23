package endpoints

import (
	"github.com/codemicro/workboat/workboat/db"
	"sync"
	"time"

	"github.com/codemicro/workboat/workboat/db/models"
	"github.com/codemicro/workboat/workboat/paths"
	"github.com/codemicro/workboat/workboat/util"
	"github.com/codemicro/workboat/workboat/views"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
)

var sessionCookieKey = "workboat_session"

type loginStateManager struct {
	states       map[string]*loginState
	lock         sync.RWMutex
	stateTimeout time.Duration
}

type loginState struct {
	createdAt time.Time
	nextURL   string
}

func newLoginStateManager() *loginStateManager {
	lsm := &loginStateManager{
		states:       make(map[string]*loginState),
		stateTimeout: time.Minute * 5,
	}
	go lsm.Worker()
	return lsm
}

func (lsm *loginStateManager) New(nextURL string) (string, error) {
	lsm.lock.Lock()
	defer lsm.lock.Unlock()

	dat, err := util.GenerateRandomDataString(30)
	if err != nil {
		return "", errors.WithStack(err)
	}

	lsm.states[dat] = &loginState{
		createdAt: time.Now(),
		nextURL:   nextURL,
	}

	return dat, nil
}

func (lsm *loginStateManager) IsValid(state string) bool {
	lsm.lock.RLock()
	defer lsm.lock.RUnlock()

	retrievedState, found := lsm.states[state]
	if !found || time.Since(retrievedState.createdAt) > lsm.stateTimeout {
		return false
	}

	return true
}

func (lsm *loginStateManager) Use(state string) (*loginState, error) {
	lsm.lock.Lock()
	defer lsm.lock.Unlock()

	retrievedState, found := lsm.states[state]
	if !found {
		return nil, errors.New("unknown state")
	}

	delete(lsm.states, state)

	if time.Since(retrievedState.createdAt) > lsm.stateTimeout {
		return nil, errors.New("expired state")
	}

	return retrievedState, nil
}

func (lsm *loginStateManager) Worker() {
	for {
		time.Sleep(lsm.stateTimeout)

		lsm.lock.Lock()

		var expired []string
		now := time.Now()
		for key, state := range lsm.states {
			if now.Sub(state.createdAt) > lsm.stateTimeout {
				expired = append(expired, key)
			}
		}

		for _, expiredKey := range expired {
			delete(lsm.states, expiredKey)
		}

		lsm.lock.Unlock()
	}
}

func (e *Endpoints) AuthLogin(ctx *fiber.Ctx) error {
	ctx.Type("html")
	return ctx.SendString(
		views.LoginPage(),
	)
}

func (e *Endpoints) AuthOauthOutbound(ctx *fiber.Ctx) error {
	var state string
	if st := ctx.Query("state"); st != "" && e.loginStateManager.IsValid(st) {
		state = st
	} else {
		st, err := e.loginStateManager.New(paths.Make(paths.Index))
		if err != nil {
			return errors.WithStack(err)
		}
		state = st
	}

	return ctx.Redirect(e.giteaClient.OauthAuthCodeURL(state))
}

func (e *Endpoints) AuthOauthInbound(ctx *fiber.Ctx) error {
	stateFromRequest := ctx.Query("state")
	var nextURL string
	if state, err := e.loginStateManager.Use(stateFromRequest); err != nil {
		return util.NewRichError(fiber.StatusBadRequest, "invalid state", err.Error())
	} else {
		nextURL = state.nextURL
	}

	session, err := e.createSessionFromOauthExchange(ctx.Query("code"))
	if err != nil {
		return errors.WithStack(err)
	}

	setCookieWithSession(ctx, session)

	return ctx.Redirect(nextURL)
}

func (e *Endpoints) createSessionFromOauthExchange(code string) (*models.Session, error) {
	token, err := e.giteaClient.OauthExchange(code)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	session, err := models.NewSession(token)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if err := e.db.InsertSession(session); err != nil {
		return nil, errors.WithStack(err)
	}

	return session, nil
}

func (e *Endpoints) getSession(ctx *fiber.Ctx) (*models.Session, bool, error) {
	if cookieValue := ctx.Cookies(sessionCookieKey); cookieValue == "" {
		return nil, false, nil
	} else {
		sess, err := e.db.GetSession(cookieValue)
		if err != nil {
			if errors.Is(err, db.ErrNotFound) {
				return nil, false, nil
			}
			return nil, false, errors.WithStack(err)
		}
		return sess, true, nil
	}
}

func setCookieWithSession(ctx *fiber.Ctx, session *models.Session) {
	ctx.Cookie(&fiber.Cookie{
		Name:     sessionCookieKey,
		Value:    session.Token,
		Expires:  session.ExpiresAt,
		HTTPOnly: true,
	})
}
