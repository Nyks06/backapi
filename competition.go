package webcore

import "time"

// Sport defines all the fields that represent a group of pronostics
type Competition struct {
	ID      string
	SportID string

	Name string

	StartAt time.Time
	EndAt   time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}

// SportStore defines all the methods usable to write data in the store for pronostics group
type CompetitionsStore interface {
	Create(*Competition) (*Competition, error)
	Delete(ID string) error
}

// TicketFinder defines all the methods usable to read data in the store for pronostics group
type CompetitionsFinder interface {
	ByID(id string) (*Competition, error)
	List() ([]Competition, error)
}
