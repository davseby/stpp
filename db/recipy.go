package db

import (
	"context"
	"database/sql"
	"foodie/core"

	"github.com/Masterminds/squirrel"
	"github.com/go-sql-driver/mysql"
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
	db *sql.DB,
	uid xid.ID,
	rc core.RecipyCore,
) (*core.Recipy, error) {

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	recipy := core.Recipy{
		ID:         xid.New(),
		UserID:     uid,
		RecipyCore: rc,
	}

	_, err = squirrel.ExecContextWith(
		ctx,
		tx,
		squirrel.Insert("recipy").SetMap(map[string]interface{}{
			"recipy.id":          recipy.ID,
			"recipy.user_id":     recipy.UserID,
			"recipy.name":        recipy.Name,
			"recipy.description": recipy.Description,
		}),
	)
	if err != nil {
		return nil, err
	}

	for _, rp := range rc.Products {
		if err := upsertRecipyProduct(
			ctx,
			tx,
			rp,
		); err != nil {
			if merr, ok := err.(*mysql.MySQLError); ok && merr.Number == 1452 {
				return nil, ErrNotFound
			}

			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &recipy, nil
}

func UpdateRecipyByID(
	ctx context.Context,
	db *sql.DB,
	id xid.ID,
	rc core.RecipyCore,
) error {

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil
	}
	defer tx.Rollback()

	oldRecipyProducts, err := getRecipyProductsByRecipyID(ctx, tx, id)
	if err != nil {
		return err
	}

	for _, recipyProduct := range rc.Products {
		if err := upsertRecipyProduct(
			ctx,
			tx,
			recipyProduct,
		); err != nil {
			if merr, ok := err.(*mysql.MySQLError); ok && merr.Number == 1452 {
				return ErrNotFound
			}

			return err
		}
	}

	for _, oldRecipyProduct := range oldRecipyProducts {
		var found bool
		for _, recipyProduct := range rc.Products {
			if recipyProduct.ProductID.Compare(oldRecipyProduct.ProductID) == 0 {
				found = true
			}
		}

		if !found {
			if err := deleteRecipyProduct(
				ctx,
				tx,
				oldRecipyProduct,
			); err != nil {
				return err
			}
		}
	}

	_, err = squirrel.ExecContextWith(
		ctx,
		tx,
		squirrel.Update("recipy").SetMap(map[string]interface{}{
			"recipy.name":        rc.Name,
			"recipy.description": rc.Description,
		}).Where(
			squirrel.Eq{"recipy.id": id},
		),
	)

	if err := tx.Commit(); err != nil {
		return nil
	}

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
		var recipy core.Recipy
		if err := rows.Scan(
			&recipy.ID,
			&recipy.UserID,
			&recipy.Name,
			&recipy.Description,
		); err != nil {
			return nil, err
		}

		recipyProducts, err := getRecipyProductsByRecipyID(ctx, qc, recipy.ID)
		if err != nil {
			return nil, err
		}

		recipy.Products = recipyProducts
		recipies = append(recipies, recipy)
	}

	return recipies, nil
}

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

func deleteRecipyProduct(
	ctx context.Context,
	ec squirrel.ExecerContext,
	rp core.RecipyProduct,
) error {

	_, err := squirrel.ExecContextWith(
		ctx,
		ec,
		squirrel.Delete("recipy_product").Where(
			squirrel.Eq{"recipy_product.recipy_id": rp.RecipyID},
			squirrel.Eq{"recipy_product.product_id": rp.ProductID},
		),
	)
	return err
}

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

	recipyProducts := make([]core.RecipyProduct, 0)
	for rows.Next() {
		var recipyProduct core.RecipyProduct
		if err := rows.Scan(
			&recipyProduct.RecipyID,
			&recipyProduct.ProductID,
			&recipyProduct.Quantity,
		); err != nil {
			return nil, err
		}

		recipyProducts = append(recipyProducts, recipyProduct)
	}

	return recipyProducts, nil
}
