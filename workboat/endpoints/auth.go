package endpoints

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/codemicro/workboat/workboat/util"
	"github.com/codemicro/workboat/workboat/views"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"sync"
	"time"
)

type loginStateManager struct {
	states       map[string]time.Time
	lock         sync.Mutex
	stateTimeout time.Duration
}

func newLoginStateManager() *loginStateManager {
	lsm := &loginStateManager{
		states:       make(map[string]time.Time),
		stateTimeout: time.Minute * 5,
	}
	go lsm.Worker()
	return lsm
}

func (lsm *loginStateManager) New() (string, error) {
	lsm.lock.Lock()
	defer lsm.lock.Unlock()

	randData := make([]byte, 30)
	if _, err := rand.Read(randData); err != nil {
		return "", errors.WithStack(err)
	}
	asBase64 := base64.URLEncoding.EncodeToString(randData)
	lsm.states[asBase64] = time.Now()
	return asBase64, nil
}

func (lsm *loginStateManager) Use(state string) error {
	lsm.lock.Lock()
	defer lsm.lock.Unlock()

	createdAt, found := lsm.states[state]
	if !found {
		return errors.New("unknown state")
	}

	delete(lsm.states, state)

	if time.Now().Sub(createdAt) > lsm.stateTimeout {
		return errors.New("expired state")
	}

	return nil
}

func (lsm *loginStateManager) Worker() {
	for {
		time.Sleep(lsm.stateTimeout)

		lsm.lock.Lock()

		var expired []string
		now := time.Now()
		for key, createdAt := range lsm.states {
			if now.Sub(createdAt) > lsm.stateTimeout {
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
	state, err := e.login.stateManager.New()
	if err != nil {
		return errors.WithStack(err)
	}

	return ctx.Redirect(e.login.oauthConfig.AuthCodeURL(state))
}

func (e *Endpoints) AuthOauthInbound(ctx *fiber.Ctx) error {
	stateFromRequest := ctx.Query("state")
	if err := e.login.stateManager.Use(stateFromRequest); err != nil {
		return util.NewRichError(fiber.StatusBadRequest, "invalid state", err.Error())
	}

	exchangeCtx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	token, err := e.login.oauthConfig.Exchange(exchangeCtx, ctx.Query("code"))
	if err != nil {
		return errors.WithStack(err)
	}

	// TODO: databaseify
	return ctx.SendString(fmt.Sprintf("%#v", token))
}
