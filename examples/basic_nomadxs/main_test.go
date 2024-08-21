package main

import (
	"bytes"
	"log"
	"strings"
	"testing"
)

func TestMain(t *testing.T) {
	// Create a buffer to capture the output
	var buf bytes.Buffer
	log.SetOutput(&buf)

	// Run the main function
	main()

	// Check if the expected output is present in the buffer
	expectedOutput := `47.041811`
	if !strings.Contains(buf.String(), expectedOutput) {
		t.Errorf("expected output %q not found", expectedOutput)
	}
}
