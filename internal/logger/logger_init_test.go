package logger

import (
	"testing"
)

func TestLoggerSync(t *testing.T) {
	// Initialize the logger
	NewLogger()
	
	// Check if logger is initialized
	if Logger == nil {
		t.Fatal("Logger should be initialized")
	}
	
	// Test logging
	Logger.Info("Test log message")
} 