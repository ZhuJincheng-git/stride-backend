package apperror

import (
	"errors"
	"fmt"
)

type Code string

const (
	CodeInvalidArgument  Code = "invalid_argument"
	CodeUnauthenticated  Code = "unauthenticated"
	CodePermissionDenied Code = "permission_denied"
	CodeNotFound         Code = "not_found"
	CodeInternal         Code = "internal"
)

type Error struct {
	Code Code // machine-readable
	Message string // human-readable, show to API clients
	cause error // optional underlying error (NOT exposed to clients)
}

func New(code Code, message string) *Error {
	return &Error{Code: code, Message: message}
}

func Wrap(code Code, message string, cause error) *Error {
	return &Error{Code: code, Message: message, cause: cause}
}

// Error is used for logging.
func (e *Error) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("%s: %s: %v", e.Code, e.Message, e.cause)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *Error) Unwrap() error { return e.cause }

func AsAppError(err error) (*Error, bool) {
	var ae *Error
	if errors.As(err, &ae) {
		return ae, true
	}
	return  nil, false
}