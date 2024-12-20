package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

type Option func(*loggerConfig)

type loggerConfig struct {
	Level   zapcore.Level
	Encoder zapcore.Encoder
}

func NewLogger(options ...Option) {
	// create a custom encoder
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "", // disable stack traces
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
	}

	// default options
	config := &loggerConfig{
		Level:   zapcore.InfoLevel,
		Encoder: zapcore.NewConsoleEncoder(encoderConfig),
	}

	for _, opt := range options {
		opt(config)
	}

	core := zapcore.NewCore(
		config.Encoder,
		zapcore.AddSync(os.Stdout),
		config.Level,
	)

	Logger = zap.New(core)
}

func WithDebug() Option {
	return func(c *loggerConfig) {
		c.Level = zapcore.DebugLevel
	}
}

func WithInfo() Option {
	return func(c *loggerConfig) {
		c.Level = zapcore.InfoLevel
	}
}

func WithEncoder(encoder zapcore.Encoder) Option {
	return func(c *loggerConfig) {
		c.Encoder = encoder
	}
}

func Sync() {
	if Logger != nil {
		_ = Logger.Sync()
	}
}
