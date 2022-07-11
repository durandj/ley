package errortypes

import "fmt"

// UserError is an error created by the user.
type UserError struct {
	SafeMessage  string
	WrappedError error
}

func (err UserError) Error() string {
	return err.SafeMessage
}

// Unwrap returns the inner error that this error wraps.
func (err UserError) Unwrap() error {
	return err.WrappedError
}

var _ error = (*UserError)(nil)

// ValidationError is an error type created by the user submitting
// invalid data.
type ValidationError struct {
	UserError
}

// NewValidationError creates a validation error instance.
func NewValidationError(message string, values ...any) ValidationError {
	return NewWrappedValidationError(nil, message, values...)
}

// NewWrappedValidationError creates a validation error that wraps another
// error.
func NewWrappedValidationError(err error, message string, values ...any) ValidationError {
	// TODO: err should only be a safe error
	return ValidationError{
		UserError: UserError{
			SafeMessage:  fmt.Sprintf(message, values...),
			WrappedError: err,
		},
	}
}

// NotFoundError is returned when requested data could not be found.
type NotFoundError struct {
	UserError
}

var _ error = (*NotFoundError)(nil)
