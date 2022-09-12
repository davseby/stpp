package server

import (
	"encoding/json"
	"foodie/core"
	"foodie/db"
	"foodie/server/apierr"
	"io"
	"net/http"
	"time"

	"github.com/rs/xid"
	"golang.org/x/crypto/bcrypt"
)

// CreateAdminUser creates an admin user with provided credentials.
func (s *Server) CreateAdminUser(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		apierr.MalformedDataInput(apierr.DataTypeRequestBody).Respond(w)
		return
	}

	var ui core.UserInput
	if err := json.Unmarshal(data, &ui); err != nil {
		apierr.MalformedDataInput(apierr.DataTypeJSON).Respond(w)
		return
	}

	if aerr := ui.Validate(); aerr != nil {
		aerr.Respond(w)
		return
	}

	ph, err := bcrypt.GenerateFromPassword([]byte(ui.Password), bcrypt.DefaultCost)
	if err != nil {
		s.log.WithError(err).Error("generating bcrypt hash")
		apierr.Internal().Respond(w)
		return
	}

	usr, err := db.InsertUser(
		r.Context(),
		s.db,
		ui.Name,
		ph,
		true,
	)
	switch err {
	case nil:
		// OK.
	case r.Context().Err():
		apierr.Context().Respond(w)
		return
	default:
		s.log.WithError(err).Error("creating a new user")
		apierr.Internal().Respond(w)
		return
	}

	s.respondJSON(w, usr)
}

// GetUsers retrieves all users.
func (s *Server) GetUsers(w http.ResponseWriter, r *http.Request) {
	uu, err := db.GetUsers(r.Context(), s.db)
	switch err {
	case nil:
		// OK.
	case r.Context().Err():
		apierr.Context().Respond(w)
		return
	default:
		s.log.WithError(err).Error("fetching users")
		apierr.Internal().Respond(w)
		return
	}

	s.respondJSON(w, uu)
}

// GetUser retrieves user by the user id.
func (s *Server) GetUser(w http.ResponseWriter, r *http.Request) {
	uid, aerr := s.extractContextUserID(r)
	if aerr != nil {
		aerr.Respond(w)
		return
	}

	usr, err := db.GetUserByID(r.Context(), s.db, uid)
	switch err {
	case nil:
		// OK.
	case r.Context().Err():
		apierr.Context().Respond(w)
		return
	case db.ErrNotFound:
		apierr.NotFound("user").Respond(w)
		return
	default:
		s.log.WithError(err).Error("fetching user by id")
		apierr.Internal().Respond(w)
		return
	}

	s.respondJSON(w, usr)
}

