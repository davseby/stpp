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

func Test_InsertProduct(t *testing.T) {
	dbh := dbFn(t)
	cleanUpTables(t, dbh)

	pc := core.ProductCore{
		Name: "123",
		Serving: core.Serving{
			Type:     "units",
			Size:     decimal.NewFromInt(1),
			Calories: 2,
		},
	}

	prd, err := InsertProduct(
		context.Background(),
		dbh,
		pc,
	)
	require.NoError(t, err)
	assert.NotEmpty(t, prd.ID)
	assert.NotEmpty(t, prd.CreatedAt)
	assert.Equal(t, pc, prd.ProductCore)
}

func Test_GetProducts(t *testing.T) {
	dbh := dbFn(t)
	cleanUpTables(t, dbh)

	pp := []core.Product{
		{
			ID:        xid.New(),
			CreatedAt: time.Now().UTC().Truncate(time.Second),
			ProductCore: core.ProductCore{
				Name: "123",
				Serving: core.Serving{
					Type:     "units",
					Size:     decimal.New(3000, -4),
					Calories: 2,
				},
			},
		},
		{
			ID:        xid.New(),
			CreatedAt: time.Now().UTC().Truncate(time.Second),
			ProductCore: core.ProductCore{
				Name: "12",
				Serving: core.Serving{
					Type:     "grams",
					Size:     decimal.New(5000, -4),
					Calories: 4,
				},
			},
		},
		{
			ID:        xid.New(),
			CreatedAt: time.Now().UTC().Truncate(time.Second),
			ProductCore: core.ProductCore{
				Name: "125",
				Serving: core.Serving{
					Type:     "milliliters",
					Size:     decimal.New(9000, -4),
					Calories: 7,
				},
			},
		},
	}

	for _, prd := range pp {
		mockProduct(t, dbh, prd)
	}

	res, err := GetProducts(context.Background(), dbh)
	require.NoError(t, err)
	assert.Equal(t, pp, res)
}

func Test_GetProductByID(t *testing.T) {
	dbh := dbFn(t)
	cleanUpTables(t, dbh)

	pp := []core.Product{
		{
			ID:        xid.New(),
			CreatedAt: time.Now().UTC().Truncate(time.Second),
			ProductCore: core.ProductCore{
				Name: "123",
				Serving: core.Serving{
					Type:     "units",
					Size:     decimal.New(1000, -4),
					Calories: 2,
				},
			},
		},
		{
			ID:        xid.New(),
			CreatedAt: time.Now().UTC().Truncate(time.Second),
			ProductCore: core.ProductCore{
				Name: "12",
				Serving: core.Serving{
					Type:     "grams",
					Size:     decimal.New(5000, -4),
					Calories: 9,
				},
			},
		},
		{
			ID:        xid.New(),
			CreatedAt: time.Now().UTC().Truncate(time.Second),
			ProductCore: core.ProductCore{
				Name: "125",
				Serving: core.Serving{
					Type:     "milliliters",
					Size:     decimal.New(3000, -4),
					Calories: 4,
				},
			},
		},
	}

	for _, prd := range pp {
		mockProduct(t, dbh, prd)
	}

	t.Run("not found", func(t *testing.T) {
		res, err := GetProductByID(context.Background(), dbh, xid.New())
		assert.Empty(t, res)
		require.Equal(t, ErrNotFound, err)
	})

	t.Run("successfully retrieved a product by id", func(t *testing.T) {
		res, err := GetProductByID(context.Background(), dbh, pp[1].ID)
		require.NoError(t, err)
		assert.Equal(t, &pp[1], res)
	})
}

