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

type SportHandler struct {
	// Defines required Services here
	APIService *webcore.APIService
}

type SportPayload struct {
	Name string `json:"name"`
}

func (h *SportHandler) CreateSport(c echo.Context) error {
	ctx := c.Get(_FilledContext).(context.Context)

	sportPld := new(SportPayload)
	if err := c.Bind(sportPld); err != nil {
		return c.JSON(http.StatusBadRequest, HandleHTTPResponse(http.StatusBadRequest, _messageBindError, nil))
	}

	sportPld.Name = strings.TrimSpace(sportPld.Name)

	validator := validator.NewValidator()
	if errs := validator.Validate(sportPld); errs != nil {
		return c.JSON(http.StatusBadRequest, HandleHTTPResponse(http.StatusBadRequest, _messageValidationError, nil))
	}

	sport, err := h.APIService.SportCreate(ctx, &webcore.Sport{
		Name: sportPld.Name,
	})
	if err != nil {
		if webcore.IsInternalServerError(err) {
			return c.JSON(http.StatusInternalServerError, HandleHTTPResponse(http.StatusInternalServerError, _messageInternalServerError, nil))
		}
		return c.JSON(http.StatusInternalServerError, HandleHTTPResponse(http.StatusInternalServerError, _messageInternalServerError, nil))
	}

	return c.JSON(http.StatusCreated, HandleHTTPResponse(http.StatusCreated, _messageSuccess, sport))
}

func (h *SportHandler) GetSport(c echo.Context) error {
	ctx := c.Get(_FilledContext).(context.Context)

	ID := c.Param("id")
	if ID == "" {
		return c.JSON(http.StatusBadRequest, HandleHTTPResponse(http.StatusBadRequest, _messageBindError, nil))
	}

	pGroup, err := h.APIService.SportGetByID(ctx, ID)
	if err != nil {
		if webcore.IsInternalServerError(err) {
			return c.JSON(http.StatusInternalServerError, HandleHTTPResponse(http.StatusInternalServerError, _messageInternalServerError, nil))
		}
		return c.JSON(http.StatusBadRequest, HandleHTTPResponse(http.StatusBadRequest, _messageGenericBadRequestError, nil))
	}

	return c.JSON(http.StatusOK, HandleHTTPResponse(http.StatusOK, _messageSuccess, pGroup))
}

func (h *SportHandler) DeleteSport(c echo.Context) error {
	ctx := c.Get(_FilledContext).(context.Context)

	ID := c.Param("id")
	if ID == "" {
		return c.JSON(http.StatusBadRequest, HandleHTTPResponse(http.StatusBadRequest, _messageBindError, nil))
	}

	err := h.APIService.SportDelete(ctx, ID)
	if err != nil {
		if webcore.IsInternalServerError(err) {
			return c.JSON(http.StatusInternalServerError, HandleHTTPResponse(http.StatusInternalServerError, _messageInternalServerError, nil))
		}
		return c.JSON(http.StatusBadRequest, HandleHTTPResponse(http.StatusBadRequest, _messageGenericBadRequestError, nil))
	}

	return c.JSON(http.StatusOK, HandleHTTPResponse(http.StatusOK, _messageSuccess, nil))
}

func (h *SportHandler) ListSport(c echo.Context) error {
	ctx := c.Get(_FilledContext).(context.Context)

	pGroups, err := h.APIService.SportList(ctx)
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
