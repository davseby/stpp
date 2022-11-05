package server

import (
	"bytes"
	"foodie/core"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/rs/xid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Server_Register(t *testing.T) {
	stubAuthorizer := func(token []byte, err error) *AuthorizerMock {
		return &AuthorizerMock{
			IssueFunc: func(_ xid.ID, _ bool, _ time.Time) ([]byte, error) {
				return token, err
			},
		}
	}

	dbh := dbFn(t)
	cleanUpTables(t, dbh)

	usr := core.User{
		ID:           xid.New(),
		Name:         "1321",
		PasswordHash: []byte{1, 2, 3, 4, 6},
		CreatedAt:    time.Now(),
	}

	_, err := squirrel.ExecWith(
		dbh,
		squirrel.Insert("users").SetMap(map[string]interface{}{
			"users.id":            usr.ID,
			"users.name":          usr.Name,
			"users.password_hash": usr.PasswordHash,
			"users.admin":         usr.Admin,
			"users.created_at":    usr.CreatedAt,
		}),
	)
	require.NoError(t, err)

	tests := map[string]struct {
		Auth       *AuthorizerMock
		Body       []byte
		Response   string
		StatusCode int
	}{
		"Invalid body JSON": {
			Auth:       stubAuthorizer([]byte{1}, nil),
			Body:       []byte("{"),
			Response:   "malformed json",
			StatusCode: http.StatusBadRequest,
		},
		"Invalid user input": {
			Auth:       stubAuthorizer([]byte{1}, nil),
			Body:       []byte("{}"),
			Response:   "name: cannot be empty",
			StatusCode: http.StatusBadRequest,
		},
		"Conflicting user": {
			Auth:       stubAuthorizer([]byte{1}, nil),
			Body:       []byte(`{"name":"1321", "password":"1248abc"}`),
			Response:   "user",
			StatusCode: http.StatusConflict,
		},
		"Authorizer returns an error": {
			Auth:       stubAuthorizer([]byte{1}, assert.AnError),
			Body:       []byte(`{"name":"1325", "password":"1248abc"}`),
			Response:   "",
			StatusCode: http.StatusInternalServerError,
		},
		"Successfully registred a user": {
			Auth:       stubAuthorizer([]byte{1}, nil),
			Body:       []byte(`{"name":"1323", "password":"1248abc"}`),
			Response:   `"name":"1323"`,
			StatusCode: http.StatusOK,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			req := httptest.NewRequest(
				http.MethodPost,
				"http://test.com/123",
				bytes.NewBuffer(test.Body),
			)

			resp := httptest.NewRecorder()

			server := &Server{
				log:  logrus.New(),
				db:   dbh,
				auth: test.Auth,
			}

			server.Register(resp, req)
			assert.Equal(t, test.StatusCode, resp.Code)
			assert.Regexp(t, test.Response, string(resp.Body.Bytes()))
		})
	}
}
