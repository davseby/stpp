package server

import (
	"encoding/json"
	"foodie/core"
	"foodie/db"
	"io"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func (s *Server) Register(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var uc core.UserCore
	if err := json.Unmarshal(data, &uc); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid JSON object"))
		return
	}

	if err := uc.Validate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	_, err = db.GetUserByName(r.Context(), s.db, uc.Name)
	switch err {
	case db.ErrNotFound:
		// OK.
	case nil:
		w.WriteHeader(http.StatusConflict)
		return
	case r.Context().Err():
		w.WriteHeader(http.StatusBadRequest)
		return
	default:
		s.log.WithError(err).Error("fetching user by name")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ph, err := bcrypt.GenerateFromPassword([]byte(uc.Password), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user, err := db.InsertUser(
		r.Context(),
		s.db,
		uc.Name,
		ph,
		true,
	)
	switch err {
	case nil:
		// OK.
	case r.Context().Err():
		w.WriteHeader(http.StatusBadRequest)
		return
	default:
		s.log.WithError(err).Error("creating a new user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	signedJWT, err := core.IssueJWT(s.secret, user.ID, user.Admin, time.Now())
	if err != nil {
		s.log.WithError(err).Error("creating jwt")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	s.respondJSON(w, struct {
		User        *core.User `json:"user"`
		AccessToken string     `json:"access_token"`
	}{
		User:        user,
		AccessToken: string(signedJWT),
	})
}

func (s *Server) Login(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var uc core.UserCore
	if err := json.Unmarshal(data, &uc); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid JSON object"))
		return
	}

	user, err := db.GetUserByName(r.Context(), s.db, uc.Name)
	switch err {
	case nil:
		// OK.
	case r.Context().Err(), db.ErrNotFound:
		w.WriteHeader(http.StatusBadRequest)
		return
	default:
		s.log.WithError(err).Error("fetching user by name")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(uc.Password))
	switch err {
	case nil:
		// OK.
	case bcrypt.ErrMismatchedHashAndPassword:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid password"))
		return
	default:
		s.log.WithError(err).Error("comparing bcrypt password")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	signedJWT, err := core.IssueJWT(s.secret, user.ID, user.Admin, time.Now())
	if err != nil {
		s.log.WithError(err).Error("creating jwt")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	s.respondJSON(w, struct {
		User        *core.User `json:"user"`
		AccessToken string     `json:"access_token"`
	}{
		User:        user,
		AccessToken: string(signedJWT),
	})
}

func (s *Server) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := db.GetUsers(r.Context(), s.db)
	switch err {
	case nil:
		// OK.
	case r.Context().Err():
		w.WriteHeader(http.StatusBadRequest)
		return
	default:
		s.log.WithError(err).Error("fetching users")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	s.respondJSON(w, users)
}

func (s *Server) CreateAdminUser(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var uc core.UserCore
	if err := json.Unmarshal(data, &uc); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid JSON object"))
		return
	}

	if err := uc.Validate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	ph, err := bcrypt.GenerateFromPassword([]byte(uc.Password), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user, err := db.InsertUser(
		r.Context(),
		s.db,
		uc.Name,
		ph,
		true,
	)
	switch err {
	case nil:
		// OK.
	case r.Context().Err():
		w.WriteHeader(http.StatusBadRequest)
		return
	default:
		s.log.WithError(err).Error("creating a new user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	s.respondJSON(w, user)
}
