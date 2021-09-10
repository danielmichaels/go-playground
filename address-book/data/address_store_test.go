package data

import (
	"testing"
	"time"
)

func TestAddAddressToStore(t *testing.T) {
	t.Run("validation is working", func(t *testing.T) {
		t.Helper()
		a := &AddressBook{
			ID:          1,
			FirstName:   "Michael",
			LastName:    "Jordan",
			Email:       "twotrey@hornets.nba",
			PhoneNumber: 23232323,
			CreatedON:   time.Now().UTC().String(),
			UpdatedON:   time.Now().UTC().String(),
			DeletedON:   time.Now().UTC().String(),
		}

		err := a.Validate()

		if err != nil {
			t.Fatal(err)
		}
	})
	t.Run("validation fails on no data", func(t *testing.T) {
		t.Helper()
		a := &AddressBook{}
		err := a.Validate()
		if err == nil {
			t.Errorf("passing in no data will fail validation")
		}
	})
	t.Run("validation fails on invalid key type", func(t *testing.T) {
		t.Helper()
		a := &AddressBook{
			FirstName:   "person",
			LastName:    "surname",
			Email:       "not-valid",
			PhoneNumber: 111,
		}
		err := a.Validate()
		if err == nil {
			t.Errorf("passing in invalid data will fail validation")
		}
	})
}
