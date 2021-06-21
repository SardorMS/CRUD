package app

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/SardorMS/CRUD/cmd/app/middleware"
	"github.com/SardorMS/CRUD/pkg/customers"
	"github.com/SardorMS/CRUD/pkg/managers"
	"github.com/gorilla/mux"
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

	managersSubrouter.HandleFunc("", s.handleManagerRegistration).Methods(POST) // right
	managersSubrouter.HandleFunc("/token", s.handleManagerGetToken).Methods(POST) //right
	managersSubrouter.HandleFunc("/sales", s.handleManagerGetSales).Methods(GET) //right
	managersSubrouter.HandleFunc("/sales", s.handleManagerMakeSale).Methods(POST) //right
	managersSubrouter.HandleFunc("/products", s.handleManagerGetProducts).Methods(GET)
	managersSubrouter.HandleFunc("/products", s.handleManagerChangeProduct).Methods(POST) //right
	managersSubrouter.HandleFunc("/products/{id:[0-9]+}", s.handleManagerRemoveProductByID).Methods(DELETE)
	managersSubrouter.HandleFunc("/customers", s.handleManagerGetCustomers).Methods(GET)
	managersSubrouter.HandleFunc("/customers", s.handleManagerChangeCustomer).Methods(POST)
	managersSubrouter.HandleFunc("/customers/{id:[0-9]+}", s.handleManagerRemoveCustomerByID).Methods(DELETE)

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
