package pg

import (
	"database/sql"
	"time"

	"github.com/davecgh/go-spew/spew"

	"github.com/jmoiron/sqlx"
	"github.com/nyks06/backapi"
)

const userColumns = `id,
customer_id,
firstname,
lastname,
username,
email,
password,
admin,
created_at,
updated_at,
confirmed_at,
deleted_at`

// UserRecord defines fields of a user for pg usage
type UserRecord struct {
	ID          string    `db:"id"`
	CustomerID  string    `db:"customer_id"`
	Firstname   string    `db:"firstname"`
	Lastname    string    `db:"lastname"`
	Username    string    `db:"username"`
	Email       string    `db:"email"`
	Password    string    `db:"password"`
	Admin       bool      `db:"admin"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
	ConfirmedAt time.Time `db:"confirmed_at"`
	DeletedAt   time.Time `db:"deleted_at"`
}

func (u *UserRecord) toUser() *backapi.User {
	return &backapi.User{
		ID:          u.ID,
		CustomerID:  u.CustomerID,
		Firstname:   u.Firstname,
		Lastname:    u.Lastname,
		Username:    u.Username,
		Email:       u.Email,
		Password:    u.Password,
		Admin:       u.Admin,
		CreatedAt:   u.CreatedAt,
		ConfirmedAt: u.ConfirmedAt,
		DeletedAt:   u.DeletedAt,
	}
}

func newUserRecord(u *backapi.User) *UserRecord {
	return &UserRecord{
		ID:          u.ID,
		CustomerID:  u.CustomerID,
		Firstname:   u.Firstname,
		Lastname:    u.Lastname,
		Username:    u.Username,
		Email:       u.Email,
		Password:    u.Password,
		Admin:       u.Admin,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
		ConfirmedAt: u.ConfirmedAt,
		DeletedAt:   u.DeletedAt,
	}
}

// UserStore is the PG implementation of the interface defined in the domain that allows to insert and update users information
type UserStore struct {
	DB *sqlx.DB
}

// Create tries to insert in DB a user
func (s *UserStore) Create(u *backapi.User) (*backapi.User, error) {
	rec := newUserRecord(u)

	namedStmt, err := s.DB.PrepareNamed(`
		INSERT INTO users (
			customer_id,
			firstname,
			lastname,
			username,
			email,
			password,
			admin,
			confirmed_at,
			deleted_at)
		VALUES (
			:customer_id,
			:firstname,
			:lastname,
			:username,
			:email,
			:password,
			:admin,
			:confirmed_at,
			:deleted_at)
		RETURNING ` + userColumns)

	if err != nil {
		return nil, backapi.NewInternalServerError(err.Error())
	}

	defer namedStmt.Close()

	user := new(UserRecord)
	err = namedStmt.Get(user, *rec)
	if err != nil {
		return nil, backapi.NewInternalServerError(err.Error())
	}

	return user.toUser(), nil
}

func (s *UserStore) Confirm(id string) error {
	_, err := s.DB.Exec(`UPDATE users SET confirmed_at=$1 WHERE id=$2`, time.Now(), id)
	if err != nil {
		return backapi.NewInternalServerError(err.Error())
	}

	return nil
}

func (s *UserStore) Delete(id string) error {
	_, err := s.DB.Exec(`UPDATE users SET deleted_at=$1 WHERE id=$2`, time.Now(), id)
	if err != nil {
		return backapi.NewInternalServerError(err.Error())
	}

	return nil
}

func (s *UserStore) SetPassword(id, password string) error {
	_, err := s.DB.Exec(`UPDATE users SET password=$1 WHERE id=$2`, password, id)
	if err != nil {
		return backapi.NewInternalServerError(err.Error())
	}

	return nil
}

func (s *UserStore) SetEmail(id, email string) error {
	_, err := s.DB.Exec(`UPDATE users SET email=$1 WHERE id=$2`, email, id)
	if err != nil {
		return backapi.NewInternalServerError(err.Error())
	}

	return nil
}

// UserFinder is the PG implementation of the interface defined in the domain that allows to find user's information based on various parameters
type UserFinder struct {
	DB *sqlx.DB
}

// ByEmail returns the associated user stored with, as email, the parameter given
func (s *UserFinder) ByEmail(email string) (*backapi.User, error) {
	user := new(UserRecord)

	err := s.DB.Get(user, `SELECT `+userColumns+` FROM users WHERE email=$1 LIMIT 1`, email)
	spew.Dump(err)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, backapi.NewNotFoundError(err.Error())
		}
		return nil, backapi.NewInternalServerError(err.Error())
	}
	return user.toUser(), nil
}

// ByID returns the associated user stored with, as ID, the parameter given
func (s *UserFinder) ByID(id string) (*backapi.User, error) {
	user := new(UserRecord)

	err := s.DB.Get(user, `SELECT `+userColumns+` FROM users WHERE id=$1 LIMIT 1`, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, backapi.NewNotFoundError(err.Error())
		}
		return nil, backapi.NewInternalServerError(err.Error())
	}
	return user.toUser(), nil
}

// List returns all the pronostics group stored
func (s *UserFinder) List() ([]backapi.User, error) {
	uRec := make([]UserRecord, 0)

	err := s.DB.Select(&uRec, `SELECT id, firstname, lastname, username, email FROM users ORDER BY created_at DESC`)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, backapi.NewNotFoundError(err.Error())
		}
		return nil, backapi.NewInternalServerError(err.Error())
	}

	users := make([]backapi.User, 0)
	for _, user := range uRec {
		users = append(users, backapi.User{
			ID:        user.ID,
			Firstname: user.Firstname,
			Lastname:  user.Lastname,
			Username:  user.Username,
			Email:     user.Email,
		})
	}

	return users, nil
}
