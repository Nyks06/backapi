package http

import (
	"context"
	"net/http"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/labstack/echo"

	webcore "github.com/nyks06/backapi"
	"github.com/nyks06/backapi/http/validator"
)

type TicketHandler struct {
	// Defines required Services here
	APIService *webcore.APIService
}

type ticketPayload struct {
	Title  string  `json:"title"`
	Stake  float64 `json:"stake"`
	Public bool    `json:"is_public"`
	Live   bool    `json:"live"`
	Risk   string  `json:"risk"`
	Pack   string  `json:"pack"`
	Odd    float64 `json:"odd"`
	Status string  `json:"status"`
}

func (h *TicketHandler) CreateTicket(c echo.Context) error {
	ctx := c.Get(_FilledContext).(context.Context)

	ticketPld := new(ticketPayload)
	if err := c.Bind(ticketPld); err != nil {
		return c.JSON(http.StatusBadRequest, HandleHTTPResponse(http.StatusBadRequest, _messageBindError, nil))
	}

	ticketPld.Title = strings.TrimSpace(ticketPld.Title)

	validator := validator.NewValidator()
	if errs := validator.Validate(ticketPld); errs != nil {
		return c.JSON(http.StatusBadRequest, HandleHTTPResponse(http.StatusBadRequest, _messageValidationError, nil))
	}

	ticket, err := h.APIService.TicketCreate(ctx, &webcore.Ticket{
		Title:  ticketPld.Title,
		Stake:  ticketPld.Stake,
		Public: ticketPld.Public,
		Live:   ticketPld.Live,
		Risk:   ticketPld.Risk,
		Pack:   ticketPld.Pack,
	})
	if err != nil {
		if webcore.IsInternalServerError(err) {
			return c.JSON(http.StatusInternalServerError, HandleHTTPResponse(http.StatusInternalServerError, _messageInternalServerError, nil))
		}
		return c.JSON(http.StatusInternalServerError, HandleHTTPResponse(http.StatusInternalServerError, _messageInternalServerError, nil))
	}

	return c.JSON(http.StatusCreated, HandleHTTPResponse(http.StatusCreated, _messageSuccess, ticket))
}

func (h *TicketHandler) GetTicket(c echo.Context) error {
	ctx := c.Get(_FilledContext).(context.Context)

	ID := c.Param("id")
	if ID == "" {
		return c.JSON(http.StatusBadRequest, HandleHTTPResponse(http.StatusBadRequest, _messageBindError, nil))
	}

	pGroup, err := h.APIService.TicketGetByID(ctx, ID)
	if err != nil {
		spew.Dump(err)
		if webcore.IsInternalServerError(err) {
			return c.JSON(http.StatusInternalServerError, HandleHTTPResponse(http.StatusInternalServerError, _messageInternalServerError, nil))
		}
		return c.JSON(http.StatusBadRequest, HandleHTTPResponse(http.StatusBadRequest, _messageGenericBadRequestError, nil))
	}

	return c.JSON(http.StatusOK, HandleHTTPResponse(http.StatusOK, _messageSuccess, pGroup))
}

func (h *TicketHandler) DeleteTicket(c echo.Context) error {
	ctx := c.Get(_FilledContext).(context.Context)

	ID := c.Param("id")
	if ID == "" {
		return c.JSON(http.StatusBadRequest, HandleHTTPResponse(http.StatusBadRequest, _messageBindError, nil))
	}

	err := h.APIService.TicketDelete(ctx, ID)
	if err != nil {
		if webcore.IsInternalServerError(err) {
			return c.JSON(http.StatusInternalServerError, HandleHTTPResponse(http.StatusInternalServerError, _messageInternalServerError, nil))
		}
		return c.JSON(http.StatusBadRequest, HandleHTTPResponse(http.StatusBadRequest, _messageGenericBadRequestError, nil))
	}

	return c.JSON(http.StatusOK, HandleHTTPResponse(http.StatusOK, _messageSuccess, nil))
}
func (h *TicketHandler) UpdateTicket(c echo.Context) error {
	ctx := c.Get(_FilledContext).(context.Context)

	ID := c.Param("id")
	if ID == "" {
		return c.JSON(http.StatusBadRequest, HandleHTTPResponse(http.StatusBadRequest, _messageBindError, nil))
	}

	ticketPayload := new(ticketPayload)
	if err := c.Bind(ticketPayload); err != nil {
		return c.JSON(http.StatusBadRequest, HandleHTTPResponse(http.StatusBadRequest, _messageBindError, nil))
	}

	ticketPayload.Title = strings.TrimSpace(ticketPayload.Title)

	err := h.APIService.TicketUpdate(ctx, ID, &webcore.Ticket{
		Title:  ticketPayload.Title,
		Stake:  ticketPayload.Stake,
		Public: ticketPayload.Public,
		Live:   ticketPayload.Live,
		Risk:   ticketPayload.Risk,
	})
	if err != nil {
		if webcore.IsInternalServerError(err) {
			spew.Dump(err)
			return c.JSON(http.StatusInternalServerError, HandleHTTPResponse(http.StatusInternalServerError, _messageInternalServerError, nil))
		}
		return c.JSON(http.StatusBadRequest, HandleHTTPResponse(http.StatusBadRequest, _messageGenericBadRequestError, nil))
	}

	return c.JSON(http.StatusOK, HandleHTTPResponse(http.StatusOK, _messageSuccess, nil))
}

func (h *TicketHandler) ListTicket(c echo.Context) error {
	ctx := c.Get(_FilledContext).(context.Context)

	pGroups, err := h.APIService.TicketList(ctx)
	if err != nil {
		spew.Dump(err)
		if webcore.IsInternalServerError(err) {
			return c.JSON(http.StatusInternalServerError, HandleHTTPResponse(http.StatusInternalServerError, _messageInternalServerError, nil))
		}
		if webcore.IsUnauthorizedError(err) {
			return c.JSON(http.StatusUnauthorized, HandleHTTPResponse(http.StatusUnauthorized, _messageUnauthorizedError, nil))
		}
		return c.JSON(http.StatusInternalServerError, HandleHTTPResponse(http.StatusInternalServerError, _messageInternalServerError, nil))
	}

	return c.JSON(http.StatusOK, HandleHTTPResponse(http.StatusOK, _messageSuccess, pGroups))
}
