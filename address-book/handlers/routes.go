package handlers

import (
	"net/http"
)

func (s *apiServer) routes() {
	api := s.router.PathPrefix("/api").Subrouter()
	api.Use(s.jsonContentTypeMiddleware)

	api.HandleFunc("/address", s.readAddressAll()).Methods(http.MethodGet)
	api.HandleFunc("/address", s.createAddress()).Methods(http.MethodPost)
	api.HandleFunc("/address/{id}", s.updateAddress()).Methods(http.MethodPut)
	api.HandleFunc("/address/{id}", s.readAddress()).Methods(http.MethodGet)
	api.HandleFunc("/address/{id}", s.deleteAddress()).Methods(http.MethodDelete)
	api.HandleFunc("/address/export", s.exportAddressToCSV()).Methods(http.MethodGet)

	// catchall -> anything that does not match a path is caught here.
	api.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Invalid path", http.StatusBadRequest)
	})

}
