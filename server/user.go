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
		s.log.WithError(err).Error("generating bcrypt password")
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

func (s *Server) GetUser(w http.ResponseWriter, r *http.Request) {
	id, err := extractIDFromPath(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := db.GetUserByID(r.Context(), s.db, id)
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
		s.log.WithError(err).Error("fetching user by id")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	s.respondJSON(w, user)
}

func (s *Server) UpdateUserPassword(w http.ResponseWriter, r *http.Request) {
	id, err := extractContextData(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var input struct {
		Password    string `json:"password"`
		OldPassword string `json:"old_password"`
	}

	if err := json.Unmarshal(data, &input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid JSON object"))
		return
	}

	user, err := db.GetUserByID(r.Context(), s.db, id)
	switch err {
	case nil:
		// OK.
	case r.Context().Err():
		w.WriteHeader(http.StatusBadRequest)
		return
	default:
		s.log.WithError(err).Error("fetching user by id for password update")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(input.OldPassword))
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

	if err := core.ValidatePassword(input.Password); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	ph, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		s.log.WithError(err).Error("generating bcrypt password")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = db.UpdateUserPasswordByID(r.Context(), s.db, id, ph)
	switch err {
	case nil:
		// OK.
	case r.Context().Err():
		w.WriteHeader(http.StatusBadRequest)
		return
	default:
		s.log.WithError(err).Error("updating user password")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
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

func (s *Server) DeleteUser(admin bool) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := extractContextData(r)
		if err != nil {
			s.log.WithError(err).Error("extracting context data")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if admin {
			oid, err := extractIDFromPath(r)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			if id.Compare(oid) == 0 {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("cannot self delete if the user is admin"))
				return
			}

			id = oid
		}

		err = db.DeleteUserByID(r.Context(), s.db, id)
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

		w.WriteHeader(http.StatusNoContent)
	}
}
