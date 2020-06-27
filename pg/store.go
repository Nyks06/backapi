package pg

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Store struct {
	*sqlx.DB
}

// NewStore returns a new store using postgresql.
// The returned struct could be mainly used to query the database and store / retrieve data on it.
func NewStore(url string) (*Store, error) {
	db, err := sqlx.Connect("postgres", url)
	if err != nil {
		return nil, err
	}

	return &Store{
		DB: db,
	}, nil
}
