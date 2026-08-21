package logger

import "go.uber.org/zap/zapcore"

var allowedFieldKeys = map[string]struct{}{
	"category":   {},
	"current":    {},
	"decoder":    {},
	"duration":   {},
	"from":       {},
	"hint":       {},
	"latency":    {},
	"latest":     {},
	"method":     {},
	"note":       {},
	"outcome":    {},
	"port":       {},
	"status":     {},
	"threshold":  {},
	"to":         {},
	"type":       {},
	"error_type": {},
}

type filteringCore struct {
	zapcore.Core
}

func newFilteringCore(core zapcore.Core) zapcore.Core {
	return &filteringCore{Core: core}
}

// WrapCore applies the fail-closed field allowlist to an existing zap core.
func WrapCore(core zapcore.Core) zapcore.Core {
	return newFilteringCore(core)
}

func (c *filteringCore) With(fields []zapcore.Field) zapcore.Core {
	return &filteringCore{Core: c.Core.With(filterFields(fields))}
}

func (c *filteringCore) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if !c.Enabled(ent.Level) {
		return ce
	}
	return ce.AddCore(ent, c)
}

func (c *filteringCore) Write(ent zapcore.Entry, fields []zapcore.Field) error {
	return c.Core.Write(ent, filterFields(fields))
}

func filterFields(fields []zapcore.Field) []zapcore.Field {
	if len(fields) == 0 {
		return fields
	}
	filtered := make([]zapcore.Field, 0, len(fields))
	for _, field := range fields {
		if allowedField(field) {
			filtered = append(filtered, field)
		}
	}
	return filtered
}

func allowedField(field zapcore.Field) bool {
	if _, ok := allowedFieldKeys[field.Key]; !ok {
		return false
	}
	switch field.Type { //nolint:exhaustive // fail-closed allowlist accepts only primitive field types
	case zapcore.BoolType,
		zapcore.Float32Type,
		zapcore.Float64Type,
		zapcore.Int16Type,
		zapcore.Int32Type,
		zapcore.Int64Type,
		zapcore.Int8Type,
		zapcore.StringType,
		zapcore.Uint16Type,
		zapcore.Uint32Type,
		zapcore.Uint64Type,
		zapcore.Uint8Type,
		zapcore.DurationType,
		zapcore.TimeType:
		return true
	default:
		return false
	}
}
