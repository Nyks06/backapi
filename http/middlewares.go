package http

import (
	"context"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/nyks06/backapi"

	"github.com/labstack/echo"
)

const _FilledContext = "_FilledContext"
const _Authorization = "_Authorization"
const _CUser = "_CUser"
const _Session = "_Session"

// middlewares is the struct responsible of handling and being used as receiver for custom middlewares
// the struct contains all the services used in the middlwares and potential vars that middlwares do want to interact with
type middlewares struct {
	APIService *webcore.APIService
}

// Defines all the middlewares here
func (m *middlewares) SetContext(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		authorization := c.Request().Header.Get("Authorization")
		spew.Dump(authorization)
		if authorization != "" {
			ctx = context.WithValue(ctx, _Authorization, authorization)
		}

		c.Set(_FilledContext, ctx)
		return next(c)
	}
}

func (m *middlewares) Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Get(_FilledContext).(context.Context)

		authorization, ok := ctx.Value(_Authorization).(string)
		if !ok {
			spew.Dump("No authorization")
			return next(c)
		}
		session, err := m.APIService.GetSessionByID(ctx, authorization)
		if err != nil {
			// Should handle an error here
			spew.Dump("cannot find session", err)
			return next(c)
		}
		if session.ExpiresAt.Before(time.Now()) {
			// Should handle an error here
			spew.Dump("expired one")
			return next(c)
		}
		user, err := m.APIService.GetUserByID(ctx, session.UserID)
		if err != nil {
			// Should handle an error here
			spew.Dump("cannot get user", err)
			return next(c)
		}

		ctx = context.WithValue(ctx, _CUser, user)
		ctx = context.WithValue(ctx, _Session, session)

		c.Set(_FilledContext, ctx)
		return next(c)
	}
}

func (m *middlewares) MustBeAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Get(_FilledContext).(context.Context)

		_, ok := ctx.Value(_CUser).(*webcore.User)
		if !ok {
			return echo.ErrUnauthorized
		}

		return next(c)
	}
}

// func (m *middlewares) MustBeAdmin(next echo.HandlerFunc) echo.HandlerFunc {
// 	return func(c echo.Context) error {
// 		cc := c.(*Context)

// 		// No user set
// 		if cc.CUser == nil || cc.CUser.Admin == false {
// 			cc.Redirect(http.StatusTemporaryRedirect, "/login")
// 		}
// 		return next(cc)
// 	}
// }

// func (m *middlewares) MustBeVIP(next echo.HandlerFunc) echo.HandlerFunc {
// 	return func(c echo.Context) error {
// 		cc := c.(*Context)

// 		// No user set
// 		if cc.CUser == nil {
// 			cc.Redirect(http.StatusTemporaryRedirect, "/login")
// 		} else if vip := m.UserService.IsVIP(cc.CUser); vip == false {
// 			cc.Redirect(http.StatusTemporaryRedirect, "/vip_info")
// 		}
// 		return next(cc)
// 	}
// }
