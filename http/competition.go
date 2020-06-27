package http

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo"

	webcore "github.com/nyks06/backapi"
	"github.com/nyks06/backapi/http/validator"
)

type CompetitionHandler struct {
	// Defines required Services here
	APIService *webcore.APIService
}

type CompetitionPayload struct {
	Name    string `json:"name"`
	SportID string `json:"sport_id"`
	StartAt string `json:"start_at"`
	EndAt   string `json:"end_at"`
}

func (h *CompetitionHandler) CreateCompetition(c echo.Context) error {
	ctx := c.Get(_FilledContext).(context.Context)

	competitionPld := new(CompetitionPayload)
	if err := c.Bind(competitionPld); err != nil {
		return c.JSON(http.StatusBadRequest, HandleHTTPResponse(http.StatusBadRequest, _messageBindError, nil))
	}

	competitionPld.Name = strings.TrimSpace(competitionPld.Name)
	startAt, err := time.Parse("2006-01-02T15:04", competitionPld.StartAt)
	if err != nil {
		return c.JSON(http.StatusBadRequest, HandleHTTPResponse(http.StatusBadRequest, _messageValidationError, nil))
	}
	endAt, err := time.Parse("2006-01-02T15:04", competitionPld.EndAt)
	if err != nil {
		return c.JSON(http.StatusBadRequest, HandleHTTPResponse(http.StatusBadRequest, _messageValidationError, nil))
	}

	validator := validator.NewValidator()
	if errs := validator.Validate(competitionPld); errs != nil {
		return c.JSON(http.StatusBadRequest, HandleHTTPResponse(http.StatusBadRequest, _messageValidationError, nil))
	}

	competition, err := h.APIService.CompetitionCreate(ctx, &webcore.Competition{
		Name: competitionPld.Name,
		Sport: &webcore.Sport{
			ID: competitionPld.SportID,
		},
		StartAt: startAt,
		EndAt:   endAt,
	})
	if err != nil {
		if webcore.IsInternalServerError(err) {
			return c.JSON(http.StatusInternalServerError, HandleHTTPResponse(http.StatusInternalServerError, _messageInternalServerError, nil))
		}
		return c.JSON(http.StatusInternalServerError, HandleHTTPResponse(http.StatusInternalServerError, _messageInternalServerError, nil))
	}

	return c.JSON(http.StatusCreated, HandleHTTPResponse(http.StatusCreated, _messageSuccess, competition))
}

func (h *CompetitionHandler) GetCompetition(c echo.Context) error {
	ctx := c.Get(_FilledContext).(context.Context)

	ID := c.Param("id")
	if ID == "" {
		return c.JSON(http.StatusBadRequest, HandleHTTPResponse(http.StatusBadRequest, _messageBindError, nil))
	}

	pGroup, err := h.APIService.CompetitionGetByID(ctx, ID)
	if err != nil {
		if webcore.IsInternalServerError(err) {
			return c.JSON(http.StatusInternalServerError, HandleHTTPResponse(http.StatusInternalServerError, _messageInternalServerError, nil))
		}
		return c.JSON(http.StatusBadRequest, HandleHTTPResponse(http.StatusBadRequest, _messageGenericBadRequestError, nil))
	}

	return c.JSON(http.StatusOK, HandleHTTPResponse(http.StatusOK, _messageSuccess, pGroup))
}

func (h *CompetitionHandler) DeleteCompetition(c echo.Context) error {
	ctx := c.Get(_FilledContext).(context.Context)

	ID := c.Param("id")
	if ID == "" {
		return c.JSON(http.StatusBadRequest, HandleHTTPResponse(http.StatusBadRequest, _messageBindError, nil))
	}

	err := h.APIService.CompetitionDelete(ctx, ID)
	if err != nil {
		if webcore.IsInternalServerError(err) {
			return c.JSON(http.StatusInternalServerError, HandleHTTPResponse(http.StatusInternalServerError, _messageInternalServerError, nil))
		}
		return c.JSON(http.StatusBadRequest, HandleHTTPResponse(http.StatusBadRequest, _messageGenericBadRequestError, nil))
	}

	return c.JSON(http.StatusOK, HandleHTTPResponse(http.StatusOK, _messageSuccess, nil))
}

func (h *CompetitionHandler) ListCompetition(c echo.Context) error {
	ctx := c.Get(_FilledContext).(context.Context)

	pGroups, err := h.APIService.CompetitionList(ctx)
	if err != nil {
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
