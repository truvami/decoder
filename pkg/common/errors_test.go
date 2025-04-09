package common

import (
	"errors"
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
