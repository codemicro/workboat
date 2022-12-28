package endpoints

import (
	"github.com/codemicro/workboat/workboat/config"
	"github.com/codemicro/workboat/workboat/paths"
	"sync"
	"time"

	"github.com/codemicro/workboat/workboat/util"
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

func (e *Endpoints) AuthOauthGetURL(ctx *fiber.Ctx) error {
	state, err := e.loginStateManager.New(ctx.Query("next", "/"))
	if err != nil {
		return errors.WithStack(err)
	}

	return ctx.JSON(e.giteaClient.OauthAuthCodeURL(state))
}

func (e *Endpoints) AuthOauthInbound(ctx *fiber.Ctx) error {
	stateFromRequest := ctx.Query("state")
	state, err := e.loginStateManager.Use(stateFromRequest)
	if err != nil {
		return util.NewRichError(fiber.StatusBadRequest, "invalid state", err.Error())
	}

	sess, err := e.createSessionFromOauthExchange(ctx.Query("code"))
	if err != nil {
		return errors.WithStack(err)
	}

	setCookieWithSession(ctx, sess)

	return ctx.Redirect(paths.JoinDomainAndPath(config.HTTP.ExternalURL, state.nextURL))
}

func setCookieWithSession(ctx *fiber.Ctx, sess *session) {
	ctx.Cookie(&fiber.Cookie{
		Name:    sessionCookieKey,
		Value:   sess.Token,
		Expires: sess.ExpiresAt,
	})
}
