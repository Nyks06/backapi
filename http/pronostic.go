package http

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/labstack/echo"

	webcore "github.com/nyks06/backapi"
	"github.com/nyks06/backapi/http/validator"
)

type PronosticHandler struct {
	// Defines required Services here
	APIService *webcore.APIService
}

type pronosticCreatePayload struct {
	TicketID      string `json:"ticket_id"`
	FirstTeam     string `json:"first_team"`
	SecondTeam    string `json:"second_team"`
	Pronostic     string `json:"pronostic"`
	CompetitionID string `json:"competition_id"`
	SportID       string `json:"sport_id"`
	Status        string `json:"status"`
	Odd           string `json:"odd"`
	EventDate     string `json:"event_date"`
}

func (h *PronosticHandler) CreatePronostic(c echo.Context) error {
	ctx := c.Get(_FilledContext).(context.Context)

	prono := new(pronosticCreatePayload)
	if err := c.Bind(prono); err != nil {
		return c.JSON(http.StatusBadRequest, HandleHTTPResponse(http.StatusBadRequest, _messageBindError, nil))
	}

	validator := validator.NewValidator()
	if errs := validator.Validate(prono); errs != nil {
		return c.JSON(http.StatusBadRequest, HandleHTTPResponse(http.StatusBadRequest, _messageValidationError, nil))
	}

	odd, err := strconv.ParseFloat(prono.Odd, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, HandleHTTPResponse(http.StatusBadRequest, _messageValidationError, nil))
	}

	eventDate, err := time.Parse("2006-01-02T15:04", prono.EventDate)
	if err != nil {
		return c.JSON(http.StatusBadRequest, HandleHTTPResponse(http.StatusBadRequest, _messageValidationError, nil))
	}

	pronoGroup, err := h.APIService.PronosticCreate(ctx, &webcore.Pronostic{
		TicketID:   prono.TicketID,
		FirstTeam:  prono.FirstTeam,
		SecondTeam: prono.SecondTeam,
		Pronostic:  prono.Pronostic,
		Competition: &webcore.Competition{
			ID: prono.CompetitionID,
		},
		Sport: &webcore.Sport{
			ID: prono.SportID,
		},
		Status:    prono.Status,
		Odd:       odd,
		EventDate: eventDate,
	})
	if err != nil {
		if webcore.IsInternalServerError(err) {
			return c.JSON(http.StatusInternalServerError, HandleHTTPResponse(http.StatusInternalServerError, _messageInternalServerError, nil))
		}
		return c.JSON(http.StatusInternalServerError, HandleHTTPResponse(http.StatusInternalServerError, _messageInternalServerError, nil))
	}

	return c.JSON(http.StatusCreated, HandleHTTPResponse(http.StatusCreated, _messageSuccess, pronoGroup))
}

func (h *PronosticHandler) GetPronostic(c echo.Context) error {
	ctx := c.Get(_FilledContext).(context.Context)

	ID := c.Param("id")
	if ID == "" {
		return c.JSON(http.StatusBadRequest, HandleHTTPResponse(http.StatusBadRequest, _messageBindError, nil))
	}

	prono, err := h.APIService.PronosticGetByID(ctx, ID)
	if err != nil {
		if webcore.IsInternalServerError(err) {
			return c.JSON(http.StatusInternalServerError, HandleHTTPResponse(http.StatusInternalServerError, _messageInternalServerError, nil))
		}
		return c.JSON(http.StatusBadRequest, HandleHTTPResponse(http.StatusBadRequest, _messageGenericBadRequestError, nil))
	}

	return c.JSON(http.StatusOK, HandleHTTPResponse(http.StatusOK, _messageSuccess, prono))
}

type pronosticUpdatePayload struct {
	Status string `json:"status"`
}

func (h *PronosticHandler) UpdatePronostic(c echo.Context) error {
	ctx := c.Get(_FilledContext).(context.Context)

	ID := c.Param("id")
	if ID == "" {
		return c.JSON(http.StatusBadRequest, HandleHTTPResponse(http.StatusBadRequest, _messageBindError, nil))
	}

	prono := new(pronosticUpdatePayload)
	if err := c.Bind(prono); err != nil {
		return c.JSON(http.StatusBadRequest, HandleHTTPResponse(http.StatusBadRequest, _messageBindError, nil))
	}

	validator := validator.NewValidator()
	if errs := validator.Validate(prono); errs != nil {
		return c.JSON(http.StatusBadRequest, HandleHTTPResponse(http.StatusBadRequest, _messageValidationError, nil))
	}

	err := h.APIService.PronosticUpdate(ctx, ID, prono.Status)
	if err != nil {
		if webcore.IsInternalServerError(err) {
			return c.JSON(http.StatusInternalServerError, HandleHTTPResponse(http.StatusInternalServerError, _messageInternalServerError, nil))
		}
		return c.JSON(http.StatusBadRequest, HandleHTTPResponse(http.StatusBadRequest, _messageGenericBadRequestError, nil))
	}

	return c.JSON(http.StatusOK, HandleHTTPResponse(http.StatusOK, _messageSuccess, nil))
}

func (h *PronosticHandler) DeletePronostic(c echo.Context) error {
	ctx := c.Get(_FilledContext).(context.Context)

	ID := c.Param("id")
	if ID == "" {
		return c.JSON(http.StatusBadRequest, HandleHTTPResponse(http.StatusBadRequest, _messageBindError, nil))
	}

	err := h.APIService.PronosticDelete(ctx, ID)
	if err != nil {
		if webcore.IsInternalServerError(err) {
			return c.JSON(http.StatusInternalServerError, HandleHTTPResponse(http.StatusInternalServerError, _messageInternalServerError, nil))
		}
		return c.JSON(http.StatusBadRequest, HandleHTTPResponse(http.StatusBadRequest, _messageGenericBadRequestError, nil))
	}

	return c.JSON(http.StatusOK, HandleHTTPResponse(http.StatusOK, _messageSuccess, nil))
}

func (h *PronosticHandler) ListPronostic(c echo.Context) error {
	ctx := c.Get(_FilledContext).(context.Context)

	pGroups, err := h.APIService.PronosticList(ctx)
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
