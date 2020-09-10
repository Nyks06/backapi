package backapi

import (
	"context"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Session struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

type SessionStore interface {
	Create(*Session) (*Session, error)
	Delete(ID string) error
	DeleteAllForUserID(userID string) error
}

type SessionFinder interface {
	ByID(ID string) (*Session, error)
}

type SessionService struct {
	StoreManager *StoreManager
}

func (s *SessionService) Create(ctx context.Context, Email string, Password string) (*Session, error) {
	user, err := s.StoreManager.UserFinder.ByEmail(Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, NewNotFoundError("user_not_found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(Password))
	if err != nil {
		return nil, NewUnauthorizedError("password_mismatch")
	}

	session := &Session{
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}
	sessionCreated, err := s.StoreManager.SessionStore.Create(session)
	if err != nil {
		return nil, err
	}
	return sessionCreated, nil
}

func (s *SessionService) GetByID(ctx context.Context, ID string) (*Session, error) {
	session, err := s.StoreManager.SessionFinder.ByID(ID)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (s *SessionService) Delete(ctx context.Context, ID string) error {
	err := s.StoreManager.SessionStore.Delete(ID)
	return err
}
