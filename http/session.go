package http

import (
	"context"
	"net/http"
	"strings"

	"github.com/nyks06/backapi"
	"github.com/nyks06/backapi/http/validator"

	"github.com/labstack/echo"
)

type SessionHandler struct {
	SessionService *backapi.SessionService
}

type sessionCreatePayload struct {
	Email    string `json:"email" form:"email" valid:"required"`
	Password string `json:"password" form:"password" valid:"required"`
}

type sessionCreateResponse struct {
	Session *backapi.Session `json:"session"`
}

func (h *SessionHandler) CreateSession(c echo.Context) error {
	ctx := context.Background()

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

	session, err := h.SessionService.Create(ctx, s.Email, s.Password)
	if err != nil {
		if backapi.IsInternalServerError(err) {
			return c.JSON(http.StatusInternalServerError, HandleHTTPResponse(http.StatusInternalServerError, _messageInternalServerError, nil))
		}
		return c.JSON(http.StatusBadRequest, HandleHTTPResponse(http.StatusBadRequest, _messageGenericBadRequestError, nil))
	}

	resp := &sessionCreateResponse{
		Session: session,
	}
	return c.JSON(http.StatusCreated, HandleHTTPResponse(http.StatusCreated, _messageSuccess, resp))
}

func (h *SessionHandler) RemoveSession(c echo.Context) error {
	ctx := context.Background()

	session, ok := c.Get(ContextKeySession).(*backapi.Session)
	if !ok {
		return c.JSON(http.StatusUnauthorized, HandleHTTPResponse(http.StatusUnauthorized, _messageUnauthorizedError, nil))
	}

	if err := h.SessionService.Delete(ctx, session.ID); err != nil {
		if backapi.IsInternalServerError(err) {
			return c.JSON(http.StatusInternalServerError, HandleHTTPResponse(http.StatusInternalServerError, _messageInternalServerError, nil))
		}
		return c.JSON(http.StatusBadRequest, HandleHTTPResponse(http.StatusBadRequest, _messageGenericBadRequestError, nil))
	}

	return c.JSON(http.StatusNoContent, nil)
}
