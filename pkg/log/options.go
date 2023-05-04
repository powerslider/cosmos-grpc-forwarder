package log

import (
	"io"

	"go.uber.org/zap/zapcore"
)

type options struct {
	Level       Level
	Format      Format
	Encoder     zapcore.Encoder
	Development bool
	Output      io.Writer
	LogToStdout bool
	AddCaller   bool
	CallerSkip  int
}

func (o options) Clone() options {
	c := options{
		Level:       o.Level,
		Format:      o.Format,
		Development: o.Development,
		Output:      o.Output,
		LogToStdout: o.LogToStdout,
		AddCaller:   o.AddCaller,
		CallerSkip:  o.CallerSkip,
	}

	if o.Encoder != nil {
		c.Encoder = o.Encoder.Clone()
	}

	return c
}

// ZapLevelEnabled checks if a zap log level is set directly or indirectly.
func (o options) ZapLevelEnabled(lvl zapcore.Level) bool {
	return o.Development || o.Level.Enabled(fromZapLevel(lvl))
}

// Option represents logger configuration options.
type Option interface {
	apply(*options)
}

type optionFunc func(*options)

func (f optionFunc) apply(log *options) {
	f(log)
}

// WithLevel sets a root log level.
func WithLevel(lvl Level) Option {
	return optionFunc(func(l *options) {
		l.Level = lvl
	})
}

// WithFormat sets a log format.
func WithFormat(format Format) Option {
	return optionFunc(func(l *options) {
		l.Format = format
	})
}

// Development turn of dev mode.
func Development() Option {
	return WithDevelopment(true)
}

// WithDevelopment sets dev mode on or off.
func WithDevelopment(development bool) Option {
	return optionFunc(func(l *options) {
		l.Development = development
	})
}

// WithEncoder sets a zap log statement encoder.
func WithEncoder(encoder zapcore.Encoder) Option {
	return optionFunc(func(l *options) {
		l.Encoder = encoder
	})
}

// WithOutput sets log output.
func WithOutput(output io.Writer) Option {
	return optionFunc(func(l *options) {
		l.Output = output
	})
}

// ToStdout sets log output to stdout.
func ToStdout() Option {
	return WithLogToStdout(true)
}

// WithLogToStdout sets log output to stdout to on or off.
func WithLogToStdout(logToStdout bool) Option {
	return optionFunc(func(l *options) {
		l.LogToStdout = logToStdout
	})
}

// WithCaller configures the Logger to annotate each message with the filename,
// line number, and function name of zap's caller, or not, depending on the
// value of enabled. This is a generalized form of AddCaller.
func WithCaller(caller bool) Option {
	return optionFunc(func(l *options) {
		l.AddCaller = caller
	})
}

// AddCaller includes the caller to the log statements as a field.
func AddCaller() Option {
	return WithCaller(true)
}

// AddCallerSkip increases the number of callers skipped by caller annotation
// (as enabled by the AddCaller option).
func AddCallerSkip(skip int) Option {
	return optionFunc(func(l *options) {
		l.CallerSkip += skip
	})
}
