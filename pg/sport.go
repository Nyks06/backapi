package pg

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	webcore "github.com/nyks06/backapi"
)

const sportsColumns = `id,
name,
created_at,
updated_at`

// TicketRecord defines fields of a pronostics group for pg usage
type SportRecord struct {
	ID        string    `db:"id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (sport *SportRecord) toSport() *webcore.Sport {
	return &webcore.Sport{
		ID:        sport.ID,
		Name:      sport.Name,
		CreatedAt: sport.CreatedAt,
		UpdatedAt: sport.UpdatedAt,
	}
}

func newSportRecord(sport *webcore.Sport) *SportRecord {
	return &SportRecord{
		ID:        sport.ID,
		Name:      sport.Name,
		CreatedAt: sport.CreatedAt,
		UpdatedAt: sport.UpdatedAt,
	}
}

// TicketStore is the PG implementation of the interface defined in the domain that allows to insert and update pronostics group information
type SportStore struct {
	DB *sqlx.DB
}

// Create tries to insert in DB a user
func (s *SportStore) Create(pGroup *webcore.Sport) (*webcore.Sport, error) {
	rec := newSportRecord(pGroup)

	namedStmt, err := s.DB.PrepareNamed(`
		INSERT INTO sports (
			name)
		VALUES (
			:name)
		RETURNING ` + sportsColumns)

	if err != nil {
		return nil, webcore.NewInternalServerError(err.Error())
	}

	defer namedStmt.Close()

	pGr := new(SportRecord)
	err = namedStmt.Get(pGr, *rec)
	if err != nil {
		return nil, webcore.NewInternalServerError(err.Error())
	}

	return pGr.toSport(), nil
}

// Delete...
func (s *SportStore) Delete(ID string) error {
	_, err := s.DB.Exec(`DELETE FROM sports WHERE id=$1`, ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return webcore.NewNotFoundError(err.Error())
		}
		return webcore.NewInternalServerError(err.Error())
	}
	return nil
}

// TicketFinder is the PG implementation of the interface defined in the domain that allows to find user's information based on various parameters
type SportFinder struct {
	DB *sqlx.DB
}

// ByID returns the associated user stored with, as id, the parameter given
func (s *SportFinder) ByID(ID string) (*webcore.Sport, error) {
	pGroupRec := new(SportRecord)

	err := s.DB.Get(pGroupRec, `SELECT `+sportsColumns+` FROM sports WHERE id=$1 LIMIT 1`, ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, webcore.NewNotFoundError(err.Error())
		}
		return nil, webcore.NewInternalServerError(err.Error())
	}
	return pGroupRec.toSport(), nil
}

// List returns all the pronostics group stored
func (s *SportFinder) List() ([]webcore.Sport, error) {
	pGrRec := make([]SportRecord, 0)

	err := s.DB.Select(&pGrRec, `SELECT `+sportsColumns+` FROM sports ORDER BY created_at DESC`)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, webcore.NewNotFoundError(err.Error())
		}
		return nil, webcore.NewInternalServerError(err.Error())
	}

	pGr := make([]webcore.Sport, 0)
	for _, sport := range pGrRec {
		pGr = append(pGr, webcore.Sport{
			ID:        sport.ID,
			Name:      sport.Name,
			CreatedAt: sport.CreatedAt,
			UpdatedAt: sport.UpdatedAt,
		})
	}

	return pGr, nil
}
