package db

import (
	"context"
	"foodie/core"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/rs/xid"
)

// InsertProduct inserts a new user into the database.
func InsertUser(
	ctx context.Context,
	ec squirrel.ExecerContext,
	name string,
	ph []byte,
	adm bool,
) (*core.User, error) {
	usr := core.User{
		ID:           xid.New(),
		Name:         name,
		PasswordHash: ph,
		CreatedAt:    time.Now(),
		Admin:        adm,
	}

	_, err := squirrel.ExecContextWith(
		ctx,
		ec,
		squirrel.Insert("users").SetMap(map[string]interface{}{
			"users.id":            usr.ID,
			"users.name":          usr.Name,
			"users.password_hash": usr.PasswordHash,
			"users.admin":         usr.Admin,
			"users.created_at":    usr.CreatedAt,
		}),
	)
	if err != nil {
		return nil, err
	}

	return &usr, nil
}

// GetUsers retrieves all users.
func GetUsers(ctx context.Context, qc squirrel.QueryerContext) ([]core.User, error) {
	return selectUsers(
		ctx,
		qc,
		func(sb squirrel.SelectBuilder) squirrel.SelectBuilder {
			return sb
		},
	)
}

// GetUserByName retrieves a user by name.
func GetUserByName(
	ctx context.Context,
	qc squirrel.QueryerContext,
	name string,
) (*core.User, error) {
	users, err := selectUsers(
		ctx,
		qc,
		func(sb squirrel.SelectBuilder) squirrel.SelectBuilder {
			return sb.Where(
				squirrel.Eq{"users.name": name},
			)
		},
	)
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, ErrNotFound
	}

	return &users[0], nil
}

// GetUserByName retrieves a user by its id.
func GetUserByID(
	ctx context.Context,
	qc squirrel.QueryerContext,
	id xid.ID,
) (*core.User, error) {
	users, err := selectUsers(
		ctx,
		qc,
		func(sb squirrel.SelectBuilder) squirrel.SelectBuilder {
			return sb.Where(
				squirrel.Eq{"users.id": id},
			)
		},
	)
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, ErrNotFound
	}

	return &users[0], nil
}

// UpdateUserPasswordByID updates user password by user id.
func UpdateUserPasswordByID(
	ctx context.Context,
	ec squirrel.ExecerContext,
	id xid.ID,
	ph []byte,
) error {
	_, err := squirrel.ExecContextWith(
		ctx,
		ec,
		squirrel.Update("users").SetMap(map[string]interface{}{
			"users.password_hash": ph,
		}).Where(
			squirrel.Eq{"users.id": id},
		),
	)

	return err
}

// DelteUserByID deletes user password by user id.
func DeleteUserByID(
	ctx context.Context,
	ec squirrel.ExecerContext,
	id xid.ID,
) error {
	_, err := squirrel.ExecContextWith(
		ctx,
		ec,
		squirrel.Delete("users").Where(
			squirrel.Eq{"users.id": id},
		),
	)

	return err
}

// selectProducts selects all users by the provided decorator function.
func selectUsers(
	ctx context.Context,
	qc squirrel.QueryerContext,
	dec func(squirrel.SelectBuilder) squirrel.SelectBuilder,
) ([]core.User, error) {
	rows, err := squirrel.QueryContextWith(ctx, qc, dec(squirrel.
		Select(
			"users.id",
			"users.name",
			"users.password_hash",
			"users.admin",
			"users.created_at",
		).From("users"),
	))
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users := make([]core.User, 0)

	for rows.Next() {
		var user core.User

		if err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.PasswordHash,
			&user.Admin,
			&user.CreatedAt,
		); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}
