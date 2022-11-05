package db

import (
	"context"
	"foodie/core"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/rs/xid"
)

// InsertProduct inserts a new product into the database.
func InsertProduct(
	ctx context.Context,
	ec squirrel.ExecerContext,
	pc core.ProductCore,
) (*core.Product, error) {

	product := core.Product{
		ID:          xid.New(),
		CreatedAt:   time.Now(),
		ProductCore: pc,
	}

	_, err := squirrel.ExecContextWith(
		ctx,
		ec,
		squirrel.Insert("products").SetMap(map[string]interface{}{
			"products.id":               product.ID,
			"products.name":             product.Name,
			"products.description":      product.Description,
			"products.image_url":        product.ImageURL,
			"products.serving_type":     product.Serving.Type,
			"products.serving_size":     product.Serving.Size,
			"products.serving_calories": product.Serving.Calories,
			"products.created_at":       product.CreatedAt,
		}),
	)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

// GetProducts retrieves all products.
func GetProducts(ctx context.Context, qc squirrel.QueryerContext) ([]core.Product, error) {
	return selectProducts(
		ctx,
		qc,
		func(sb squirrel.SelectBuilder) squirrel.SelectBuilder {
			return sb
		},
	)
}

// GetProductByID retrieves a product by the product id.
func GetProductByID(
	ctx context.Context,
	qc squirrel.QueryerContext,
	id xid.ID,
) (*core.Product, error) {

	products, err := selectProducts(
		ctx,
		qc,
		func(sb squirrel.SelectBuilder) squirrel.SelectBuilder {
			return sb.Where(
				squirrel.Eq{"products.id": id},
			)
		},
	)
	if err != nil {
		return nil, err
	}

	if len(products) == 0 {
		return nil, ErrNotFound
	}

	return &products[0], nil
}

// UpdateProductByID updates an existing recipe by its id. An updated product
// is returned.
func UpdateProductByID(
	ctx context.Context,
	ssc squirrel.StdSqlCtx,
	id xid.ID,
	pc core.ProductCore,
) (*core.Product, error) {

	_, err := squirrel.ExecContextWith(
		ctx,
		ssc,
		squirrel.Update("products").SetMap(map[string]interface{}{
			"products.name":             pc.Name,
			"products.description":      pc.Description,
			"products.image_url":        pc.ImageURL,
			"products.serving_type":     pc.Serving.Type,
			"products.serving_size":     pc.Serving.Size,
			"products.serving_calories": pc.Serving.Calories,
		}).Where(
			squirrel.Eq{"products.id": id},
		),
	)
	if err != nil {
		return nil, err
	}

	prd, err := GetProductByID(ctx, ssc, id)
	if err != nil {
		return nil, err
	}

	return prd, nil
}

// DeleteProductByID deletes a product by its id.
func DeleteProductByID(
	ctx context.Context,
	ec squirrel.ExecerContext,
	id xid.ID,
) error {

	_, err := squirrel.ExecContextWith(
		ctx,
		ec,
		squirrel.Delete("products").Where(
			squirrel.Eq{"products.id": id},
		),
	)
	return err
}

// selectProducts selects all products by the provided decorator function.
func selectProducts(
	ctx context.Context,
	qc squirrel.QueryerContext,
	dec func(squirrel.SelectBuilder) squirrel.SelectBuilder,
) ([]core.Product, error) {

	rows, err := squirrel.QueryContextWith(ctx, qc, dec(squirrel.
		Select(
			"products.id",
			"products.name",
			"COALESCE(products.image_url, '')",
			"products.description",
			"products.serving_type",
			"products.serving_size",
			"products.serving_calories",
			"products.created_at",
		).From("products"),
	))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]core.Product, 0)
	for rows.Next() {
		var product core.Product
		if err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.ImageURL,
			&product.Description,
			&product.Serving.Type,
			&product.Serving.Size,
			&product.Serving.Calories,
			&product.CreatedAt,
		); err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	return products, nil
}
