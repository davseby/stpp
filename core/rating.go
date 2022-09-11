package core

import (
	"errors"

	"github.com/rs/xid"
	"github.com/shopspring/decimal"
)

type Rating struct {
	RecipyID xid.ID `json:"recipy_id"`
	UserID   xid.ID `json:"user_id"`
	RatingCore
}

type RatingCore struct {
	Score   decimal.Decimal `json:"score"`
	Comment string          `json:"comment"`
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
