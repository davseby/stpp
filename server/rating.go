package server

import (
	"encoding/json"
	"foodie/core"
	"foodie/db"
	"io"
	"net/http"
)

func (s *Server) GetRecipyRatings(w http.ResponseWriter, r *http.Request) {
	rid, ok := extractPathID(r, "recipyId")
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ratings, err := db.GetRatings(r.Context(), s.db, rid)
	switch err {
	case nil:
		// OK.
	case r.Context().Err():
		w.WriteHeader(http.StatusBadRequest)
		return
	default:
		s.log.WithError(err).Error("fetching ratings")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	s.respondJSON(w, ratings)
}

func (s *Server) CreateRating(w http.ResponseWriter, r *http.Request) {
	rid, ok := extractPathID(r, "recipyId")
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	uid, ok := extractContextUserID(r)
	if !ok {
		s.log.Error("extracting context user id data")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var rc core.RatingCore
	if err := json.Unmarshal(data, &rc); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid JSON object"))
		return
	}

	if err := rc.Validate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	recipy, err := db.GetRecipyByID(r.Context(), s.db, rid)
	switch err {
	case nil:
		// OK.
	case r.Context().Err():
		w.WriteHeader(http.StatusBadRequest)
		return
	case db.ErrNotFound:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("recipy"))
		return
	default:
		s.log.WithError(err).Error("fetching recipy by id")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if recipy.UserID.Compare(uid) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("cannot create a rating for your recipy"))
		return
	}

	rating, err := db.InsertRating(r.Context(), s.db, rid, uid, rc)
	switch err {
	case nil:
		// OK.
	case r.Context().Err():
		w.WriteHeader(http.StatusBadRequest)
		return
	case db.ErrDuplicate:
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte("rating"))
		return
	default:
		s.log.WithError(err).Error("creating a new rating")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	s.respondJSON(w, rating)
}

func (s *Server) UpdateRating(w http.ResponseWriter, r *http.Request) {
	rid, ok := extractPathID(r, "recipyId")
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	uid, ok := extractContextUserID(r)
	if !ok {
		s.log.Error("extracting context user id data")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var rc core.RatingCore
	if err := json.Unmarshal(data, &rc); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid JSON object"))
		return
	}

	rating, err := db.GetRating(r.Context(), s.db, rid, uid)
	switch err {
	case nil:
		// OK.
	case r.Context().Err():
		w.WriteHeader(http.StatusBadRequest)
		return
	case db.ErrNotFound:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("rating"))
		return
	default:
		s.log.WithError(err).Error("fetching rating by id")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = db.UpdateRating(r.Context(), s.db, rid, uid, rc)
	switch err {
	case nil:
		// OK.
	case r.Context().Err():
		w.WriteHeader(http.StatusBadRequest)
		return
	default:
		s.log.WithError(err).Error("updating rating")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	rating, err = db.GetRating(r.Context(), s.db, rid, uid)
	switch err {
	case nil:
		// OK.
	case r.Context().Err():
		w.WriteHeader(http.StatusBadRequest)
		return
	case db.ErrNotFound:
		s.log.WithError(err).Error("fetching rating by id after its update")
		w.WriteHeader(http.StatusInternalServerError)
		return
	default:
		s.log.WithError(err).Error("fetching rating by id")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	s.respondJSON(w, rating)
}
func (s *Server) DeleteRating(w http.ResponseWriter, r *http.Request) {
	rid, ok := extractPathID(r, "recipyId")
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	uid, ok := extractContextUserID(r)
	if !ok {
		s.log.Error("extracting context user id data")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err := db.GetRating(r.Context(), s.db, rid, uid)
	switch err {
	case nil:
		// OK.
	case r.Context().Err():
		w.WriteHeader(http.StatusBadRequest)
		return
	case db.ErrNotFound:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("rating"))
		return
	default:
		s.log.WithError(err).Error("fetching rating by id")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = db.DeleteRating(r.Context(), s.db, rid, uid)
	switch err {
	case nil:
		// OK.
	case r.Context().Err():
		w.WriteHeader(http.StatusBadRequest)
		return
	default:
		s.log.WithError(err).Error("deleting rating by id")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
