package db

import (
	"context"
	"database/sql"
	"foodie/core"
	"testing"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_InsertUser(t *testing.T) {
	dbh := dbFn(t)
	cleanUpTables(t, dbh)

	t.Run("database returns an error", func(t *testing.T) {
		usr, err := InsertUser(
			context.Background(),
			dbh,
			"333",
			nil,
			true,
		)
		assert.Nil(t, usr)
		assert.Error(t, err)
	})

	t.Run("successful new user insert", func(t *testing.T) {
		usr, err := InsertUser(
			context.Background(),
			dbh,
			"333",
			[]byte{1},
			true,
		)
		require.NoError(t, err)
		assert.NotEmpty(t, usr.ID)
		assert.Equal(t, "333", usr.Name)
		assert.Equal(t, []byte{1}, usr.PasswordHash)
		assert.NotEmpty(t, usr.CreatedAt)
		assert.True(t, usr.Admin)
	})
}

func Test_GetUsers(t *testing.T) {
	dbh := dbFn(t)
	cleanUpTables(t, dbh)

	uu := []core.User{
		{
			ID:           xid.New(),
			Name:         "1",
			PasswordHash: []byte{1},
			CreatedAt:    time.Now().UTC().Truncate(time.Second),
			Admin:        true,
		},
		{
			ID:           xid.New(),
			Name:         "2",
			PasswordHash: []byte{2},
			CreatedAt:    time.Now().UTC().Truncate(time.Second),
			Admin:        false,
		},
		{
			ID:           xid.New(),
			Name:         "4",
			PasswordHash: []byte{5},
			CreatedAt:    time.Now().UTC().Truncate(time.Second),
			Admin:        false,
		},
	}

	for _, usr := range uu {
		mockUser(t, dbh, usr)
	}

	res, err := GetUsers(context.Background(), dbh)
	require.NoError(t, err)
	assert.Equal(t, uu, res)
}

func Test_GetUserByName(t *testing.T) {
	dbh := dbFn(t)
	cleanUpTables(t, dbh)

	uu := []core.User{
		{
			ID:           xid.New(),
			Name:         "1",
			PasswordHash: []byte{1},
			CreatedAt:    time.Now().UTC().Truncate(time.Second),
			Admin:        true,
		},
		{
			ID:           xid.New(),
			Name:         "2",
			PasswordHash: []byte{2},
			CreatedAt:    time.Now().UTC().Truncate(time.Second),
			Admin:        false,
		},
		{
			ID:           xid.New(),
			Name:         "4",
			PasswordHash: []byte{5},
			CreatedAt:    time.Now().UTC().Truncate(time.Second),
			Admin:        false,
		},
	}

	for _, usr := range uu {
		mockUser(t, dbh, usr)
	}

	t.Run("not found", func(t *testing.T) {
		res, err := GetUserByName(context.Background(), dbh, "3")
		assert.Empty(t, res)
		require.Equal(t, ErrNotFound, err)
	})

	t.Run("successfully retrieved a user by name", func(t *testing.T) {
		res, err := GetUserByName(context.Background(), dbh, "2")
		require.NoError(t, err)
		assert.Equal(t, &uu[1], res)
	})
}

func Test_GetUserByID(t *testing.T) {
	dbh := dbFn(t)
	cleanUpTables(t, dbh)

	uu := []core.User{
		{
			ID:           xid.New(),
			Name:         "1",
			PasswordHash: []byte{1},
			CreatedAt:    time.Now().UTC().Truncate(time.Second),
			Admin:        true,
		},
		{
			ID:           xid.New(),
			Name:         "2",
			PasswordHash: []byte{2},
			CreatedAt:    time.Now().UTC().Truncate(time.Second),
			Admin:        false,
		},
		{
			ID:           xid.New(),
			Name:         "4",
			PasswordHash: []byte{5},
			CreatedAt:    time.Now().UTC().Truncate(time.Second),
			Admin:        false,
		},
	}

	for _, usr := range uu {
		mockUser(t, dbh, usr)
	}

	t.Run("not found", func(t *testing.T) {
		res, err := GetUserByID(context.Background(), dbh, xid.New())
		assert.Empty(t, res)
		require.Equal(t, ErrNotFound, err)
	})

	t.Run("successfully retrieved a user by name", func(t *testing.T) {
		res, err := GetUserByID(context.Background(), dbh, uu[1].ID)
		require.NoError(t, err)
		assert.Equal(t, &uu[1], res)
	})
}

func Test_UpdateUserPasswordByID(t *testing.T) {
	dbh := dbFn(t)
	cleanUpTables(t, dbh)

	usr := core.User{
		ID:           xid.New(),
		Name:         "4",
		PasswordHash: []byte{5},
		CreatedAt:    time.Now().UTC().Truncate(time.Second),
		Admin:        false,
	}

	mockUser(t, dbh, usr)

	err := UpdateUserPasswordByID(context.Background(), dbh, usr.ID, []byte{2})
	require.NoError(t, err)

	uu := retrieveUsers(t, dbh)
	require.Len(t, uu, 1)

	usr.PasswordHash = []byte{2}
	assert.Equal(t, usr, uu[0])
}

func Test_DeleteUserByID(t *testing.T) {
	dbh := dbFn(t)
	cleanUpTables(t, dbh)

	usr := core.User{
		ID:           xid.New(),
		Name:         "4",
		PasswordHash: []byte{5},
		CreatedAt:    time.Now().UTC().Truncate(time.Second),
		Admin:        false,
	}

	mockUser(t, dbh, usr)

	uu := retrieveUsers(t, dbh)
	require.Len(t, uu, 1)
	assert.Equal(t, usr, uu[0])

	require.NoError(t, DeleteUserByID(context.Background(), dbh, usr.ID))
	require.Len(t, retrieveUsers(t, dbh), 0)
}

func Test_selectUsers(t *testing.T) {
	dbh := dbFn(t)
	cleanUpTables(t, dbh)

	uu := []core.User{
		{
			ID:           xid.New(),
			Name:         "1",
			PasswordHash: []byte{1},
			CreatedAt:    time.Now().UTC().Truncate(time.Second),
			Admin:        true,
		},
		{
			ID:           xid.New(),
			Name:         "2",
			PasswordHash: []byte{2},
			CreatedAt:    time.Now().UTC().Truncate(time.Second),
			Admin:        false,
		},
		{
			ID:           xid.New(),
			Name:         "4",
			PasswordHash: []byte{5},
			CreatedAt:    time.Now().UTC().Truncate(time.Second),
			Admin:        false,
		},
	}

	for _, usr := range uu {
		mockUser(t, dbh, usr)
	}

	res, err := selectUsers(context.Background(), dbh, func(sb squirrel.SelectBuilder) squirrel.SelectBuilder {
		return sb
	})
	require.NoError(t, err)
	assert.Equal(t, retrieveUsers(t, dbh), res)
}

func mockUser(t *testing.T, dbh *sql.DB, usr core.User) {
	t.Helper()

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
}

func retrieveUsers(t *testing.T, dbh *sql.DB) []core.User {
	rows, err := squirrel.QueryWith(dbh, squirrel.
		Select(
			"users.id",
			"users.name",
			"users.password_hash",
			"users.admin",
			"users.created_at",
		).From("users"),
	)
	require.NoError(t, err)
	defer rows.Close()

	users := make([]core.User, 0)
	for rows.Next() {
		var user core.User
		require.NoError(t, rows.Scan(
			&user.ID,
			&user.Name,
			&user.PasswordHash,
			&user.Admin,
			&user.CreatedAt,
		))

		users = append(users, user)
	}

	return users
}
