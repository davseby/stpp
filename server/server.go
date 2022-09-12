package server

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"foodie/core"
	"foodie/db"
	"foodie/server/apierr"
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
	contextKeyAdmin
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
		sr.Get("/{productId}", s.GetProduct)

		sr.Group(func(ssr chi.Router) {
			ssr.Use(s.authorize(true))
			ssr.Post("/", s.CreateProduct)
			ssr.Patch("/{productId}", s.UpdateProduct)
			ssr.Delete("/{productId}", s.DeleteProduct)
		})
	})

	r.Route("/recipes", func(sr chi.Router) {
		sr.Get("/", s.GetPublicRecipes)
		sr.Get("/{userId}", s.GetUserRecipes)

		sr.Route("/{recipyId}", func(ssr chi.Router) {
			ssr.Get("/", s.GetRecipy)
			ssr.Route("/ratings", func(sssr chi.Router) {
				sssr.Get("/", s.GetRecipyRatings)

				sssr.Group(func(ssssr chi.Router) {
					ssssr.Use(s.authorize(false))
					ssssr.Post("/", s.CreateRating)
					ssssr.Patch("/", s.UpdateRating)
					ssssr.Delete("/", s.DeleteRating)
				})
			})
		})

		sr.Group(func(ssr chi.Router) {
			ssr.Use(s.authorize(false))
			ssr.Post("/", s.CreateRecipy)
			ssr.Patch("/{recipyId}", s.UpdateRecipy)
			ssr.Delete("/{recipyId}", s.DeleteRecipy)
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
			ssr.Get("/{userId}", s.GetUser)
			ssr.Delete("/{userId}", s.DeleteUser(true))
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
			signedJWT, aerr := s.extractAuthorizationToken(r)
			if aerr != nil {
				aerr.Respond(w)
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
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			_, err = db.GetUserByID(r.Context(), s.db, id)
			switch err {
			case nil:
				// OK.
			case db.ErrNotFound:
				w.WriteHeader(http.StatusUnauthorized)
				return
			case r.Context().Err():
				w.WriteHeader(http.StatusBadRequest)
				return
			default:
				s.log.WithError(err).Error("fetching user by name")
				w.WriteHeader(http.StatusInternalServerError)
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
						context.WithValue(
							r.Context(),
							contextKeyID,
							id,
						),
						contextKeyAdmin,
						admin,
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

func (s *Server) extractPathID(r *http.Request, key string) (xid.ID, *apierr.Error) {
	sid := chi.URLParam(r, key)
	if sid == "" {
		return xid.NilID(), apierr.BadRequest("invalid object path identification")
	}

	id, err := xid.FromString(sid)
	if err != nil {
		return xid.NilID(), apierr.BadRequest("incorrect object identification format")
	}

	return id, nil
}

func (s *Server) extractAuthorizationToken(r *http.Request) ([]byte, *apierr.Error) {
	value := r.Header.Get("Authorization")

	parts := strings.SplitN(value, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return nil, apierr.Unauthorized()
	}

	return []byte(parts[1]), nil
}

func (s *Server) extractContextUserID(r *http.Request) (xid.ID, *apierr.Error) {
	vid := r.Context().Value(contextKeyID)

	id, ok := vid.(xid.ID)
	if vid == nil || !ok {
		s.log.Error("missing context user id value")
		return xid.NilID(), apierr.Internal()
	}

	return id, nil
}

func (s *Server) extractContextAdmin(r *http.Request) (bool, *apierr.Error) {
	vid := r.Context().Value(contextKeyAdmin)

	admin, ok := vid.(bool)
	if vid == nil || !ok {
		s.log.Error("missing context admin value")
		return false, apierr.Internal()
	}

	return admin, nil
}
