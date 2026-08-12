package logger

import "go.uber.org/zap/zapcore"

var allowedFieldKeys = map[string]struct{}{
	"category":   {},
	"current":    {},
	"decoder":    {},
	"duration":   {},
	"from":       {},
	"hint":       {},
	"host":       {},
	"latency":    {},
	"latest":     {},
	"method":     {},
	"note":       {},
	"outcome":    {},
	"path":       {},
	"port":       {},
	"route":      {},
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
		if _, ok := allowedFieldKeys[field.Key]; ok {
			filtered = append(filtered, field)
		}
	}
	return filtered
}
