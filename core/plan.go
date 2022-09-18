package core

import (
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

	// Recipies contains plan recipies.
	Recipies []PlanRecipy `json:"recipes"`
}

// Validate checks whether plan core contains valid attributes.
func (pc *PlanCore) Validate() *apierr.Error {
	if pc.Name == "" {
		return apierr.InvalidAttribute("name", "cannot be empty")
	}

	if pc.Description == "" {
		return apierr.InvalidAttribute("description", "cannot be empty")
	}

	if len(pc.Recipies) < 1 {
		return apierr.InvalidAttribute("products", "must contains at least one element")
	}

	return nil
}

// PlanRecipy maps plan recipes with the actual recipes stored in the
// system.
type PlanRecipy struct {
	// PlanID specifies the plan id.
	PlanID xid.ID `json:"-"`

	// RecipyID specifies the recipy id.
	RecipyID xid.ID `json:"recipy_id"`

	// Quantity specifies recipy count.
	Quantity uint64 `json:"quantity"`
}

// FindMatching finds the matching recipy based on the id.
func (pr *PlanRecipy) FindMatching(recipes []Recipy) (Recipy, bool) {
	for _, recipy := range recipes {
		if recipy.ID == pr.RecipyID {
			return recipy, true
		}
	}

	return Recipy{}, false
}
