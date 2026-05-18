package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	// ErrExpiredToken is returned when the token is expired.
	ErrExpiredToken = errors.New("token has invalid claims: token is expired")
	// ErrInvalidToken is returned when the token is invalid.
	ErrInvalidToken = errors.New("token is invalid")
)

// Payload contains the payload data ot the token.
type Payload struct {
	jwt.RegisteredClaims

	ID        uuid.UUID `json:"id"`
	UserID    int64     `json:"user_id"`
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

// NewPayload creates a new token payload with a specific username and duration.
func NewPayload(userID int64, username string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	issuedAt := time.Now()
	expiredAt := time.Now().Add(duration)

	payload := &Payload{
		ID:        tokenID,
		UserID:    userID,
		Username:  username,
		IssuedAt:  issuedAt,
		ExpiredAt: expiredAt,

		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    username,
			Subject:   "",
			Audience:  nil,
			ExpiresAt: jwt.NewNumericDate(expiredAt),
			NotBefore: jwt.NewNumericDate(issuedAt),
			IssuedAt:  jwt.NewNumericDate(issuedAt),
			ID:        tokenID.String(),
		},
	}

	return payload, nil
}

// Valid checks if the token is valid.
func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}
