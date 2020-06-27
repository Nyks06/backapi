package pg

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	webcore "github.com/nyks06/backapi"
)

const competitionsColumns = `id,
sport_id,
name,
start_at,
end_at,
created_at,
updated_at`

// TicketRecord defines fields of a pronostics group for pg usage
type CompetitionRecord struct {
	ID        string    `db:"id"`
	SportID   string    `db:"sport_id"`
	Name      string    `db:"name"`
	StartAt   time.Time `db:"start_at"`
	EndAt     time.Time `db:"end_at"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (competition *CompetitionRecord) toCompetition() *webcore.Competition {
	return &webcore.Competition{
		ID:        competition.ID,
		SportID:   competition.SportID,
		Name:      competition.Name,
		StartAt:   competition.StartAt,
		EndAt:     competition.EndAt,
		CreatedAt: competition.CreatedAt,
		UpdatedAt: competition.UpdatedAt,
	}
}

func newCompetitionRecord(competition *webcore.Competition) *CompetitionRecord {
	return &CompetitionRecord{
		ID:        competition.ID,
		SportID:   competition.SportID,
		Name:      competition.Name,
		StartAt:   competition.StartAt,
		EndAt:     competition.EndAt,
		CreatedAt: competition.CreatedAt,
		UpdatedAt: competition.UpdatedAt,
	}
}

// TicketStore is the PG implementation of the interface defined in the domain that allows to insert and update pronostics group information
type CompetitionStore struct {
	DB *sqlx.DB
}

// Create tries to insert in DB a user
func (s *CompetitionStore) Create(pGroup *webcore.Competition) (*webcore.Competition, error) {
	rec := newCompetitionRecord(pGroup)

	namedStmt, err := s.DB.PrepareNamed(`
		INSERT INTO competitions (
			name,
			sport_id,
			start_at,
			end_at)
		VALUES (
			:name,
			:sport_id,
			:start_at,
			:end_at)
		RETURNING ` + competitionsColumns)

	if err != nil {
		return nil, webcore.NewInternalServerError(err.Error())
	}

	defer namedStmt.Close()

	pGr := new(CompetitionRecord)
	err = namedStmt.Get(pGr, *rec)
	if err != nil {
		return nil, webcore.NewInternalServerError(err.Error())
	}

	return pGr.toCompetition(), nil
}

// Delete...
func (s *CompetitionStore) Delete(ID string) error {
	_, err := s.DB.Exec(`DELETE FROM competitions WHERE id=$1`, ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return webcore.NewNotFoundError(err.Error())
		}
		return webcore.NewInternalServerError(err.Error())
	}
	return nil
}

// TicketFinder is the PG implementation of the interface defined in the domain that allows to find user's information based on various parameters
type CompetitionFinder struct {
	DB *sqlx.DB
}

// ByID returns the associated user stored with, as id, the parameter given
func (s *CompetitionFinder) ByID(ID string) (*webcore.Competition, error) {
	pGroupRec := new(CompetitionRecord)

	err := s.DB.Get(pGroupRec, `SELECT `+competitionsColumns+` FROM competitions WHERE id=$1 LIMIT 1`, ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, webcore.NewNotFoundError(err.Error())
		}
		return nil, webcore.NewInternalServerError(err.Error())
	}
	return pGroupRec.toCompetition(), nil
}

// List returns all the pronostics group stored
func (s *CompetitionFinder) List() ([]webcore.Competition, error) {
	pGrRec := make([]CompetitionRecord, 0)

	err := s.DB.Select(&pGrRec, `SELECT `+competitionsColumns+` FROM competitions ORDER BY created_at DESC`)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, webcore.NewNotFoundError(err.Error())
		}
		return nil, webcore.NewInternalServerError(err.Error())
	}

	pGr := make([]webcore.Competition, 0)
	for _, competition := range pGrRec {
		pGr = append(pGr, webcore.Competition{
			ID:        competition.ID,
			SportID:   competition.SportID,
			Name:      competition.Name,
			StartAt:   competition.StartAt,
			EndAt:     competition.EndAt,
			CreatedAt: competition.CreatedAt,
			UpdatedAt: competition.UpdatedAt,
		})
	}

	return pGr, nil
}
