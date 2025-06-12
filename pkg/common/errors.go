package common

import (
	"errors"
	"fmt"
)

var (
	ErrPortNotSupported = errors.New("port not supported")

	ErrInvalidPayloadLength = errors.New("invalid payload length")
	ErrPayloadTooShort      = errors.New("payload too short")
	ErrPayloadTooLong       = errors.New("payload too long")

	ErrValidationFailed = errors.New("validation failed")
)

// WrapError wraps two errors into a single error, combining the parent and child errors.
// It uses the %w verb to allow the wrapped errors to be unwrapped later.
//
// Parameters:
//   - parent: The outer error that provides context.
//   - child: The inner error that provides additional details.
//
// Returns:
//
//	A new error that combines the parent and child errors.
func WrapError(parent error, child error) error {
	return fmt.Errorf("%w: %w", parent, child)
}

// WrapErrorWithMessage wraps two errors (parent and child) with an additional
// custom message, creating a new error. The resulting error includes the
// provided message followed by the parent and child errors in a nested format.
//
// Parameters:
//   - parent: The outer error to be wrapped.
//   - child: The inner error to be wrapped.
//   - message: A custom message to provide additional context.
//
// Returns:
//   - An error that combines the custom message, parent error, and child error
//     in a formatted string.
func WrapErrorWithMessage(parent error, child error, message string) error {
	return fmt.Errorf("%s: %w: %w", message, parent, child)
}

func UnwrapError(err error) []error {
	var errs []error = []error{}
	if err, ok := err.(interface{ Unwrap() []error }); ok {
		errs = append(errs, err.Unwrap()...)
	}
	return errs
}
