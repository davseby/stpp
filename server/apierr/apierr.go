package apierr

import (
	"fmt"
	"net/http"
)

type Data int

const (
	RequestData Data = iota
	JSONData
)

type Error struct {
	statusCode int
	message    string
}

func (e *Error) Respond(w http.ResponseWriter) {
	w.WriteHeader(e.statusCode)
	if e.message != "" {
		w.Write([]byte(e.message))
	}
}

func Unauthorized() *Error {
	return &Error{
		statusCode: http.StatusUnauthorized,
	}
}

func Forbidden() *Error {
	return &Error{
		statusCode: http.StatusForbidden,
	}
}

func Conflict(object string) *Error {
	return &Error{
		statusCode: http.StatusBadRequest,
		message:    object,
	}
}

func Context() *Error {
	return &Error{
		statusCode: http.StatusBadRequest,
	}
}

func NotFound(object string) *Error {
	return &Error{
		statusCode: http.StatusNotFound,
		message:    object,
	}
}

func Database() *Error {
	return &Error{
		statusCode: http.StatusServiceUnavailable,
	}
}

func Internal() *Error {
	return &Error{
		statusCode: http.StatusInternalServerError,
	}
}

func Attribute(attribute, message string) *Error {
	return &Error{
		statusCode: http.StatusBadRequest,
		message:    fmt.Sprintf("%s: %s", attribute, message),
	}
}

func BadRequest(message string) *Error {
	return &Error{
		statusCode: http.StatusBadRequest,
		message:    message,
	}
}

func DataFormat(dt Data) *Error {
	var message string
	switch dt {
	case RequestData:
		message = "invalid body"
	case JSONData:
		message = "malformed json"
	default:
	}

	return &Error{
		statusCode: http.StatusBadRequest,
		message:    message,
	}
}
