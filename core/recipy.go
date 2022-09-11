package core

import (
	"errors"

	"github.com/rs/xid"
	"github.com/shopspring/decimal"
)

type Recipy struct {
	ID     xid.ID `json:"id"`
	UserID xid.ID `json:"user_id"`
	RecipyCore
}

type RecipyCore struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Products    []RecipyProduct `json:"products"`
}

func (rc *RecipyCore) Validate() error {
	if rc.Name == "" {
		return errors.New("name cannot be empty")
	}

	if rc.Description == "" {
		return errors.New("description cannot be empty")
	}

	if len(rc.Products) < 2 {
		return errors.New("recipy must contain at least two product")
	}

	return nil
}

type RecipyProduct struct {
	RecipyID  xid.ID          `json:"-"`
	ProductID xid.ID          `json:"product_id"`
	Quantity  decimal.Decimal `json:"quantity"`
}
