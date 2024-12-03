package logger

import (
	"log/slog"
	"os"
	"regexp"
	"strings"
	"testing"
)

type captureStream struct {
	lines [][]byte
}

func (cs *captureStream) Write(bytes []byte) (int, error) {
	cs.lines = append(cs.lines, bytes)
	return len(bytes), nil
}

func TestWritesToProvidedStream(t *testing.T) {
	cs := &captureStream{}
	handler := New(nil, WithDestinationWriter(cs), WithOutputEmptyAttrs())
	logger := slog.New(handler)

	logger.Info("testing logger")
	if len(cs.lines) != 1 {
		t.Errorf("expected 1 lines logged, got: %d", len(cs.lines))
	}

	lineMatcher := regexp.MustCompile(`\[\d{2}:\d{2}:\d{2}\.\d{3}\] INFO: testing logger {}`)
	line := string(cs.lines[0])
	if lineMatcher.MatchString(line) == false {
		t.Errorf("expected `testing logger` but found `%s`", line)
	}
	if !strings.HasSuffix(line, "\n") {
		t.Errorf("expected line to be terminated with `\\n` but found `%s`", line[len(line)-1:])
	}
}

func TestSkipEmptyAttributes(t *testing.T) {
	cs := &captureStream{}
	handler := New(nil, WithDestinationWriter(cs))
	logger := slog.New(handler)

	logger.Info("testing logger")
	if len(cs.lines) != 1 {
		t.Errorf("expected 1 lines logged, got: %d", len(cs.lines))
	}

	lineMatcher := regexp.MustCompile(`\[\d{2}:\d{2}:\d{2}\.\d{3}\] INFO: testing logger`)
	line := string(cs.lines[0])
	if lineMatcher.MatchString(line) == false {
		t.Errorf("expected `testing logger` but found `%s`", line)
	}
	if !strings.HasSuffix(line, "\n") {
		t.Errorf("expected line to be terminated with `\\n` but found `%s`", line[len(line)-1:])
	}
}

func TestColorizer(t *testing.T) {
	// Test case 1: color code 31 (red)
	expected1 := "\033[31mHello\033[0m"
	result1 := colorizer(31, "Hello")
	if result1 != expected1 {
		t.Errorf("Expected: %s, but got: %s", expected1, result1)
	}

	// Test case 2: color code 92 (light green)
	expected2 := "\033[92mWorld\033[0m"
	result2 := colorizer(92, "World")
	if result2 != expected2 {
		t.Errorf("Expected: %s, but got: %s", expected2, result2)
	}

	// Test case 3: color code 97 (white)
	expected3 := "\033[97m123\033[0m"
	result3 := colorizer(97, "123")
	if result3 != expected3 {
		t.Errorf("Expected: %s, but got: %s", expected3, result3)
	}
}

func TestNewHandler(t *testing.T) {
	// Test case 1: Verify that the handler is created with the correct writer
	writer := os.Stdout
	handler := NewHandler(&slog.HandlerOptions{})
	if handler.writer != writer {
		t.Errorf("Expected writer to be set to os.Stdout, but got: %v", &handler.writer)
	}

	// Test case 2: Verify that the handler is created with colorization enabled
	handler = NewHandler(&slog.HandlerOptions{})
	if !handler.colorize {
		t.Error("Expected colorization to be enabled, but it is not")
	}

	// Test case 3: Verify that the handler is created with output empty attributes enabled
	handler = NewHandler(&slog.HandlerOptions{})
	if !handler.outputEmptyAttrs {
		t.Error("Expected output empty attributes to be enabled, but it is not")
	}
}
