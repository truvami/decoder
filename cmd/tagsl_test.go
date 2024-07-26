package cmd

import (
	"bytes"
	"log"
	"strings"
	"testing"
)

func TestPrintJSON(t *testing.T) {
	// Test input
	data := map[string]interface{}{
		"key1": "value1",
		"key2": 2,
		"key3": true,
	}

	// Expected output
	expectedOutput := `{"key1":"value1","key2":2,"key3":true}`

	// Capture the output of the printJSON function
	var capturedOutput bytes.Buffer
	log.SetOutput(&capturedOutput)

	// Call the function
	printJSON(data)

	// Check if the output matches the expected output
	if !strings.Contains(capturedOutput.String(), expectedOutput) {
		t.Errorf("printJSON output = %s, expected %s", capturedOutput.String(), expectedOutput)
	}
}
