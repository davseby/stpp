package apierr

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type _mockResposeWriter struct {
	headerFn    func() http.Header
	headerCalls int

	writeFn    func([]byte) (int, error)
	writeCalls int

	writeHeaderFn    func(int)
	writeHeaderCalls int
}

func (mrw *_mockResposeWriter) Header() http.Header {
	mrw.headerCalls++
	return mrw.headerFn()
}

func (mrw *_mockResposeWriter) Write(data []byte) (int, error) {
	mrw.writeCalls++
	return mrw.writeFn(data)
}

func (mrw *_mockResposeWriter) WriteHeader(statusCode int) {
	mrw.writeHeaderCalls++
	mrw.writeHeaderFn(statusCode)
}

func Test_Error_Respond(t *testing.T) {
	stubResponseWriterMock := func(expData []byte, expStatusCode int) *_mockResposeWriter {
		return &_mockResposeWriter{
			writeFn: func(data []byte) (int, error) {
				assert.Equal(t, expData, data)
				return 0, nil
			},
			writeHeaderFn: func(statusCode int) {
				assert.Equal(t, expStatusCode, statusCode)
			},
		}
	}

	t.Run("without message", func(t *testing.T) {
		t.Parallel()

		mrw := stubResponseWriterMock(nil, http.StatusTeapot)

		(&Error{
			statusCode: http.StatusTeapot,
		}).Respond(mrw)

		assert.Equal(t, 1, mrw.writeHeaderCalls)
		assert.Equal(t, 0, mrw.writeCalls)
		assert.Equal(t, 0, mrw.headerCalls)
	})

	t.Run("with message", func(t *testing.T) {
		t.Parallel()

		mrw := stubResponseWriterMock([]byte("test21"), http.StatusTeapot)

		(&Error{
			statusCode: http.StatusTeapot,
			message:    "test21",
		}).Respond(mrw)

		assert.Equal(t, 1, mrw.writeHeaderCalls)
		assert.Equal(t, 1, mrw.writeCalls)
		assert.Equal(t, 0, mrw.headerCalls)
	})
}

func Test_Unauthorized(t *testing.T) {
	assert.Equal(
		t,
		&Error{
			statusCode: http.StatusUnauthorized,
		},
		Unauthorized(),
	)
}

func Test_Forbidden(t *testing.T) {
	assert.Equal(
		t,
		&Error{
			statusCode: http.StatusForbidden,
		},
		Forbidden(),
	)
}

func Test_Conflict(t *testing.T) {
	assert.Equal(
		t,
		&Error{
			statusCode: http.StatusConflict,
			message:    "132",
		},
		Conflict("132"),
	)
}

func Test_Context(t *testing.T) {
	assert.Equal(
		t,
		&Error{
			statusCode: http.StatusBadRequest,
		},
		Context(),
	)
}

func Test_NotFound(t *testing.T) {
	assert.Equal(
		t,
		&Error{
			statusCode: http.StatusNotFound,
			message:    "test111",
		},
		NotFound("test111"),
	)
}

func Test_Database(t *testing.T) {
	assert.Equal(
		t,
		&Error{
			statusCode: http.StatusServiceUnavailable,
		},
		Database(),
	)
}

func Test_Internal(t *testing.T) {
	assert.Equal(
		t,
		&Error{
			statusCode: http.StatusInternalServerError,
		},
		Internal(),
	)
}

func Test_InvalidAttribute(t *testing.T) {
	assert.Equal(
		t,
		&Error{
			statusCode: http.StatusBadRequest,
			message:    "11: 22",
		},
		InvalidAttribute("11", "22"),
	)
}

func Test_BadRequest(t *testing.T) {
	assert.Equal(
		t,
		&Error{
			statusCode: http.StatusBadRequest,
			message:    "1122",
		},
		BadRequest("1122"),
	)
}

func Test_MalformedDataInput(t *testing.T) {
	tests := map[string]struct {
		DataType DataType
		Result   *Error
	}{
		"Invalid data type": {
			Result: &Error{
				statusCode: http.StatusBadRequest,
			},
		},
		"Request body data type": {
			DataType: DataTypeRequestBody,
			Result: &Error{
				statusCode: http.StatusBadRequest,
				message:    "invalid body",
			},
		},
		"JSON data type": {
			DataType: DataTypeJSON,
			Result: &Error{
				statusCode: http.StatusBadRequest,
				message:    "malformed json",
			},
		},
	}

	for name, test := range tests {
		test := test

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(
				t,
				test.Result,
				MalformedDataInput(test.DataType),
			)
		})
	}
}
