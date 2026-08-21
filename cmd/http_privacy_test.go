package cmd

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/truvami/decoder/internal/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestLoggingMiddlewareEmitsOnlySafeAccessFields(t *testing.T) {
	observed, logs := observer.New(zapcore.InfoLevel)
	safeLogger := zap.New(logger.WrapCore(observed))

	handler := loggingMiddleware(safeLogger, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	req := httptest.NewRequest(http.MethodPost, "/tagsl/v1?devEui=0123456789ABCDEF&payload=deadbeef", nil)
	req.RemoteAddr = "203.0.113.10:1234"
	req.Header.Set("User-Agent", "secret-agent/1.0")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if logs.Len() != 1 {
		t.Fatalf("expected 1 access log entry, got %d", logs.Len())
	}

	entry := logs.All()[0]
	if entry.Message != "HTTP request" {
		t.Fatalf("unexpected log message %q", entry.Message)
	}

	fields := entry.ContextMap()
	for _, forbidden := range []string{"devEui", "payload", "url", "route", "requestId", "remoteAddress", "userAgent", "response"} {
		if _, ok := fields[forbidden]; ok {
			t.Fatalf("forbidden field %q present in access log", forbidden)
		}
	}

	for _, required := range []string{"method", "status", "duration"} {
		if _, ok := fields[required]; !ok {
			t.Fatalf("required field %q missing from access log", required)
		}
	}

	if fields["method"] != http.MethodPost {
		t.Fatalf("expected method POST, got %v", fields["method"])
	}
}

func TestSafeHTTPErrorLogDoesNotEmitRawErrorText(t *testing.T) {
	observed, logs := observer.New(zapcore.ErrorLevel)
	logger.Logger = zap.New(logger.WrapCore(observed))

	writer := safeHTTPErrorLog{}
	raw := []byte("listen tcp: lookup evil.example: no such host")
	if _, err := writer.Write(raw); err != nil {
		t.Fatalf("write failed: %v", err)
	}

	if logs.Len() != 1 {
		t.Fatalf("expected 1 log entry, got %d", logs.Len())
	}

	entry := logs.All()[0]
	if entry.Message != "http server error" {
		t.Fatalf("unexpected message %q", entry.Message)
	}
	if entry.ContextMap()["category"] != "internal" {
		t.Fatalf("expected category=internal, got %v", entry.ContextMap()["category"])
	}
	for _, field := range entry.Context {
		if field.Key == "error" || field.String == string(raw) {
			t.Fatal("raw server error text must not appear in logs")
		}
	}
}

func TestPrintJSONWritesToStdoutNotLogger(t *testing.T) {
	observed, logs := observer.New(zapcore.InfoLevel)
	logger.Logger = zap.New(logger.WrapCore(observed))

	originalJSON := Json
	Json = true
	defer func() { Json = originalJSON }()

	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	os.Stdout = w

	printJSON(map[string]string{"latitude": "47.0"})

	_ = w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatalf("read stdout: %v", err)
	}
	_ = r.Close()

	if logs.Len() != 0 {
		t.Fatalf("expected no zap output for JSON CLI mode, got %d entries", logs.Len())
	}
	if !bytes.Contains(buf.Bytes(), []byte(`"latitude":"47.0"`)) && !bytes.Contains(buf.Bytes(), []byte(`"latitude": "47.0"`)) {
		t.Fatalf("expected JSON on stdout, got %q", buf.String())
	}
}
