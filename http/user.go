package http

import (
	"context"
	"net/http"
	"strings"

	"github.com/labstack/echo"

	"github.com/nyks06/backapi"
	"github.com/nyks06/backapi/http/validator"
)

type UserHandler struct {
	// Defines required Services here
	UserService *backapi.UserService
}

type userCreatePayload struct {
	Firstname string `json:"firstname" form:"firstname" valid:"required"`
	Lastname  string `json:"lastname" form:"lastname" valid:"required"`
	Username  string `json:"username" form:"username" valid:"required"`
	Email     string `json:"email" form:"email" valid:"required,email"`
	Password  string `json:"password" form:"password" valid:"required"`
}

type userCreateResponse struct {
	User *backapi.User `json:"user"`
}

func (h *UserHandler) CreateUser(c echo.Context) error {
	ctx := context.Background()

	u := new(userCreatePayload)
	if err := c.Bind(u); err != nil {
		return c.JSON(http.StatusBadRequest, HandleHTTPResponse(http.StatusBadRequest, _messageBindError, nil))
	}

	u.Firstname = strings.TrimSpace(u.Firstname)
	u.Lastname = strings.TrimSpace(u.Lastname)
	u.Username = strings.TrimSpace(u.Username)
	u.Email = strings.TrimSpace(strings.ToLower(u.Email))
	u.Password = strings.TrimSpace(u.Password)

	validator := validator.NewValidator()
	if errs := validator.Validate(u); errs != nil {
		return c.JSON(http.StatusBadRequest, HandleHTTPResponse(http.StatusBadRequest, _messageValidationError, nil))
	}

	user, err := h.UserService.CreateUser(ctx, &backapi.User{
		Email:     u.Email,
		Password:  u.Password,
		Firstname: u.Firstname,
		Lastname:  u.Lastname,
		Username:  u.Username,
	})
	if err != nil {
		if backapi.IsResourceAlreadyCreatedError(err) {
			return c.JSON(http.StatusConflict, HandleHTTPResponse(http.StatusConflict, _messageUserAlreadyCreatedError, nil))
		}
		return c.JSON(http.StatusInternalServerError, HandleHTTPResponse(http.StatusInternalServerError, _messageInternalServerError, nil))
	}

	resp := &userCreateResponse{
		User: user,
	}
	return c.JSON(http.StatusCreated, HandleHTTPResponse(http.StatusCreated, _messageSuccess, resp))
}

type userGetCurrentResponse struct {
	User *backapi.User `json:"user"`
}

func (h *UserHandler) GetCUser(c echo.Context) error {
	// // Not required to verify it as the endpoint is only accessible with auth
	CUser := c.Get(ContextKeyCurrentUser).(*backapi.User)

	user := &userGetCurrentResponse{
		User: CUser,
	}
	return c.JSON(http.StatusOK, HandleHTTPResponse(http.StatusOK, _messageSuccess, user))
}

type userChangePasswordPayload struct {
	CurrentPassword string `json:"current_password" form:"current_password" valid:"required"`
	NewPassword     string `json:"new_password" form:"new_password" valid:"required"`
}

type userChangePasswordResponse struct {
	User *backapi.User `json:"user"`
}

func (h *UserHandler) UserChangePassword(c echo.Context) error {
	ctx := context.Background()

	CUser := c.Get(ContextKeyCurrentUser).(*backapi.User)

	u := new(userChangePasswordPayload)
	if err := c.Bind(u); err != nil {
		return c.JSON(http.StatusBadRequest, HandleHTTPResponse(http.StatusBadRequest, _messageBindError, nil))
	}

	validator := validator.NewValidator()
	if errs := validator.Validate(u); errs != nil {
		return c.JSON(http.StatusBadRequest, HandleHTTPResponse(http.StatusBadRequest, _messageValidationError, nil))
	}

	userModified, err := h.UserService.ChangePassword(ctx, CUser, u.CurrentPassword, u.NewPassword)
	if err != nil {
		if backapi.IsInternalServerError(err) {
			return c.JSON(http.StatusInternalServerError, HandleHTTPResponse(http.StatusInternalServerError, _messageInternalServerError, nil))
		}
		if backapi.IsUnauthorizedError(err) {
			return c.JSON(http.StatusUnauthorized, HandleHTTPResponse(http.StatusUnauthorized, _messageUnauthorizedError, nil))
		}
		return c.JSON(http.StatusBadRequest, HandleHTTPResponse(http.StatusBadRequest, _messageGenericBadRequestError, nil))
	}

	user := &userChangePasswordResponse{
		User: userModified,
	}
	return c.JSON(http.StatusOK, HandleHTTPResponse(http.StatusOK, _messageSuccess, user))
}

type userChangeEmailPayload struct {
	CurrentPassword string `json:"current_password" form:"current_password" valid:"required"`
	NewEmail        string `json:"new_email" form:"new_email" valid:"required"`
}

type userChangeEmailResponse struct {
	User *backapi.User `json:"user"`
}

func (h *UserHandler) UserChangeEmail(c echo.Context) error {
	ctx := context.Background()

	CUser := c.Get(ContextKeyCurrentUser).(*backapi.User)

	u := new(userChangeEmailPayload)
	if err := c.Bind(u); err != nil {
		return c.JSON(http.StatusBadRequest, HandleHTTPResponse(http.StatusBadRequest, _messageBindError, nil))
	}

	validator := validator.NewValidator()
	if errs := validator.Validate(u); errs != nil {
		return c.JSON(http.StatusBadRequest, HandleHTTPResponse(http.StatusBadRequest, _messageValidationError, nil))
	}

	userModified, err := h.UserService.ChangeEmail(ctx, CUser, u.CurrentPassword, u.NewEmail)
	if err != nil {
		if backapi.IsInternalServerError(err) {
			return c.JSON(http.StatusInternalServerError, HandleHTTPResponse(http.StatusInternalServerError, _messageInternalServerError, nil))
		}
		if backapi.IsUnauthorizedError(err) {
			return c.JSON(http.StatusUnauthorized, HandleHTTPResponse(http.StatusUnauthorized, _messageUnauthorizedError, nil))
		}
		return c.JSON(http.StatusBadRequest, HandleHTTPResponse(http.StatusBadRequest, _messageGenericBadRequestError, nil))
	}

	user := &userChangeEmailResponse{
		User: userModified,
	}
	return c.JSON(http.StatusOK, HandleHTTPResponse(http.StatusOK, _messageSuccess, user))
}
