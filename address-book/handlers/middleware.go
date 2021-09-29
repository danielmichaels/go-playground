package handlers

import (
	"context"
	"fmt"
	"github.com/danielmichaels/address-book/data"
	"net/http"
)

// Middleware that sets the `application/json` response type
func (s apiServer) jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func (s apiServer) addressBookValidationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		addr := data.AddressBook{}
		err := data.FromJSON(&addr, r.Body)
		if err != nil {
			return
		}

		err = addr.Validate()
		if err != nil {
			http.Error(
				w,
				fmt.Sprintf("Unable to validate Address: %s", err),
				http.StatusBadRequest,
			)
			return
		}
		ctx := context.WithValue(r.Context(), struct{}{}, addr)
		req := r.WithContext(ctx)
		next.ServeHTTP(w, req)
	})
}
