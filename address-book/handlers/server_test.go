package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/danielmichaels/address-book/data"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

// checkResponseCode testing utility for asserting the status code
func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("expected %v, actual %v", expected, actual)
	}

}

var AddressBookStore = []data.AddressBook{
	{
		ID:          1,
		FirstName:   "Michael",
		LastName:    "Jordan",
		Email:       "twotrey@hornets.nba",
		PhoneNumber: 23232323,
	},
	{
		ID:          2,
		FirstName:   "Frank",
		LastName:    "Herbert",
		Email:       "frank@arrakis.com",
		PhoneNumber: 12341234,
	},
}

func TestApiServer_Router(t *testing.T) {
	// set up the server
	router := mux.NewRouter()
	s := apiServer{
		router:       router,
		ReadTimeout:  1,
		WriteTimeout: 5,
		IdleTimeout:  15,
	}
	s.routes()

	t.Run("test server hits invalid path and returns 400", func(t *testing.T) {
		t.Helper()
		request, err := http.NewRequest(http.MethodGet, "/api/", nil)
		if err != nil {
			t.Fatal(err)
		}
		response := httptest.NewRecorder()
		s.router.ServeHTTP(response, request)
		checkResponseCode(t, response.Code, http.StatusBadRequest)
	})
	t.Run("test get all entries", func(t *testing.T) {
		t.Helper()
		request, err := http.NewRequest(http.MethodGet, "/api/address", nil)
		if err != nil {
			t.Fatal(err)
		}
		r := httptest.NewRecorder()
		s.router.ServeHTTP(r, request)
		checkResponseCode(t, r.Code, http.StatusOK)

		expectedAddressBook := AddressBookStore
		var actualAddressBook []data.AddressBook

		err = json.NewDecoder(r.Body).Decode(&actualAddressBook)
		if err != nil {
			t.Fatalf("Unable to parse response. expected %q, actual %v", r.Body, err)
		}
		if !reflect.DeepEqual(actualAddressBook, expectedAddressBook) {
			t.Errorf("expected %v, actual %v", expectedAddressBook, actualAddressBook)
		}
	})
	t.Run("test get single entry", func(t *testing.T) {
		t.Helper()
		request, err := http.NewRequest(http.MethodGet, "/api/address/1", nil)
		if err != nil {
			t.Fatal(err)
		}
		r := httptest.NewRecorder()
		s.router.ServeHTTP(r, request)
		checkResponseCode(t, r.Code, http.StatusOK)
		expected := AddressBookStore[0]
		t.Log(expected)
		var actual data.AddressBook
		err = data.FromJSON(&actual, r.Body)
		if err != nil {
			t.Fatalf("Unable to parse response. expected %q, actual %v", r.Body, err)
		}
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
	})
	t.Run("test create new entry", func(t *testing.T) {
		t.Helper()

		// table tests are helpful here because it lets you loop over several
		// endpoints rather than defining a subtest for each response/request
		var addressPostTests = []struct {
			requestBodyAddressData string
			httpStatusExpected     int
			expectedAddressData    data.AddressBook
			validate               bool
		}{
			// successful creation yields 200
			{
				`{"first_name": "Jane", "last_name":"Austin","email":"jane@austin.com","phone_number":123}`,
				http.StatusCreated,
				data.AddressBook{
					FirstName:   "Jane",
					LastName:    "Austin",
					Email:       "jane@austin.com",
					PhoneNumber: 123,
				},
				false,
			},
			// if missing a value, a validation error will return 400
			{
				`{"last_name":"Austin","email":"jane@austin.com","phone_number":123}`,
				http.StatusBadRequest,
				data.AddressBook{
					LastName:    "Austin",
					Email:       "jane@austin.com",
					PhoneNumber: 123,
				},
				false,
			},
			//// if email is not valid return 400
			{
				`{"last_name":"Austin","email":"jane@austin","phone_number":123}`,
				http.StatusBadRequest,
				data.AddressBook{
					FirstName:   "Jane",
					LastName:    "Austin",
					Email:       "jane@austin",
					PhoneNumber: 123,
				},
				false,
			},
		}
		for _, tt := range addressPostTests {
			req, err := http.NewRequest(http.MethodPost, "/api/address", bytes.NewBuffer([]byte(tt.requestBodyAddressData)))
			if err != nil {
				t.Fatal(err)
			}
			r := httptest.NewRecorder()
			s.router.ServeHTTP(r, req)

			checkResponseCode(t, r.Code, tt.httpStatusExpected)

			var actual data.AddressBook
			err = data.FromJSON(&actual, r.Body)

			if err != nil {
				if !tt.validate {
					// don't test json validation here so set the error to nil
					err = nil
				} else {
					t.Log("fromJSON error != nil")
					t.Fatalf("Unable to parse response. expected=%q, actual=%q", r.Body, err)
				}
			}

			if tt.httpStatusExpected == http.StatusOK && !reflect.DeepEqual(actual, tt.requestBodyAddressData) {
				t.Errorf("input %v, unexpected address book data. got %+v wan %+v",
					tt.requestBodyAddressData, actual, tt.expectedAddressData)
			}

		}
	})
	t.Run("test delete entry", func(t *testing.T) {
		t.Helper()
		request, err := http.NewRequest(http.MethodDelete, "/api/address/1", nil)
		if err != nil {
			t.Fatal(err)
		}
		r := httptest.NewRecorder()
		s.router.ServeHTTP(r, request)
		checkResponseCode(t, r.Code, http.StatusNoContent)
	})
	t.Run("test update entry", func(t *testing.T) {
		t.Helper()
		updateBody := bytes.NewBuffer([]byte(`{"first_name": "Jane", "last_name":"Austin","email":"jane@austin.com","phone_number":123}`))
		request, err := http.NewRequest(http.MethodPut, "/api/address/2", updateBody)
		if err != nil {
			t.Fatal(err)
		}
		r := httptest.NewRecorder()
		var actual data.AddressBook
		err = data.FromJSON(&actual, r.Body)
		t.Log(updateBody, &actual, r.Body)
		s.router.ServeHTTP(r, request)
		checkResponseCode(t, r.Code, http.StatusOK)

		if !reflect.DeepEqual(&actual, updateBody) {
			t.Errorf("unexpected address book update data. got %+v, wanted %+v", &actual, updateBody)
		}
	})

}
