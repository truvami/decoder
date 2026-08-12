package logger

import (
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestFilterFieldsStripsSensitiveKeys(t *testing.T) {
	filtered := filterFields([]zapcore.Field{
		{Key: "devEui", Type: zapcore.StringType, String: "0123456789ABCDEF"},
		{Key: "payload", Type: zapcore.StringType, String: "deadbeef"},
		{Key: "route", Type: zapcore.StringType, String: "/tagsl/v1"},
		{Key: "method", Type: zapcore.StringType, String: "POST"},
	})

	if len(filtered) != 2 {
		t.Fatalf("expected 2 allowed fields, got %d", len(filtered))
	}
	for _, field := range filtered {
		if field.Key != "route" && field.Key != "method" {
			t.Fatalf("unexpected allowed field %q", field.Key)
		}
	}
}

func TestFilteringCoreDoesNotEmitSensitiveFields(t *testing.T) {
	observed, logs := observer.New(zapcore.InfoLevel)
	Logger = zap.New(newFilteringCore(observed))

	Logger.Info("HTTP request",
		zap.String("route", "/tagsl/v1"),
		zap.String("method", "POST"),
		zap.Int("status", 200),
		zap.String("devEui", "0123456789ABCDEF"),
		zap.String("payload", "deadbeef"),
		zap.String("url", "http://localhost/tagsl/v1?secret=1"),
	)

	if logs.Len() != 1 {
		t.Fatalf("expected 1 log entry, got %d", logs.Len())
	}

	encoded := logs.All()[0]
	for _, field := range encoded.Context {
		if _, allowed := allowedFieldKeys[field.Key]; !allowed {
			t.Fatalf("unexpected field key in log output: %q", field.Key)
		}
	}

	out := encoded.ContextMap()
	if _, ok := out["devEui"]; ok {
		t.Fatal("devEui must not appear in log output")
	}
	if _, ok := out["payload"]; ok {
		t.Fatal("payload must not appear in log output")
	}
	if _, ok := out["url"]; ok {
		t.Fatal("url must not appear in log output")
	}
}

func TestFilteringCoreBlocksRawErrorField(t *testing.T) {
	observed, logs := observer.New(zapcore.ErrorLevel)
	Logger = zap.New(newFilteringCore(observed))

	Logger.Error("decode failed", zap.Error(stringsErr("payload=deadbeef devEui=0123456789ABCDEF")))

	if logs.Len() != 1 {
		t.Fatalf("expected 1 log entry, got %d", logs.Len())
	}
	for _, field := range logs.All()[0].Context {
		if field.Key == "error" {
			t.Fatal("raw error field must not appear in log output")
		}
	}
}

type stringsErr string

func (e stringsErr) Error() string { return string(e) }
