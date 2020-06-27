package webcore

import "time"

// Sport defines all the fields that represent a group of pronostics
type Sport struct {
	ID        string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// SportStore defines all the methods usable to write data in the store for pronostics group
type SportStore interface {
	Create(*Sport) (*Sport, error)
	Delete(ID string) error
}

// SportFinder defines all the methods usable to read data in the store for pronostics group
type SportFinder interface {
	ByID(id string) (*Sport, error)
	List() ([]Sport, error)
}
