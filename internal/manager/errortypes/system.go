package errortypes

// SystemError is any error created internally by the system that
// should NOT be leaked to end users.
type SystemError struct {
	SafeMessage   string
	UnsafeMessage string
	WrappedError  error
}

func (err SystemError) Error() string {
	return err.SafeMessage
}

// Unwrap attempts to pull the error that this SystemError out.
func (err SystemError) Unwrap() error {
	return err.WrappedError
}

var _ error = (*SystemError)(nil)
