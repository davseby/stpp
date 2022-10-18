package core

import (
	"foodie/server/apierr"
	"testing"

	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
)

func Test_PlanCore_Validate(t *testing.T) {
	tests := map[string]struct {
		PlanCore PlanCore
		Error    *apierr.Error
	}{
		"Invalid name": {
			PlanCore: PlanCore{
				Description: "123",
				Recipes: []PlanRecipe{
					{
						Quantity: 3,
					},
					{
						Quantity: 3,
					},
				},
			},
			Error: apierr.InvalidAttribute("name", "cannot be empty"),
		},
		"Invalid description": {
			PlanCore: PlanCore{
				Name: "123",
				Recipes: []PlanRecipe{
					{
						Quantity: 3,
					},
					{
						Quantity: 3,
					},
				},
			},
			Error: apierr.InvalidAttribute("description", "cannot be empty"),
		},
		"Invalid recipes length": {
			PlanCore: PlanCore{
				Name:        "123",
				Description: "123",
			},
			Error: apierr.InvalidAttribute("recipes", "must contain at least one element"),
		},
		"Invalid recipes quantity": {
			PlanCore: PlanCore{
				Name:        "123",
				Description: "123",
				Recipes: []PlanRecipe{
					{},
					{
						Quantity: 3,
					},
				},
			},
			Error: apierr.InvalidAttribute("recipes[0].quantity", "must be positive"),
		},
		"Valid plan core": {
			PlanCore: PlanCore{
				Name:        "123",
				Description: "123",
				Recipes: []PlanRecipe{
					{
						Quantity: 3,
					},
					{
						Quantity: 3,
					},
				},
			},
		},
	}

	for name, test := range tests {
		test := test

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, test.Error, test.PlanCore.Validate())
		})
	}
}

func Test_PlanRecipe_FindMatching(t *testing.T) {
	id := xid.New()
	pr := &PlanRecipe{
		RecipeID: id,
	}

	t.Run("not found", func(t *testing.T) {
		t.Parallel()

		rcp, found := pr.FindMatching(make([]Recipe, 4))

		assert.Empty(t, rcp)
		assert.False(t, found)
	})

	t.Run("found", func(t *testing.T) {
		t.Parallel()

		rr := []Recipe{
			{},
			{},
			{
				ID: id,
			},
			{},
		}

		rcp, found := pr.FindMatching(rr)

		assert.Equal(t, rcp, rr[2])
		assert.True(t, found)
	})
}
