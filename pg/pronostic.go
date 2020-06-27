package pg

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	webcore "github.com/nyks06/backapi"
)

const pronosticColumns = `id,
ticket_id,
first_team,
second_team,
pronostic,
competition_id,
sport_id,
status,
odd,
event_date,
created_at,
updated_at`

// PronosticsGroupRecord defines fields of a pronostics group for pg usage
type PronosticRecord struct {
	ID            string    `db:"id"`
	TicketID      string    `db:"ticket_id"`
	FirstTeam     string    `db:"first_team"`
	SecondTeam    string    `db:"second_team"`
	Pronostic     string    `db:"pronostic"`
	CompetitionID string    `db:"competition_id"`
	SportID       string    `db:"sport_id"`
	Status        string    `db:"status"`
	Odd           float64   `db:"odd"`
	EventDate     time.Time `db:"event_date"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}

func (prono *PronosticRecord) toPronostic() *webcore.Pronostic {
	return &webcore.Pronostic{
		ID:         prono.ID,
		TicketID:   prono.TicketID,
		FirstTeam:  prono.FirstTeam,
		SecondTeam: prono.SecondTeam,
		Pronostic:  prono.Pronostic,
		Competition: &webcore.Competition{
			ID: prono.CompetitionID,
		},
		Sport: &webcore.Sport{
			ID: prono.SportID,
		},
		Status:    prono.Status,
		Odd:       prono.Odd,
		EventDate: prono.EventDate,
		CreatedAt: prono.CreatedAt,
		UpdatedAt: prono.UpdatedAt,
	}
}

func newPronosticRecord(prono *webcore.Pronostic) *PronosticRecord {
	return &PronosticRecord{
		ID:            prono.ID,
		TicketID:      prono.TicketID,
		FirstTeam:     prono.FirstTeam,
		SecondTeam:    prono.SecondTeam,
		Pronostic:     prono.Pronostic,
		CompetitionID: prono.Competition.ID,
		SportID:       prono.Sport.ID,
		Status:        prono.Status,
		Odd:           prono.Odd,
		EventDate:     prono.EventDate,
		CreatedAt:     prono.CreatedAt,
		UpdatedAt:     prono.UpdatedAt,
	}
}

// PronosticsGroupStore is the PG implementation of the interface defined in the domain that allows to insert and update pronostics group information
type PronosticStore struct {
	DB *sqlx.DB
}

// Create tries to insert in DB a user
func (s *PronosticStore) Create(prono *webcore.Pronostic) (*webcore.Pronostic, error) {
	rec := newPronosticRecord(prono)

	namedStmt, err := s.DB.PrepareNamed(`
		INSERT INTO pronostics (
			ticket_id,
			first_team,
			second_team,
			pronostic,
			competition_id,
			sport_id,
			status,
			odd,
			event_date)
		VALUES (
			:ticket_id,
			:first_team,
			:second_team,
			:pronostic,
			:competition_id,
			:sport_id,
			:status,
			:odd,
			:event_date)
		RETURNING ` + pronosticColumns)

	if err != nil {
		return nil, webcore.NewInternalServerError(err.Error())
	}

	defer namedStmt.Close()

	p := new(PronosticRecord)
	err = namedStmt.Get(p, *rec)
	if err != nil {
		return nil, webcore.NewInternalServerError(err.Error())
	}

	return p.toPronostic(), nil
}

// Delete...
func (s *PronosticStore) Delete(ID string) error {
	_, err := s.DB.Exec(`DELETE FROM pronostics WHERE id=$1`, ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return webcore.NewNotFoundError(err.Error())
		}
		return webcore.NewInternalServerError(err.Error())
	}
	return nil
}

func (s *PronosticStore) Update(ID string, Status string) error {
	_, err := s.DB.Exec(`UPDATE pronostics SET status=$1 WHERE id=$2`, Status, ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return webcore.NewNotFoundError(err.Error())
		}
		return webcore.NewInternalServerError(err.Error())
	}
	return nil
}

// PronosticsGroupFinder is the PG implementation of the interface defined in the domain that allows to find user's information based on various parameters
type PronosticFinder struct {
	DB *sqlx.DB
}

// ByID returns the associated user stored with, as id, the parameter given
func (s *PronosticFinder) ByID(ID string) (*webcore.Pronostic, error) {
	pRec := new(PronosticRecord)

	err := s.DB.Get(pRec, `SELECT `+pronosticColumns+` FROM pronostics WHERE id=$1 LIMIT 1`, ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, webcore.NewNotFoundError(err.Error())
		}
		return nil, webcore.NewInternalServerError(err.Error())
	}
	return pRec.toPronostic(), nil
}

// List returns all the pronostics group stored
func (s *PronosticFinder) List() ([]webcore.Pronostic, error) {
	pRec := make([]PronosticRecord, 0)

	err := s.DB.Select(&pRec, `SELECT `+pronosticColumns+` FROM pronostics ORDER BY created_at DESC`)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, webcore.NewNotFoundError(err.Error())
		}
		return nil, webcore.NewInternalServerError(err.Error())
	}

	p := make([]webcore.Pronostic, 0)
	for _, prono := range pRec {
		p = append(p, webcore.Pronostic{
			ID:         prono.ID,
			TicketID:   prono.TicketID,
			FirstTeam:  prono.FirstTeam,
			SecondTeam: prono.SecondTeam,
			Pronostic:  prono.Pronostic,
			Competition: &webcore.Competition{
				ID: prono.CompetitionID,
			},
			Sport: &webcore.Sport{
				ID: prono.SportID,
			},
			Status:    prono.Status,
			Odd:       prono.Odd,
			EventDate: prono.EventDate,
			CreatedAt: prono.CreatedAt,
			UpdatedAt: prono.UpdatedAt,
		})
	}

	return p, nil
}

// List returns all the pronostics group stored
func (s *PronosticFinder) ListByTicketID(TicketID string) ([]webcore.Pronostic, error) {
	pRec := make([]PronosticRecord, 0)

	err := s.DB.Select(&pRec, `SELECT `+pronosticColumns+` FROM pronostics WHERE ticket_id=$1 ORDER BY created_at DESC`, TicketID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, webcore.NewNotFoundError(err.Error())
		}
		return nil, webcore.NewInternalServerError(err.Error())
	}

	p := make([]webcore.Pronostic, 0)
	for _, prono := range pRec {
		p = append(p, webcore.Pronostic{
			ID:         prono.ID,
			TicketID:   prono.TicketID,
			FirstTeam:  prono.FirstTeam,
			SecondTeam: prono.SecondTeam,
			Pronostic:  prono.Pronostic,
			Competition: &webcore.Competition{
				ID: prono.CompetitionID,
			},
			Sport: &webcore.Sport{
				ID: prono.SportID,
			},
			Status:    prono.Status,
			Odd:       prono.Odd,
			EventDate: prono.EventDate,
			CreatedAt: prono.CreatedAt,
			UpdatedAt: prono.UpdatedAt,
		})
	}

	return p, nil
}
