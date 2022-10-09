package core

import (
	"fmt"
	"foodie/server/apierr"
	"time"

	"github.com/rs/xid"
)

// Plan contains plan data.
type Plan struct {
	PlanCore

	// ID is a unique plan identifier.
	ID xid.ID `json:"id"`

	// UserID specifies the user which created the plan.
	UserID xid.ID `json:"user_id"`

	// CreatedAt specifies a time at which the object was created.
	CreatedAt time.Time `json:"created_at"`
}

// PlanCore contains core plan information.
type PlanCore struct {
	// Plan specifies the name of the plan.
	Name string `json:"name"`

	// Description provides a brief description of the plan.
	Description string `json:"description"`

	// Recipes contains plan recipes.
	Recipes []PlanRecipe `json:"recipes"`
}

// Validate checks whether plan core contains valid attributes.
func (pc *PlanCore) Validate() *apierr.Error {
	if pc.Name == "" {
		return apierr.InvalidAttribute("name", "cannot be empty")
	}

	if pc.Description == "" {
		return apierr.InvalidAttribute("description", "cannot be empty")
	}

	if len(pc.Recipes) < 1 {
		return apierr.InvalidAttribute("recipes", "must contains at least one element")
	}

	for i, rec := range pc.Recipes {
		if rec.Quantity == 0 {
			return apierr.InvalidAttribute(fmt.Sprintf("recipes[%d].quantity", i), "must be positive")
		}
	}

	return nil
}

// PlanRecipe maps plan recipes with the actual recipes stored in the
// system.
type PlanRecipe struct {
	// PlanID specifies the plan id.
	PlanID xid.ID `json:"-"`

	// RecipeID specifies the recipe id.
	RecipeID xid.ID `json:"recipe_id"`

	// Quantity specifies recipe count.
	Quantity uint64 `json:"quantity"`
}

// FindMatching finds the matching recipe based on the id.
func (pr *PlanRecipe) FindMatching(recipes []Recipe) (Recipe, bool) {
	for _, recipe := range recipes {
		if recipe.ID == pr.RecipeID {
			return recipe, true
		}
	}

	return Recipe{}, false
}
