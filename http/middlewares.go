package http

import (
	"context"
	"time"

	"github.com/labstack/echo"
	"github.com/nyks06/backapi"
)

// middlewares is the struct responsible of handling and being used as receiver for custom middlewares
// the struct contains all the services used in the middlwares and potential vars that middlwares do want to interact with
type middlewares struct {
	SessionService *backapi.SessionService
	UserService    *backapi.UserService
}

func (m *middlewares) Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := context.Background()

		authorization := c.Request().Header.Get(HeaderKeyAuthorization)
		if authorization == "" {
			return next(c)
		}

		session, err := m.SessionService.GetByID(ctx, authorization)
		if err != nil {
			return next(c)
		}
		if session.ExpiresAt.Before(time.Now()) {
			return next(c)
		}

		user, err := m.UserService.GetByID(ctx, session.UserID)
		if err != nil {
			return next(c)
		}

		c.Set(ContextKeyCurrentUser, user)
		c.Set(ContextKeySession, session)
		return next(c)
	}
}

func (m *middlewares) MustBeAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess := c.Get(ContextKeySession)
		if sess == nil {
			return echo.ErrUnauthorized
		}
		return next(c)
	}
}

func (m *middlewares) MustBeAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cUserIfce := c.Get(ContextKeyCurrentUser)
		if cUserIfce == nil {
			return echo.ErrForbidden
		}
		cUser, ok := cUserIfce.(*backapi.User)
		if !ok {
			return echo.ErrForbidden
		}
		if cUser.Admin == false {
			return echo.ErrForbidden
		}
		return next(c)
	}
}
