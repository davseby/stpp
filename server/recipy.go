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

// CreateRecipy creates a recipy.
func (s *Server) CreateRecipy(w http.ResponseWriter, r *http.Request) {
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

	var rc core.RecipyCore
	if err := json.Unmarshal(data, &rc); err != nil {
		apierr.MalformedDataInput(apierr.DataTypeJSON).Respond(w)
		return
	}

	if aerr := s.validateRecipyCore(r.Context(), rc); aerr != nil {
		aerr.Respond(w)
		return
	}

	rec, err := db.InsertRecipy(r.Context(), s.db, uid, rc)
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
		s.log.WithError(err).Error("inserting a recipy")
		apierr.Database().Respond(w)
		return
	}

	s.respondJSON(w, rec)
}

// GetPublicRecipes retrieves all public recipes.
func (s *Server) GetPublicRecipes(w http.ResponseWriter, r *http.Request) {
	rr, err := db.GetRecipes(r.Context(), s.db, false)
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

// GetUserRecipes retrieves user recipes. If the requested recipes are from
// the user which created them, the private recipes are included.
func (s *Server) GetUserRecipes(w http.ResponseWriter, r *http.Request) {
	uid, aerr := s.extractPathID(r, "userId")
	if aerr != nil {
		aerr.Respond(w)
		return
	}

	var ip bool
	if cuid, aerr := s.extractContextUserID(r); aerr == nil && cuid.Compare(uid) == 0 {
		ip = true
	}

	rr, err := db.GetRecipesByUserID(r.Context(), s.db, uid, ip)
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

// GetRecipy retrieves a single recipy by its id. The private recipes can be
// retrieved only by the user that created them.
func (s *Server) GetRecipy(w http.ResponseWriter, r *http.Request) {
	rid, aerr := s.extractPathID(r, "recipyId")
	if aerr != nil {
		aerr.Respond(w)
		return
	}

	uid, aerr := s.extractContextUserID(r)
	if aerr != nil {
		aerr.Respond(w)
		return
	}

	rec, err := db.GetRecipyByID(r.Context(), s.db, rid, true)
	switch err {
	case nil:
		// OK.
	case r.Context().Err():
		apierr.Context().Respond(w)
		return
	case db.ErrNotFound:
		apierr.NotFound("recipy").Respond(w)
		return
	default:
		s.log.WithError(err).Error("fetching recipy by id")
		apierr.Database().Respond(w)
		return
	}

	// In case the user retrieves other users private recipy, we should
	// protect it by saying that it was not found.
	if rec.Private && rec.UserID.Compare(uid) != 0 {
		apierr.NotFound("recipy").Respond(w)
		return
	}

	s.respondJSON(w, rec)
}

// UpdateRecipy updates existing recipy by its id. The recipy can be
// updated only by the user which created it.
func (s *Server) UpdateRecipy(w http.ResponseWriter, r *http.Request) {
	rid, aerr := s.extractPathID(r, "recipyId")
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

	rec, err := db.GetRecipyByID(r.Context(), s.db, rid, true)
	switch err {
	case nil:
		// OK.
	case r.Context().Err():
		apierr.Context().Respond(w)
		return
	case db.ErrNotFound:
		apierr.NotFound("recipy").Respond(w)
		return
	default:
		s.log.WithError(err).Error("fetching recipy by id")
		apierr.Database().Respond(w)
		return
	}

	if rec.UserID.Compare(uid) != 0 {
		// In case the user retrieves other users private recipy, we
		// should protect it by saying that it was not found.
		if rec.Private {
			apierr.NotFound("recipy").Respond(w)
			return
		}

		apierr.Forbidden().Respond(w)
		return
	}

	var rc core.RecipyCore
	if err := json.Unmarshal(data, &rc); err != nil {
		apierr.MalformedDataInput(apierr.DataTypeJSON).Respond(w)
		return
	}

	if aerr := s.validateRecipyCore(r.Context(), rc); aerr != nil {
		aerr.Respond(w)
		return
	}

	rec, err = db.UpdateRecipyByID(r.Context(), s.db, rid, rc)
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
		s.log.WithError(err).Error("creating a new recipy")
		apierr.Database().Respond(w)
		return
	}

	s.respondJSON(w, rec)
}

// DeleteRecipy deletes existing recipy by its id. The recipy can be deleted
// only by an admin or the user that created it.
func (s *Server) DeleteRecipy(w http.ResponseWriter, r *http.Request) {
	rid, aerr := s.extractPathID(r, "recipyId")
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

		rec, err := db.GetRecipyByID(r.Context(), s.db, rid, true)
		switch err {
		case nil:
			// OK.
		case r.Context().Err():
			apierr.Context().Respond(w)
			return
		case db.ErrNotFound:
			apierr.NotFound("recipy").Respond(w)
			return
		default:
			s.log.WithError(err).Error("fetching recipy by id")
			apierr.Database().Respond(w)
			return
		}

		if rec.UserID.Compare(uid) != 0 {
			// In case the user retrieves other users private recipy, we
			// should protect it by saying that it was not found.
			if rec.Private {
				apierr.NotFound("recipy").Respond(w)
				return
			}

			apierr.Forbidden().Respond(w)
			return
		}
	}

	err := db.DeleteRecipyByID(r.Context(), s.db, rid)
	switch err {
	case nil:
		// OK.
	case r.Context().Err():
		apierr.Context().Respond(w)
		return
	default:
		s.log.WithError(err).Error("deleting recipy by id")
		apierr.Database().Respond(w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// validateRecipyCore validates recipy core attributes.
func (s *Server) validateRecipyCore(ctx context.Context, rc core.RecipyCore) *apierr.Error {
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
