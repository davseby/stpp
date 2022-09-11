package core

import (
	"errors"

	"github.com/rs/xid"
	"github.com/shopspring/decimal"
)

type Rating struct {
	ID       xid.ID
	RecipyID xid.ID
	UserID   xid.ID
	RatingCore
}

type RatingCore struct {
	Score   decimal.Decimal
	Comment string
}

func (rc *RatingCore) Validate() error {
	if rc.Score.IsNegative() {
		return errors.New("rating score cannot be less than 0")
	}

	if rc.Comment == "" {
		return errors.New("rating comment cannot be empty")
	}

	return nil
}
