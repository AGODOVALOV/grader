// Package common contains common error types.
package common //nolint:revive // package name is ok

import (
	"errors"

	"github.com/jackc/pgx/v5"
)

// ErrRecordNotFound is returned when a record is not found.
var ErrRecordNotFound = pgx.ErrNoRows

// ErrIncorrectPassword is returned when the password is incorrect.
var ErrIncorrectPassword = errors.New("incorrect password")
