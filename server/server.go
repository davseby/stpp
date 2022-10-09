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

// contextKey is an unique type used to store context data.
type contextKey int

const (
	// contextKeyUserID is used to store user id data.
	contextKeyUserID contextKey = iota

	// contextKeyAdmin is used to store admin flag data.
	contextKeyAdmin
)

// Server contains structures and data that is required to run a web
// server.
type Server struct {
	// log is a configured logger structure.
	log logrus.FieldLogger

	// db is an API that is used to communicate with the database.
	db *sql.DB

	// serv is a underlying structure used to establish and maintain
	// clients connections pool.
	serv *http.Server

	// secret specifies the secret that is used to generate jwt hash.
	secret []byte
}

// NewServer creates a fresh instance of the server.
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

// Start starts the server. It blocks until the server.Stop is called.
func (s *Server) Start() error {
	return s.serv.ListenAndServe()
}

// Stop shuts down the server.
func (s *Server) Stop() error {
	return s.serv.Shutdown(context.Background())
}

// router builds the server router.
func (s *Server) router() chi.Router {
	r := chi.NewRouter()

	r.Post("/login", s.Login)
	r.Post("/register", s.Register)

	r.Route("/products", func(sr chi.Router) {
		sr.Get("/", s.GetProducts)
		sr.Get("/{productID}", s.GetProduct)

		sr.Group(func(ssr chi.Router) {
			ssr.Use(s.authorize(true))
			ssr.Post("/", s.CreateProduct)
			ssr.Patch("/{productID}", s.UpdateProduct)
			ssr.Delete("/{productID}", s.DeleteProduct)
		})
	})

	r.Route("/recipes", func(sr chi.Router) {
		sr.Get("/", s.GetRecipes)
		sr.Get("/{recipeID}", s.GetRecipe)
		sr.Get("/user/{userID}", s.GetUserRecipes)

		sr.Group(func(ssr chi.Router) {
			ssr.Use(s.authorize(false))
			ssr.Post("/", s.CreateRecipe)
			ssr.Patch("/{recipeID}", s.UpdateRecipe)
			ssr.Delete("/{recipeID}", s.DeleteRecipe)
		})
	})

	r.Route("/plans", func(sr chi.Router) {
		sr.Get("/", s.GetPlans)
		sr.Get("/{planID}", s.GetPlan)
		sr.Get("/user/{userID}", s.GetUserPlans)

		sr.Group(func(ssr chi.Router) {
			ssr.Use(s.authorize(false))
			ssr.Post("/", s.CreatePlan)
			ssr.Patch("/{planID}", s.UpdatePlan)
			ssr.Delete("/{planID}", s.DeletePlan)
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
			ssr.Get("/{userID}", s.GetUser)
			ssr.Delete("/{userID}", s.DeleteUser(true))
		})
	})

	r.Get("/version", s.GetVersion)

	return r
}

// GetVersion returns server version information.
func (s *Server) GetVersion(w http.ResponseWriter, r *http.Request) {
	s.respondJSON(w, struct {
		Version string `json:"version"`
	}{
		Version: "1.0.0",
	})
}

// authorize is a middleware that authorizes incoming requests by their
// authorization token.
func (s *Server) authorize(super bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sjwt, aerr := s.extractAuthorizationToken(r)
			if aerr != nil {
				aerr.Respond(w)
				return
			}

			id, admin, aerr := core.ParseJWT(s.secret, sjwt, time.Now())
			if aerr != nil {
				if aerr == apierr.Internal() {
					s.log.WithField("token", sjwt).Error("parsing jwt token")
				}

				aerr.Respond(w)
				return
			}

			_, err := db.GetUserByID(r.Context(), s.db, id)
			switch err {
			case nil:
				// OK.
			case db.ErrNotFound:
				apierr.Unauthorized().Respond(w)
				return
			case r.Context().Err():
				apierr.Context().Respond(w)
				return
			default:
				s.log.WithError(err).Error("fetching user by name")
				apierr.Database().Respond(w)
				return
			}

			if super && !admin {
				apierr.Forbidden().Respond(w)
				return
			}

			next.ServeHTTP(
				w,
				r.WithContext(
					context.WithValue(
						context.WithValue(
							r.Context(),
							contextKeyUserID,
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

// respondJSON marshals the given object and writes its data to the response
// writer.
func (s *Server) respondJSON(w http.ResponseWriter, obj any) {
	w.Header().Add("Content-Type", "application/json")

	data, err := json.Marshal(obj)
	if err != nil {
		s.log.WithError(err).Error("marshaling response object")
		apierr.Internal().Respond(w)
		return
	}

	_, err = w.Write(data)
	if err != nil {
		s.log.WithError(err).Error("writing to client response data")
		apierr.Internal().Respond(w)
		return
	}
}

// extractPathID extracts given key from the request path.
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

// extractAuthorizationToken extract authorization token from the
// authorization header.
func (s *Server) extractAuthorizationToken(r *http.Request) ([]byte, *apierr.Error) {
	value := r.Header.Get("Authorization")

	parts := strings.SplitN(value, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return nil, apierr.Unauthorized()
	}

	return []byte(parts[1]), nil
}

// extractContextUserID extracts user id from the request context.
func (s *Server) extractContextUserID(r *http.Request) (xid.ID, *apierr.Error) {
	vid := r.Context().Value(contextKeyUserID)

	id, ok := vid.(xid.ID)
	if vid == nil || !ok {
		s.log.Error("missing context user id value")
		return xid.NilID(), apierr.Internal()
	}

	return id, nil
}

// extractContextAdmin extracts admin flag from the request context.
func (s *Server) extractContextAdmin(r *http.Request) (bool, *apierr.Error) {
	vid := r.Context().Value(contextKeyAdmin)

	admin, ok := vid.(bool)
	if vid == nil || !ok {
		s.log.Error("missing context admin value")
		return false, apierr.Internal()
	}

	return admin, nil
}
