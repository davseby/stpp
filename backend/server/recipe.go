package server

import (
	"context"
	"encoding/json"
	"foodie/core"
	"foodie/db"
	"foodie/server/apierr"
	"io"
	"net/http"
)

// CreateRecipe creates a recipe.
func (s *Server) CreateRecipe(w http.ResponseWriter, r *http.Request) {
	uid, aerr := s.extractContextUserID(r)
	if aerr != nil {
		aerr.Respond(w)
		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		apierr.MalformedDataInput(apierr.DataTypeRequestBody).Respond(w)
		return
	}

	var rc core.RecipeCore
	if err := json.Unmarshal(data, &rc); err != nil {
		apierr.MalformedDataInput(apierr.DataTypeJSON).Respond(w)
		return
	}

	if aerr := s.validateRecipeCore(r.Context(), rc); aerr != nil {
		aerr.Respond(w)
		return
	}

	rec, err := db.InsertRecipe(r.Context(), s.db, uid, rc)
	switch err {
	case nil:
		// OK.
	case r.Context().Err():
		apierr.Context().Respond(w)
		return
	case db.ErrNotFound:
		apierr.NotFound("product").Respond(w)
		return
	default:
		s.log.WithError(err).Error("inserting a recipe")
		apierr.Database().Respond(w)

		return
	}

	s.respondJSON(w, rec)
}

// GetRecipes retrieves all recipes.
func (s *Server) GetRecipes(w http.ResponseWriter, r *http.Request) {
	rr, err := db.GetRecipes(r.Context(), s.db)
	switch err {
	case nil:
		// OK.
	case r.Context().Err():
		apierr.Context().Respond(w)
		return
	default:
		s.log.WithError(err).Error("fetching recipes")
		apierr.Database().Respond(w)

		return
	}

	s.respondJSON(w, rr)
}

// GetUserRecipes retrieves user recipes.
func (s *Server) GetUserRecipes(w http.ResponseWriter, r *http.Request) {
	uid, aerr := s.extractPathID(r, "userID")
	if aerr != nil {
		aerr.Respond(w)
		return
	}

	rr, err := db.GetRecipesByUserID(r.Context(), s.db, uid)
	switch err {
	case nil:
		// OK.
	case r.Context().Err():
		apierr.Context().Respond(w)
		return
	default:
		s.log.WithError(err).Error("fetching recipes")
		apierr.Database().Respond(w)

		return
	}

	s.respondJSON(w, rr)
}

// GetRecipe retrieves a single recipe by its id.
func (s *Server) GetRecipe(w http.ResponseWriter, r *http.Request) {
	rid, aerr := s.extractPathID(r, "recipeID")
	if aerr != nil {
		aerr.Respond(w)
		return
	}

	rec, err := db.GetRecipeByID(r.Context(), s.db, rid)
	switch err {
	case nil:
		// OK.
	case r.Context().Err():
		apierr.Context().Respond(w)
		return
	case db.ErrNotFound:
		apierr.NotFound("recipe").Respond(w)
		return
	default:
		s.log.WithError(err).Error("fetching recipe by id")
		apierr.Database().Respond(w)

		return
	}

	s.respondJSON(w, rec)
}

// UpdateRecipe updates existing recipe by its id. The recipe can be
// updated only by the user which created it.
func (s *Server) UpdateRecipe(w http.ResponseWriter, r *http.Request) {
	rid, aerr := s.extractPathID(r, "recipeID")
	if aerr != nil {
		aerr.Respond(w)
		return
	}

	uid, aerr := s.extractContextUserID(r)
	if aerr != nil {
		aerr.Respond(w)
		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		apierr.MalformedDataInput(apierr.DataTypeRequestBody).Respond(w)
		return
	}

	rec, err := db.GetRecipeByID(r.Context(), s.db, rid)
	switch err {
	case nil:
		// OK.
	case r.Context().Err():
		apierr.Context().Respond(w)
		return
	case db.ErrNotFound:
		apierr.NotFound("recipe").Respond(w)
		return
	default:
		s.log.WithError(err).Error("fetching recipe by id")
		apierr.Database().Respond(w)

		return
	}

	if rec.UserID.Compare(uid) != 0 {
		apierr.Forbidden().Respond(w)
		return
	}

	var rc core.RecipeCore
	if err := json.Unmarshal(data, &rc); err != nil {
		apierr.MalformedDataInput(apierr.DataTypeJSON).Respond(w)
		return
	}

	if aerr := s.validateRecipeCore(r.Context(), rc); aerr != nil {
		aerr.Respond(w)
		return
	}

	rec, err = db.UpdateRecipeByID(r.Context(), s.db, rid, rc)
	switch err {
	case nil:
		// OK.
	case r.Context().Err():
		apierr.Context().Respond(w)
		return
	case db.ErrNotFound:
		apierr.NotFound("product").Respond(w)
		return
	default:
		s.log.WithError(err).Error("creating a new recipe")
		apierr.Database().Respond(w)

		return
	}

	s.respondJSON(w, rec)
}

// DeleteRecipe deletes existing recipe by its id. The recipe can be deleted
// only by an admin or the user that created it.
func (s *Server) DeleteRecipe(w http.ResponseWriter, r *http.Request) {
	rid, aerr := s.extractPathID(r, "recipeID")
	if aerr != nil {
		aerr.Respond(w)
		return
	}

	adm, aerr := s.extractContextAdmin(r)
	if aerr != nil {
		aerr.Respond(w)
		return
	}

	if !adm {
		uid, aerr := s.extractContextUserID(r)
		if aerr != nil {
			aerr.Respond(w)
			return
		}

		rec, err := db.GetRecipeByID(r.Context(), s.db, rid)
		switch err {
		case nil:
			// OK.
		case r.Context().Err():
			apierr.Context().Respond(w)
			return
		case db.ErrNotFound:
			apierr.NotFound("recipe").Respond(w)
			return
		default:
			s.log.WithError(err).Error("fetching recipe by id")
			apierr.Database().Respond(w)

			return
		}

		if rec.UserID.Compare(uid) != 0 {
			apierr.Forbidden().Respond(w)
			return
		}
	}

	plans, err := db.GetPlanRecipesByRecipeID(r.Context(), s.db, rid)
	switch err {
	case nil:
		// OK.
	case r.Context().Err():
		apierr.Context().Respond(w)
		return
	default:
		s.log.WithError(err).Error("getting plan recipes by recipe id")
		apierr.Database().Respond(w)

		return
	}

	if len(plans) > 0 {
		apierr.Conflict("recipe in use").Respond(w)
		return
	}

	err = db.DeleteRecipeByID(r.Context(), s.db, rid)
	switch err {
	case nil:
		// OK.
	case r.Context().Err():
		apierr.Context().Respond(w)
		return
	default:
		s.log.WithError(err).Error("deleting recipe by id")
		apierr.Database().Respond(w)

		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// validateRecipeCore validates recipe core attributes.
func (s *Server) validateRecipeCore(ctx context.Context, rc core.RecipeCore) *apierr.Error {
	pp, err := db.GetProducts(ctx, s.db)
	switch err {
	case nil:
		// OK.
	case ctx.Err():
		return apierr.Context()
	default:
		s.log.WithError(err).Error("fetching products")
		return apierr.Database()
	}

	for _, rp := range rc.Products {
		if _, ok := rp.FindMatching(pp); !ok {
			return apierr.NotFound("product")
		}
	}

	return rc.Validate()
}
