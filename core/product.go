package core

import (
	"errors"

	"github.com/rs/xid"
	"github.com/shopspring/decimal"
)

type ServingType string

const (
	ServingTypeGrams       ServingType = "grams"
	ServingTypeMilliliters ServingType = "milliliters"
	ServingTypeUnits       ServingType = "units"
)

type Product struct {
	ID xid.ID `json:"id"`
	ProductCore
}

type ProductCore struct {
	Name    string  `json:"name"`
	Serving Serving `json:"serving"`
}

type Serving struct {
	Type     ServingType     `json:"type"`
	Size     decimal.Decimal `json:"size"`
	Calories decimal.Decimal `json:"calories"`
}

func (pc *ProductCore) Validate() error {
	if pc.Name == "" {
		return errors.New("name cannot be empty")
	}

	switch pc.Serving.Type {
	case ServingTypeGrams, ServingTypeMilliliters, ServingTypeUnits:
	default:
		return errors.New("invalid serving type")
	}

	if pc.Serving.Size.Cmp(decimal.Zero) <= 0 {
		return errors.New("serving size cannot be less or equal to 0")
	}

	if pc.Serving.Calories.IsNegative() {
		return errors.New("serving calories cannot be less than 0")
	}

	return nil
}
