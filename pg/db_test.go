package pg_test

import (
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func testDatabase(t *testing.T) (*sqlx.DB, func()) {
	uri := os.Getenv("POSTGRESQL_DATASOURCE")
	if uri == "" {
		t.Error("Env var POSTGRESQL_DATASOURCE not found")
	}

	db, err := sqlx.Connect("postgres", uri)
	require.NoError(t, err)

	return db, func() {
		db.Close()
	}
}