// UpdateUserPassword updates user password.
func (s *Server) UpdateUserPassword(w http.ResponseWriter, r *http.Request) {
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

	var inp struct {
		Password    string `json:"password"`
		OldPassword string `json:"old_password"`
	}

	if err := json.Unmarshal(data, &inp); err != nil {
		apierr.MalformedDataInput(apierr.DataTypeJSON).Respond(w)
		return
	}

	user, err := db.GetUserByID(r.Context(), s.db, uid)
	switch err {
	case nil:
		// OK.
	case r.Context().Err():
		apierr.Context().Respond(w)
		return
	default:
		s.log.WithError(err).Error("fetching user by id for password update")
		apierr.Database().Respond(w)
		return
	}

	err = bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(inp.OldPassword))
	switch err {
	case nil:
		// OK.
	case bcrypt.ErrMismatchedHashAndPassword:
		apierr.BadRequest("invalid password")
		return
	default:
		s.log.WithError(err).Error("comparing bcrypt hash to a password")
		apierr.Internal().Respond(w)
		return
	}

	if aerr := core.ValidatePassword(inp.Password); aerr != nil {
		aerr.Respond(w)
		return
	}

	ph, err := bcrypt.GenerateFromPassword([]byte(inp.Password), bcrypt.DefaultCost)
	if err != nil {
		s.log.WithError(err).Error("generating bcrypt password hash")
		apierr.Internal().Respond(w)
		return
	}

	err = db.UpdateUserPasswordByID(r.Context(), s.db, uid, ph)
	switch err {
	case nil:
		// OK.
	case r.Context().Err():
		apierr.Context().Respond(w)
		return
	default:
		s.log.WithError(err).Error("updating user password")
		apierr.Database().Respond(w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeleteUser deletes user. If the super user initiated the request the user
// id is extracted from the path, otherwise the user id is extracted from
// the context.
func (s *Server) DeleteUser(super bool) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			uid  xid.ID
			aerr *apierr.Error
		)

		if !super {
			uid, aerr = s.extractContextUserID(r)
			if aerr != nil {
				aerr.Respond(w)
				return
			}
		} else {
			uid, aerr = s.extractPathID(r, "userId")
			if aerr != nil {
				aerr.Respond(w)
				return
			}
		}

		adm, aerr := s.extractContextAdmin(r)
		if aerr != nil {
			aerr.Respond(w)
			return
		}

		if adm {
			usr, err := db.GetUserByID(r.Context(), s.db, uid)
			switch err {
			case nil:
				// OK.
			case r.Context().Err():
				apierr.Context().Respond(w)
				return
			case db.ErrNotFound:
				apierr.NotFound("user").Respond(w)
				return
			default:
				s.log.WithError(err).Error("fetching user by id")
				apierr.Database().Respond(w)
				return
			}

			if usr.Name == core.RootAdminName {
				apierr.BadRequest("cannot delete root admin").Respond(w)
				return
			}
		}

		err := db.DeleteUserByID(r.Context(), s.db, uid)
		switch err {
		case nil:
			// OK.
		case r.Context().Err():
			apierr.Context().Respond(w)
			return
		default:
			s.log.WithError(err).Error("creating a new user")
			apierr.Database().Respond(w)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// Login authenticates the user by its credentials and returns a JWT token.
func (s *Server) Login(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		apierr.MalformedDataInput(apierr.DataTypeRequestBody).Respond(w)
		return
	}

	var ui core.UserInput
	if err := json.Unmarshal(data, &ui); err != nil {
		apierr.MalformedDataInput(apierr.DataTypeJSON).Respond(w)
		return
	}

	usr, err := db.GetUserByName(r.Context(), s.db, ui.Name)
	switch err {
	case nil:
		// OK.
	case r.Context().Err(), db.ErrNotFound:
		apierr.Unauthorized().Respond(w)
		return
	default:
		s.log.WithError(err).Error("fetching user by name")
		apierr.Database().Respond(w)
		return
	}

	err = bcrypt.CompareHashAndPassword(usr.PasswordHash, []byte(ui.Password))
	switch err {
	case nil:
		// OK.
	case bcrypt.ErrMismatchedHashAndPassword:
		apierr.BadRequest("invalid password")
		return
	default:
		s.log.WithError(err).Error("comparing bcrypt hash to a password")
		apierr.Internal().Respond(w)
		return
	}

	sjwt, err := core.IssueJWT(s.secret, usr.ID, usr.Admin, time.Now())
	if err != nil {
		s.log.WithError(err).Error("creating jwt")
		apierr.Internal().Respond(w)
		return
	}

	s.respondJSON(w, struct {
		User        *core.User `json:"user"`
		AccessToken string     `json:"access_token"`
	}{
		User:        usr,
		AccessToken: string(sjwt),
	})
}

// Register creates a new user by its credentials and retrieves a JWT token.
func (s *Server) Register(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		apierr.MalformedDataInput(apierr.DataTypeRequestBody).Respond(w)
		return
	}

	var ui core.UserInput
	if err := json.Unmarshal(data, &ui); err != nil {
		apierr.MalformedDataInput(apierr.DataTypeJSON).Respond(w)
		return
	}

	if aerr := ui.Validate(); aerr != nil {
		aerr.Respond(w)
		return
	}

	_, err = db.GetUserByName(r.Context(), s.db, ui.Name)
	switch err {
	case db.ErrNotFound:
		// OK.
	case nil:
		apierr.Conflict("user").Respond(w)
		return
	case r.Context().Err():
		apierr.Context().Respond(w)
		return
	default:
		s.log.WithError(err).Error("fetching user by name")
		apierr.Database().Respond(w)
		return
	}

	ph, err := bcrypt.GenerateFromPassword([]byte(ui.Password), bcrypt.DefaultCost)
	if err != nil {
		s.log.WithError(err).Error("generating bcrypt password hash")
		apierr.Internal().Respond(w)
		return
	}

	usr, err := db.InsertUser(
		r.Context(),
		s.db,
		ui.Name,
		ph,
		false,
	)
	switch err {
	case nil:
		// OK.
	case r.Context().Err():
		apierr.Context().Respond(w)
		return
	default:
		s.log.WithError(err).Error("creating a new user")
		apierr.Database().Respond(w)
		return
	}

	sjwt, err := core.IssueJWT(s.secret, usr.ID, usr.Admin, time.Now())
	if err != nil {
		s.log.WithError(err).Error("creating jwt")
		apierr.Internal().Respond(w)
		return
	}

	s.respondJSON(w, struct {
		User        *core.User `json:"user"`
		AccessToken string     `json:"access_token"`
	}{
		User:        usr,
		AccessToken: string(sjwt),
	})
}
