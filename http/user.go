package http

import (
	"context"
	"net/http"
	"strings"

	"github.com/labstack/echo"

	webcore "github.com/nyks06/backapi"
	"github.com/nyks06/backapi/http/validator"
)

type UserResponse struct {
	ID            string `json:"id"`
	CustomerID    string `json:"customer_id"`
	SponsorID     string `json:"sponsor_id"`
	SponsorshipID string `json:"sponsorship_id"`
	VIP           bool   `json:"vip"`
	Admin         bool   `json:"admin"`
	Firstname     string `json:"firstname"`
	Lastname      string `json:"lastname"`
	Username      string `json:"username"`
	Email         string `json:"email"`
	PhoneNumber   string `json:"phone"`
	Telegram      string `json:"telegram"`
	AccessToken   string `json:"accessToken"`
	SubIDPrive    string `json:"sub_id_prive"`
	SubIDFun      string `json:"sub_id_fun"`
	SubIDChampion string `json:"sub_id_champion"`
	Pic           string `json:"pic"`
}

type UserHandler struct {
	// Defines required Services here
	APIService *webcore.APIService
}

type userCreatePayload struct {
	Firstname   string `json:"firstname"`
	Lastname    string `json:"lastname"`
	Username    string `json:"username"`
	SponsorID   string `json:"sponsor_id"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone"`
	Telegram    string `json:"telegram"`
	Password    string `json:"password"`
}

func (h *UserHandler) CreateUser(c echo.Context) error {
	ctx := c.Get(_FilledContext).(context.Context)

	u := new(userCreatePayload)
	if err := c.Bind(u); err != nil {
		return c.JSON(http.StatusBadRequest, HandleHTTPResponse(http.StatusBadRequest, _messageBindError, nil))
	}

	u.Firstname = strings.TrimSpace(u.Firstname)
	u.Lastname = strings.TrimSpace(u.Lastname)
	u.Username = strings.TrimSpace(u.Username)
	u.SponsorID = strings.TrimSpace(u.SponsorID)
	u.Email = strings.TrimSpace(strings.ToLower(u.Email))
	u.PhoneNumber = strings.TrimSpace(u.PhoneNumber)
	u.Telegram = strings.TrimSpace(u.Telegram)
	u.Password = strings.TrimSpace(u.Password)

	validator := validator.NewValidator()
	if errs := validator.Validate(u); errs != nil {
		return c.JSON(http.StatusBadRequest, HandleHTTPResponse(http.StatusBadRequest, _messageValidationError, nil))
	}

	user, err := h.APIService.CreateUser(ctx, &webcore.User{
		Firstname:   u.Firstname,
		Lastname:    u.Lastname,
		Username:    u.Username,
		SponsorID:   u.SponsorID,
		Email:       u.Email,
		PhoneNumber: u.PhoneNumber,
		Telegram:    u.Telegram,
		Password:    u.Password,
	})
	if err != nil {
		if webcore.IsInternalServerError(err) {
			return c.JSON(http.StatusInternalServerError, HandleHTTPResponse(http.StatusInternalServerError, _messageInternalServerError, nil))
		}
		if webcore.IsResourceAlreadyCreatedError(err) {
			return c.JSON(http.StatusConflict, HandleHTTPResponse(http.StatusConflict, _messageUserAlreadyCreatedError, nil))
		}
		return c.JSON(http.StatusInternalServerError, HandleHTTPResponse(http.StatusInternalServerError, _messageInternalServerError, nil))
	}

	return c.JSON(http.StatusCreated, HandleHTTPResponse(http.StatusCreated, _messageSuccess, user))
}

func (h *UserHandler) ListUsers(c echo.Context) error {
	ctx := c.Get(_FilledContext).(context.Context)

	users, err := h.APIService.ListUsers(ctx)
	if err != nil {
		if webcore.IsInternalServerError(err) {
			return c.JSON(http.StatusInternalServerError, HandleHTTPResponse(http.StatusInternalServerError, _messageInternalServerError, nil))
		}
		return c.JSON(http.StatusBadRequest, HandleHTTPResponse(http.StatusBadRequest, _messageGenericBadRequestError, nil))
	}

	return c.JSON(http.StatusOK, HandleHTTPResponse(http.StatusOK, _messageSuccess, users))
}

func (h *UserHandler) GetUser(c echo.Context) error {
	ctx := c.Get(_FilledContext).(context.Context)

	UserID := c.Param("id")
	if UserID == "" {
		return c.JSON(http.StatusBadRequest, HandleHTTPResponse(http.StatusBadRequest, _messageBindError, nil))
	}

	user, err := h.APIService.GetUserByID(ctx, UserID)
	if err != nil {
		if webcore.IsInternalServerError(err) {
			return c.JSON(http.StatusInternalServerError, HandleHTTPResponse(http.StatusInternalServerError, _messageInternalServerError, nil))
		}
		return c.JSON(http.StatusBadRequest, HandleHTTPResponse(http.StatusBadRequest, _messageGenericBadRequestError, nil))
	}

	return c.JSON(http.StatusOK, HandleHTTPResponse(http.StatusOK, _messageSuccess, user))
}

