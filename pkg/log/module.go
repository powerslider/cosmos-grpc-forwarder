package log

// InitializeLogger configures a sane default logger module.
func InitializeLogger(logLevel string, logFormat string) *StructuredLogger {
	lf, err := ParseFormat(logFormat)
	if err != nil {
		panic(err)
	}

	ll, err := ParseLevel(logLevel)
	if err != nil {
		panic(err)
	}

	return New(
		WithFormat(lf),
		WithLevel(ll),
		AddCaller(),
		ToStdout(),
	)
}
