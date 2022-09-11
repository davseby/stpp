package core

import (
	"github.com/rs/xid"
	"github.com/shopspring/decimal"
)

type Rating struct {
	ID       xid.ID
	RecipyID xid.ID
	UserID   xid.ID
	Score    decimal.Decimal
	Comment  string
}
