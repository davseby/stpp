package db

import (
	"context"
	"foodie/core"

	"github.com/Masterminds/squirrel"
	"github.com/rs/xid"
)

func GetProducts(ctx context.Context, qc squirrel.QueryerContext) ([]core.Product, error) {
	return selectProducts(
		ctx,
		qc,
		func(sb squirrel.SelectBuilder) squirrel.SelectBuilder {
			return sb
		},
	)
}

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
				squirrel.Eq{"product.id": id},
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

func InsertProduct(
	ctx context.Context,
	ec squirrel.ExecerContext,
	aid xid.ID,
	pc core.ProductCore,
) (*core.Product, error) {

	product := core.Product{
		ID:          xid.New(),
		AdminID:     aid,
		ProductCore: pc,
	}

	_, err := squirrel.ExecContextWith(
		ctx,
		ec,
		squirrel.Insert("product").SetMap(map[string]interface{}{
			"product.id":               product.ID,
			"product.admin_id":         product.AdminID,
			"product.name":             product.Name,
			"product.serving_type":     product.Serving.Type,
			"product.serving_size":     product.Serving.Size,
			"product.serving_calories": product.Serving.Calories,
		}),
	)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

func UpdateProduct(
	ctx context.Context,
	ec squirrel.ExecerContext,
	id xid.ID,
	pc core.ProductCore,
) error {

	_, err := squirrel.ExecContextWith(
		ctx,
		ec,
		squirrel.Update("product").SetMap(map[string]interface{}{
			"product.name":             pc.Name,
			"product.serving_type":     pc.Serving.Type,
			"product.serving_size":     pc.Serving.Size,
			"product.serving_calories": pc.Serving.Calories,
		}).Where(
			squirrel.Eq{"product.id": id},
		),
	)
	return err
}

func DeleteProduct(
	ctx context.Context,
	ec squirrel.ExecerContext,
	id xid.ID,
) error {

	_, err := squirrel.ExecContextWith(
		ctx,
		ec,
		squirrel.Delete("product").Where(
			squirrel.Eq{"product.id": id},
		),
	)
	return err
}

func selectProducts(
	ctx context.Context,
	qc squirrel.QueryerContext,
	dec func(squirrel.SelectBuilder) squirrel.SelectBuilder,
) ([]core.Product, error) {

	rows, err := squirrel.QueryContextWith(ctx, qc, dec(squirrel.
		Select(
			"product.id",
			"product.admin_id",
			"product.name",
			"product.serving_type",
			"product.serving_size",
			"product.serving_calories",
		).From("product"),
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
			&product.AdminID,
			&product.Name,
			&product.Serving.Type,
			&product.Serving.Size,
			&product.Serving.Calories,
		); err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	return products, nil
}
