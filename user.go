package webcore

import "time"

// User defines all the fields that represent a user
type User struct {
	// IDs
	ID            string
	SponsorshipID string
	SponsorID     string
	CustomerID    string

	// Login information
	Email    string
	Password string

	// Personal information
	Firstname string
	Lastname  string
	Username  string

	// Status
	Admin bool
	VIP   bool

	// Steps date
	CreatedAt   time.Time
	ConfirmedAt time.Time
	UpdatedAt   time.Time
	DeletedAt   time.Time
}

// UserStore defines all the methods usable to write data in the store for users
type UserStore interface {
	Create(*User) (*User, error)
	Confirm(ID string) error
	Delete(ID string) error
	// SetPassword(string, string) error
	// UpdateDetails(string, string, string) error
	// UpdateCustomerID(string, string) error
	// UpdateSubscription(string, string, string) error
	// UpdateContactSettings(string, string, string, string) error
}

// UserFinder defines all the methods usable to read data in the store for users
type UserFinder interface {
	ByEmail(email string) (*User, error)
	ByID(id string) (*User, error)
	List() ([]User, error)
}
