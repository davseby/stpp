package db

import (
	"context"
	"foodie/core"

	"github.com/Masterminds/squirrel"
	"github.com/go-sql-driver/mysql"
	"github.com/rs/xid"
)

func GetRatings(ctx context.Context, qc squirrel.QueryerContext) ([]core.Rating, error) {
	return selectRatings(
		ctx,
		qc,
		func(sb squirrel.SelectBuilder) squirrel.SelectBuilder {
			return sb
		},
	)
}

func GetRatingByID(
	ctx context.Context,
	qc squirrel.QueryerContext,
	id xid.ID,
) (*core.Rating, error) {

	ratings, err := selectRatings(
		ctx,
		qc,
		func(sb squirrel.SelectBuilder) squirrel.SelectBuilder {
			return sb.Where(
				squirrel.Eq{"rating.id": id},
			)
		},
	)
	if err != nil {
		return nil, err
	}

	if len(ratings) == 0 {
		return nil, ErrNotFound
	}

	return &ratings[0], nil
}

func InsertRating(
	ctx context.Context,
	ec squirrel.ExecerContext,
	rid xid.ID,
	uid xid.ID,
	rc core.RatingCore,
) (*core.Rating, error) {

	rating := core.Rating{
		ID:         xid.New(),
		RecipyID:   rid,
		UserID:     uid,
		RatingCore: rc,
	}

	_, err := squirrel.ExecContextWith(
		ctx,
		ec,
		squirrel.Insert("rating").SetMap(map[string]interface{}{
			"rating.id":        rating.ID,
			"rating.recipy_id": rating.RecipyID,
			"rating.user_id":   rating.UserID,
			"rating.score":     rating.Score,
			"rating.comment":   rating.Comment,
		}),
	)
	if err != nil {
		if merr, ok := err.(*mysql.MySQLError); ok && merr.Number == 1452 {
			return nil, ErrNotFound
		}

		return nil, err
	}

	return &rating, nil
}

func UpdateRatingByID(
	ctx context.Context,
	ec squirrel.ExecerContext,
	id xid.ID,
	rc core.RatingCore,
) error {

	_, err := squirrel.ExecContextWith(
		ctx,
		ec,
		squirrel.Update("rating").SetMap(map[string]interface{}{
			"rating.score":   rc.Score,
			"rating.comment": rc.Comment,
		}).Where(
			squirrel.Eq{"rating.id": id},
		),
	)
	return err
}

func DeleteRatingByID(
	ctx context.Context,
	ec squirrel.ExecerContext,
	id xid.ID,
) error {

	_, err := squirrel.ExecContextWith(
		ctx,
		ec,
		squirrel.Delete("rating").Where(
			squirrel.Eq{"rating.id": id},
		),
	)
	return err
}

func selectRatings(
	ctx context.Context,
	qc squirrel.QueryerContext,
	dec func(squirrel.SelectBuilder) squirrel.SelectBuilder,
) ([]core.Rating, error) {

	rows, err := squirrel.QueryContextWith(ctx, qc, dec(squirrel.
		Select(
			"rating.id",
			"rating.recipy_id",
			"rating.user_id",
			"rating.score",
			"rating.comment",
		).From("rating"),
	))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ratings := make([]core.Rating, 0)
	for rows.Next() {
		var rating core.Rating
		if err := rows.Scan(
			&rating.ID,
			&rating.RecipyID,
			&rating.UserID,
			&rating.Score,
			&rating.Comment,
		); err != nil {
			return nil, err
		}

		ratings = append(ratings, rating)
	}

	return ratings, nil
}
