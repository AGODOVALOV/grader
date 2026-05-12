package handler

import (
	"errors"
	"net/http"
	"net/http/httputil"

	"github.com/AGODOVALOV/grader/pkg/logger"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

//23505 unique violation
//23503 foreign key violation
//23502 not null violation
//22P02 invalid input syntax
//42703 undefined column
//42P01 undefined table
//08006 connection failure
//53300 too many connections

var (
	ErrNotFound                = errors.New("not found")
	ErrDuplicate               = errors.New("duplicate value")
	ErrInvalidData             = errors.New("invalid data")
	ErrDatabaseError           = errors.New("database error")
	ErrDatabaseConnectionError = errors.New("database connection error")
)

func mapDBError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}

	if _, ok := errors.AsType[*pgconn.ConnectError](err); ok {
		return ErrDatabaseConnectionError
	}

	if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
		switch pgErr.Code {
		case "23505":
			return ErrDuplicate
		case "23503", "23502", "22P02":
			return ErrInvalidData
		default:
			return ErrDatabaseError
		}
	}

	return err
}

func logErrorRequestWithDump(r *http.Request, err error) {
	dumpReq, _ := httputil.DumpRequest(r, false)
	logger.Z(r.Context()).Error(r.Context(), "service error", err.Error(),
		map[string]string{
			"dump request": string(dumpReq),
		})
}

func writeHTTPError(r *http.Request, w http.ResponseWriter, err error) {
	logErrorRequestWithDump(r, err)

	err = mapDBError(err)

	switch {
	case errors.Is(err, ErrNotFound):
		http.Error(w, "not found", http.StatusNotFound)
	case errors.Is(err, ErrDuplicate), errors.Is(err, ErrInvalidData):
		http.Error(w, "bad request", http.StatusBadRequest)
	case errors.Is(err, ErrDatabaseError):
		http.Error(w, "database error", http.StatusInternalServerError)
	case errors.Is(err, ErrDatabaseConnectionError):
		http.Error(w, "database unavailable", http.StatusServiceUnavailable)
	default:
		http.Error(w, "internal error", http.StatusInternalServerError)
	}

}
