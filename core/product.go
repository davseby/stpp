package core

import (
	"foodie/server/apierr"
	"time"

	"github.com/rs/xid"
	"github.com/shopspring/decimal"
)

// ServingType specifies the product single serving type.
type ServingType string

const (
	// ServingTypeGrams specifies the serving in grams.
	ServingTypeGrams ServingType = "grams"

	// ServingTypeMilliliters specifies the serving in milliliters.
	ServingTypeMilliliters ServingType = "milliliters"

	// ServingTypeUnits specifies the serving in units.
	ServingTypeUnits ServingType = "units"
)

// Product contains product data.
type Product struct {
	ProductCore

	// ID is a unique product identifier.
	ID xid.ID `json:"id"`

	// CreatedAt specifies a time at which the object was created.
	CreatedAt time.Time `json:"created_at"`
}

// ProductCore contains core product information.
type ProductCore struct {
	// Name specifies the name of the product.
	Name string `json:"name"`

	// Serving specifies the serving information of the product.
	Serving Serving `json:"serving"`
}

// Serving specifies the serving information of the product.
type Serving struct {
	// Types specifies the serving measurement type.
	Type ServingType `json:"type"`

	// Size specifies the amount in a single serving.
	Size decimal.Decimal `json:"size"`

	// Calories specifies how many calories are in the single serving.
	Calories int `json:"calories"`
}

// Validate checks whether product core contains valid attributes.
func (pc *ProductCore) Validate() *apierr.Error {
	if pc.Name == "" {
		return apierr.InvalidAttribute("name", "cannot be empty")
	}

	switch pc.Serving.Type {
	case ServingTypeGrams, ServingTypeMilliliters, ServingTypeUnits:
	default:
		return apierr.InvalidAttribute("type", "must be of a valid type")
	}

	if pc.Serving.Size.Cmp(decimal.Zero) <= 0 {
		return apierr.InvalidAttribute("size", "cannot be less or equal to 0")
	}

	if pc.Serving.Calories < 0 {
		return apierr.InvalidAttribute("calories", "cannot be less than 0")
	}

	return nil
}
