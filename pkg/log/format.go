package log

import (
	"fmt"
	"strings"
)

// Format is an enum representing log formats.
type Format int

const (
	// FormatConsole is a setting for a key-value pair formatted statements.
	FormatConsole Format = iota
	// FormatJSON is a setting for json formatted statements.
	FormatJSON
)

// ParseFormat parses to Format a passed string value.
func ParseFormat(format string) (Format, error) {
	switch strings.ToLower(format) {
	case "console":
		return FormatConsole, nil
	case "json":
		return FormatJSON, nil
	}

	return FormatConsole, fmt.Errorf("not a valid log format: %q", format)
}
