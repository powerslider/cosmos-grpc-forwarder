package log

import (
	"fmt"
	"strings"

	"go.uber.org/zap/zapcore"
)

// Level represents a log level.
type Level int

const (
	// DebugLevel debug log level
	DebugLevel Level = iota - 1
	// InfoLevel info log level
	InfoLevel
	// WarnLevel warn log level
	WarnLevel
	// ErrorLevel error log level
	ErrorLevel
	// DPanicLevel dpanic log level
	DPanicLevel
	// PanicLevel panic log level
	PanicLevel
	// FatalLevel fatal log level
	FatalLevel
)

// Enabled checks if a log level is set.
func (l Level) Enabled(lvl Level) bool {
	return lvl >= l
}

// ParseLevel parses a log level from a string value.
func ParseLevel(lvl string) (Level, error) {
	switch strings.ToLower(lvl) {
	case "fatal":
		return FatalLevel, nil
	case "panic":
		return PanicLevel, nil
	case "dpanic":
		return DPanicLevel, nil
	case "error":
		return ErrorLevel, nil
	case "warn", "warning":
		return WarnLevel, nil
	case "info":
		return InfoLevel, nil
	case "debug":
		return DebugLevel, nil
	}

	return InfoLevel, fmt.Errorf("not a valid Level: %q", lvl)
}

func fromZapLevel(lvl zapcore.Level) Level {
	switch lvl {
	case zapcore.DebugLevel:
		return DebugLevel
	case zapcore.InfoLevel:
		return InfoLevel
	case zapcore.WarnLevel:
		return WarnLevel
	case zapcore.ErrorLevel:
		return ErrorLevel
	case zapcore.DPanicLevel:
		return DPanicLevel
	case zapcore.PanicLevel:
		return PanicLevel
	case zapcore.FatalLevel:
		return FatalLevel
	}

	return InfoLevel
}

//nolint:unused
func toZapLevel(lvl Level) zapcore.Level {
	switch lvl {
	case DebugLevel:
		return zapcore.DebugLevel
	case InfoLevel:
		return zapcore.InfoLevel
	case WarnLevel:
		return zapcore.WarnLevel
	case ErrorLevel:
		return zapcore.ErrorLevel
	case DPanicLevel:
		return zapcore.DPanicLevel
	case PanicLevel:
		return zapcore.PanicLevel
	case FatalLevel:
		return zapcore.FatalLevel
	}

	return zapcore.InfoLevel
}
