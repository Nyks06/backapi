package pg_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/nyks06/backapi"
	"github.com/nyks06/backapi/pg"
)

func TestSession_Create(t *testing.T) {
	db, cleanup := testDatabase(t)
	defer cleanup()

	store := pg.SessionStore{
		DB: db,
	}

	s := &backapi.Session{
		UserID:    "user-123",
		ExpiresAt: time.Now().Add(48 * time.Hour),
	}

	sessStored, err := store.Create(s)
	require.NoError(t, err)
	require.Equal(t, "user-123", sessStored.UserID)
	require.NotZero(t, sessStored.ID)
	require.NotZero(t, sessStored.CreatedAt)
	require.NotZero(t, sessStored.UpdatedAt)
	require.NotZero(t, sessStored.ExpiresAt)
}

func TestSession_FindByID(t *testing.T) {
	db, cleanup := testDatabase(t)
	defer cleanup()

	store := pg.SessionStore{
		DB: db,
	}
	finder := pg.SessionFinder{
		DB: db,
	}

	s := &backapi.Session{
		UserID:    "user-123",
		ExpiresAt: time.Now().Add(48 * time.Hour),
	}
	sessionStored, err := store.Create(s)
	require.NoError(t, err)
	require.NotNil(t, sessionStored)

	s2 := &backapi.Session{
		UserID:    "user-123",
		ExpiresAt: time.Now().Add(48 * time.Hour),
	}
	sessionStored2, err := store.Create(s2)
	require.NoError(t, err)
	require.NotNil(t, sessionStored2)

	s3 := &backapi.Session{
		UserID:    "user-124",
		ExpiresAt: time.Now().Add(48 * time.Hour),
	}
	sessionStored3, err := store.Create(s3)
	require.NoError(t, err)
	require.NotNil(t, sessionStored3)

	found, err := finder.ByID(sessionStored3.ID)
	require.Nil(t, err)
	require.NotNil(t, found)
	require.Equal(t, sessionStored3, found)

	notFound, err := finder.ByID("notfound_id")
	require.NotNil(t, err)
	require.Nil(t, notFound)

	err = store.Delete(sessionStored3.ID)
	require.Nil(t, err)

	found, err = finder.ByID(sessionStored3.ID)
	require.NotNil(t, err)
	require.Nil(t, found)
}
