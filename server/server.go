package server

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"foodie/core"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/rs/xid"
	"github.com/sirupsen/logrus"
)

type contextKey int

const (
	contextKeyID contextKey = iota
)

type Server struct {
	log    logrus.FieldLogger
	db     *sql.DB
	serv   *http.Server
	secret []byte
}

func NewServer(db *sql.DB, port string, secret []byte) *Server {
	s := &Server{
		log:    logrus.New(),
		db:     db,
		secret: secret,
	}

	s.serv = &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: s.router(),
	}

	return s
}

func (s *Server) Start() error {
	return s.serv.ListenAndServe()
}

func (s *Server) Stop() error {
	return s.serv.Shutdown(context.Background())
}

func (s *Server) router() chi.Router {
	r := chi.NewRouter()

	r.Post("/register", s.Register)
	r.Post("/login", s.Login)

	r.Route("/products", func(sr chi.Router) {
		sr.Get("/", s.GetProducts)
		sr.Get("/{id}", s.GetProduct)

		sr.Group(func(ssr chi.Router) {
			ssr.Use(s.authorize(true))
			ssr.Post("/", s.CreateProduct)
			ssr.Patch("/{id}", s.UpdateProduct)
			ssr.Delete("/{id}", s.DeleteProduct)
		})
	})

	r.Route("/users", func(sr chi.Router) {
		sr.Group(func(ssr chi.Router) {
			ssr.Use(s.authorize(false))
			ssr.Delete("/", s.DeleteUser(false))
			ssr.Patch("/", s.UpdateUserPassword)
		})

		sr.Group(func(ssr chi.Router) {
			ssr.Use(s.authorize(true))
			ssr.Get("/", s.GetUsers)
			ssr.Post("/", s.CreateAdminUser)
			ssr.Get("/{id}", s.GetUser)
			ssr.Delete("/{id}", s.DeleteUser(true))
		})
	})

	r.Get("/version", s.Version)

	return r
}

func (s *Server) Version(w http.ResponseWriter, r *http.Request) {
	s.respondJSON(w, struct {
		Version string `json:"version"`
	}{
		Version: "1.0.0",
	})
}

func (s *Server) authorize(super bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			signedJWT, ok := extractAuthorizationToken(r)
			if !ok {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			id, admin, err := core.ParseJWT(s.secret, signedJWT, time.Now())
			switch err {
			case nil:
				// OK.
			case core.ErrExpiredToken:
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(err.Error()))
				return
			default:
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			if super && !admin {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			next.ServeHTTP(
				w,
				r.WithContext(
					context.WithValue(
						r.Context(),
						contextKeyID,
						id,
					),
				),
			)
		})
	}
}

func (s *Server) respondJSON(w http.ResponseWriter, obj any) {
	data, err := json.Marshal(obj)
	if err != nil {
		s.log.WithError(err).Error("marshalling response object")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = w.Write(data)
	if err != nil {
		s.log.WithError(err).Error("writing to client response data")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func extractContextData(r *http.Request) (xid.ID, error) {
	vid := r.Context().Value(contextKeyID)

	id, ok := vid.(xid.ID)
	if vid == nil || !ok {
		return xid.NilID(), errors.New("cannot extract id")
	}

	return id, nil
}

func extractIDFromPath(r *http.Request) (xid.ID, error) {
	sid := chi.URLParam(r, "id")
	if sid == "" {
		return xid.NilID(), nil
	}

	return xid.FromString(sid)
}

func extractAuthorizationToken(r *http.Request) ([]byte, bool) {
	value := r.Header.Get("Authorization")

	parts := strings.SplitN(value, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return nil, false
	}

	return []byte(parts[1]), true
}
