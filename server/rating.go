package server

import (
	"encoding/json"
	"foodie/core"
	"foodie/db"
	"io"
	"net/http"
)

func (s *Server) GetRatings(w http.ResponseWriter, r *http.Request) {
	ratings, err := db.GetRatings(r.Context(), s.db)
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
	rid, ok := extractPathID(r, "ratingId")
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

	rating, err := db.InsertRating(r.Context(), s.db, rid, uid, rc)
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
		s.log.WithError(err).Error("creating a new rating")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	s.respondJSON(w, rating)
}

func (s *Server) UpdateRating(w http.ResponseWriter, r *http.Request) {
	rid, ok := extractPathID(r, "ratingId")
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

	rating, err := db.GetRatingByID(r.Context(), s.db, rid)
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

	if rating.UserID.Compare(uid) != 0 {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	err = db.UpdateRatingByID(r.Context(), s.db, rid, rc)
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

	rating, err = db.GetRatingByID(r.Context(), s.db, rid)
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
	rid, ok := extractPathID(r, "ratingId")
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

	rating, err := db.GetRatingByID(r.Context(), s.db, rid)
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

	if rating.UserID.Compare(uid) != 0 {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	err = db.DeleteRatingByID(r.Context(), s.db, rid)
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

func (s *Server) GetRating(w http.ResponseWriter, r *http.Request) {
	rid, ok := extractPathID(r, "ratingId")
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	rating, err := db.GetRatingByID(r.Context(), s.db, rid)
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
		s.log.WithError(err).Error("fetching rating by id")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	s.respondJSON(w, rating)
}
