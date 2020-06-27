package webcore

import "time"

type Session struct {
	ID        string
	UserID    string
	CreatedAt time.Time
	UpdatedAt time.Time
	ExpiresAt time.Time
}

type SessionStore interface {
	Create(*Session) (*Session, error)
	Delete(ID string) error
	DeleteAllForUserID(userID string) error
}

type SessionFinder interface {
	ByID(ID string) (*Session, error)
}