func Test_UpdateProductByID(t *testing.T) {
	dbh := dbFn(t)
	cleanUpTables(t, dbh)

	prd := core.Product{
		ID:        xid.New(),
		CreatedAt: time.Now().UTC().Truncate(time.Second),
		ProductCore: core.ProductCore{
			Name: "11",
			Serving: core.Serving{
				Type:     "units",
				Size:     decimal.New(4000, -4),
				Calories: 9,
			},
		},
	}

	mockProduct(t, dbh, prd)

	prd.Name = "12"
	prd.Serving.Calories = 200
	prd.Serving.Size = decimal.New(5000, -4)

	res, err := UpdateProductByID(context.Background(), dbh, prd.ID, prd.ProductCore)
	require.NoError(t, err)
	assert.Equal(t, &prd, res)
}

func Test_DeleteProductByID(t *testing.T) {
	dbh := dbFn(t)
	cleanUpTables(t, dbh)

	prd := core.Product{
		ID:        xid.New(),
		CreatedAt: time.Now().UTC().Truncate(time.Second),
		ProductCore: core.ProductCore{
			Name: "11",
			Serving: core.Serving{
				Type:     "units",
				Size:     decimal.New(4000, -4),
				Calories: 9,
			},
		},
	}

	mockProduct(t, dbh, prd)

	pp := retrieveProducts(t, dbh)
	require.Len(t, pp, 1)
	assert.Equal(t, prd, pp[0])

	require.NoError(t, DeleteProductByID(context.Background(), dbh, prd.ID))
	require.Len(t, retrieveProducts(t, dbh), 0)
}

func Test_selectProduct(t *testing.T) {
	dbh := dbFn(t)
	cleanUpTables(t, dbh)

	pp := []core.Product{
		{
			ID:        xid.New(),
			CreatedAt: time.Now().UTC().Truncate(time.Second),
			ProductCore: core.ProductCore{
				Name: "123",
				Serving: core.Serving{
					Type:     "units",
					Size:     decimal.New(3000, -4),
					Calories: 2,
				},
			},
		},
		{
			ID:        xid.New(),
			CreatedAt: time.Now().UTC().Truncate(time.Second),
			ProductCore: core.ProductCore{
				Name: "12",
				Serving: core.Serving{
					Type:     "grams",
					Size:     decimal.New(5000, -4),
					Calories: 4,
				},
			},
		},
		{
			ID:        xid.New(),
			CreatedAt: time.Now().UTC().Truncate(time.Second),
			ProductCore: core.ProductCore{
				Name: "125",
				Serving: core.Serving{
					Type:     "milliliters",
					Size:     decimal.New(9000, -4),
					Calories: 7,
				},
			},
		},
	}

	for _, prd := range pp {
		mockProduct(t, dbh, prd)
	}

	res, err := selectProducts(context.Background(), dbh, func(sb squirrel.SelectBuilder) squirrel.SelectBuilder {
		return sb
	})
	require.NoError(t, err)
	assert.Equal(t, retrieveProducts(t, dbh), res)
}

func mockProduct(t *testing.T, dbh *sql.DB, prd core.Product) {
	t.Helper()

	_, err := squirrel.ExecWith(
		dbh,
		squirrel.Insert("products").SetMap(map[string]interface{}{
			"products.id":               prd.ID,
			"products.name":             prd.Name,
			"products.serving_type":     prd.Serving.Type,
			"products.serving_size":     prd.Serving.Size,
			"products.serving_calories": prd.Serving.Calories,
			"products.created_at":       prd.CreatedAt,
		}),
	)
	require.NoError(t, err)
}

func retrieveProducts(t *testing.T, dbh *sql.DB) []core.Product {
	rows, err := squirrel.QueryWith(dbh, squirrel.
		Select(
			"products.id",
			"products.name",
			"products.serving_type",
			"products.serving_size",
			"products.serving_calories",
			"products.created_at",
		).From("products"),
	)
	require.NoError(t, err)
	defer rows.Close()

	products := make([]core.Product, 0)
	for rows.Next() {
		var product core.Product
		require.NoError(t, rows.Scan(
			&product.ID,
			&product.Name,
			&product.Serving.Type,
			&product.Serving.Size,
			&product.Serving.Calories,
			&product.CreatedAt,
		))

		products = append(products, product)
	}

	return products
}
