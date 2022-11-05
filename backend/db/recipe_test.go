package db

import (
	"context"
	"database/sql"
	"foodie/core"
	"testing"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/rs/xid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_InsertRecipe(t *testing.T) {
	dbh := dbFn(t)
	cleanUpTables(t, dbh)

	uid := xid.New()

	mockUsers(t, dbh, core.User{
		ID:           uid,
		Name:         "1",
		PasswordHash: []byte{1},
		Admin:        true,
	})

	rcp := core.RecipeCore{
		Name:        "123",
		Description: "123",
		Products: []core.RecipeProduct{
			{
				ProductID: xid.New(),
				Quantity:  decimal.New(4000, -4),
			},
		},
	}

	t.Run("foreign key constraint fails", func(t *testing.T) {
		rec, err := InsertRecipe(
			context.Background(),
			dbh,
			uid,
			rcp,
		)
		require.Nil(t, rec)
		require.Error(t, err)
	})

	pid := xid.New()

	mockProducts(t, dbh, core.Product{
		ID: pid,
		ProductCore: core.ProductCore{
			Name: "123",
			Serving: core.Serving{
				Type:     "units",
				Size:     decimal.NewFromInt(1),
				Calories: 2,
			},
		},
	})

	rcp.Products = []core.RecipeProduct{
		{
			ProductID: pid,
			Quantity:  decimal.New(4000, -4),
		},
	}

	t.Run("successfully inserted a new recipe", func(t *testing.T) {
		rec, err := InsertRecipe(
			context.Background(),
			dbh,
			uid,
			rcp,
		)
		require.NoError(t, err)
		assert.NotEmpty(t, rec.ID)
		assert.NotEmpty(t, rec.CreatedAt)
		assert.Equal(t, rcp, rec.RecipeCore)
	})
}

func Test_GetRecipes(t *testing.T) {
	dbh := dbFn(t)
	cleanUpTables(t, dbh)

	uid1 := xid.New()
	uid2 := xid.New()

	mockUsers(t, dbh, []core.User{
		{
			ID:           uid1,
			Name:         "1",
			PasswordHash: []byte{1},
			Admin:        true,
		},
		{
			ID:           uid2,
			Name:         "2",
			PasswordHash: []byte{2},
			Admin:        false,
		},
	}...)

	pid1 := xid.New()
	pid2 := xid.New()
	pid3 := xid.New()

	mockProducts(t, dbh, []core.Product{
		{
			ID: pid1,
			ProductCore: core.ProductCore{
				Name: "123",
				Serving: core.Serving{
					Type:     "units",
					Size:     decimal.NewFromInt(1),
					Calories: 2,
				},
			},
		},
		{
			ID: pid2,
			ProductCore: core.ProductCore{
				Name: "223",
				Serving: core.Serving{
					Type:     "grams",
					Size:     decimal.NewFromInt(1),
					Calories: 5,
				},
			},
		},
		{
			ID: pid3,
			ProductCore: core.ProductCore{
				Name: "143",
				Serving: core.Serving{
					Type:     "units",
					Size:     decimal.NewFromInt(1),
					Calories: 9,
				},
			},
		},
	}...)

	rid1 := xid.New()
	rid2 := xid.New()
	rid3 := xid.New()

	rr := []core.Recipe{
		{
			ID:        rid1,
			CreatedAt: time.Now().UTC().Truncate(time.Second),
			UserID:    uid1,
			RecipeCore: core.RecipeCore{
				Name:        "1",
				Description: "test1",
				Products: []core.RecipeProduct{
					{
						RecipeID:  rid1,
						ProductID: pid1,
						Quantity:  decimal.New(1000, -4),
					},
					{
						RecipeID:  rid1,
						ProductID: pid3,
						Quantity:  decimal.New(3000, -4),
					},
				},
			},
		},
		{
			ID:        rid2,
			CreatedAt: time.Now().UTC().Truncate(time.Second),
			UserID:    uid2,
			RecipeCore: core.RecipeCore{
				Name:        "0",
				Description: "test9",
				Products: []core.RecipeProduct{
					{
						RecipeID:  rid2,
						ProductID: pid2,
						Quantity:  decimal.New(2000, -4),
					},
					{
						RecipeID:  rid2,
						ProductID: pid3,
						Quantity:  decimal.New(2000, -4),
					},
				},
			},
		},
		{
			ID:        rid3,
			CreatedAt: time.Now().UTC().Truncate(time.Second),
			UserID:    uid2,
			RecipeCore: core.RecipeCore{
				Name:        "4",
				Description: "test3",
				Products: []core.RecipeProduct{
					{
						RecipeID:  rid3,
						ProductID: pid1,
						Quantity:  decimal.New(9000, -4),
					},
					{
						RecipeID:  rid3,
						ProductID: pid2,
						Quantity:  decimal.New(8000, -4),
					},
					{
						RecipeID:  rid3,
						ProductID: pid3,
						Quantity:  decimal.New(7000, -4),
					},
				},
			},
		},
	}

	mockRecipes(t, dbh, rr...)

	res, err := GetRecipes(context.Background(), dbh)
	require.NoError(t, err)
	assert.Equal(t, rr, res)
}

func Test_GetRecipesByUserID(t *testing.T) {
	dbh := dbFn(t)
	cleanUpTables(t, dbh)

	uid1 := xid.New()
	uid2 := xid.New()

	mockUsers(t, dbh, []core.User{
		{
			ID:           uid1,
			Name:         "1",
			PasswordHash: []byte{1},
			Admin:        true,
		},
		{
			ID:           uid2,
			Name:         "2",
			PasswordHash: []byte{2},
			Admin:        false,
		},
	}...)

	pid1 := xid.New()
	pid2 := xid.New()
	pid3 := xid.New()

	mockProducts(t, dbh, []core.Product{
		{
			ID: pid1,
			ProductCore: core.ProductCore{
				Name: "123",
				Serving: core.Serving{
					Type:     "units",
					Size:     decimal.NewFromInt(1),
					Calories: 2,
				},
			},
		},
		{
			ID: pid2,
			ProductCore: core.ProductCore{
				Name: "223",
				Serving: core.Serving{
					Type:     "grams",
					Size:     decimal.NewFromInt(1),
					Calories: 5,
				},
			},
		},
		{
			ID: pid3,
			ProductCore: core.ProductCore{
				Name: "143",
				Serving: core.Serving{
					Type:     "units",
					Size:     decimal.NewFromInt(1),
					Calories: 9,
				},
			},
		},
	}...)

	rid1 := xid.New()
	rid2 := xid.New()
	rid3 := xid.New()

	rr := []core.Recipe{
		{
			ID:        rid1,
			CreatedAt: time.Now().UTC().Truncate(time.Second),
			UserID:    uid1,
			RecipeCore: core.RecipeCore{
				Name:        "1",
				Description: "test1",
				Products: []core.RecipeProduct{
					{
						RecipeID:  rid1,
						ProductID: pid1,
						Quantity:  decimal.New(1000, -4),
					},
					{
						RecipeID:  rid1,
						ProductID: pid3,
						Quantity:  decimal.New(3000, -4),
					},
				},
			},
		},
		{
			ID:        rid2,
			CreatedAt: time.Now().UTC().Truncate(time.Second),
			UserID:    uid2,
			RecipeCore: core.RecipeCore{
				Name:        "0",
				Description: "test9",
				Products: []core.RecipeProduct{
					{
						RecipeID:  rid2,
						ProductID: pid2,
						Quantity:  decimal.New(2000, -4),
					},
					{
						RecipeID:  rid2,
						ProductID: pid3,
						Quantity:  decimal.New(2000, -4),
					},
				},
			},
		},
		{
			ID:        rid3,
			CreatedAt: time.Now().UTC().Truncate(time.Second),
			UserID:    uid2,
			RecipeCore: core.RecipeCore{
				Name:        "4",
				Description: "test3",
				Products: []core.RecipeProduct{
					{
						RecipeID:  rid3,
						ProductID: pid1,
						Quantity:  decimal.New(9000, -4),
					},
					{
						RecipeID:  rid3,
						ProductID: pid2,
						Quantity:  decimal.New(8000, -4),
					},
					{
						RecipeID:  rid3,
						ProductID: pid3,
						Quantity:  decimal.New(7000, -4),
					},
				},
			},
		},
	}

	mockRecipes(t, dbh, rr...)

	res, err := GetRecipesByUserID(context.Background(), dbh, uid2)
	require.NoError(t, err)
	assert.Equal(t, rr[1:], res)
}

func Test_GetRecipeByID(t *testing.T) {
	dbh := dbFn(t)
	cleanUpTables(t, dbh)

	uid1 := xid.New()

	mockUsers(t, dbh, core.User{
		ID:           uid1,
		Name:         "1",
		PasswordHash: []byte{1},
		Admin:        true,
	})

	pid1 := xid.New()
	pid2 := xid.New()

	mockProducts(t, dbh, []core.Product{
		{
			ID: pid1,
			ProductCore: core.ProductCore{
				Name: "123",
				Serving: core.Serving{
					Type:     "units",
					Size:     decimal.NewFromInt(1),
					Calories: 2,
				},
			},
		},
		{
			ID: pid2,
			ProductCore: core.ProductCore{
				Name: "223",
				Serving: core.Serving{
					Type:     "grams",
					Size:     decimal.NewFromInt(1),
					Calories: 5,
				},
			},
		},
	}...)

	rid1 := xid.New()

	rcp := core.Recipe{
		ID:        rid1,
		CreatedAt: time.Now().UTC().Truncate(time.Second),
		UserID:    uid1,
		RecipeCore: core.RecipeCore{
			Name:        "1",
			Description: "test1",
			Products: []core.RecipeProduct{
				{
					RecipeID:  rid1,
					ProductID: pid1,
					Quantity:  decimal.New(1000, -4),
				},
				{
					RecipeID:  rid1,
					ProductID: pid2,
					Quantity:  decimal.New(3000, -4),
				},
			},
		},
	}

	mockRecipes(t, dbh, rcp)

	t.Run("not found", func(t *testing.T) {
		res, err := GetRecipeByID(context.Background(), dbh, xid.New())
		assert.Empty(t, res)
		require.Equal(t, ErrNotFound, err)
	})

	t.Run("successfully retrieved a recipe by id", func(t *testing.T) {
		res, err := GetRecipeByID(context.Background(), dbh, rcp.ID)
		require.NoError(t, err)
		assert.Equal(t, &rcp, res)
	})
}

func Test_UpdateRecipeByID(t *testing.T) {
	dbh := dbFn(t)
	cleanUpTables(t, dbh)

	uid1 := xid.New()

	mockUsers(t, dbh, core.User{
		ID:           uid1,
		Name:         "1",
		PasswordHash: []byte{1},
		Admin:        true,
	})

	pid1 := xid.New()
	pid2 := xid.New()

	mockProducts(t, dbh, []core.Product{
		{
			ID: pid1,
			ProductCore: core.ProductCore{
				Name: "123",
				Serving: core.Serving{
					Type:     "units",
					Size:     decimal.NewFromInt(1),
					Calories: 2,
				},
			},
		},
		{
			ID: pid2,
			ProductCore: core.ProductCore{
				Name: "223",
				Serving: core.Serving{
					Type:     "grams",
					Size:     decimal.NewFromInt(1),
					Calories: 5,
				},
			},
		},
	}...)

	rid1 := xid.New()

	rcp := core.Recipe{
		ID:        rid1,
		CreatedAt: time.Now().UTC().Truncate(time.Second),
		UserID:    uid1,
		RecipeCore: core.RecipeCore{
			Name:        "1",
			Description: "test1",
			Products: []core.RecipeProduct{
				{
					RecipeID:  rid1,
					ProductID: pid1,
					Quantity:  decimal.New(1000, -4),
				},
			},
		},
	}

	mockRecipes(t, dbh, rcp)

	rcp.Name = "999"
	rcp.Description = "34124"
	rcp.Products = []core.RecipeProduct{
		{
			RecipeID:  rid1,
			ProductID: pid1,
			Quantity:  decimal.New(3000, -4),
		},
		{
			RecipeID:  rid1,
			ProductID: pid2,
			Quantity:  decimal.New(9000, -4),
		},
	}

	res, err := UpdateRecipeByID(context.Background(), dbh, rcp.ID, rcp.RecipeCore)
	require.NoError(t, err)
	assert.Equal(t, &rcp, res)
}

func Test_DeleteRecipeByID(t *testing.T) {
	dbh := dbFn(t)
	cleanUpTables(t, dbh)

	uid1 := xid.New()

	mockUsers(t, dbh, core.User{
		ID:           uid1,
		Name:         "1",
		PasswordHash: []byte{1},
		Admin:        true,
	})

	pid1 := xid.New()
	pid2 := xid.New()

	mockProducts(t, dbh, []core.Product{
		{
			ID: pid1,
			ProductCore: core.ProductCore{
				Name: "123",
				Serving: core.Serving{
					Type:     "units",
					Size:     decimal.NewFromInt(1),
					Calories: 2,
				},
			},
		},
		{
			ID: pid2,
			ProductCore: core.ProductCore{
				Name: "223",
				Serving: core.Serving{
					Type:     "grams",
					Size:     decimal.NewFromInt(1),
					Calories: 5,
				},
			},
		},
	}...)

	rid1 := xid.New()

	rcp := core.Recipe{
		ID:        rid1,
		CreatedAt: time.Now().UTC().Truncate(time.Second),
		UserID:    uid1,
		RecipeCore: core.RecipeCore{
			Name:        "1",
			Description: "test1",
			Products: []core.RecipeProduct{
				{
					RecipeID:  rid1,
					ProductID: pid1,
					Quantity:  decimal.New(1000, -4),
				},
			},
		},
	}

	mockRecipes(t, dbh, rcp)

	rr := retrieveRecipes(t, dbh)
	require.Len(t, rr, 1)
	assert.Equal(t, rcp, rr[0])

	require.NoError(t, DeleteRecipeByID(context.Background(), dbh, rcp.ID))
	require.Len(t, retrieveRecipes(t, dbh), 0)
}

func Test_GetRecipeProductsByProductID(t *testing.T) {
	dbh := dbFn(t)
	cleanUpTables(t, dbh)

	uid1 := xid.New()
	uid2 := xid.New()

	mockUsers(t, dbh, []core.User{
		{
			ID:           uid1,
			Name:         "1",
			PasswordHash: []byte{1},
			Admin:        true,
		},
		{
			ID:           uid2,
			Name:         "2",
			PasswordHash: []byte{2},
			Admin:        false,
		},
	}...)

	pid1 := xid.New()
	pid2 := xid.New()
	pid3 := xid.New()

	mockProducts(t, dbh, []core.Product{
		{
			ID: pid1,
			ProductCore: core.ProductCore{
				Name: "123",
				Serving: core.Serving{
					Type:     "units",
					Size:     decimal.NewFromInt(1),
					Calories: 2,
				},
			},
		},
		{
			ID: pid2,
			ProductCore: core.ProductCore{
				Name: "223",
				Serving: core.Serving{
					Type:     "grams",
					Size:     decimal.NewFromInt(1),
					Calories: 5,
				},
			},
		},
		{
			ID: pid3,
			ProductCore: core.ProductCore{
				Name: "143",
				Serving: core.Serving{
					Type:     "units",
					Size:     decimal.NewFromInt(1),
					Calories: 9,
				},
			},
		},
	}...)

	rid1 := xid.New()
	rid2 := xid.New()
	rid3 := xid.New()

	rr := []core.Recipe{
		{
			ID:        rid1,
			CreatedAt: time.Now().UTC().Truncate(time.Second),
			UserID:    uid1,
			RecipeCore: core.RecipeCore{
				Name:        "1",
				Description: "test1",
				Products: []core.RecipeProduct{
					{
						RecipeID:  rid1,
						ProductID: pid1,
						Quantity:  decimal.New(1000, -4),
					},
					{
						RecipeID:  rid1,
						ProductID: pid3,
						Quantity:  decimal.New(3000, -4),
					},
				},
			},
		},
		{
			ID:        rid2,
			CreatedAt: time.Now().UTC().Truncate(time.Second),
			UserID:    uid2,
			RecipeCore: core.RecipeCore{
				Name:        "0",
				Description: "test9",
				Products: []core.RecipeProduct{
					{
						RecipeID:  rid2,
						ProductID: pid2,
						Quantity:  decimal.New(2000, -4),
					},
					{
						RecipeID:  rid2,
						ProductID: pid3,
						Quantity:  decimal.New(2000, -4),
					},
				},
			},
		},
		{
			ID:        rid3,
			CreatedAt: time.Now().UTC().Truncate(time.Second),
			UserID:    uid2,
			RecipeCore: core.RecipeCore{
				Name:        "4",
				Description: "test3",
				Products: []core.RecipeProduct{
					{
						RecipeID:  rid3,
						ProductID: pid1,
						Quantity:  decimal.New(9000, -4),
					},
					{
						RecipeID:  rid3,
						ProductID: pid2,
						Quantity:  decimal.New(8000, -4),
					},
					{
						RecipeID:  rid3,
						ProductID: pid3,
						Quantity:  decimal.New(7000, -4),
					},
				},
			},
		},
	}

	mockRecipes(t, dbh, rr...)

	res, err := GetRecipeProductsByProductID(context.Background(), dbh, pid2)
	require.NoError(t, err)
	assert.Equal(t, []core.RecipeProduct{
		{
			RecipeID:  rid2,
			ProductID: pid2,
			Quantity:  decimal.New(2000, -4),
		},
		{
			RecipeID:  rid3,
			ProductID: pid2,
			Quantity:  decimal.New(8000, -4),
		},
	}, res)
}

func mockRecipes(t *testing.T, dbh *sql.DB, rr ...core.Recipe) {
	t.Helper()

	for _, rcp := range rr {
		_, err := squirrel.ExecWith(
			dbh,
			squirrel.Insert("recipes").SetMap(map[string]interface{}{
				"recipes.id":          rcp.ID,
				"recipes.user_id":     rcp.UserID,
				"recipes.name":        rcp.Name,
				"recipes.description": rcp.Description,
				"recipes.created_at":  rcp.CreatedAt,
			}),
		)
		require.NoError(t, err)

		for _, rp := range rcp.Products {
			_, err = squirrel.ExecWith(
				dbh,
				squirrel.Insert("recipe_products").SetMap(map[string]interface{}{
					"recipe_products.recipe_id":  rp.RecipeID,
					"recipe_products.product_id": rp.ProductID,
					"recipe_products.quantity":   rp.Quantity,
				}).Suffix("ON DUPLICATE KEY UPDATE recipe_products.quantity = VALUES(recipe_products.quantity)"),
			)
			require.NoError(t, err)
		}
	}
}

func retrieveRecipes(t *testing.T, dbh *sql.DB) []core.Recipe {
	rows, err := squirrel.QueryWith(dbh, squirrel.
		Select(
			"recipes.id",
			"recipes.user_id",
			"recipes.name",
			"recipes.description",
			"recipes.created_at",
		).From("recipes"),
	)
	require.NoError(t, err)

	defer rows.Close()

	rr := make([]core.Recipe, 0)

	for rows.Next() {
		var rec core.Recipe

		require.NoError(t, rows.Scan(
			&rec.ID,
			&rec.UserID,
			&rec.Name,
			&rec.Description,
			&rec.CreatedAt,
		))

		rec.Products = retrieveRecipeProducts(t, dbh, rec.ID)
		rr = append(rr, rec)
	}

	return rr
}

func retrieveRecipeProducts(t *testing.T, dbh *sql.DB, rid xid.ID) []core.RecipeProduct {
	rows, err := squirrel.QueryWith(dbh, squirrel.
		Select(
			"recipe_products.recipe_id",
			"recipe_products.product_id",
			"recipe_products.quantity",
		).From("recipe_products").
		Where(squirrel.Eq{
			"recipe_products.recipe_id": rid,
		}),
	)
	require.NoError(t, err)

	defer rows.Close()

	rps := make([]core.RecipeProduct, 0)

	for rows.Next() {
		var rp core.RecipeProduct

		require.NoError(t, rows.Scan(
			&rp.RecipeID,
			&rp.ProductID,
			&rp.Quantity,
		))

		rps = append(rps, rp)
	}

	return rps
}
