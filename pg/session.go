package pg

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/nyks06/backapi"
)

const sessionColumns = `id,
user_id,
created_at,
updated_at,
expires_at`

// SessionRecord defines fields of a session for pg usage
type SessionRecord struct {
	ID        string    `db:"id"`
	UserID    string    `db:"user_id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	ExpiresAt time.Time `db:"expires_at"`
}

func (s *SessionRecord) toSession() *webcore.Session {
	return &webcore.Session{
		ID:        s.ID,
		UserID:    s.UserID,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
		ExpiresAt: s.ExpiresAt,
	}
}

func newSessionRecord(s *webcore.Session) *SessionRecord {
	return &SessionRecord{
		ID:        s.ID,
		UserID:    s.UserID,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
		ExpiresAt: s.ExpiresAt,
	}
}

// SessionStore is the PG implementation of the interface defined in the domain that allows to insert and update session information
type SessionStore struct {
	DB *sqlx.DB
}

// Create tries to insert in DB a user
func (s *SessionStore) Create(sess *webcore.Session) (*webcore.Session, error) {
	rec := newSessionRecord(sess)

	namedStmt, err := s.DB.PrepareNamed(`
		INSERT INTO sessions (
			user_id,
			expires_at)
		VALUES (
			:user_id,
			:expires_at)
		RETURNING ` + sessionColumns)

	if err != nil {
		return nil, webcore.NewInternalServerError(err.Error())
	}

	defer namedStmt.Close()

	session := new(SessionRecord)
	err = namedStmt.Get(session, *rec)
	if err != nil {
		return nil, webcore.NewInternalServerError(err.Error())
	}

	return session.toSession(), nil
}

func (s *SessionStore) Delete(ID string) error {
	_, err := s.DB.Exec(`DELETE FROM sessions WHERE id=$1`, ID)
	if err != nil {
		return webcore.NewInternalServerError(err.Error())
	}

	return nil
}

func (s *SessionStore) DeleteAllForUserID(UserID string) error {
	_, err := s.DB.Exec(`DELETE FROM sessions WHERE user_id=$1`, UserID)
	if err != nil {
		return webcore.NewInternalServerError(err.Error())
	}

	return nil
}

// SessionFinder is the PG implementation of the interface defined in the domain that allows to find session's information based on various parameters
type SessionFinder struct {
	DB *sqlx.DB
}

// ByID returns the associated session stored with, as ID, the parameter given
func (s *SessionFinder) ByID(id string) (*webcore.Session, error) {
	sess := new(SessionRecord)

	err := s.DB.Get(sess, `SELECT `+sessionColumns+` FROM sessions WHERE id=$1 LIMIT 1`, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, webcore.NewNotFoundError(err.Error())
		}
		return nil, webcore.NewInternalServerError(err.Error())
	}
	return sess.toSession(), nil
}
