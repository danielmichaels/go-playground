package handlers

import (
	"fmt"
	"github.com/danielmichaels/address-book/data"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"time"
)

// NewServer create a new server which returns a server struct that implements a Router interface
func NewServer() Server {
	s := apiServer{
		router:       mux.NewRouter(),
		ReadTimeout:  5 * time.Second,   // max time to read from client
		WriteTimeout: 10 * time.Second,  // max time to write to client
		IdleTimeout:  120 * time.Second, // max time for TCP Keep-Alive

	}
	s.routes()
	return s
}

// apiServer struct which contains dependencies
type apiServer struct {
	router       *mux.Router
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// Server interface implements a Router type. Handles all HTTP requests for server.
type Server interface {
	Router() *mux.Router
}

func (s apiServer) Router() *mux.Router {
	return s.router
}

func (s *apiServer) respondWithJSON(w http.ResponseWriter, i interface{}, status int) error {
	w.WriteHeader(status)
	e := data.ToJSON(i, w)
	return e
}

// validateAddress validate the address has the required json fields
func (s apiServer) validateAddress(address data.AddressBook) error {
	err := address.Validate()
	if err != nil {
		return fmt.Errorf("unable to validate address: %s", err)
	}
	return nil
}

func (s *apiServer) readAddressAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		addressBook := data.GetAddressBook()
		err := s.respondWithJSON(w, &addressBook, http.StatusOK)
		if err != nil {
			http.Error(w, "Unable to marshal json", http.StatusInternalServerError)
			return
		}
	}
}
func (s *apiServer) createAddress() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var address data.AddressBook
		err := data.FromJSON(&address, r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = s.validateAddress(address)
		if err != nil {
			http.Error(
				w,
				fmt.Sprintf("%s", err),
				http.StatusBadRequest,
			)
			return
		}
		data.AddAddressToStore(&address)
		err = s.respondWithJSON(w, &address, http.StatusCreated)
	}
}

func (s *apiServer) updateAddress() http.HandlerFunc {
	var address data.AddressBook
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = data.FromJSON(&address, r.Body)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = s.validateAddress(address)
		if err != nil {
			log.Println(err)
			http.Error(
				w,
				fmt.Sprintf("%s", err),
				http.StatusBadRequest,
			)
			return
		}

		err = data.UpdateAddress(id, &address)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = s.respondWithJSON(w, &address, http.StatusOK)
	}
}

func (s *apiServer) readAddress() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		result, _, _ := data.FindAddress(id)
		err = s.respondWithJSON(w, &result, http.StatusOK)
	}
}
func (s *apiServer) deleteAddress() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		deleteAddress := data.DeleteAddress(id)
		err = s.respondWithJSON(w, &deleteAddress, http.StatusNoContent)
	}
}
func (s *apiServer) exportAddressToCSV() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, "export to csv")
	}
}
