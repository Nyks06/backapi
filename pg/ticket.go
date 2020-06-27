package pg

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	webcore "github.com/nyks06/backapi"
)

const ticketsColumns = `id,
title,
stake,
public,
live,
risk,
pack,
created_at,
updated_at`

// TicketRecord defines fields of a pronostics group for pg usage
type TicketRecord struct {
	ID        string    `db:"id"`
	Title     string    `db:"title"`
	Stake     float64   `db:"stake"`
	Public    bool      `db:"public"`
	Live      bool      `db:"live"`
	Risk      string    `db:"risk"`
	Pack      string    `db:"pack"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (ticket *TicketRecord) toTicket() *webcore.Ticket {
	return &webcore.Ticket{
		ID:        ticket.ID,
		Title:     ticket.Title,
		Stake:     ticket.Stake,
		Public:    ticket.Public,
		Live:      ticket.Live,
		Risk:      ticket.Risk,
		Pack:      ticket.Pack,
		CreatedAt: ticket.CreatedAt,
		UpdatedAt: ticket.UpdatedAt,
	}
}

func newTicketRecord(u *webcore.Ticket) *TicketRecord {
	return &TicketRecord{
		ID:        u.ID,
		Title:     u.Title,
		Stake:     u.Stake,
		Public:    u.Public,
		Live:      u.Live,
		Risk:      u.Risk,
		Pack:      u.Pack,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// TicketStore is the PG implementation of the interface defined in the domain that allows to insert and update pronostics group information
type TicketStore struct {
	DB *sqlx.DB
}

// Create tries to insert in DB a user
func (s *TicketStore) Create(pGroup *webcore.Ticket) (*webcore.Ticket, error) {
	rec := newTicketRecord(pGroup)

	namedStmt, err := s.DB.PrepareNamed(`
		INSERT INTO tickets (
			title,
			stake,
			public,
			live,
			risk,
			pack)
		VALUES (
			:title,
			:stake,
			:public,
			:live,
			:risk,
			:pack)
		RETURNING ` + ticketsColumns)

	if err != nil {
		return nil, webcore.NewInternalServerError(err.Error())
	}

	defer namedStmt.Close()

	pGr := new(TicketRecord)
	err = namedStmt.Get(pGr, *rec)
	if err != nil {
		return nil, webcore.NewInternalServerError(err.Error())
	}

	return pGr.toTicket(), nil
}

// Delete...
func (s *TicketStore) Delete(ID string) error {
	_, err := s.DB.Exec(`DELETE FROM tickets WHERE id=$1`, ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return webcore.NewNotFoundError(err.Error())
		}
		return webcore.NewInternalServerError(err.Error())
	}
	return nil
}

func (s *TicketStore) Update(ID string, ticket *webcore.Ticket) error {
	_, err := s.DB.Exec(`UPDATE tickets SET title=$1, stake=$2, public=$3, live=$4, risk=$5 WHERE id=$6`, ticket.Title, ticket.Stake, ticket.Public, ticket.Live, ticket.Risk, ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return webcore.NewNotFoundError(err.Error())
		}
		return webcore.NewInternalServerError(err.Error())
	}
	return nil
}

// TicketFinder is the PG implementation of the interface defined in the domain that allows to find user's information based on various parameters
type TicketFinder struct {
	DB *sqlx.DB
}

// ByID returns the associated user stored with, as id, the parameter given
func (s *TicketFinder) ByID(ID string) (*webcore.Ticket, error) {
	pGroupRec := new(TicketRecord)

	err := s.DB.Get(pGroupRec, `SELECT `+ticketsColumns+` FROM tickets WHERE id=$1 LIMIT 1`, ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, webcore.NewNotFoundError(err.Error())
		}
		return nil, webcore.NewInternalServerError(err.Error())
	}
	return pGroupRec.toTicket(), nil
}

// List returns all the pronostics group stored
func (s *TicketFinder) List() ([]webcore.Ticket, error) {
	pGrRec := make([]TicketRecord, 0)

	err := s.DB.Select(&pGrRec, `SELECT `+ticketsColumns+` FROM tickets ORDER BY created_at DESC`)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, webcore.NewNotFoundError(err.Error())
		}
		return nil, webcore.NewInternalServerError(err.Error())
	}

	pGr := make([]webcore.Ticket, 0)
	for _, prono := range pGrRec {
		pGr = append(pGr, webcore.Ticket{
			ID:        prono.ID,
			Title:     prono.Title,
			Stake:     prono.Stake,
			Public:    prono.Public,
			Live:      prono.Live,
			Risk:      prono.Risk,
			Pack:      prono.Pack,
			CreatedAt: prono.CreatedAt,
			UpdatedAt: prono.UpdatedAt,
		})
	}

	return pGr, nil
}
