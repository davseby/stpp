package db

import (
	"context"
	"database/sql"
	"foodie/core"

	"github.com/Masterminds/squirrel"
	"github.com/rs/xid"
)

// InsertRecipy inserts a new recipy into the database.
func InsertRecipy(
	ctx context.Context,
	db *sql.DB,
	uid xid.ID,
	rc core.RecipyCore,
) (*core.Recipy, error) {

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	rec := core.Recipy{
		ID:         xid.New(),
		UserID:     uid,
		RecipyCore: rc,
	}

	_, err = squirrel.ExecContextWith(
		ctx,
		tx,
		squirrel.Insert("recipy").SetMap(map[string]interface{}{
			"recipy.id":          rec.ID,
			"recipy.user_id":     rec.UserID,
			"recipy.name":        rec.Name,
			"recipy.private":     rec.Private,
			"recipy.description": rec.Description,
		}),
	)
	if err != nil {
		return nil, err
	}

	for _, rp := range rc.Products {
		rp.RecipyID = rec.ID

		if err := upsertRecipyProduct(
			ctx,
			tx,
			rp,
		); err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &rec, nil
}

// GetRecipes retrieves all recipes. The ip paratemer
// specifies whether the private recipes should be retrieved.
func GetRecipes(
	ctx context.Context,
	qc squirrel.QueryerContext,
	ip bool,
) ([]core.Recipy, error) {

	return selectRecipes(
		ctx,
		qc,
		func(sb squirrel.SelectBuilder) squirrel.SelectBuilder {
			if !ip {
				return sb.Where(
					squirrel.Eq{"recipy.private": false},
				)
			}

			return sb
		},
	)
}

// GetRecipesByUserID retrieves recipes by the user id. The ip paratemer
// specifies whether the private recipes should be retrieved.
func GetRecipesByUserID(
	ctx context.Context,
	qc squirrel.QueryerContext,
	uid xid.ID,
	ip bool,
) ([]core.Recipy, error) {

	return selectRecipes(
		ctx,
		qc,
		func(sb squirrel.SelectBuilder) squirrel.SelectBuilder {
			if !ip {
				sb = sb.Where(
					squirrel.Eq{"recipy.private": false},
				)
			}

			return sb.Where(
				squirrel.Eq{"recipy.user_id": uid},
			)
		},
	)
}

// GetRecipyByID retrieves a recipy by its id. The ip paratemer
// specifies whether the private recipes should be retrieved.
func GetRecipyByID(
	ctx context.Context,
	qc squirrel.QueryerContext,
	id xid.ID,
	ip bool,
) (*core.Recipy, error) {

	rr, err := selectRecipes(
		ctx,
		qc,
		func(sb squirrel.SelectBuilder) squirrel.SelectBuilder {
			if !ip {
				sb = sb.Where(
					squirrel.Eq{"recipy.private": false},
				)
			}

			return sb.Where(
				squirrel.Eq{"recipy.id": id},
			)
		},
	)
	if err != nil {
		return nil, err
	}

	if len(rr) == 0 {
		return nil, ErrNotFound
	}

	return &rr[0], nil
}

// UpdateRecipyByID updates an existing recipy by its id. An updated recipy
// is returned.
func UpdateRecipyByID(
	ctx context.Context,
	db *sql.DB,
	id xid.ID,
	rc core.RecipyCore,
) (*core.Recipy, error) {

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, nil
	}
	defer tx.Rollback()

	if err := deleteRecipyProducts(
		ctx,
		tx,
		id,
	); err != nil {
		return nil, err
	}

	for _, rp := range rc.Products {
		rp.RecipyID = id

		if err := upsertRecipyProduct(
			ctx,
			tx,
			rp,
		); err != nil {
			return nil, err
		}
	}

	_, err = squirrel.ExecContextWith(
		ctx,
		tx,
		squirrel.Update("recipy").SetMap(map[string]interface{}{
			"recipy.name":        rc.Name,
			"recipy.private":     rc.Private,
			"recipy.description": rc.Description,
		}).Where(
			squirrel.Eq{"recipy.id": id},
		),
	)

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	rec, err := GetRecipyByID(ctx, db, id, true)
	if err != nil {
		return nil, err
	}

	return rec, nil
}

