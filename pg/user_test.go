package pg_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/nyks06/backapi"
	"github.com/nyks06/backapi/pg"
)

func TestUser_Create(t *testing.T) {
	db, cleanup := testDatabase(t)
	defer cleanup()

	store := pg.UserStore{
		DB: db,
	}

	u := &webcore.User{
		Firstname:   "Vincent",
		Lastname:    "Vielle",
		Username:    "Nykxs",
		Email:       "vincent_create@gmail.com",
		PhoneNumber: "+33611223344",
		Password:    "123pass",
	}

	userStored, err := store.Create(u)
	require.NoError(t, err)
	require.Equal(t, "Vincent", userStored.Firstname)
	require.Equal(t, "Vielle", userStored.Lastname)
	require.Equal(t, "Nykxs", userStored.Username)
	require.Equal(t, "vincent_create@gmail.com", userStored.Email)
	require.Equal(t, "+33611223344", userStored.PhoneNumber)
	require.Equal(t, "123pass", userStored.Password)
	require.NotZero(t, userStored.ID)
	require.NotZero(t, userStored.CreatedAt)
	require.NotZero(t, userStored.UpdatedAt)
}

func TestUser_FindByEmail(t *testing.T) {
	db, cleanup := testDatabase(t)
	defer cleanup()

	store := pg.UserStore{
		DB: db,
	}
	finder := pg.UserFinder{
		DB: db,
	}

	u := &webcore.User{
		Firstname:   "Vincent",
		Lastname:    "Vielle",
		Username:    "Nykxs",
		Email:       "vincent_email@gmail.com",
		PhoneNumber: "+33611223344",
		Password:    "123pass",
	}
	userStored, err := store.Create(u)
	require.NoError(t, err)
	require.NotNil(t, userStored)

	u2 := &webcore.User{
		Firstname:   "Vincent",
		Lastname:    "Vielle",
		Username:    "Nykxs",
		Email:       "vincent_email2@gmail.com",
		PhoneNumber: "+33611223355",
		Password:    "123pass",
	}
	userStored2, err := store.Create(u2)
	require.NoError(t, err)
	require.NotNil(t, userStored2)

	u3 := &webcore.User{
		Firstname:   "Vincent",
		Lastname:    "Vielle",
		Username:    "Nykxs",
		Email:       "vincent_email3@gmail.com",
		PhoneNumber: "+33611223366",
		Password:    "123pass",
	}
	userStored3, err := store.Create(u3)
	require.NoError(t, err)
	require.NotNil(t, userStored3)

	found, err := finder.ByEmail("vincent_email3@gmail.com")
	require.Nil(t, err)
	require.NotNil(t, found)
	require.Equal(t, u3.PhoneNumber, found.PhoneNumber)
	require.Equal(t, u3.Password, found.Password)
	require.NotZero(t, found.CreatedAt)
	require.NotZero(t, found.UpdatedAt)

	notFound, err := finder.ByEmail("notfound@gmail.com")
	require.NotNil(t, err)
	require.Nil(t, notFound)
}

func TestUser_FindByID(t *testing.T) {
	db, cleanup := testDatabase(t)
	defer cleanup()

	store := pg.UserStore{
		DB: db,
	}
	finder := pg.UserFinder{
		DB: db,
	}

	u := &webcore.User{
		Email:       "vincent_phone@gmail.com",
		PhoneNumber: "+33611223344",
		Password:    "123pass",
	}
	userStored, err := store.Create(u)
	require.NoError(t, err)
	require.NotNil(t, userStored)

	u2 := &webcore.User{
		Email:       "vincent_phone2@gmail.com",
		PhoneNumber: "+33611223355",
		Password:    "123pass",
	}
	userStored2, err := store.Create(u2)
	require.NoError(t, err)
	require.NotNil(t, userStored2)

	u3 := &webcore.User{
		Email:       "vincent_phone3@gmail.com",
		PhoneNumber: "+33611223366",
		Password:    "123pass",
	}
	userStored3, err := store.Create(u3)
	require.NoError(t, err)
	require.NotNil(t, userStored3)

	found, err := finder.ByID(userStored3.ID)
	require.Nil(t, err)
	require.NotNil(t, found)
	require.Equal(t, userStored3, found)

	notFound, err := finder.ByID("notfound_id")
	require.NotNil(t, err)
	require.Nil(t, notFound)
}

func TestUser_ListEmails(t *testing.T) {
	db, cleanup := testDatabase(t)
	defer cleanup()

	store := pg.UserStore{
		DB: db,
	}
	finder := pg.UserFinder{
		DB: db,
	}

	u := &webcore.User{
		Email:       "vincent_test@gmail.com",
		PhoneNumber: "+33611223344",
		Password:    "123pass",
	}
	userStored, err := store.Create(u)
	require.NoError(t, err)
	require.NotNil(t, userStored)

	u2 := &webcore.User{
		Email:       "vincent_test@gmail.com",
		PhoneNumber: "+33611223355",
		Password:    "123pass",
	}
	userStored2, err := store.Create(u2)
	require.NoError(t, err)
	require.NotNil(t, userStored2)

	u3 := &webcore.User{
		Email:       "vincent_test@gmail.com",
		PhoneNumber: "+33611223366",
		Password:    "123pass",
	}
	userStored3, err := store.Create(u3)
	require.NoError(t, err)
	require.NotNil(t, userStored3)

	emails, err := finder.ListEmails()
	require.Nil(t, err)
	require.Len(t, emails, 10)
}

func TestUser_ListPhoneNumbers(t *testing.T) {
	db, cleanup := testDatabase(t)
	defer cleanup()

	store := pg.UserStore{
		DB: db,
	}
	finder := pg.UserFinder{
		DB: db,
	}

	u := &webcore.User{
		Email:       "vincent_test@gmail.com",
		PhoneNumber: "+33611223344",
		Password:    "123pass",
	}
	userStored, err := store.Create(u)
	require.NoError(t, err)
	require.NotNil(t, userStored)

	u2 := &webcore.User{
		Email:       "vincent_test@gmail.com",
		PhoneNumber: "+33611223355",
		Password:    "123pass",
	}
	userStored2, err := store.Create(u2)
	require.NoError(t, err)
	require.NotNil(t, userStored2)

	u3 := &webcore.User{
		Email:       "vincent_test@gmail.com",
		PhoneNumber: "+33611223366",
		Password:    "123pass",
	}
	userStored3, err := store.Create(u3)
	require.NoError(t, err)
	require.NotNil(t, userStored3)

	emails, err := finder.ListEmails()
	require.Nil(t, err)
	require.Len(t, emails, 13)
}
