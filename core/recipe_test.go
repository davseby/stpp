package core

import (
	"foodie/server/apierr"
	"testing"

	"github.com/rs/xid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func Test_RecipeCore_Validate(t *testing.T) {
	tests := map[string]struct {
		RecipeCore RecipeCore
		Error      *apierr.Error
	}{
		"Invalid name": {
			RecipeCore: RecipeCore{
				Description: "123",
				Products: []RecipeProduct{
					{
						Quantity: decimal.NewFromInt(3),
					},
					{
						Quantity: decimal.NewFromInt(3),
					},
				},
			},
			Error: apierr.InvalidAttribute("name", "cannot be empty"),
		},
		"Invalid description": {
			RecipeCore: RecipeCore{
				Name: "333",
				Products: []RecipeProduct{
					{
						Quantity: decimal.NewFromInt(3),
					},
					{
						Quantity: decimal.NewFromInt(3),
					},
				},
			},
			Error: apierr.InvalidAttribute("description", "cannot be empty"),
		},
		"Invalid products length": {
			RecipeCore: RecipeCore{
				Name:        "333",
				Description: "123",
			},
			Error: apierr.InvalidAttribute("products", "must contains at least two elements"),
		},
		"Invalid products quantity": {
			RecipeCore: RecipeCore{
				Name:        "333",
				Description: "123",
				Products: []RecipeProduct{
					{
						Quantity: decimal.NewFromInt(3),
					},
					{
						Quantity: decimal.NewFromInt(0),
					},
				},
			},
			Error: apierr.InvalidAttribute("products[1].quantity", "must be positive"),
		},
		"Valid recipe core": {
			RecipeCore: RecipeCore{
				Name:        "123",
				Description: "123",
				Products: []RecipeProduct{
					{
						Quantity: decimal.NewFromInt(3),
					},
					{
						Quantity: decimal.NewFromInt(3),
					},
				},
			},
		},
	}

	for name, test := range tests {
		test := test

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, test.Error, test.RecipeCore.Validate())
		})
	}
}

func Test_RecipeProduct_FindMatching(t *testing.T) {
	id := xid.New()
	rp := &RecipeProduct{
		ProductID: id,
	}

	t.Run("not found", func(t *testing.T) {
		t.Parallel()

		prd, found := rp.FindMatching(make([]Product, 4))

		assert.Empty(t, prd)
		assert.False(t, found)
	})

	t.Run("found", func(t *testing.T) {
		t.Parallel()

		pp := []Product{
			{},
			{},
			{
				ID: id,
			},
			{},
		}

		prd, found := rp.FindMatching(pp)

		assert.Equal(t, prd, pp[2])
		assert.True(t, found)
	})
}
