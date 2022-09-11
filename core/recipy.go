package core

import (
	"errors"

	"github.com/rs/xid"
	"github.com/shopspring/decimal"
)

type Recipy struct {
	ID     xid.ID
	UserID xid.ID
	RecipyCore
}

type RecipyCore struct {
	Name        string
	Description string
	Products    []RecipyProduct
}

func (rc *RecipyCore) Validate() error {
	if rc.Name == "" {
		return errors.New("name cannot be empty")
	}

	if rc.Description == "" {
		return errors.New("description cannot be empty")
	}

	return nil
}

type RecipyProduct struct {
	RecipyID  xid.ID
	ProductID xid.ID
	Quantity  decimal.Decimal
}