type userChangePasswordPayload struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

func (h *UserHandler) UserChangePassword(c echo.Context) error {
	ctx := c.Get(_FilledContext).(context.Context)
	CUser := ctx.Value(_CUser).(*webcore.User)

	u := new(userChangePasswordPayload)
	if err := c.Bind(u); err != nil {
		return c.JSON(http.StatusBadRequest, HandleHTTPResponse(http.StatusBadRequest, _messageBindError, nil))
	}

	user, err := h.APIService.ChangePassword(ctx, CUser, u.CurrentPassword, u.NewPassword)
	if err != nil {
		if webcore.IsInternalServerError(err) {
			return c.JSON(http.StatusInternalServerError, HandleHTTPResponse(http.StatusInternalServerError, _messageInternalServerError, nil))
		}
		if webcore.IsUnauthorizedError(err) {
			return c.JSON(http.StatusUnauthorized, HandleHTTPResponse(http.StatusUnauthorized, _messageUnauthorizedError, nil))
		}
		return c.JSON(http.StatusInternalServerError, HandleHTTPResponse(http.StatusInternalServerError, _messageInternalServerError, nil))
	}

	return c.JSON(http.StatusOK, HandleHTTPResponse(http.StatusOK, _messageSuccess, user))
}

type userUpdateDetailsPayload struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

func (h *UserHandler) UserUpdateDetails(c echo.Context) error {
	ctx := c.Get(_FilledContext).(context.Context)
	CUser := ctx.Value(_CUser).(*webcore.User)

	u := new(userUpdateDetailsPayload)
	if err := c.Bind(u); err != nil {
		return c.JSON(http.StatusBadRequest, HandleHTTPResponse(http.StatusBadRequest, _messageBindError, nil))
	}

	user, err := h.APIService.UpdateDetails(ctx, CUser, u.Firstname, u.Lastname)
	if err != nil {
		if webcore.IsInternalServerError(err) {
			return c.JSON(http.StatusInternalServerError, HandleHTTPResponse(http.StatusInternalServerError, _messageInternalServerError, nil))
		}
		if webcore.IsUnauthorizedError(err) {
			return c.JSON(http.StatusUnauthorized, HandleHTTPResponse(http.StatusUnauthorized, _messageUnauthorizedError, nil))
		}
		return c.JSON(http.StatusInternalServerError, HandleHTTPResponse(http.StatusInternalServerError, _messageInternalServerError, nil))
	}

	return c.JSON(http.StatusOK, HandleHTTPResponse(http.StatusOK, _messageSuccess, user))
}

type userUpdateContactSettingsPayload struct {
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Telegram string `json:"telegram"`
}

func (h *UserHandler) UserUpdateContactSettings(c echo.Context) error {
	ctx := c.Get(_FilledContext).(context.Context)
	CUser := ctx.Value(_CUser).(*webcore.User)

	u := new(userUpdateContactSettingsPayload)
	if err := c.Bind(u); err != nil {
		return c.JSON(http.StatusBadRequest, HandleHTTPResponse(http.StatusBadRequest, _messageBindError, nil))
	}

	user, err := h.APIService.UpdateContactSettings(ctx, CUser, u.Email, u.Phone, u.Telegram)
	if err != nil {
		if webcore.IsInternalServerError(err) {
			return c.JSON(http.StatusInternalServerError, HandleHTTPResponse(http.StatusInternalServerError, _messageInternalServerError, nil))
		}
		if webcore.IsUnauthorizedError(err) {
			return c.JSON(http.StatusUnauthorized, HandleHTTPResponse(http.StatusUnauthorized, _messageUnauthorizedError, nil))
		}
		return c.JSON(http.StatusInternalServerError, HandleHTTPResponse(http.StatusInternalServerError, _messageInternalServerError, nil))
	}

	return c.JSON(http.StatusOK, HandleHTTPResponse(http.StatusOK, _messageSuccess, user))
}

func (h *UserHandler) GetCUser(c echo.Context) error {
	ctx := c.Get(_FilledContext).(context.Context)

	// // Not required to verify it as the endpoint is only accessible with auth
	CUser := ctx.Value(_CUser).(*webcore.User)
	Session := ctx.Value(_Session).(*webcore.Session)

	user := &UserResponse{
		ID:            CUser.ID,
		CustomerID:    CUser.CustomerID,
		SponsorID:     CUser.SponsorID,
		SponsorshipID: CUser.SponsorshipID,
		VIP:           CUser.VIP,
		Admin:         CUser.Admin,
		Firstname:     CUser.Firstname,
		Lastname:      CUser.Lastname,
		Username:      CUser.Username,
		Email:         CUser.Email,
		PhoneNumber:   CUser.PhoneNumber,
		Telegram:      CUser.Telegram,
		AccessToken:   Session.ID,
		SubIDPrive:    CUser.SubIDPrive,
		SubIDFun:      CUser.SubIDFun,
		SubIDChampion: CUser.SubIDChampion,
		Pic:           "./assets/media/users/default.jpg",
	}
	return c.JSON(http.StatusOK, HandleHTTPResponse(http.StatusOK, _messageSuccess, user))
}