// DeleteRecipyByID deletes a recipy by its id.
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

// selectRecipes selects all recipes by the provided decorator function.
func selectRecipes(
	ctx context.Context,
	qc squirrel.QueryerContext,
	dec func(squirrel.SelectBuilder) squirrel.SelectBuilder,
) ([]core.Recipy, error) {

	rows, err := squirrel.QueryContextWith(ctx, qc, dec(squirrel.
		Select(
			"recipy.id",
			"recipy.user_id",
			"recipy.name",
			"recipy.private",
			"recipy.description",
		).From("recipy"),
	))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rr := make([]core.Recipy, 0)
	for rows.Next() {
		var rec core.Recipy
		if err := rows.Scan(
			&rec.ID,
			&rec.UserID,
			&rec.Name,
			&rec.Private,
			&rec.Description,
		); err != nil {
			return nil, err
		}

		rps, err := getRecipyProductsByRecipyID(ctx, qc, rec.ID)
		if err != nil {
			return nil, err
		}

		rec.Products = rps
		rr = append(rr, rec)
	}

	return rr, nil
}

// getRecipyProductsByRecipyID selects recipy products by the recipy id.
func getRecipyProductsByRecipyID(
	ctx context.Context,
	qc squirrel.QueryerContext,
	id xid.ID,
) ([]core.RecipyProduct, error) {

	return selectRecipyProducts(
		ctx,
		qc,
		func(sb squirrel.SelectBuilder) squirrel.SelectBuilder {
			return sb.Where(
				squirrel.Eq{"recipy_product.recipy_id": id},
			)
		},
	)
}

// deleteRecipyProducts deletes all recipy products.
func deleteRecipyProducts(
	ctx context.Context,
	ec squirrel.ExecerContext,
	rid xid.ID,
) error {

	_, err := squirrel.ExecContextWith(
		ctx,
		ec,
		squirrel.Delete("recipy_product").Where(
			squirrel.Eq{"recipy_product.recipy_id": rid},
		),
	)
	return err
}

// upsertRecipyProduct upserts recipy products.
func upsertRecipyProduct(
	ctx context.Context,
	ec squirrel.ExecerContext,
	rp core.RecipyProduct,
) error {

	_, err := squirrel.ExecContextWith(
		ctx,
		ec,
		squirrel.Insert("recipy_product").SetMap(map[string]interface{}{
			"recipy_product.recipy_id":  rp.RecipyID,
			"recipy_product.product_id": rp.ProductID,
			"recipy_product.quantity":   rp.Quantity,
		}).Suffix("ON DUPLICATE KEY UPDATE recipy_product.quantity = VALUES(recipy_product.quantity)"),
	)
	return err
}

// selectRecipyProducts selects all recipy products by the provided decorator
// function.
func selectRecipyProducts(
	ctx context.Context,
	qc squirrel.QueryerContext,
	dec func(squirrel.SelectBuilder) squirrel.SelectBuilder,
) ([]core.RecipyProduct, error) {

	rows, err := squirrel.QueryContextWith(ctx, qc, dec(squirrel.
		Select(
			"recipy_product.recipy_id",
			"recipy_product.product_id",
			"recipy_product.quantity",
		).From("recipy_product"),
	))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rps := make([]core.RecipyProduct, 0)
	for rows.Next() {
		var rp core.RecipyProduct
		if err := rows.Scan(
			&rp.RecipyID,
			&rp.ProductID,
			&rp.Quantity,
		); err != nil {
			return nil, err
		}

		rps = append(rps, rp)
	}

	return rps, nil
}
