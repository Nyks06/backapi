package http

import "net/http"

const (
	_messageSuccess                 = "Success"
	_messageBindError               = "Unable to bind data - malformated parameters"
	_messageValidationError         = "Unable to validate data"
	_messageInternalServerError     = "Internal Server Error - please try again later"
	_messageUserAlreadyCreatedError = "Conflict - User Already Exist"
	_messageGenericBadRequestError  = "Bad Request Error - please confirm your request and retry"
	_messageUnauthorizedError       = "Unauthorized Error - you are not authorized to do that"
)

type HTTPResponseStatus struct {
	Error   bool   `json:"error"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type HTTPResponse struct {
	Status HTTPResponseStatus `json:"status"`
	Data   interface{}        `json:"data"`
}

func HandleHTTPResponse(statusCode int, message string, data interface{}) HTTPResponse {
	resp := HTTPResponse{
		Status: HTTPResponseStatus{
			Code:    statusCode,
			Message: message,
		},
		Data: data,
	}

	if statusCode != http.StatusOK && statusCode != http.StatusCreated && statusCode != http.StatusNoContent {
		resp.Status.Error = true
	}

	return resp
}
