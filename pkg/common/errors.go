package common

import (
	"fmt"

	"github.com/jackc/pgx/v5"
)

var ErrRecordNotFound = pgx.ErrNoRows

var ErrIncorrectPassword = fmt.Errorf("incorrect password")
