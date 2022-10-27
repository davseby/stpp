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

func Test_InsertPlan(t *testing.T) {
	dbh := dbFn(t)
	cleanUpTables(t, dbh)

	uid := xid.New()
	mockUsers(t, dbh, core.User{
		ID:           uid,
		Name:         "1",
		PasswordHash: []byte{1},
		Admin:        true,
	})

	pc := core.PlanCore{
		Name:        "123",
		Description: "123",
		Recipes: []core.PlanRecipe{
			{
				RecipeID: xid.New(),
				Quantity: 2,
			},
		},
	}

	t.Run("foreign key constraint fails", func(t *testing.T) {
		rec, err := InsertPlan(
			context.Background(),
			dbh,
			uid,
			pc,
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

	rid := xid.New()
	mockRecipes(t, dbh, core.Recipe{
		ID: rid,
		RecipeCore: core.RecipeCore{
			Name:        "123124",
			Description: "48924",
			Products: []core.RecipeProduct{
				{
					RecipeID:  rid,
					ProductID: pid,
					Quantity:  decimal.New(4000, -4),
				},
			},
		},
	})

	pc.Recipes = []core.PlanRecipe{
		{
			RecipeID: rid,
			Quantity: 9,
		},
	}

	t.Run("succesfully inserted a new plan", func(t *testing.T) {
		rec, err := InsertPlan(
			context.Background(),
			dbh,
			uid,
			pc,
		)
		require.NoError(t, err)
		assert.NotEmpty(t, rec.ID)
		assert.NotEmpty(t, rec.CreatedAt)
		assert.Equal(t, pc, rec.PlanCore)
	})
}

func Test_GetPlans(t *testing.T) {
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

	plid1 := xid.New()
	plid2 := xid.New()
	plid3 := xid.New()
	pp := []core.Plan{
		{
			ID:        plid1,
			CreatedAt: time.Now().UTC().Truncate(time.Second),
			UserID:    uid1,
			PlanCore: core.PlanCore{
				Name:        "1",
				Description: "test1",
				Recipes: []core.PlanRecipe{
					{
						PlanID:   plid1,
						RecipeID: rid1,
						Quantity: 3,
					},
					{
						PlanID:   plid1,
						RecipeID: rid3,
						Quantity: 2,
					},
				},
			},
		},
		{
			ID:        plid2,
			CreatedAt: time.Now().UTC().Truncate(time.Second),
			UserID:    uid2,
			PlanCore: core.PlanCore{
				Name:        "0",
				Description: "test9",
				Recipes: []core.PlanRecipe{
					{
						PlanID:   plid2,
						RecipeID: rid1,
						Quantity: 2,
					},
					{
						PlanID:   plid2,
						RecipeID: rid2,
						Quantity: 3,
					},
				},
			},
		},
		{
			ID:        plid3,
			CreatedAt: time.Now().UTC().Truncate(time.Second),
			UserID:    uid2,
			PlanCore: core.PlanCore{
				Name:        "4",
				Description: "test3",
				Recipes: []core.PlanRecipe{
					{
						PlanID:   plid3,
						RecipeID: rid1,
						Quantity: 4,
					},
					{
						PlanID:   plid3,
						RecipeID: rid2,
						Quantity: 2,
					},
					{
						PlanID:   plid3,
						RecipeID: rid3,
						Quantity: 5,
					},
				},
			},
		},
	}

	mockPlans(t, dbh, pp...)

	res, err := GetPlans(context.Background(), dbh)
	require.NoError(t, err)
	assert.Equal(t, pp, res)
}

func Test_GetPlansByUserID(t *testing.T) {
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

	plid1 := xid.New()
	plid2 := xid.New()
	plid3 := xid.New()
	pp := []core.Plan{
		{
			ID:        plid1,
			CreatedAt: time.Now().UTC().Truncate(time.Second),
			UserID:    uid1,
			PlanCore: core.PlanCore{
				Name:        "1",
				Description: "test1",
				Recipes: []core.PlanRecipe{
					{
						PlanID:   plid1,
						RecipeID: rid1,
						Quantity: 3,
					},
					{
						PlanID:   plid1,
						RecipeID: rid3,
						Quantity: 2,
					},
				},
			},
		},
		{
			ID:        plid2,
			CreatedAt: time.Now().UTC().Truncate(time.Second),
			UserID:    uid2,
			PlanCore: core.PlanCore{
				Name:        "0",
				Description: "test9",
				Recipes: []core.PlanRecipe{
					{
						PlanID:   plid2,
						RecipeID: rid1,
						Quantity: 2,
					},
					{
						PlanID:   plid2,
						RecipeID: rid2,
						Quantity: 3,
					},
				},
			},
		},
		{
			ID:        plid3,
			CreatedAt: time.Now().UTC().Truncate(time.Second),
			UserID:    uid2,
			PlanCore: core.PlanCore{
				Name:        "4",
				Description: "test3",
				Recipes: []core.PlanRecipe{
					{
						PlanID:   plid3,
						RecipeID: rid1,
						Quantity: 4,
					},
					{
						PlanID:   plid3,
						RecipeID: rid2,
						Quantity: 2,
					},
					{
						PlanID:   plid3,
						RecipeID: rid3,
						Quantity: 5,
					},
				},
			},
		},
	}

	mockPlans(t, dbh, pp...)

	res, err := GetPlansByUserID(context.Background(), dbh, uid2)
	require.NoError(t, err)
	assert.Equal(t, pp[1:], res)
}

func Test_GetPlanByID(t *testing.T) {
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

	plid1 := xid.New()
	pln := core.Plan{
		ID:        plid1,
		CreatedAt: time.Now().UTC().Truncate(time.Second),
		UserID:    uid1,
		PlanCore: core.PlanCore{
			Name:        "1",
			Description: "test1",
			Recipes: []core.PlanRecipe{
				{
					PlanID:   plid1,
					RecipeID: rid1,
					Quantity: 3,
				},
			},
		},
	}

	mockPlans(t, dbh, pln)

	t.Run("not found", func(t *testing.T) {
		res, err := GetPlanByID(context.Background(), dbh, xid.New())
		assert.Empty(t, res)
		require.Equal(t, ErrNotFound, err)
	})

	t.Run("successfully retrieved a plan by id", func(t *testing.T) {
		res, err := GetPlanByID(context.Background(), dbh, pln.ID)
		require.NoError(t, err)
		assert.Equal(t, &pln, res)
	})
}

func Test_UpdatePlanByID(t *testing.T) {
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
	rid2 := xid.New()
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
						ProductID: pid2,
						Quantity:  decimal.New(3000, -4),
					},
				},
			},
		},
		{
			ID:        rid2,
			CreatedAt: time.Now().UTC().Truncate(time.Second),
			UserID:    uid1,
			RecipeCore: core.RecipeCore{
				Name:        "0",
				Description: "test9",
				Products: []core.RecipeProduct{
					{
						RecipeID:  rid2,
						ProductID: pid2,
						Quantity:  decimal.New(2000, -4),
					},
				},
			},
		},
	}

	mockRecipes(t, dbh, rr...)

	plid1 := xid.New()
	pln := core.Plan{
		ID:        plid1,
		CreatedAt: time.Now().UTC().Truncate(time.Second),
		UserID:    uid1,
		PlanCore: core.PlanCore{
			Name:        "1",
			Description: "test1",
			Recipes: []core.PlanRecipe{
				{
					PlanID:   plid1,
					RecipeID: rid1,
					Quantity: 3,
				},
			},
		},
	}

	mockPlans(t, dbh, pln)

	pln.Name = "another"
	pln.Description = "test"
	pln.Recipes = []core.PlanRecipe{
		{
			PlanID:   plid1,
			RecipeID: rid1,
			Quantity: 1,
		},
		{
			PlanID:   plid1,
			RecipeID: rid2,
			Quantity: 2,
		},
	}

	res, err := UpdatePlanByID(context.Background(), dbh, pln.ID, pln.PlanCore)
	require.NoError(t, err)
	assert.Equal(t, &pln, res)
}

