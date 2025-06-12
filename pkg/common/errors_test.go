package common

import (
	"errors"
	"fmt"
	"testing"
)

func TestWrapError(t *testing.T) {
	parentErr := errors.New("parent error")
	childErr := errors.New("child error")

	wrappedErr := WrapError(parentErr, childErr)

	// Check if the wrapped error is not nil
	if wrappedErr == nil {
		t.Fatalf("expected wrapped error to be non-nil")
	}

	// Check if the wrapped error contains the parent error
	if !errors.Is(wrappedErr, parentErr) {
		t.Errorf("expected wrapped error to contain parent error, but it does not")
	}

	// Check if the wrapped error contains the child error
	if !errors.Is(wrappedErr, childErr) {
		t.Errorf("expected wrapped error to contain child error, but it does not")
	}

	// Check the error message format
	expectedMessage := "parent error: child error"
	if wrappedErr.Error() != expectedMessage {
		t.Errorf("expected error message '%s', got '%s'", expectedMessage, wrappedErr.Error())
	}
}

func TestWrapErrorWithMessage(t *testing.T) {
	parentErr := errors.New("parent error")
	childErr := errors.New("child error")
	message := "custom message"

	wrappedErr := WrapErrorWithMessage(parentErr, childErr, message)

	// Check if the wrapped error is not nil
	if wrappedErr == nil {
		t.Fatalf("expected wrapped error to be non-nil")
	}

	// Check if the wrapped error contains the parent error
	if !errors.Is(wrappedErr, parentErr) {
		t.Errorf("expected wrapped error to contain parent error, but it does not")
	}

	// Check if the wrapped error contains the child error
	if !errors.Is(wrappedErr, childErr) {
		t.Errorf("expected wrapped error to contain child error, but it does not")
	}

	// Check the error message format
	expectedMessage := "custom message: parent error: child error"
	if wrappedErr.Error() != expectedMessage {
		t.Errorf("expected error message '%s', got '%s'", expectedMessage, wrappedErr.Error())
	}
}

func TestUnwrapError_Nil(t *testing.T) {
	var err error = nil
	result := UnwrapError(err)
	if len(result) != 0 {
		t.Fatalf("expected empty slice, got %v", result)
	}
}

type multiError struct {
	errs []error
}

func (m multiError) Error() string {
	return "multi error"
}

func (m multiError) Unwrap() []error {
	return m.errs
}

func TestUnwrapError_WithUnwrap(t *testing.T) {
	err1 := fmt.Errorf("error 1")
	err2 := fmt.Errorf("error 2")
	merr := multiError{errs: []error{err1, err2}}
	result := UnwrapError(merr)
	if len(result) != 2 || result[0] != err1 || result[1] != err2 {
		t.Fatalf("expected [%v %v], got %v", err1, err2, result)
	}
}

func TestUnwrapError_NoUnwrap(t *testing.T) {
	err := fmt.Errorf("plain error")
	result := UnwrapError(err)
	if len(result) != 0 {
		t.Fatalf("expected empty slice, got %v", result)
	}
}
