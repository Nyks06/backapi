package http

import (
	"context"
	"net/http"
	"strings"

	"github.com/davecgh/go-spew/spew"
	webcore "github.com/nyks06/backapi"
	"github.com/nyks06/backapi/http/validator"

	"github.com/labstack/echo"
)

type SessionHandler struct {
	APIService *webcore.APIService
}

type sessionCreatePayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *SessionHandler) CreateSession(c echo.Context) error {
	ctx := c.Get(_FilledContext).(context.Context)

	s := new(sessionCreatePayload)
	if err := c.Bind(s); err != nil {
		return c.JSON(http.StatusBadRequest, HandleHTTPResponse(http.StatusBadRequest, _messageBindError, nil))
	}

	s.Email = strings.TrimSpace(strings.ToLower(s.Email))
	s.Password = strings.TrimSpace(s.Password)

	validator := validator.NewValidator()
	if errs := validator.Validate(s); errs != nil {
		return c.JSON(http.StatusBadRequest, HandleHTTPResponse(http.StatusBadRequest, _messageValidationError, nil))
	}

	session, user, err := h.APIService.CreateSession(ctx, s.Email, s.Password)
	if err != nil {
		if webcore.IsInternalServerError(err) {
			return c.JSON(http.StatusInternalServerError, HandleHTTPResponse(http.StatusInternalServerError, _messageInternalServerError, nil))
		}
		return c.JSON(http.StatusBadRequest, HandleHTTPResponse(http.StatusBadRequest, _messageGenericBadRequestError, nil))
	}
	spew.Dump(user)

	u := UserResponse{
		ID:          user.ID,
		CustomerID:  user.CustomerID,
		Firstname:   user.Firstname,
		Lastname:    user.Lastname,
		Username:    user.Username,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		Admin:       user.Admin,
		VIP:         user.VIP,
		AccessToken: session.ID,
		Pic:         "./assets/media/users/default.jpg",
	}
	return c.JSON(http.StatusCreated, &u)
}

func (h *SessionHandler) RemoveSession(c echo.Context) error {
	ctx := c.Get(_FilledContext).(context.Context)

	session, ok := ctx.Value(_Session).(*webcore.Session)
	if !ok {
		return c.JSON(http.StatusUnauthorized, HandleHTTPResponse(http.StatusUnauthorized, _messageUnauthorizedError, nil))
	}

	if err := h.APIService.RemoveSession(ctx, session.ID); err != nil {
		if webcore.IsInternalServerError(err) {
			return c.JSON(http.StatusInternalServerError, HandleHTTPResponse(http.StatusInternalServerError, _messageInternalServerError, nil))
		}
		return c.JSON(http.StatusBadRequest, HandleHTTPResponse(http.StatusBadRequest, _messageGenericBadRequestError, nil))
	}

	return c.JSON(http.StatusNoContent, nil)
}
