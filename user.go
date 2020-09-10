package backapi

import (
	"context"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User defines all the fields that represent a user
type User struct {
	// IDs
	ID string `json:"id"`
	// CustomerID is related to Stripe
	CustomerID string `json:"customer_id"`

	// Login information
	Email    string `json:"email"`
	Password string `json:"-"`

	// Personal information
	Firstname   string `json:"firstname"`
	Lastname    string `json:"lastname"`
	Username    string `json:"username"`
	PhoneNumber string `json:"phone_number"`

	// Status
	Admin bool `json:"admin"`

	// Steps date
	CreatedAt   time.Time `json:"created_at"`
	ConfirmedAt time.Time `json:"confirmed_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   time.Time `json:"deleted_at"`
}

// UserStore defines all the methods usable to write data in the store for users
type UserStore interface {
	Create(*User) (*User, error)
	Confirm(ID string) error
	Delete(ID string) error
	SetPassword(ID string, passwd string) error
	SetEmail(ID string, email string) error
}

// UserFinder defines all the methods usable to read data in the store for users
type UserFinder interface {
	ByEmail(email string) (*User, error)
	ByID(id string) (*User, error)
	List() ([]User, error)
}

type UserService struct {
	StoreManager *StoreManager
}

func (s *UserService) CreateUser(ctx context.Context, user *User) (*User, error) {
	u, err := s.StoreManager.UserFinder.ByEmail(user.Email)
	if u != nil {
		return nil, NewResourceAlreadyCreatedError("The user already exist")
	} else if err != nil {
		if !IsNotFoundError(err) {
			return nil, err
		}
	}

	password, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, NewInternalServerError(err.Error())
	}
	user.Password = string(password)

	// Ensure that the user is created not being an admin
	user.Admin = false

	createdUser, err := s.StoreManager.UserStore.Create(user)
	if err != nil {
		return nil, NewInternalServerError(err.Error())
	}

	return createdUser, nil
}

func (s *UserService) ChangePassword(ctx context.Context, user *User, currentPassword string, newPassword string) (*User, error) {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(currentPassword))
	if err != nil {
		return nil, NewUnauthorizedError("password_mismatch")
	}

	newPasswordHashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, NewInternalServerError(err.Error())
	}

	if err = s.StoreManager.UserStore.SetPassword(user.ID, string(newPasswordHashed)); err != nil {
		return nil, err
	}

	return s.StoreManager.UserFinder.ByID(user.ID)
}

func (s *UserService) ChangeEmail(ctx context.Context, user *User, currentPassword string, newEmail string) (*User, error) {
	u, err := s.StoreManager.UserFinder.ByEmail(newEmail)
	if u != nil {
		return nil, NewResourceAlreadyCreatedError("The user already exist")
	} else if err != nil {
		if !IsNotFoundError(err) {
			return nil, err
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(currentPassword))
	if err != nil {
		return nil, NewUnauthorizedError("password_mismatch")
	}

	if err = s.StoreManager.UserStore.SetEmail(user.ID, newEmail); err != nil {
		return nil, err
	}

	return s.StoreManager.UserFinder.ByID(user.ID)
}

func (s *UserService) GetByID(ctx context.Context, ID string) (*User, error) {
	u, err := s.StoreManager.UserFinder.ByID(ID)
	if err != nil {
		return nil, err
	}

	return u, nil
}
