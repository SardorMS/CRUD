package app

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/SardorMS/CRUD/cmd/app/middleware"
	"github.com/SardorMS/CRUD/pkg/managers"
	"github.com/SardorMS/CRUD/pkg/types"
	"github.com/gorilla/mux"
)

// handleManagerRegistration - ...
func (s *Server) handleManagerRegistration(writer http.ResponseWriter, request *http.Request) {

	var item *types.ManagerRegister

	id, err := middleware.Authentication(request.Context())
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if id == 0 {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	if err = json.NewDecoder(request.Body).Decode(&item); err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// if admin := s.managersSvc.IsAdmin(request.Context(), id); admin != true {
	// 	log.Println(err)
	// 	http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	// 	return
	// }

	items := &types.Managers{
		ID:    item.ID,
		Name:  item.Name,
		Phone: item.Phone,
	}

	for _, role := range item.Roles {
		if role == "ADMIN" {
			items.IsAdmin = true
			break
		}
	}
	token, err := s.managersSvc.Register(request.Context(), items)
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	respondJSON(writer, map[string]interface{}{"token": token})
}

// handleManagerGetToken - ...
func (s *Server) handleManagerGetToken(writer http.ResponseWriter, request *http.Request) {
	var item *types.Managers

	if err := json.NewDecoder(request.Body).Decode(&item); err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	token, err := s.managersSvc.Token(request.Context(), item.Phone, item.Password)
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	respondJSON(writer, map[string]interface{}{"token": token})
}

// handleManagerGetSales - ...
func (s *Server) handleManagerGetSales(writer http.ResponseWriter, request *http.Request) {

	id, err := middleware.Authentication(request.Context())
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if id == 0 {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	total, err := s.managersSvc.GetSales(request.Context(), id)
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	respondJSON(writer, map[string]interface{}{"manager_id": id, "total": total})
}

// handleManagerMakeSale - ...
func (s *Server) handleManagerMakeSale(writer http.ResponseWriter, request *http.Request) {

	id, err := middleware.Authentication(request.Context())
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if id == 0 {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	sale := &types.Sale{}
	sale.ManagerID = id

	if err := json.NewDecoder(request.Body).Decode(&sale); err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	sale, err = s.managersSvc.MakeSales(request.Context(), sale)
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	respondJSON(writer, sale)
}

// handleManagerGetProducts - ...
func (s *Server) handleManagerGetProducts(writer http.ResponseWriter, request *http.Request) {

	items, err := s.managersSvc.Products(request.Context())
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	respondJSON(writer, items)
}

// handleManagerChangeProduct - ...
func (s *Server) handleManagerChangeProduct(writer http.ResponseWriter, request *http.Request) {
	id, err := middleware.Authentication(request.Context())
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if id == 0 {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	product := &types.Products{}

	if err := json.NewDecoder(request.Body).Decode(&product); err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	product, err = s.managersSvc.ChangeProduct(request.Context(), product)
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	respondJSON(writer, product)
}


// handleManagerRemoveProductByID - ...
func (s *Server) handleRehandleManagerRemoveProductByIDmoveCustomerByID(writer http.ResponseWriter, request *http.Request) {
	id, err := middleware.Authentication(request.Context())
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if id == 0 {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	idParam, ok := mux.Vars(request)["id"]
	if !ok {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	
	productID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// НАЧНи ЗДЕСЬ
	item, err := s.managersSvc.RemoveProductByID(request.Context(), productID)
	if errors.Is(err, managers.ErrNotFound) {
		http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusInternalServerError)
		return
	}

	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	respondJSON(writer, item)
}