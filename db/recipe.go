package db

import (
	"context"
	"database/sql"
	"foodie/core"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/rs/xid"
)

// InsertRecipe inserts a new recipe into the database.
func InsertRecipe(
	ctx context.Context,
	db *sql.DB,
	uid xid.ID,
	rc core.RecipeCore,
) (*core.Recipe, error) {

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	rec := core.Recipe{
		ID:         xid.New(),
		UserID:     uid,
		CreatedAt:  time.Now(),
		RecipeCore: rc,
	}

	_, err = squirrel.ExecContextWith(
		ctx,
		tx,
		squirrel.Insert("recipes").SetMap(map[string]interface{}{
			"recipes.id":          rec.ID,
			"recipes.user_id":     rec.UserID,
			"recipes.name":        rec.Name,
			"recipes.description": rec.Description,
			"recipes.created_at":  rec.CreatedAt,
		}),
	)
	if err != nil {
		return nil, err
	}

	for _, rp := range rc.Products {
		rp.RecipeID = rec.ID

		if err := upsertRecipeProduct(
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

// GetRecipes retrieves all recipes.
func GetRecipes(
	ctx context.Context,
	qc squirrel.QueryerContext,
) ([]core.Recipe, error) {

	return selectRecipes(
		ctx,
		qc,
		func(sb squirrel.SelectBuilder) squirrel.SelectBuilder {
			return sb
		},
	)
}

// GetRecipesByUserID retrieves recipes by the user id.
func GetRecipesByUserID(
	ctx context.Context,
	qc squirrel.QueryerContext,
	uid xid.ID,
) ([]core.Recipe, error) {

	return selectRecipes(
		ctx,
		qc,
		func(sb squirrel.SelectBuilder) squirrel.SelectBuilder {
			return sb.Where(
				squirrel.Eq{"recipes.user_id": uid},
			)
		},
	)
}

// GetRecipeByID retrieves a recipe by its id.
func GetRecipeByID(
	ctx context.Context,
	qc squirrel.QueryerContext,
	id xid.ID,
) (*core.Recipe, error) {

	rr, err := selectRecipes(
		ctx,
		qc,
		func(sb squirrel.SelectBuilder) squirrel.SelectBuilder {
			return sb.Where(
				squirrel.Eq{"recipes.id": id},
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

// UpdateRecipeByID updates an existing recipe by its id. An updated recipe
// is returned.
func UpdateRecipeByID(
	ctx context.Context,
	db *sql.DB,
	id xid.ID,
	rc core.RecipeCore,
) (*core.Recipe, error) {

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, nil
	}
	defer tx.Rollback()

	if err := deleteRecipeProducts(
		ctx,
		tx,
		id,
	); err != nil {
		return nil, err
	}

	for _, rp := range rc.Products {
		rp.RecipeID = id

		if err := upsertRecipeProduct(
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
		squirrel.Update("recipes").SetMap(map[string]interface{}{
			"recipes.name":        rc.Name,
			"recipes.description": rc.Description,
		}).Where(
			squirrel.Eq{"recipes.id": id},
		),
	)

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	rec, err := GetRecipeByID(ctx, db, id)
	if err != nil {
		return nil, err
	}

	return rec, nil
}

// DeleteRecipeByID deletes a recipe by its id.
func DeleteRecipeByID(
	ctx context.Context,
	ec squirrel.ExecerContext,
	id xid.ID,
) error {

	_, err := squirrel.ExecContextWith(
		ctx,
		ec,
		squirrel.Delete("recipes").Where(
			squirrel.Eq{"recipes.id": id},
		),
	)
	return err
}

// selectRecipes selects all recipes by the provided decorator function.
func selectRecipes(
	ctx context.Context,
	qc squirrel.QueryerContext,
	dec func(squirrel.SelectBuilder) squirrel.SelectBuilder,
) ([]core.Recipe, error) {

	rows, err := squirrel.QueryContextWith(ctx, qc, dec(squirrel.
		Select(
			"recipes.id",
			"recipes.user_id",
			"recipes.name",
			"recipes.description",
			"recipes.created_at",
		).From("recipes"),
	))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rr := make([]core.Recipe, 0)
	for rows.Next() {
		var rec core.Recipe
		if err := rows.Scan(
			&rec.ID,
			&rec.UserID,
			&rec.Name,
			&rec.Description,
			&rec.CreatedAt,
		); err != nil {
			return nil, err
		}

		rps, err := getRecipeProductsByRecipeID(ctx, qc, rec.ID)
		if err != nil {
			return nil, err
		}

		rec.Products = rps
		rr = append(rr, rec)
	}

	return rr, nil
}

// GetRecipeProductsByPorudctID selects recipe products by the product id.
func GetRecipeProductsByProductID(
	ctx context.Context,
	qc squirrel.QueryerContext,
	id xid.ID,
) ([]core.RecipeProduct, error) {

	return selectRecipeProducts(
		ctx,
		qc,
		func(sb squirrel.SelectBuilder) squirrel.SelectBuilder {
			return sb.Where(
				squirrel.Eq{"recipe_products.product_id": id},
			)
		},
	)
}

// getRecipeProductsByRecipeID selects recipe products by the recipe id.
func getRecipeProductsByRecipeID(
	ctx context.Context,
	qc squirrel.QueryerContext,
	id xid.ID,
) ([]core.RecipeProduct, error) {

	return selectRecipeProducts(
		ctx,
		qc,
		func(sb squirrel.SelectBuilder) squirrel.SelectBuilder {
			return sb.Where(
				squirrel.Eq{"recipe_products.recipe_id": id},
			)
		},
	)
}

// deleteRecipeProducts deletes all recipe products.
func deleteRecipeProducts(
	ctx context.Context,
	ec squirrel.ExecerContext,
	rid xid.ID,
) error {

	_, err := squirrel.ExecContextWith(
		ctx,
		ec,
		squirrel.Delete("recipe_products").Where(
			squirrel.Eq{"recipe_products.recipe_id": rid},
		),
	)
	return err
}

// upsertRecipeProduct upserts recipe products.
func upsertRecipeProduct(
	ctx context.Context,
	ec squirrel.ExecerContext,
	rp core.RecipeProduct,
) error {

	_, err := squirrel.ExecContextWith(
		ctx,
		ec,
		squirrel.Insert("recipe_products").SetMap(map[string]interface{}{
			"recipe_products.recipe_id":  rp.RecipeID,
			"recipe_products.product_id": rp.ProductID,
			"recipe_products.quantity":   rp.Quantity,
		}).Suffix("ON DUPLICATE KEY UPDATE recipe_products.quantity = VALUES(recipe_products.quantity)"),
	)
	return err
}

// selectRecipeProducts selects all recipe products by the provided decorator
// function.
func selectRecipeProducts(
	ctx context.Context,
	qc squirrel.QueryerContext,
	dec func(squirrel.SelectBuilder) squirrel.SelectBuilder,
) ([]core.RecipeProduct, error) {

	rows, err := squirrel.QueryContextWith(ctx, qc, dec(squirrel.
		Select(
			"recipe_products.recipe_id",
			"recipe_products.product_id",
			"recipe_products.quantity",
		).From("recipe_products"),
	))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rps := make([]core.RecipeProduct, 0)
	for rows.Next() {
		var rp core.RecipeProduct
		if err := rows.Scan(
			&rp.RecipeID,
			&rp.ProductID,
			&rp.Quantity,
		); err != nil {
			return nil, err
		}

		rps = append(rps, rp)
	}

	return rps, nil
}
