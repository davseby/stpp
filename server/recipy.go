package server

import (
	"encoding/json"
	"foodie/core"
	"foodie/db"
	"io"
	"net/http"
)

func (s *Server) GetRecipies(w http.ResponseWriter, r *http.Request) {
	recipies, err := db.GetRecipies(r.Context(), s.db)
	switch err {
	case nil:
		// OK.
	case r.Context().Err():
		w.WriteHeader(http.StatusBadRequest)
		return
	default:
		s.log.WithError(err).Error("fetching recipies")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	s.respondJSON(w, recipies)
}

func (s *Server) CreateRecipy(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	uid, ok := extractContextUserID(r)
	if !ok {
		s.log.Error("extracting context user id data")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var rc core.RecipyCore
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

	recipy, err := db.InsertRecipy(r.Context(), s.db, uid, rc)
	switch err {
	case nil:
		// OK.
	case r.Context().Err():
		w.WriteHeader(http.StatusBadRequest)
		return
	case core.ErrInvalidProduct:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	default:
		s.log.WithError(err).Error("creating a new recipy")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	s.respondJSON(w, recipy)
}

func (s *Server) UpdateRecipy(w http.ResponseWriter, r *http.Request) {
	id, ok := extractPathID(r)
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

	uid, ok := extractContextUserID(r)
	if !ok {
		s.log.Error("extracting context user id data")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	recipy, err := db.GetRecipyByID(r.Context(), s.db, id)
	switch err {
	case nil:
		// OK.
	case r.Context().Err():
		w.WriteHeader(http.StatusBadRequest)
		return
	case db.ErrNotFound:
		w.WriteHeader(http.StatusNotFound)
		return
	default:
		s.log.WithError(err).Error("fetching recipy by id")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if recipy.UserID.Compare(uid) != 0 {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	var rc core.RecipyCore
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

	err = db.UpdateRecipyByID(r.Context(), s.db, id, rc)
	switch err {
	case nil:
		// OK.
	case r.Context().Err():
		w.WriteHeader(http.StatusBadRequest)
		return
	case core.ErrInvalidProduct:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	default:
		s.log.WithError(err).Error("creating a new recipy")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	recipy, err = db.GetRecipyByID(r.Context(), s.db, id)
	switch err {
	case nil:
		// OK.
	case r.Context().Err():
		w.WriteHeader(http.StatusBadRequest)
		return
	case db.ErrNotFound:
		w.WriteHeader(http.StatusNotFound)
		return
	default:
		s.log.WithError(err).Error("fetching recipy by id")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	s.respondJSON(w, recipy)
}

func (s *Server) DeleteRecipy(w http.ResponseWriter, r *http.Request) {
	id, ok := extractPathID(r)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	admin, ok := extractContextAdmin(r)
	if !ok {
		s.log.Error("extracting context admin data")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !admin {
		uid, ok := extractContextUserID(r)
		if !ok {
			s.log.Error("extracting context user id data")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		recipy, err := db.GetRecipyByID(r.Context(), s.db, id)
		switch err {
		case nil:
			// OK.
		case r.Context().Err():
			w.WriteHeader(http.StatusBadRequest)
			return
		case db.ErrNotFound:
			w.WriteHeader(http.StatusNotFound)
			return
		default:
			s.log.WithError(err).Error("fetching recipy by id")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if recipy.UserID.Compare(uid) != 0 {
			w.WriteHeader(http.StatusForbidden)
			return
		}
	}

	err := db.DeleteRecipyByID(r.Context(), s.db, id)
	switch err {
	case nil:
		// OK.
	case r.Context().Err():
		w.WriteHeader(http.StatusBadRequest)
		return
	default:
		s.log.WithError(err).Error("deleting recipy by id")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) GetRecipy(w http.ResponseWriter, r *http.Request) {
	id, ok := extractPathID(r)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	recipy, err := db.GetRecipyByID(r.Context(), s.db, id)
	switch err {
	case nil:
		// OK.
	case r.Context().Err():
		w.WriteHeader(http.StatusBadRequest)
		return
	case db.ErrNotFound:
		w.WriteHeader(http.StatusNotFound)
		return
	default:
		s.log.WithError(err).Error("fetching recipy by id")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	s.respondJSON(w, recipy)
}
