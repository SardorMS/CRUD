package app

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/SardorMS/CRUD/cmd/app/middleware"
	"github.com/SardorMS/CRUD/pkg/customers"
	"github.com/SardorMS/CRUD/pkg/managers"
	"github.com/SardorMS/CRUD/pkg/types"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

const (
	GET    = "GET"
	POST   = "POST"
	DELETE = "DELETE"
)

// Server - represents the logical server of application.
type Server struct {
	mux          *mux.Router
	customersSvc *customers.Service
	managersSvc  *managers.Service
}

// NewServer - constructor function to create a new server.
func NewServer(mux *mux.Router, customersSvc *customers.Service, managersSvc *managers.Service) *Server {
	return &Server{
		mux:          mux,
		customersSvc: customersSvc,
		managersSvc:  managersSvc,
	}
}

// ServeHTTP - method to start the server.
func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	s.mux.ServeHTTP(writer, request)
}

// Init - initializes the server (register all handlers).
func (s *Server) Init() {

	// s.mux.Handle("/customers", middleware.Logger(http.HandlerFunc(s.handleGetAllCustomers)))
	// s.mux.Use(middleware.Logger)
	// chMD := middleware.CheckHeader("Content-Type", "application/json")
	// s.mux.Handle("customers", chMD(http.HandlerFunc(s.handleSaveCustomer))).Methods(POST)
	// s.mux.Use(middleware.CheckHeader("Content-Type", "application/json"))
	// s.mux.Use(middleware.Basic(s.securitySvc.Auth))

	// s.mux.HandleFunc("/api/customers", s.handleSaveCustomer).Methods(POST)
	// s.mux.HandleFunc("/api/customers/token", s.handleGenerateToken).Methods(POST)
	// s.mux.HandleFunc("/api/customers/token/validate", s.handleValidateToken).Methods(POST)
	// managersSubrouter.HandleFunc("/customers/active", s.handleGetAllActive).Methods(GET)
	// managersSubrouter.HandleFunc("/customers/{id:[0-9]+}", s.handleGetCustomerByID).Methods(GET)
	// managersSubrouter.HandleFunc("/customers/{id:[0-9]+}/block", s.handlePostBlock).Methods(POST)
	// managersSubrouter.HandleFunc("/customers/{id:[0-9]+}/block", s.handleDeleteBlock).Methods(DELETE)

	s.mux.Use(middleware.Logger)

	customerAuthenticateMd := middleware.Authenticate(s.customersSvc.IDByToken)
	customersSubrouter := s.mux.PathPrefix("/api/customers").Subrouter()
	customersSubrouter.Use(customerAuthenticateMd)

	customersSubrouter.HandleFunc("", s.handleCustomerRegistration).Methods(POST)
	customersSubrouter.HandleFunc("/token", s.handleCustomerGetToken).Methods(POST)
	customersSubrouter.HandleFunc("/products", s.handleCustomerGetProducts).Methods(GET)
	customersSubrouter.HandleFunc("/purchases", s.handleCustomerGetPurchases).Methods(GET)
	// customersSubrouter.HandleFunc("/purchases", s.handleCustomerMakePurchase).Methods(POST)


	managerAuthenticateMd := middleware.Authenticate(s.managersSvc.IDByToken)
	managersSubrouter := s.mux.PathPrefix("/api/managers").Subrouter()
	managersSubrouter.Use(managerAuthenticateMd)

	managersSubrouter.HandleFunc("", s.handleManagerRegistration).Methods(POST)
	managersSubrouter.HandleFunc("/token", s.handleManagerGetToken).Methods(POST)
	managersSubrouter.HandleFunc("/sales", s.handleManagerGetSales).Methods(GET)
	managersSubrouter.HandleFunc("/sales", s.handleManagerMakeSale).Methods(POST)
	managersSubrouter.HandleFunc("/products", s.handleManagerGetProducts).Methods(GET)
	managersSubrouter.HandleFunc("/products", s.handleManagerChangeProduct).Methods(POST)
	managersSubrouter.HandleFunc("/products/{id:[0-9]+}", s.handleManagerRemoveProductByID).Methods(DELETE)
	managersSubrouter.HandleFunc("/customers", s.handleManagerGetCustomers).Methods(GET)
	managersSubrouter.HandleFunc("/customers", s.handleManagerChangeCustomer).Methods(POST)
	managersSubrouter.HandleFunc("/customers/{id:[0-9]+}", s.handleManagerRemoveCustomerByID).Methods(DELETE)

}

