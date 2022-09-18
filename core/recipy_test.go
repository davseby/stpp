package core

import (
	"foodie/server/apierr"
	"testing"

	"github.com/rs/xid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func Test_RecipyCore_Validate(t *testing.T) {
	tests := map[string]struct {
		RecipyCore RecipyCore
		Error      *apierr.Error
	}{
		"Invalid name": {
			RecipyCore: RecipyCore{
				Description: "123",
				Products: []RecipyProduct{
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
			RecipyCore: RecipyCore{
				Name: "333",
				Products: []RecipyProduct{
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
			RecipyCore: RecipyCore{
				Name:        "333",
				Description: "123",
			},
			Error: apierr.InvalidAttribute("products", "must contains at least two elements"),
		},
		"Invalid products quantity": {
			RecipyCore: RecipyCore{
				Name:        "333",
				Description: "123",
				Products: []RecipyProduct{
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
		"Valid recipy core": {
			RecipyCore: RecipyCore{
				Name:        "123",
				Description: "123",
				Products: []RecipyProduct{
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

			assert.Equal(t, test.Error, test.RecipyCore.Validate())
		})
	}
}

func Test_RecipyProduct_FindMatching(t *testing.T) {
	id := xid.New()
	rp := &RecipyProduct{
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
