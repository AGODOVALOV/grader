// Package session provides session management functionality.
package session

import (
	"errors"

	"github.com/google/uuid"
)

type ctxKey int

// SessionKey is the key used to store the session in the context.
const SessionKey ctxKey = 1

// ErrNoAuth is returned when no session is found.
var ErrNoAuth = errors.New("no session found")

// Session represents a user session.
type Session struct {
	ID     uuid.UUID
	UserID int64
}