// handleSaveCustomer - handler for creating or saving.
func (s *Server) handleSaveCustomer(writer http.ResponseWriter, request *http.Request) {

	var item *types.Customer
	if err := json.NewDecoder(request.Body).Decode(&item); err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(item.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	item.Password = string(hash)

	customer, err := s.customersSvc.Save(request.Context(), item) //item
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	respondJSON(writer, customer)
}

// handleGenerateToken - handler to generate token.
func (s *Server) handleGenerateToken(writer http.ResponseWriter, request *http.Request) {

	var item *struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(request.Body).Decode(&item); err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	token, err := s.securitySvc.TokenForCustomer(request.Context(), item.Login, item.Password)
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	respondJSON(writer, map[string]interface{}{"status": "ok", "token": token})
}

// handleValidateToken - handler to check generated token.
func (s *Server) handleValidateToken(writer http.ResponseWriter, request *http.Request) {
	var item *struct {
		Token string `json:"token"`
	}

	if err := json.NewDecoder(request.Body).Decode(&item); err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	id, cerr := s.securitySvc.AuthenticateCustomer(request.Context(), item.Token)
	if cerr != nil {
		status := http.StatusInternalServerError
		text := http.StatusText(http.StatusInternalServerError)
		if cerr == security.ErrNoSuchUser {
			status = http.StatusNotFound
			text = "not found"
		}
		if cerr == security.ErrExpired {
			status = http.StatusBadRequest
			text = "expired"
		}
		writer.WriteHeader(status)
		respondJSON(writer, map[string]interface{}{"status": "fail", "reason": text})
		return

	}

	writer.WriteHeader(http.StatusOK)
	respondJSON(writer, map[string]interface{}{"status": "ok", "customerId": id})
}

// handleGetAllCustomers - handler to get information about all customers.
func (s *Server) handleGetAllCustomers(writer http.ResponseWriter, request *http.Request) {

	item, err := s.customersSvc.All(request.Context())
	if errors.Is(err, customers.ErrNotFound) {
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

// handleGetAllActive - handler to get information about all active customers.
func (s *Server) handleGetAllActive(writer http.ResponseWriter, request *http.Request) {

	item, err := s.customersSvc.AllActive(request.Context())
	if errors.Is(err, customers.ErrNotFound) {
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

// handleGetCustomerByID - handler to get information about customers by ID.
func (s *Server) handleGetCustomerByID(writer http.ResponseWriter, request *http.Request) {

	idParam, ok := mux.Vars(request)["id"]
	if !ok {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	item, err := s.customersSvc.ByID(request.Context(), id)
	if errors.Is(err, customers.ErrNotFound) {
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

// handlePostBlock - handler to block customers.
func (s *Server) handlePostBlock(writer http.ResponseWriter, request *http.Request) {

	idParam, ok := mux.Vars(request)["id"]
	if !ok {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	item, err := s.customersSvc.Block(request.Context(), id, false)
	if errors.Is(err, customers.ErrNotFound) {
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

// handleDeleteBlock - handler to unblock customers.
func (s *Server) handleDeleteBlock(writer http.ResponseWriter, request *http.Request) {

	idParam, ok := mux.Vars(request)["id"]
	if !ok {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	item, err := s.customersSvc.Block(request.Context(), id, true)
	if errors.Is(err, customers.ErrNotFound) {
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

// handleRemoveCustomerByID - handler to remove customers.
func (s *Server) handleRemoveCustomerByID(writer http.ResponseWriter, request *http.Request) {

	idParam, ok := mux.Vars(request)["id"]
	if !ok {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	item, err := s.customersSvc.Delete(request.Context(), id)
	if errors.Is(err, customers.ErrNotFound) {
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

// respondJSON - response from JSON.
func respondJSON(w http.ResponseWriter, item interface{}) {

	data, err := json.Marshal(item)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		log.Println(err)
	}
}
