package session

import (
	"errors"

	"github.com/google/uuid"
)

type ctxKey int

const SessionKey ctxKey = 1

var (
	ErrNoAuth = errors.New("No session found")
)

type Session struct {
	ID     uuid.UUID
	UserID int64
}
