package repository

import "errors"

var (
	ErrDuplicateKey           = errors.New("duplicate key violation")
	ErrForeignKeyViolation    = errors.New("foreign key violation")
	ErrNotNullViolation       = errors.New("not null constraint violation")
	ErrInvalidInput           = errors.New("invalid input syntax")
	ErrConnectionDoesNotExist = errors.New("connection does not exist")
	ErrServerShutdown         = errors.New("server is shutting down")
	ErrSyntaxError            = errors.New("syntax error")
	ErrDeadlockDetected       = errors.New("deadlock detected")
	ErrInternalServer         = errors.New("internal server error")
)
