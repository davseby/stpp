package core

import (
	"fmt"
	"foodie/server/apierr"
	"time"

	"github.com/rs/xid"
	"github.com/shopspring/decimal"
)

// Recipe contains recipe data.
type Recipe struct {
	RecipeCore

	// ID is a unique recipe identifier.
	ID xid.ID `json:"id"`

	// UserID specifies the user which created the recipe.
	UserID xid.ID `json:"user_id"`

	// CreatedAt specifies a time at which the object was created.
	CreatedAt time.Time `json:"created_at"`
}

// RecipeCore contains core recipe information.
type RecipeCore struct {
	// Name specifies the name of the recipe.
	Name string `json:"name"`

	// ImageURL specifies the image url for the recipe.
	ImageURL string `json:"image_url"`

	// Description provides a brief description of the recipe.
	Description string `json:"description"`

	// Products contains recipe products.
	Products []RecipeProduct `json:"products"`
}

// Validate checks whether recipe core contains valid attributes.
func (rc *RecipeCore) Validate() *apierr.Error {
	if rc.Name == "" {
		return apierr.InvalidAttribute("name", "cannot be empty")
	}

	if rc.Description == "" {
		return apierr.InvalidAttribute("description", "cannot be empty")
	}

	if len(rc.Products) < 2 {
		return apierr.InvalidAttribute("products", "must contains at least two elements")
	}

	for i, prod := range rc.Products {
		if !prod.Quantity.IsPositive() {
			return apierr.InvalidAttribute(fmt.Sprintf("products[%d].quantity", i), "must be positive")
		}
	}

	return nil
}

// RecipeProduct maps recipe product with the actual product stored in the
// system.
type RecipeProduct struct {
	// RecipeID specifies the recipe id of the recipe that it belongs to.
	RecipeID xid.ID `json:"-"`

	// ProductID specifies the product id of the product.
	ProductID xid.ID `json:"product_id"`

	// Quantity specifies how much of a product should be used in the
	// recipe.
	Quantity decimal.Decimal `json:"quantity"`
}

// FindMatching finds the matching product based on the id.
func (rp *RecipeProduct) FindMatching(products []Product) (Product, bool) {
	for _, product := range products {
		if product.ID == rp.ProductID {
			return product, true
		}
	}

	return Product{}, false
}
