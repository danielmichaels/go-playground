package data

import (
	"fmt"
	"time"
)

// ErrAddressNotFound is an error raised when a product can not be found in the database
var ErrAddressNotFound = fmt.Errorf("address not found")

type AddressBook struct {
	ID          int    `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	PhoneNumber int    `json:"phone_number,string"`
	CreatedON   string `json:"-"` // "-" omit from response
	UpdatedON   string `json:"-"` // "-" omit from response
	DeletedON   string `json:"-"` // "-" omit from response
}

var AddressBookStore = []*AddressBook{
	{
		ID:          1,
		FirstName:   "Michael",
		LastName:    "Jordan",
		Email:       "twotrey@hornets.nba",
		PhoneNumber: 23232323,
		CreatedON:   time.Now().UTC().String(),
		UpdatedON:   time.Now().UTC().String(),
		DeletedON:   time.Now().UTC().String(),
	},
	{
		ID:          2,
		FirstName:   "Frank",
		LastName:    "Herbert",
		Email:       "frank@arrakis.com",
		PhoneNumber: 12341234,
		CreatedON:   time.Now().UTC().String(),
		UpdatedON:   time.Now().UTC().String(),
		DeletedON:   time.Now().UTC().String(),
	},
}

// GetAddressBook is a helper which returns all the books in the store
func GetAddressBook() []*AddressBook {
	return AddressBookStore
}

// findIndexByAddressID finds the index of a product in the database
// returns -1 when no product can be found
func findIndexByAddressID(id int) int {
	for i, p := range AddressBookStore {
		if p.ID == id {
			return i
		}
	}
	return -1
}

// FindAddress iterates over the AddressBookStore and looks for a matching ID
func FindAddress(id int) (*AddressBook, int, error) {
	for i, p := range AddressBookStore {
		if p.ID == id {
			return p, i, nil
		}
	}
	return nil, -1, ErrAddressNotFound
}
func DeleteAddress(id int) error {
	addr := findIndexByAddressID(id)
	if addr == -1 {
		return ErrAddressNotFound
	}
	AddressBookStore = append(AddressBookStore[:addr], AddressBookStore[addr+1])
	return nil
}

// getNextID is a utility function that will create an ID value based on the number of addresses
// in the store plus 1.
func getNextID() int {
	lastAddress := AddressBookStore[len(AddressBookStore)-1]
	return lastAddress.ID + 1
}

// AddAddressToStore is a helper which creates a new entry in the address store
func AddAddressToStore(a *AddressBook) {
	a.ID = getNextID()
	AddressBookStore = append(AddressBookStore, a)
}
