package webcore

import "time"

// Pronostic defines all the fields that represent a pronostic
type Pronostic struct {
	// IDs
	ID       string
	TicketID string

	// Ticket information
	FirstTeam   string
	SecondTeam  string
	Title       string
	Competition *Competition
	Sport       *Sport
	Status      string
	Odd         float64
	StartsAt    time.Time

	// Steps Date
	CreatedAt time.Time
	UpdatedAt time.Time
}

// PronosticStore defines all the methods usable to write data in the store for pronostic
type PronosticsStore interface {
	Create(*Pronostic) (*Pronostic, error)
	Delete(string) error
	UpdateStatus(id string, status string) error
}

// PronosticsFinder defines all the methods usable to read data in the store for pronostic
type PronosticsFinder interface {
	ByID(id string) (*Pronostic, error)
	List() ([]Pronostic, error)
	ListByTicketID(ID string) ([]Pronostic, error)
}