func Test_DeletePlanByID(t *testing.T) {
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

	plid1 := xid.New()
	pln := core.Plan{
		ID:        plid1,
		CreatedAt: time.Now().UTC().Truncate(time.Second),
		UserID:    uid1,
		PlanCore: core.PlanCore{
			Name:        "1",
			Description: "test1",
			Recipes: []core.PlanRecipe{
				{
					PlanID:   plid1,
					RecipeID: rid1,
					Quantity: 3,
				},
			},
		},
	}

	mockPlans(t, dbh, pln)

	pp := retrievePlans(t, dbh)
	require.Len(t, pp, 1)
	assert.Equal(t, pln, pp[0])

	require.NoError(t, DeletePlanByID(context.Background(), dbh, pln.ID))
	require.Len(t, retrievePlans(t, dbh), 0)
}

func Test_GetPlanRecipesByRecipeID(t *testing.T) {
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

	plid1 := xid.New()
	plid2 := xid.New()
	plid3 := xid.New()
	pp := []core.Plan{
		{
			ID:        plid1,
			CreatedAt: time.Now().UTC().Truncate(time.Second),
			UserID:    uid1,
			PlanCore: core.PlanCore{
				Name:        "1",
				Description: "test1",
				Recipes: []core.PlanRecipe{
					{
						PlanID:   plid1,
						RecipeID: rid1,
						Quantity: 3,
					},
					{
						PlanID:   plid1,
						RecipeID: rid3,
						Quantity: 2,
					},
				},
			},
		},
		{
			ID:        plid2,
			CreatedAt: time.Now().UTC().Truncate(time.Second),
			UserID:    uid2,
			PlanCore: core.PlanCore{
				Name:        "0",
				Description: "test9",
				Recipes: []core.PlanRecipe{
					{
						PlanID:   plid2,
						RecipeID: rid1,
						Quantity: 2,
					},
					{
						PlanID:   plid2,
						RecipeID: rid2,
						Quantity: 3,
					},
				},
			},
		},
		{
			ID:        plid3,
			CreatedAt: time.Now().UTC().Truncate(time.Second),
			UserID:    uid2,
			PlanCore: core.PlanCore{
				Name:        "4",
				Description: "test3",
				Recipes: []core.PlanRecipe{
					{
						PlanID:   plid3,
						RecipeID: rid1,
						Quantity: 4,
					},
					{
						PlanID:   plid3,
						RecipeID: rid2,
						Quantity: 2,
					},
					{
						PlanID:   plid3,
						RecipeID: rid3,
						Quantity: 5,
					},
				},
			},
		},
	}

	mockPlans(t, dbh, pp...)

	res, err := GetPlanRecipesByRecipeID(context.Background(), dbh, rid2)
	require.NoError(t, err)
	assert.Equal(t, []core.PlanRecipe{
		{
			PlanID:   plid2,
			RecipeID: rid2,
			Quantity: 3,
		},
		{
			PlanID:   plid3,
			RecipeID: rid2,
			Quantity: 2,
		},
	}, res)
}

