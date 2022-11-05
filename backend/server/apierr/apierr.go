package apierr

import (
	"fmt"
	"net/http"
)

// DataType specifies the data type.
type DataType int

const (
	// DataTypeRequestBody specifies that the data type is the request body.
	DataTypeRequestBody DataType = iota + 1

	// DataTypeJSON specifies that the data type is JSON.
	DataTypeJSON
)

// Error contains http error information.
type Error struct {
	statusCode int
	message    string
}

// Respond responds to the response writer with the provided error.
func (e *Error) Respond(w http.ResponseWriter) {
	w.WriteHeader(e.statusCode)
	if e.message != "" {
		w.Write([]byte(e.message)) //nolint:errcheck // cannot recover from this.
	}
}

// Unauthorized creates a new unauthorized error.
func Unauthorized() *Error {
	return &Error{
		statusCode: http.StatusUnauthorized,
	}
}

// Forbidden creates a new forbidden error.
func Forbidden() *Error {
	return &Error{
		statusCode: http.StatusForbidden,
	}
}

// Conflict creates a new conflict error.
func Conflict(object string) *Error {
	return &Error{
		statusCode: http.StatusConflict,
		message:    object,
	}
}

// Context creates a new context error.
func Context() *Error {
	return &Error{
		statusCode: http.StatusBadRequest,
	}
}

// NotFound creates a new not found error.
func NotFound(object string) *Error {
	return &Error{
		statusCode: http.StatusNotFound,
		message:    object,
	}
}

// Database creates a new database error.
func Database() *Error {
	return &Error{
		statusCode: http.StatusServiceUnavailable,
	}
}

// Internal creates a new internal error.
func Internal() *Error {
	return &Error{
		statusCode: http.StatusInternalServerError,
	}
}

// InvalidAttribute creates a new attribute error.
func InvalidAttribute(attribute, message string) *Error {
	return &Error{
		statusCode: http.StatusBadRequest,
		message:    fmt.Sprintf("%s: %s", attribute, message),
	}
}

// BadRequest creates a new bad request error.
func BadRequest(message string) *Error {
	return &Error{
		statusCode: http.StatusBadRequest,
		message:    message,
	}
}

// MalformedDataInput creates a new malformed error.
func MalformedDataInput(dt DataType) *Error {
	var message string
	switch dt {
	case DataTypeRequestBody:
		message = "invalid body"
	case DataTypeJSON:
		message = "malformed json"
	default:
	}

	return &Error{
		statusCode: http.StatusBadRequest,
		message:    message,
	}
}
