package log

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger represents all methods for logging at different levels.
type Logger interface {
	Print(msg string, args ...Field)
	Debug(msg string, args ...Field)
	Info(msg string, args ...Field)
	Warn(msg string, args ...Field)
	Error(msg string, args ...Field)
	DPanic(msg string, args ...Field)
	Panic(msg string, args ...Field)
	Fatal(msg string, args ...Field)
}

// Field is a type alias for zap.Field.
type Field = zap.Field

// Re-mapped zap methods for statically typed log fields.
var (
	Skip        = zap.Skip
	Binary      = zap.Binary
	Bool        = zap.Bool
	Boolp       = zap.Boolp
	ByteString  = zap.ByteString
	Complex128  = zap.Complex128
	Complex128p = zap.Complex128p
	Complex64   = zap.Complex64
	Complex64p  = zap.Complex64p
	Error       = zap.Error
	Float64     = zap.Float64
	Float64p    = zap.Float64p
	Float32     = zap.Float32
	Float32p    = zap.Float32p
	Int         = zap.Int
	Intp        = zap.Intp
	Int64       = zap.Int64
	Int64p      = zap.Int64p
	Int32       = zap.Int32
	Int32p      = zap.Int32p
	Int16       = zap.Int16
	Int16p      = zap.Int16p
	Int8        = zap.Int8
	Int8p       = zap.Int8p
	String      = zap.String
	Stringp     = zap.Stringp
	Uint        = zap.Uint
	Uintp       = zap.Uintp
	Uint64      = zap.Uint64
	Uint64p     = zap.Uint64p
	Uint32      = zap.Uint32
	Uint32p     = zap.Uint32p
	Uint16      = zap.Uint16
	Uint16p     = zap.Uint16p
	Uint8       = zap.Uint8
	Uint8p      = zap.Uint8p
	Uintptr     = zap.Uintptr
	Uintptrp    = zap.Uintptrp
	Reflect     = zap.Reflect
	Namespace   = zap.Namespace
	Stringer    = zap.Stringer
	Time        = zap.Time
	Timep       = zap.Timep
	Stack       = zap.Stack
	StackSkip   = zap.StackSkip
	Duration    = zap.Duration
	Durationp   = zap.Durationp
	Any         = zap.Any
)

type logFunc func(logger *zap.Logger, msg string, fields ...Field)

// StructuredLogger is an implementation of the Logger inteface for structured logging.
type StructuredLogger struct {
	base    *zap.Logger
	options options

	print  logFunc
	debug  logFunc
	info   logFunc
	warn   logFunc
	error  logFunc
	dpanic logFunc
	panic  logFunc
	fatal  logFunc
}

var _ Logger = (*StructuredLogger)(nil)

var _defaultOptions = options{
	Development: false,
	Format:      FormatJSON,
	Level:       InfoLevel,
	LogToStdout: true,
	AddCaller:   false,
	CallerSkip:  1,
}

// New constructs a new StructuredLogger.
func New(opt ...Option) *StructuredLogger {
	opts := _defaultOptions

	for _, o := range opt {
		o.apply(&opts)
	}

	return newLogger(opts)
}

func newLogger(opts options) *StructuredLogger {
	var encoder zapcore.Encoder

	if opts.Encoder != nil {
		encoder = opts.Encoder
	} else {
		var encoderCfg zapcore.EncoderConfig

		if opts.Development {
			encoderCfg = zap.NewDevelopmentEncoderConfig()
			encoderCfg.EncodeTime = zapcore.RFC3339TimeEncoder
			encoderCfg.EncodeCaller = zapcore.FullCallerEncoder
		} else {
			encoderCfg = zap.NewProductionEncoderConfig()
			encoderCfg.TimeKey = "time"
			encoderCfg.EncodeTime = zapcore.RFC3339TimeEncoder
		}

		switch opts.Format {
		case FormatJSON:
			encoder = zapcore.NewJSONEncoder(encoderCfg)
		default:
			encoder = zapcore.NewConsoleEncoder(encoderCfg)
		}
	}

	cores := make([]zapcore.Core, 0)

	// add stdout log
	if opts.LogToStdout {
		stdoutCore := zapcore.NewCore(
			encoder,
			zapcore.Lock(os.Stdout),
			zap.LevelEnablerFunc(opts.ZapLevelEnabled),
		)
		cores = append(cores, stdoutCore)
	}

	// add output core
	if opts.Output != nil {
		outputCore := zapcore.NewCore(
			encoder,
			zapcore.Lock(zapcore.AddSync(opts.Output)),
			zap.LevelEnablerFunc(opts.ZapLevelEnabled),
		)
		cores = append(cores, outputCore)
	}

	zapOptions := []zap.Option{
		zap.WithCaller(opts.AddCaller),
		zap.AddCallerSkip(opts.CallerSkip),
	}

	if opts.Development {
		zapOptions = append(zapOptions, zap.Development())
	}

	l := &StructuredLogger{
		base:    zap.New(zapcore.NewTee(cores...), zapOptions...),
		options: opts,

		debug:  (*zap.Logger).Debug,
		info:   (*zap.Logger).Info,
		warn:   (*zap.Logger).Warn,
		error:  (*zap.Logger).Error,
		dpanic: (*zap.Logger).DPanic,
		panic:  (*zap.Logger).Panic,
		fatal:  (*zap.Logger).Fatal,
	}

	if opts.Development {
		l.print = (*zap.Logger).Debug
	} else {
		l.print = (*zap.Logger).Info
	}

	return l
}

// WithOptions allows configuring a logger instance with pre-defined settings.
func (l *StructuredLogger) WithOptions(opt ...Option) *StructuredLogger {
	opts := l.options.Clone()

	for _, o := range opt {
		o.apply(&opts)
	}

	return newLogger(opts)
}

// Print logs a log statement with either Debug on Info log levels.
func (l *StructuredLogger) Print(msg string, fields ...Field) {
	l.print(l.base, msg, fields...)
}

// Debug logs a log statement with Debug log level.
func (l *StructuredLogger) Debug(msg string, fields ...Field) {
	l.debug(l.base, msg, fields...)
}

// Info logs a log statement with Info log level.
func (l *StructuredLogger) Info(msg string, fields ...Field) {
	l.info(l.base, msg, fields...)
}

// Warn logs a log statement with Warn log level.
func (l *StructuredLogger) Warn(msg string, fields ...Field) {
	l.warn(l.base, msg, fields...)
}

// Error logs a log statement with Error log level.
func (l *StructuredLogger) Error(msg string, fields ...Field) {
	l.error(l.base, msg, fields...)
}

// DPanic logs a log statement with DPanic log level.
func (l *StructuredLogger) DPanic(msg string, fields ...Field) {
	l.dpanic(l.base, msg, fields...)
}

// Panic logs a log statement with Panic log level.
func (l *StructuredLogger) Panic(msg string, fields ...Field) {
	l.panic(l.base, msg, fields...)
}

// Fatal logs a log statement with Fatal log level.
func (l *StructuredLogger) Fatal(msg string, fields ...Field) {
	l.fatal(l.base, msg, fields...)
}
