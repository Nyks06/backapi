package http

import (
	"context"
	"net/http"

	"github.com/labstack/echo"

	webcore "github.com/nyks06/backapi"
)

type SubscriptionHandler struct {
	APIService *webcore.APIService
	// Defines required Services here
}

type subscriptionCreatePayload struct {
	PlanID   string `json:"plan_id" form:"plan_id" query:"plan_id" valid:"required"`
	SourceID string `json:"source_id" form:"source_id" query:"source_id" valid:"required"`
}

func (h *SubscriptionHandler) Create(c echo.Context) error {
	ctx := c.Get(_FilledContext).(context.Context)
	CUser := ctx.Value(_CUser).(*webcore.User)

	m := new(subscriptionCreatePayload)
	if err := c.Bind(m); err != nil {
		return c.JSON(http.StatusBadRequest, HandleHTTPResponse(http.StatusBadRequest, _messageBindError, nil))
	}

	if err := h.APIService.SubscriptionCreate(ctx, CUser, m.PlanID, m.SourceID); err != nil {
		if webcore.IsInternalServerError(err) {
			return c.JSON(http.StatusInternalServerError, HandleHTTPResponse(http.StatusInternalServerError, _messageInternalServerError, nil))
		}
		if webcore.IsResourceAlreadyCreatedError(err) {
			return c.JSON(http.StatusConflict, HandleHTTPResponse(http.StatusConflict, _messageUserAlreadyCreatedError, nil))
		}
		return c.JSON(http.StatusInternalServerError, HandleHTTPResponse(http.StatusInternalServerError, _messageInternalServerError, nil))
	}

	return c.JSON(http.StatusCreated, HandleHTTPResponse(http.StatusCreated, _messageSuccess, nil))
}

func (h *SubscriptionHandler) Get(c echo.Context) error {
	ctx := c.Get(_FilledContext).(context.Context)
	CUser := ctx.Value(_CUser).(*webcore.User)

	subs, err := h.APIService.SubscriptionGet(ctx, CUser.CustomerID)
	if err != nil {
		if webcore.IsInternalServerError(err) {
			return c.JSON(http.StatusInternalServerError, HandleHTTPResponse(http.StatusInternalServerError, _messageInternalServerError, nil))
		}
		if webcore.IsResourceAlreadyCreatedError(err) {
			return c.JSON(http.StatusConflict, HandleHTTPResponse(http.StatusConflict, _messageUserAlreadyCreatedError, nil))
		}
		return c.JSON(http.StatusInternalServerError, HandleHTTPResponse(http.StatusInternalServerError, _messageInternalServerError, nil))
	}

	return c.JSON(http.StatusCreated, HandleHTTPResponse(http.StatusCreated, _messageSuccess, subs))
}

// func (h *SubscriptionHandler) Cancel(c echo.Context) error {
// 	ctx := c.Get(_FilledContext).(context.Context)
// 	CUser := ctx.Value(_CUser).(*webcore.User)

// 	SubID := c.Param("id")

// 	err := h.APIService.SubscriptionCancel(ctx, CUser.CustomerID, SubID)
// 	if err != nil {
// 		if webcore.IsInternalServerError(err) {
// 			return c.JSON(http.StatusInternalServerError, HandleHTTPResponse(http.StatusInternalServerError, _messageInternalServerError, nil))
// 		}
// 		if webcore.IsResourceAlreadyCreatedError(err) {
// 			return c.JSON(http.StatusConflict, HandleHTTPResponse(http.StatusConflict, _messageUserAlreadyCreatedError, nil))
// 		}
// 		return c.JSON(http.StatusInternalServerError, HandleHTTPResponse(http.StatusInternalServerError, _messageInternalServerError, nil))
// 	}

// 	return c.JSON(http.StatusNoContent, HandleHTTPResponse(http.StatusNoContent, _messageSuccess, nil))
// }
