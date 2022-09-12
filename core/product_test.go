package core

import (
	"foodie/server/apierr"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func Test_ProductCore_Validate(t *testing.T) {
	tests := map[string]struct {
		ProductCore ProductCore
		Error       *apierr.Error
	}{
		"Invalid name": {
			ProductCore: ProductCore{
				Serving: Serving{
					Type:     ServingTypeMilliliters,
					Size:     decimal.NewFromInt(10),
					Calories: decimal.NewFromInt(50),
				},
			},
			Error: apierr.InvalidAttribute("name", "cannot be empty"),
		},
		"Invalid serving type": {
			ProductCore: ProductCore{
				Name: "123",
				Serving: Serving{
					Size:     decimal.NewFromInt(10),
					Calories: decimal.NewFromInt(50),
				},
			},
			Error: apierr.InvalidAttribute("type", "must be of a valid type"),
		},
		"Invalid serving size": {
			ProductCore: ProductCore{
				Name: "123",
				Serving: Serving{
					Type:     ServingTypeMilliliters,
					Calories: decimal.NewFromInt(50),
				},
			},
			Error: apierr.InvalidAttribute("size", "cannot be less or equal to 0"),
		},
		"Invalid serving calories": {
			ProductCore: ProductCore{
				Name: "123",
				Serving: Serving{
					Type:     ServingTypeMilliliters,
					Size:     decimal.NewFromInt(10),
					Calories: decimal.NewFromInt(-50),
				},
			},
			Error: apierr.InvalidAttribute("name", "cannot be less than 0"),
		},
		"Valid product core": {
			ProductCore: ProductCore{
				Name: "123",
				Serving: Serving{
					Type:     ServingTypeMilliliters,
					Size:     decimal.NewFromInt(10),
					Calories: decimal.NewFromInt(50),
				},
			},
		},
	}

	for name, test := range tests {
		test := test

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, test.Error, test.ProductCore.Validate())
		})
	}
}
