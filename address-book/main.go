package main

import (
	"fmt"
	"github.com/danielmichaels/address-book/handlers"
	muxHandlers "github.com/gorilla/handlers"
	"log"
	"net/http"
	"os"
)

func main() {
	log.Println("Starting Address Book API Server")

	s := handlers.NewServer()
	// use mux for logging
	http.Handle("/", muxHandlers.LoggingHandler(os.Stdout, s.Router()))
	// we would use environment variables here and a config struct
	port := "9090"
	log.Printf("Listening on port: %s", port)
	// trap any start up errors and fail
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
	// next I would handle graceful shutdown using context.WithTimeout, goroutines and channels
}
