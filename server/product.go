package server

import (
	"encoding/json"
	"foodie/core"
	"foodie/db"
	"io"
	"net/http"
)

// GetProducts retrieves all products.
func (s *Server) GetProducts(w http.ResponseWriter, r *http.Request) {
	products, err := db.GetProducts(r.Context(), s.db)
	switch err {
	case nil:
		// OK.
	case r.Context().Err():
		w.WriteHeader(http.StatusBadRequest)
		return
	default:
		s.log.WithError(err).Error("fetching products")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	s.respondJSON(w, products)
}

func (s *Server) CreateProduct(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var pc core.ProductCore
	if err := json.Unmarshal(data, &pc); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid JSON object"))
		return
	}

	if err := pc.Validate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	product, err := db.InsertProduct(r.Context(), s.db, pc)
	switch err {
	case nil:
		// OK.
	case r.Context().Err():
		w.WriteHeader(http.StatusBadRequest)
		return
	default:
		s.log.WithError(err).Error("creating a new product")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	s.respondJSON(w, product)
}

func (s *Server) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	pid, aerr := s.extractPathID(r, "productId")
	if aerr != nil {
		aerr.Respond(w)
		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var pc core.ProductCore
	if err := json.Unmarshal(data, &pc); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid JSON object"))
		return
	}

	err = db.UpdateProductByID(r.Context(), s.db, pid, pc)
	switch err {
	case nil:
		// OK.
	case r.Context().Err():
		w.WriteHeader(http.StatusBadRequest)
		return
	default:
		s.log.WithError(err).Error("updating product")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	product, err := db.GetProductByID(r.Context(), s.db, pid)
	switch err {
	case nil:
		// OK.
	case r.Context().Err():
		w.WriteHeader(http.StatusBadRequest)
		return
	case db.ErrNotFound:
		s.log.WithError(err).Error("fetching product by id after its update")
		w.WriteHeader(http.StatusInternalServerError)
		return
	default:
		s.log.WithError(err).Error("fetching product by id")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	s.respondJSON(w, product)
}
func (s *Server) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	pid, aerr := s.extractPathID(r, "productId")
	if aerr != nil {
		aerr.Respond(w)
		return
	}

	err := db.DeleteProductByID(r.Context(), s.db, pid)
	switch err {
	case nil:
		// OK.
	case r.Context().Err():
		w.WriteHeader(http.StatusBadRequest)
		return
	default:
		s.log.WithError(err).Error("deleting product by id")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) GetProduct(w http.ResponseWriter, r *http.Request) {
	pid, aerr := s.extractPathID(r, "productId")
	if aerr != nil {
		aerr.Respond(w)
		return
	}

	product, err := db.GetProductByID(r.Context(), s.db, pid)
	switch err {
	case nil:
		// OK.
	case r.Context().Err():
		w.WriteHeader(http.StatusBadRequest)
		return
	case db.ErrNotFound:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("product"))
		return
	default:
		s.log.WithError(err).Error("fetching product by id")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	s.respondJSON(w, product)
}
