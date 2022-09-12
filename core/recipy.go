package core

import (
	"foodie/server/apierr"

	"github.com/rs/xid"
	"github.com/shopspring/decimal"
)

// Recipy contains recipy data.
type Recipy struct {
	RecipyCore

	// ID is a unique recipy identifier.
	ID xid.ID `json:"id"`

	// UserID specifies the user which created the recipy.
	UserID xid.ID `json:"user_id"`
}

// RecipyCore contains core recipy information.
type RecipyCore struct {
	// Name specifies the name of the recipy.
	Name string `json:"name"`

	// Private specifies whether the recipy is private.
	Private bool `json:"private"`

	// Description provides a brief description of the recipy.
	Description string `json:"description"`

	// Products contains recipy products.
	Products []RecipyProduct `json:"products"`
}

// Validate checks whether recipy core contains valid attributes.
func (rc *RecipyCore) Validate() *apierr.Error {
	if rc.Name == "" {
		return apierr.InvalidAttribute("name", "cannot be empty")
	}

	if rc.Description == "" {
		return apierr.InvalidAttribute("description", "cannot be empty")
	}

	if len(rc.Products) < 2 {
		return apierr.InvalidAttribute("products", "must contains at least two elements")
	}

	return nil
}

// RecipyProduct maps recipy product with the actual product stored in the
// system.
type RecipyProduct struct {
	// RecipyID specifies the recipy id of the recipy that it belongs to.
	RecipyID xid.ID `json:"-"`

	// ProductID specifies the product id of the product.
	ProductID xid.ID `json:"product_id"`

	// Quantity specifies how much of a product should be used in the
	// recipy.
	Quantity decimal.Decimal `json:"quantity"`
}

// FindMatching finds the matching product based on the id.
func (rp *RecipyProduct) FindMatching(products []Product) (Product, bool) {
	for _, product := range products {
		if product.ID == rp.ProductID {
			return product, true
		}
	}

	return Product{}, false
}
