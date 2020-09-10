package http

import "net/http"

const (
	_messageSuccess                 = "SUCCESS"
	_messageBindError               = "MALFORMATED_PARAMETERS"
	_messageValidationError         = "INVALID_PARAMETERS"
	_messageInternalServerError     = "INTERNAL_SERVER_ERROR"
	_messageUserAlreadyCreatedError = "CONFLICT_ALREADY_EXIST"
	_messageGenericBadRequestError  = "BAD_REQUEST_ERROR"
	_messageUnauthorizedError       = "UNAUTHORIZED_ERROR"
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

	// TODO : Better handling - Currently relying only on the HTTP status code
	// 200, 201, 204 are equivalent to no error
	// all other status code are equivalent to error
	if statusCode != http.StatusOK && statusCode != http.StatusCreated && statusCode != http.StatusNoContent {
		resp.Status.Error = true
	}

	return resp
}
