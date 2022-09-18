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

// CreatePlan creates a plan.
func (s *Server) CreatePlan(w http.ResponseWriter, r *http.Request) {
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

	var pc core.PlanCore
	if err := json.Unmarshal(data, &pc); err != nil {
		apierr.MalformedDataInput(apierr.DataTypeJSON).Respond(w)
		return
	}

	if aerr := s.validatePlanCore(r.Context(), pc); aerr != nil {
		aerr.Respond(w)
		return
	}

	pl, err := db.InsertPlan(r.Context(), s.db, uid, pc)
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
		s.log.WithError(err).Error("inserting a plan")
		apierr.Database().Respond(w)
		return
	}

	s.respondJSON(w, pl)
}

// GetPlans retrieves all plans.
func (s *Server) GetPlans(w http.ResponseWriter, r *http.Request) {
	pp, err := db.GetPlans(r.Context(), s.db)
	switch err {
	case nil:
		// OK.
	case r.Context().Err():
		apierr.Context().Respond(w)
		return
	default:
		s.log.WithError(err).Error("fetching plans")
		apierr.Database().Respond(w)
		return
	}

	s.respondJSON(w, pp)
}

// GetUserPlans retrieves user plans.
func (s *Server) GetUserPlans(w http.ResponseWriter, r *http.Request) {
	uid, aerr := s.extractPathID(r, "userID")
	if aerr != nil {
		aerr.Respond(w)
		return
	}

	pp, err := db.GetPlansByUserID(r.Context(), s.db, uid)
	switch err {
	case nil:
		// OK.
	case r.Context().Err():
		apierr.Context().Respond(w)
		return
	default:
		s.log.WithError(err).Error("fetching plans")
		apierr.Database().Respond(w)
		return
	}

	s.respondJSON(w, pp)
}

// GetPlan retrieves a single plan by its id.
func (s *Server) GetPlan(w http.ResponseWriter, r *http.Request) {
	pid, aerr := s.extractPathID(r, "planID")
	if aerr != nil {
		aerr.Respond(w)
		return
	}

	pl, err := db.GetPlanByID(r.Context(), s.db, pid)
	switch err {
	case nil:
		// OK.
	case r.Context().Err():
		apierr.Context().Respond(w)
		return
	case db.ErrNotFound:
		apierr.NotFound("plan").Respond(w)
		return
	default:
		s.log.WithError(err).Error("fetching planby id")
		apierr.Database().Respond(w)
		return
	}

	s.respondJSON(w, pl)
}

// UpdatePlan updates existing plan by its id. The plan can be
// updated only by the user which created it.
func (s *Server) UpdatePlan(w http.ResponseWriter, r *http.Request) {
	pid, aerr := s.extractPathID(r, "planID")
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

	pl, err := db.GetPlanByID(r.Context(), s.db, pid)
	switch err {
	case nil:
		// OK.
	case r.Context().Err():
		apierr.Context().Respond(w)
		return
	case db.ErrNotFound:
		apierr.NotFound("plan").Respond(w)
		return
	default:
		s.log.WithError(err).Error("fetching plan by id")
		apierr.Database().Respond(w)
		return
	}

	if pl.UserID.Compare(uid) != 0 {
		apierr.Forbidden().Respond(w)
		return
	}

	var pc core.PlanCore
	if err := json.Unmarshal(data, &pc); err != nil {
		apierr.MalformedDataInput(apierr.DataTypeJSON).Respond(w)
		return
	}

	if aerr := s.validatePlanCore(r.Context(), pc); aerr != nil {
		aerr.Respond(w)
		return
	}

	pl, err = db.UpdatePlanByID(r.Context(), s.db, pid, pc)
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
		s.log.WithError(err).Error("creating a new plan")
		apierr.Database().Respond(w)
		return
	}

	s.respondJSON(w, pl)
}

// DeletePlan deletes existing plan by its id. The plan can be deleted
// only by an admin or the user that created it.
func (s *Server) DeletePlan(w http.ResponseWriter, r *http.Request) {
	pid, aerr := s.extractPathID(r, "planID")
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

		pl, err := db.GetPlanByID(r.Context(), s.db, pid)
		switch err {
		case nil:
			// OK.
		case r.Context().Err():
			apierr.Context().Respond(w)
			return
		case db.ErrNotFound:
			apierr.NotFound("plan").Respond(w)
			return
		default:
			s.log.WithError(err).Error("fetching plan by id")
			apierr.Database().Respond(w)
			return
		}

		if pl.UserID.Compare(uid) != 0 {
			apierr.Forbidden().Respond(w)
			return
		}
	}

	err := db.DeletePlanByID(r.Context(), s.db, pid)
	switch err {
	case nil:
		// OK.
	case r.Context().Err():
		apierr.Context().Respond(w)
		return
	default:
		s.log.WithError(err).Error("deleting plan by id")
		apierr.Database().Respond(w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// validatePlanCore validates plan core attributes.
func (s *Server) validatePlanCore(ctx context.Context, pc core.PlanCore) *apierr.Error {
	rr, err := db.GetRecipes(ctx, s.db)
	switch err {
	case nil:
		// OK.
	case ctx.Err():
		return apierr.Context()
	default:
		s.log.WithError(err).Error("fetching recipes")
		return apierr.Database()
	}

	for _, pr := range pc.Recipes {
		if _, ok := pr.FindMatching(rr); !ok {
			return apierr.NotFound("recipy")
		}
	}

	return pc.Validate()
}