func mockPlans(t *testing.T, dbh *sql.DB, pp ...core.Plan) {
	t.Helper()

	for _, pl := range pp {
		_, err := squirrel.ExecWith(
			dbh,
			squirrel.Insert("plans").SetMap(map[string]interface{}{
				"plans.id":          pl.ID,
				"plans.user_id":     pl.UserID,
				"plans.name":        pl.Name,
				"plans.description": pl.Description,
				"plans.created_at":  pl.CreatedAt,
			}),
		)
		require.NoError(t, err)

		for _, pr := range pl.Recipes {
			_, err = squirrel.ExecWith(
				dbh,
				squirrel.Insert("plan_recipes").SetMap(map[string]interface{}{
					"plan_recipes.plan_id":   pr.PlanID,
					"plan_recipes.recipe_id": pr.RecipeID,
					"plan_recipes.quantity":  pr.Quantity,
				}).Suffix("ON DUPLICATE KEY UPDATE plan_recipes.quantity = VALUES(plan_recipes.quantity)"),
			)
			require.NoError(t, err)
		}
	}
}

func retrievePlans(t *testing.T, dbh *sql.DB) []core.Plan {
	rows, err := squirrel.QueryWith(dbh, squirrel.
		Select(
			"plans.id",
			"plans.user_id",
			"plans.name",
			"plans.description",
			"plans.created_at",
		).From("plans"),
	)
	require.NoError(t, err)
	defer rows.Close()

	pp := make([]core.Plan, 0)
	for rows.Next() {
		var pl core.Plan
		require.NoError(t, rows.Scan(
			&pl.ID,
			&pl.UserID,
			&pl.Name,
			&pl.Description,
			&pl.CreatedAt,
		))

		pl.Recipes = retrievePlanRecipes(t, dbh, pl.ID)
		pp = append(pp, pl)
	}

	return pp
}

func retrievePlanRecipes(t *testing.T, dbh *sql.DB, pid xid.ID) []core.PlanRecipe {
	rows, err := squirrel.QueryWith(dbh, squirrel.
		Select(
			"plan_recipes.plan_id",
			"plan_recipes.recipe_id",
			"plan_recipes.quantity",
		).From("plan_recipes").
		Where(squirrel.Eq{
			"plan_recipes.plan_id": pid,
		}),
	)
	require.NoError(t, err)
	defer rows.Close()

	prs := make([]core.PlanRecipe, 0)
	for rows.Next() {
		var pr core.PlanRecipe
		require.NoError(t, rows.Scan(
			&pr.PlanID,
			&pr.RecipeID,
			&pr.Quantity,
		))

		prs = append(prs, pr)
	}

	return prs
}
