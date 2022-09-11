package core

import (
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
	Items       []RecipyItem
}

type RecipyItem struct {
	ID       xid.ID
	Product  Product
	Quantity decimal.Decimal
}
