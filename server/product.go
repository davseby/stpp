package server

import (
	"encoding/json"
	"foodie/core"
	"foodie/db"
	"foodie/server/apierr"
	"io"
	"net/http"
)

// CreateProduct creates a product.
func (s *Server) CreateProduct(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		apierr.MalformedDataInput(apierr.DataTypeRequestBody).Respond(w)
		return
	}

	var pc core.ProductCore
	if err := json.Unmarshal(data, &pc); err != nil {
		apierr.MalformedDataInput(apierr.DataTypeJSON).Respond(w)
		return
	}

	if aerr := pc.Validate(); aerr != nil {
		aerr.Respond(w)
		return
	}

	prd, err := db.InsertProduct(r.Context(), s.db, pc)
	switch err {
	case nil:
		// OK.
	case r.Context().Err():
		apierr.Context().Respond(w)
		return
	default:
		s.log.WithError(err).Error("creating a new product")
		apierr.Database().Respond(w)
		return
	}

	s.respondJSON(w, prd)
}

// GetProducts retrieves all products.
func (s *Server) GetProducts(w http.ResponseWriter, r *http.Request) {
	pp, err := db.GetProducts(r.Context(), s.db)
	switch err {
	case nil:
		// OK.
	case r.Context().Err():
		apierr.Context().Respond(w)
		return
	default:
		s.log.WithError(err).Error("fetching products")
		apierr.Database().Respond(w)
		return
	}

	s.respondJSON(w, pp)
}

// GetProduct retrieves a single product by its id.
func (s *Server) GetProduct(w http.ResponseWriter, r *http.Request) {
	pid, aerr := s.extractPathID(r, "productID")
	if aerr != nil {
		aerr.Respond(w)
		return
	}

	prd, err := db.GetProductByID(r.Context(), s.db, pid)
	switch err {
	case nil:
		// OK.
	case r.Context().Err():
		apierr.Context().Respond(w)
		return
	case db.ErrNotFound:
		apierr.NotFound("product").Respond(w)
		return
	default:
		s.log.WithError(err).Error("fetching product by id")
		apierr.Database().Respond(w)
		return
	}

	s.respondJSON(w, prd)
}

// UpdateProduct updates existing product by its id.
func (s *Server) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	pid, aerr := s.extractPathID(r, "productID")
	if aerr != nil {
		aerr.Respond(w)
		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		apierr.MalformedDataInput(apierr.DataTypeRequestBody).Respond(w)
		return
	}

	var pc core.ProductCore
	if err := json.Unmarshal(data, &pc); err != nil {
		apierr.MalformedDataInput(apierr.DataTypeJSON).Respond(w)
		return
	}

	prd, err := db.UpdateProductByID(r.Context(), s.db, pid, pc)
	switch err {
	case nil:
		// OK.
	case r.Context().Err():
		apierr.Context().Respond(w)
		return
	default:
		s.log.WithError(err).Error("updating product")
		apierr.Database().Respond(w)
		return
	}

	s.respondJSON(w, prd)
}

// DeleteProduct deletes existing product by its id.
func (s *Server) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	pid, aerr := s.extractPathID(r, "productID")
	if aerr != nil {
		aerr.Respond(w)
		return
	}

	recipes, err := db.GetRecipyProductsByProductID(r.Context(), s.db, pid)
	switch err {
	case nil:
		// OK.
	case r.Context().Err():
		apierr.Context().Respond(w)
		return
	default:
		s.log.WithError(err).Error("getting recipy products by product id")
		apierr.Database().Respond(w)
		return
	}

	if len(recipes) > 0 {
		apierr.Conflict("product in use").Respond(w)
		return
	}

	err = db.DeleteProductByID(r.Context(), s.db, pid)
	switch err {
	case nil:
		// OK.
	case r.Context().Err():
		apierr.Context().Respond(w)
		return
	default:
		s.log.WithError(err).Error("deleting product by id")
		apierr.Database().Respond(w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
