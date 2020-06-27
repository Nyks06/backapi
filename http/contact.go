package http

import (
	"context"
	"net/http"

	"github.com/labstack/echo"

	webcore "github.com/nyks06/backapi"
)

type ContactHandler struct {
	APIService *webcore.APIService
	// Defines required Services here
}

type sendMessagePayload struct {
	Name    string `json:"name" form:"name" query:"name" valid:"required"`
	Phone   string `json:"phone" form:"phone" query:"phone" valid:"required"`
	Email   string `json:"email" form:"email" query:"email" valid:"required"`
	Message string `json:"message" form:"message" query:"message" valid:"required"`
}

func (h *ContactHandler) SendMessage(c echo.Context) error {
	ctx := c.Get(_FilledContext).(context.Context)

	m := new(sendMessagePayload)
	if err := c.Bind(m); err != nil {
		return c.JSON(http.StatusBadRequest, HandleHTTPResponse(http.StatusBadRequest, _messageBindError, nil))
	}

	if err := h.APIService.MessageCreate(ctx, m.Name, m.Phone, m.Email, m.Message); err != nil {
		if webcore.IsInternalServerError(err) {
			return c.JSON(http.StatusInternalServerError, HandleHTTPResponse(http.StatusInternalServerError, _messageInternalServerError, nil))
		}
		if webcore.IsResourceAlreadyCreatedError(err) {
			return c.JSON(http.StatusConflict, HandleHTTPResponse(http.StatusConflict, _messageUserAlreadyCreatedError, nil))
		}
		return c.JSON(http.StatusInternalServerError, HandleHTTPResponse(http.StatusInternalServerError, _messageInternalServerError, nil))
	}

	return c.JSON(http.StatusCreated, HandleHTTPResponse(http.StatusOK, _messageSuccess, nil))
}
