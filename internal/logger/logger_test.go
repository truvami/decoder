package logger

import (
	"testing"

	"go.uber.org/zap/zapcore"
)

func TestWithInfoSetsLevelToInfo(t *testing.T) {
	cfg := &loggerConfig{
		Level: zapcore.DebugLevel, // set to something else first
	}
	opt := WithInfo()
	opt(cfg)
	if cfg.Level != zapcore.InfoLevel {
		t.Errorf("WithInfo() did not set Level to InfoLevel, got %v", cfg.Level)
	}
}

func TestWithEncoderSetsEncoder(t *testing.T) {
	// Create a dummy encoder
	encoderConfig := zapcore.EncoderConfig{}
	dummyEncoder := zapcore.NewConsoleEncoder(encoderConfig)

	cfg := &loggerConfig{
		Encoder: nil,
	}
	opt := WithEncoder(dummyEncoder)
	opt(cfg)
	if cfg.Encoder != dummyEncoder {
		t.Errorf("WithEncoder() did not set Encoder correctly")
	}
}
