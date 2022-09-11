package server

import (
	"encoding/json"
	"foodie/core"
	"foodie/db"
	"io"
	"net/http"
	"time"

	"github.com/rs/xid"
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
		w.Write([]byte("user"))
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
	uid, ok := extractPathID(r, "userId")
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := db.GetUserByID(r.Context(), s.db, uid)
	switch err {
	case nil:
		// OK.
	case r.Context().Err():
		w.WriteHeader(http.StatusBadRequest)
		return
	case db.ErrNotFound:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("user"))
		return
	default:
		s.log.WithError(err).Error("fetching user by id")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	s.respondJSON(w, user)
}

func (s *Server) UpdateUserPassword(w http.ResponseWriter, r *http.Request) {
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

	var input struct {
		Password    string `json:"password"`
		OldPassword string `json:"old_password"`
	}

	if err := json.Unmarshal(data, &input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid JSON object"))
		return
	}

	user, err := db.GetUserByID(r.Context(), s.db, uid)
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

	err = db.UpdateUserPasswordByID(r.Context(), s.db, uid, ph)
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

func (s *Server) DeleteUser(super bool) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var uid xid.ID

		var ok bool

		if !super {
			uid, ok = extractContextUserID(r)
			if !ok {
				s.log.Error("extracting context user id data")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else {
			uid, ok = extractPathID(r, "userId")
			if !ok {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		admin, ok := extractContextAdmin(r)
		if !ok {
			s.log.Error("extracting context admin data")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if admin {
			user, err := db.GetUserByID(r.Context(), s.db, uid)
			switch err {
			case nil:
				// OK.
			case r.Context().Err():
				w.WriteHeader(http.StatusBadRequest)
				return
			case db.ErrNotFound:
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte("user"))
				return
			default:
				s.log.WithError(err).Error("fetching user by id")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			if user.Name == core.RootAdminName {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("cannot delete root admin"))
				return
			}
		}

		err := db.DeleteUserByID(r.Context(), s.db, uid)
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
