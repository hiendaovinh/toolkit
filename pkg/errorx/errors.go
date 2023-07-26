package errorx

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/ory/ladon"
)

type Kind uint8

const (
	Other Kind = iota
	Invalid
	Validation
	NotExist
	Exist
	RateLimiting
	Authn
	Authz
	Captcha
	Database
	Service
)

func (k Kind) String() string {
	switch k {
	case Invalid:
		return "invalid-request"
	case Validation:
		return "validation"
	case NotExist:
		return "resource-not-found"
	case Exist:
		return "resource-already-exists"
	case RateLimiting:
		return "rate-limiting"
	case Authn:
		return "authentication"
	case Authz:
		return "authorization"
	case Captcha:
		return "captcha"
	case Database:
		return "database-query"
	case Service:
		return "internal-service-failure"
	}

	return "unknown"
}

func (k Kind) HTTPStatus() int {
	switch k {
	case Invalid:
		return http.StatusBadRequest
	case Validation, Captcha:
		return http.StatusUnprocessableEntity
	case NotExist:
		return http.StatusNotFound
	case Exist:
		return http.StatusConflict
	case RateLimiting:
		return http.StatusTooManyRequests
	case Authn:
		return http.StatusUnauthorized
	case Authz:
		return http.StatusForbidden
	case Database:
		return http.StatusInternalServerError
	}

	return http.StatusInternalServerError
}

type Error struct {
	err     error
	kind    Kind
	message string
}

func (e *Error) Error() string {
	return e.message
}

// Unwrap offers the ability to check the cause using errors.Is
func (e *Error) Unwrap() error {
	return e.err
}

func (e *Error) Of(k Kind) bool {
	return e.kind == k
}

func (e *Error) Code() string {
	return e.kind.String()
}

func (e *Error) Status() int {
	return e.kind.HTTPStatus()
}

func Wrap(err error, kind Kind) *Error {
	if IsNoRows(err) {
		return &Error{err: err, kind: NotExist, message: "not found"}
	}

	if IsForbidden(err) {
		return &Error{err: err, kind: Authz, message: "unauthorized"}
	}

	if v, ok := IsDuplicated(err); ok {
		return &Error{err: err, kind: Exist, message: v.Detail}
	}

	if _, ok := err.(*pgconn.PgError); ok {
		return &Error{err: err, kind: Database, message: err.Error()}
	}

	return &Error{err: err, kind: kind, message: err.Error()}
}

func IsNoRows(err error) bool {
	if err == sql.ErrNoRows || err == pgx.ErrNoRows {
		return true
	}

	return false
}

func IsForbidden(err error) bool {
	return errors.Is(err, ladon.ErrRequestDenied) || errors.Is(err, ladon.ErrRequestForcefullyDenied)
}

func IsDuplicated(err error) (*pgconn.PgError, bool) {
	if err == nil {
		return nil, false
	}

	if v, ok := err.(*pgconn.PgError); ok {
		return v, v.SQLState() == "23505"
	}

	return nil, false
}

func MaskErrorMessage(err error) string {
	var target *Error

	if !errors.As(err, &target) {
		return "unexpected error occurred"
	}

	if target.Of(Database) || target.Of(Service) {
		return "unable to process"
	}

	return target.Error()
}
