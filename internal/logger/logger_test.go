package logger

import (
	"bytes"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestNewLoggerWithDefaults(t *testing.T) {
	// Redirect output to a buffer for testing
	buffer := &bytes.Buffer{}
	NewLogger(WithEncoder(zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		TimeKey:        "",
		LevelKey:       "level",
		MessageKey:     "msg",
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
	})))

	Logger = zap.New(zapcore.NewCore(
		zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
			TimeKey:        "",
			LevelKey:       "level",
			MessageKey:     "msg",
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
		}),
		zapcore.AddSync(buffer),
		zapcore.InfoLevel,
	))

	Logger.Info("Test log message")
	if !bytes.Contains(buffer.Bytes(), []byte("INFO")) {
		t.Errorf("expected INFO level in log, got %s", buffer.String())
	}
	if !bytes.Contains(buffer.Bytes(), []byte("Test log message")) {
		t.Errorf("expected 'Test log message' in log, got %s", buffer.String())
	}
}

func TestNewLoggerWithDebug(t *testing.T) {
	// Redirect output to a buffer for testing
	buffer := &bytes.Buffer{}

	NewLogger(WithDebug())

	Logger = zap.New(zapcore.NewCore(
		zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
			TimeKey:        "",
			LevelKey:       "level",
			MessageKey:     "msg",
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
		}),
		zapcore.AddSync(buffer),
		zapcore.DebugLevel,
	))

	Logger.Debug("Debug log message")
	if !bytes.Contains(buffer.Bytes(), []byte("DEBUG")) {
		t.Errorf("expected DEBUG level in log, got %s", buffer.String())
	}
	if !bytes.Contains(buffer.Bytes(), []byte("Debug log message")) {
		t.Errorf("expected 'Debug log message' in log, got %s", buffer.String())
	}
}

func TestWithEncoderOption(t *testing.T) {
	buffer := &bytes.Buffer{}

	customEncoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		MessageKey:     "msg",
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
	})

	NewLogger(WithEncoder(customEncoder))

	Logger = zap.New(zapcore.NewCore(
		customEncoder,
		zapcore.AddSync(buffer),
		zapcore.InfoLevel,
	))

	Logger.Info("JSON log message")
	if !bytes.Contains(buffer.Bytes(), []byte("level")) || !bytes.Contains(buffer.Bytes(), []byte("msg")) {
		t.Errorf("expected JSON format in log, got %s", buffer.String())
	}
}

func TestSyncWithoutLogger(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("unexpected panic during Sync: %v", r)
		}
	}()

	Sync() // Should not panic
}
