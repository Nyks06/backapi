package webcore

import "time"

// Ticket defines all the fields that represent a group of pronostics
type Ticket struct {
	ID     string
	PackID string

	Title  string
	Stake  float64
	Public bool
	Live   bool
	Risk   string

	CreatedAt time.Time
	UpdatedAt time.Time
}

// FilledTicket ...
type FilledTicket struct {
	ID     string
	PackID string

	Title  string
	Stake  float64
	Public bool
	Live   bool
	Risk   string

	CreatedAt time.Time
	UpdatedAt time.Time
}

// TicketStore defines all the methods usable to write data in the store for pronostics group
type TicketStore interface {
	Create(*Ticket) (*Ticket, error)
	Delete(ID string) error
}

// TicketFinder defines all the methods usable to read data in the store for pronostics group
type TicketFinder interface {
	ByID(id string) (*Ticket, error)
	List() ([]Ticket, error)
	ListByPackID(string) ([]Ticket, error)
}
