package pg

import (
	"database/sql"
	"time"

	"github.com/davecgh/go-spew/spew"

	"github.com/jmoiron/sqlx"
	webcore "github.com/nyks06/backapi"
)

const userColumns = `id,
customer_id,
sponsor_id,
sponsorship_id,
firstname,
lastname,
username,
email,
phone_number,
telegram,
password,
confirmed,
admin,
sub_id_prive,
sub_id_champion,
sub_id_fun,
created_at,
updated_at`

// UserRecord defines fields of a user for pg usage
type UserRecord struct {
	ID            string    `db:"id"`
	CustomerID    string    `db:"customer_id"`
	SponsorID     string    `db:"sponsor_id"`
	SponsorshipID string    `db:"sponsorship_id"`
	Firstname     string    `db:"firstname"`
	Lastname      string    `db:"lastname"`
	Username      string    `db:"username"`
	Email         string    `db:"email"`
	PhoneNumber   string    `db:"phone_number"`
	Telegram      string    `db:"telegram"`
	Password      string    `db:"password"`
	Confirmed     bool      `db:"confirmed"`
	Admin         bool      `db:"admin"`
	SubIDPrive    string    `db:"sub_id_prive"`
	SubIDChampion string    `db:"sub_id_champion"`
	SubIDFun      string    `db:"sub_id_fun"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}

func (u *UserRecord) toUser() *webcore.User {
	return &webcore.User{
		ID:            u.ID,
		CustomerID:    u.CustomerID,
		SponsorID:     u.SponsorID,
		SponsorshipID: u.SponsorshipID,
		Firstname:     u.Firstname,
		Lastname:      u.Lastname,
		Username:      u.Username,
		Email:         u.Email,
		PhoneNumber:   u.PhoneNumber,
		Telegram:      u.Telegram,
		Password:      u.Password,
		Confirmed:     u.Confirmed,
		Admin:         u.Admin,
		SubIDPrive:    u.SubIDPrive,
		SubIDChampion: u.SubIDChampion,
		SubIDFun:      u.SubIDFun,
		CreatedAt:     u.CreatedAt,
		UpdatedAt:     u.UpdatedAt,
	}
}

func newUserRecord(u *webcore.User) *UserRecord {
	return &UserRecord{
		ID:            u.ID,
		CustomerID:    u.CustomerID,
		SponsorID:     u.SponsorID,
		SponsorshipID: u.SponsorshipID,
		Firstname:     u.Firstname,
		Lastname:      u.Lastname,
		Username:      u.Username,
		Email:         u.Email,
		PhoneNumber:   u.PhoneNumber,
		Telegram:      u.Telegram,
		Password:      u.Password,
		Confirmed:     u.Confirmed,
		Admin:         u.Admin,
		SubIDPrive:    u.SubIDPrive,
		SubIDChampion: u.SubIDChampion,
		SubIDFun:      u.SubIDFun,
		CreatedAt:     u.CreatedAt,
		UpdatedAt:     u.UpdatedAt,
	}
}

// UserStore is the PG implementation of the interface defined in the domain that allows to insert and update users information
type UserStore struct {
	DB *sqlx.DB
}

// Create tries to insert in DB a user
func (s *UserStore) Create(u *webcore.User) (*webcore.User, error) {
	rec := newUserRecord(u)

	namedStmt, err := s.DB.PrepareNamed(`
		INSERT INTO users (
			customer_id,
			sponsor_id,
			sponsorship_id,
			firstname,
			lastname,
			username,
			email,
			phone_number,
			telegram,
			password)
		VALUES (
			:customer_id,
			:sponsor_id,
			:sponsorship_id,
			:firstname,
			:lastname,
			:username,
			:email,
			:phone_number,
			:telegram,
			:password)
		RETURNING ` + userColumns)

	if err != nil {
		return nil, webcore.NewInternalServerError(err.Error())
	}

	defer namedStmt.Close()

	user := new(UserRecord)
	err = namedStmt.Get(user, *rec)
	if err != nil {
		return nil, webcore.NewInternalServerError(err.Error())
	}

	return user.toUser(), nil
}

func (s *UserStore) Confirm(id string) error {
	_, err := s.DB.Exec(`UPDATE users SET confirmed=$1 WHERE id=$2`, true, id)
	if err != nil {
		return webcore.NewInternalServerError(err.Error())
	}

	return nil
}

func (s *UserStore) SetPassword(id, password string) error {
	_, err := s.DB.Exec(`UPDATE users SET password=$1 WHERE id=$2`, password, id)
	if err != nil {
		return webcore.NewInternalServerError(err.Error())
	}

	return nil
}

func (s *UserStore) UpdateDetails(id, firstname string, lastname string) error {
	_, err := s.DB.Exec(`UPDATE users SET firstname=$1, lastname=$2 WHERE id=$3`, firstname, lastname, id)
	if err != nil {
		return webcore.NewInternalServerError(err.Error())
	}

	return nil
}

func (s *UserStore) UpdateCustomerID(id, cid string) error {
	_, err := s.DB.Exec(`UPDATE users SET customer_id=$1 WHERE id=$2`, cid, id)
	if err != nil {
		return webcore.NewInternalServerError(err.Error())
	}

	return nil
}

func (s *UserStore) UpdateSubscription(id string, subscriptionID string, pack string) error {

	_, err := s.DB.Exec(`UPDATE users SET `+pack+`=$1 WHERE id=$2`, subscriptionID, id)
	if err != nil {
		return webcore.NewInternalServerError(err.Error())
	}

	return nil
}

func (s *UserStore) UpdateContactSettings(id, email string, phone string, telegram string) error {
	_, err := s.DB.Exec(`UPDATE users SET email=$1, phone_number=$2, telegram=$3 WHERE id=$4`, email, phone, telegram, id)
	if err != nil {
		return webcore.NewInternalServerError(err.Error())
	}

	return nil
}

// UserFinder is the PG implementation of the interface defined in the domain that allows to find user's information based on various parameters
type UserFinder struct {
	DB *sqlx.DB
}

// ByEmail returns the associated user stored with, as email, the parameter given
func (s *UserFinder) ByEmail(email string) (*webcore.User, error) {
	user := new(UserRecord)

	err := s.DB.Get(user, `SELECT `+userColumns+` FROM users WHERE email=$1 LIMIT 1`, email)
	spew.Dump(err)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, webcore.NewNotFoundError(err.Error())
		}
		return nil, webcore.NewInternalServerError(err.Error())
	}
	return user.toUser(), nil
}

// ByID returns the associated user stored with, as ID, the parameter given
func (s *UserFinder) ByID(id string) (*webcore.User, error) {
	user := new(UserRecord)

	err := s.DB.Get(user, `SELECT `+userColumns+` FROM users WHERE id=$1 LIMIT 1`, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, webcore.NewNotFoundError(err.Error())
		}
		return nil, webcore.NewInternalServerError(err.Error())
	}
	return user.toUser(), nil
}

// ListEmails returns all the email stored
func (s *UserFinder) ListEmails() ([]string, error) {
	emails := make([]string, 0)

	err := s.DB.Select(&emails, "SELECT email FROM users")
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, webcore.NewNotFoundError(err.Error())
		}
		return nil, webcore.NewInternalServerError(err.Error())
	}
	return emails, nil
}

// List returns all the pronostics group stored
func (s *UserFinder) ListInfo() ([]webcore.User, error) {
	uRec := make([]UserRecord, 0)

	err := s.DB.Select(&uRec, `SELECT id, firstname, lastname, username, email FROM users ORDER BY created_at DESC`)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, webcore.NewNotFoundError(err.Error())
		}
		return nil, webcore.NewInternalServerError(err.Error())
	}

	users := make([]webcore.User, 0)
	for _, user := range uRec {
		users = append(users, webcore.User{
			ID:        user.ID,
			Firstname: user.Firstname,
			Lastname:  user.Lastname,
			Username:  user.Username,
			Email:     user.Email,
		})
	}

	return users, nil
}

// ListPhoneNumbers returns all the phone_number stored
func (s *UserFinder) ListPhoneNumbers() ([]string, error) {
	numbers := make([]string, 0)

	err := s.DB.Select(&numbers, "SELECT phone_number FROM users")
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, webcore.NewNotFoundError(err.Error())
		}
		return nil, webcore.NewInternalServerError(err.Error())
	}
	return numbers, nil
}
