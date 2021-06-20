package app

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/SardorMS/CRUD/cmd/app/middleware"
	"github.com/SardorMS/CRUD/pkg/types"
)

// handleCustomerRegistration - ...
func (s *Server) handleCustomerRegistration(writer http.ResponseWriter, request *http.Request) {
	var item *types.Registration

	err := json.NewDecoder(request.Body).Decode(&item)
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	saved, err := s.customersSvc.Register(request.Context(), item)
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	respondJSON(writer, saved)
}

//handleCustomerGetToken - ...
func (s *Server) handleCustomerGetToken(writer http.ResponseWriter, request *http.Request) {
	var item *types.Auth

	err := json.NewDecoder(request.Body).Decode(&item)
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	token, err := s.customersSvc.Token(request.Context(), item.Login, item.Password)
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	respondJSON(writer, &types.Token{Token: token})
}

// - handleCustomerGetProducts - ...
func (s *Server) handleCustomerGetProducts(writer http.ResponseWriter, request *http.Request) {
	items, err := s.customersSvc.Products(request.Context())
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	respondJSON(writer, items)
}

// handleCustomerGetPurchases - ...
func (s *Server) handleCustomerGetPurchases(writer http.ResponseWriter, request *http.Request) {
	id, err := middleware.Authentication(request.Context())
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	items, err := s.customersSvc.Purchases(request.Context(), id)
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	respondJSON(writer, items)
}

