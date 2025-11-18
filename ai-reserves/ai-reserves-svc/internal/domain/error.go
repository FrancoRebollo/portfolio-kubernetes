package domain

import (
	"errors"
	"fmt"
)

const (
	ErrCodeDuplicateKey            = "duplicate_key"
	ErrCodeForeignKeyViolation     = "foreign_key"
	ErrCodeNotNullViolation        = "not_null_constraint"
	ErrCodeInvalidInput            = "invalid_input"
	ErrCodeDeadlockDetected        = "deadlock_detected"
	ErrCodeTableOrViewDoesNotExist = "table_view_not_exist"
	ErrCodeConnectionTimeout       = "connection_timeout"
	ErrCodeEndOfCommunication      = "end_of_communication"
	ErrCodeInternalServer          = "internal_server"
	ErrCodeRouteNotFound           = "not_found"
	ErrCodeRequestTimeout          = "request_cancelled"
	ErrCodeMediaType               = "Tipo de contenido no soportado"
	ErrCodeForbidden               = "Prohibido"
	ErrCodeUnauthorized            = "No autorizado"
	ErrCodeServiceUnavailable      = "Servicio no disponible"
)

var (
	ErrDuplicateKey            = errors.New("duplicate key violation")
	ErrForeignKeyViolation     = errors.New("foreign key violation")
	ErrNotNullViolation        = errors.New("not null constraint violation")
	ErrInvalidInput            = errors.New("invalid input syntax")
	ErrDeadlockDetected        = errors.New("deadlock detected")
	ErrTableOrViewDoesNotExist = errors.New("table or view does not exist")
	ErrConnectionTimeout       = errors.New("connection timeout")
	ErrEndOfCommunication      = errors.New("end of communication channel")
	ErrInternalServer          = errors.New("internal server error")
)

type HealthcheckError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *HealthcheckError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

var ErrDuplicateEvent = errors.New("duplicate event ignored")
