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

func (o options) ZapLevelEnabled(lvl zapcore.Level) bool {
	return o.Development || o.Level.Enabled(fromZapLevel(lvl))
}

type Option interface {
	apply(*options)
}

type optionFunc func(*options)

func (f optionFunc) apply(log *options) {
	f(log)
}

func WithLevel(lvl Level) Option {
	return optionFunc(func(l *options) {
		l.Level = lvl
	})
}

func WithFormat(format Format) Option {
	return optionFunc(func(l *options) {
		l.Format = format
	})
}

func Development() Option {
	return WithDevelopment(true)
}

func WithDevelopment(development bool) Option {
	return optionFunc(func(l *options) {
		l.Development = development
	})
}

func WithEncoder(encoder zapcore.Encoder) Option {
	return optionFunc(func(l *options) {
		l.Encoder = encoder
	})
}

func WithOutput(output io.Writer) Option {
	return optionFunc(func(l *options) {
		l.Output = output
	})
}

func LogToStdout() Option {
	return WithLogToStdout(true)
}

func WithLogToStdout(logToStdout bool) Option {
	return optionFunc(func(l *options) {
		l.LogToStdout = logToStdout
	})
}

func WithCaller(caller bool) Option {
	return optionFunc(func(l *options) {
		l.AddCaller = caller
	})
}

func AddCaller() Option {
	return WithCaller(true)
}

func AddCallerSkip(skip int) Option {
	return optionFunc(func(l *options) {
		l.CallerSkip += skip
	})
}
