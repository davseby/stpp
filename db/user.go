package db

import (
	"context"
	"foodie/core"

	"github.com/Masterminds/squirrel"
	"github.com/rs/xid"
)

func GetUsers(ctx context.Context, qc squirrel.QueryerContext) ([]core.User, error) {
	return selectUsers(
		ctx,
		qc,
		func(sb squirrel.SelectBuilder) squirrel.SelectBuilder {
			return sb
		},
	)
}

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
				squirrel.Eq{"user.name": name},
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
				squirrel.Eq{"user.id": id},
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

func InsertUser(
	ctx context.Context,
	ec squirrel.ExecerContext,
	name string,
	ph []byte,
	adm bool,
) (*core.User, error) {

	user := core.User{
		ID:           xid.New(),
		Name:         name,
		PasswordHash: ph,
		Admin:        adm,
	}

	_, err := squirrel.ExecContextWith(
		ctx,
		ec,
		squirrel.Insert("user").SetMap(map[string]interface{}{
			"user.id":            user.ID,
			"user.name":          user.Name,
			"user.password_hash": user.PasswordHash,
			"user.admin":         user.Admin,
		}),
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func UpdateUserPasswordByID(
	ctx context.Context,
	ec squirrel.ExecerContext,
	id xid.ID,
	ph []byte,
) error {

	_, err := squirrel.ExecContextWith(
		ctx,
		ec,
		squirrel.Update("user").SetMap(map[string]interface{}{
			"user.password_hash": ph,
		}).Where(
			squirrel.Eq{"user.id": id},
		),
	)
	return err
}

func DeleteUserByID(
	ctx context.Context,
	ec squirrel.ExecerContext,
	id xid.ID,
) error {

	_, err := squirrel.ExecContextWith(
		ctx,
		ec,
		squirrel.Delete("user").Where(
			squirrel.Eq{"user.id": id},
		),
	)
	return err
}

func selectUsers(
	ctx context.Context,
	qc squirrel.QueryerContext,
	dec func(squirrel.SelectBuilder) squirrel.SelectBuilder,
) ([]core.User, error) {

	rows, err := squirrel.QueryContextWith(ctx, qc, dec(squirrel.
		Select(
			"user.id",
			"user.name",
			"user.password_hash",
			"user.admin",
		).From("user"),
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
		); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}
