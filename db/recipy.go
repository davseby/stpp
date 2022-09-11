package db

import (
	"context"
	"encoding/json"
	"foodie/core"

	"github.com/Masterminds/squirrel"
	"github.com/rs/xid"
)

func GetRecipies(ctx context.Context, qc squirrel.QueryerContext) ([]core.Recipy, error) {
	return selectRecipies(
		ctx,
		qc,
		func(sb squirrel.SelectBuilder) squirrel.SelectBuilder {
			return sb
		},
	)
}

func GetRecipyByID(
	ctx context.Context,
	qc squirrel.QueryerContext,
	id xid.ID,
) (*core.Recipy, error) {

	recipies, err := selectRecipies(
		ctx,
		qc,
		func(sb squirrel.SelectBuilder) squirrel.SelectBuilder {
			return sb.Where(
				squirrel.Eq{"recipy.id": id},
			)
		},
	)
	if err != nil {
		return nil, err
	}

	if len(recipies) == 0 {
		return nil, ErrNotFound
	}

	return &recipies[0], nil
}

func InsertRecipy(
	ctx context.Context,
	ec squirrel.ExecerContext,
	uid xid.ID,
	rc core.RecipyCore,
) (*core.Recipy, error) {

	recipy := core.Recipy{
		ID:         xid.New(),
		UserID:     uid,
		RecipyCore: rc,
	}

	data, err := json.Marshal(rc.Items)
	if err != nil {
		return nil, err
	}

	_, err = squirrel.ExecContextWith(
		ctx,
		ec,
		squirrel.Insert("recipy").SetMap(map[string]interface{}{
			"recipy.id":          recipy.ID,
			"recipy.user_id":     recipy.UserID,
			"recipy.name":        recipy.Name,
			"recipy.description": recipy.Description,
			"recipy.items":       data,
		}),
	)
	if err != nil {
		return nil, err
	}

	return &recipy, nil
}

func UpdateRecipyByID(
	ctx context.Context,
	ec squirrel.ExecerContext,
	id xid.ID,
	rc core.RecipyCore,
) error {

	data, err := json.Marshal(rc.Items)
	if err != nil {
		return err
	}

	_, err = squirrel.ExecContextWith(
		ctx,
		ec,
		squirrel.Update("recipy").SetMap(map[string]interface{}{
			"recipy.name":        rc.Name,
			"recipy.description": rc.Description,
			"recipy.items":       data,
		}).Where(
			squirrel.Eq{"recipy.id": id},
		),
	)
	return err
}

func DeleteRecipyByID(
	ctx context.Context,
	ec squirrel.ExecerContext,
	id xid.ID,
) error {

	_, err := squirrel.ExecContextWith(
		ctx,
		ec,
		squirrel.Delete("recipy").Where(
			squirrel.Eq{"recipy.id": id},
		),
	)
	return err
}

func selectRecipies(
	ctx context.Context,
	qc squirrel.QueryerContext,
	dec func(squirrel.SelectBuilder) squirrel.SelectBuilder,
) ([]core.Recipy, error) {

	rows, err := squirrel.QueryContextWith(ctx, qc, dec(squirrel.
		Select(
			"recipy.id",
			"recipy.user_id",
			"recipy.name",
			"recipy.description",
			"recipy.items",
		).From("recipy"),
	))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	recipies := make([]core.Recipy, 0)
	for rows.Next() {
		var (
			recipy core.Recipy
			data   []byte
		)

		if err := rows.Scan(
			&recipy.ID,
			&recipy.UserID,
			&recipy.Name,
			&recipy.Description,
			&data,
		); err != nil {
			return nil, err
		}

		if err := json.Unmarshal(data, &recipy.Items); err != nil {
			return nil, err
		}

		recipies = append(recipies, recipy)
	}

	return recipies, nil
}
