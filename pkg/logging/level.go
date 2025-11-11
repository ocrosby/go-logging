package logging

import "strings"

// Level represents the severity level of a log message.
// Lower numeric values indicate more verbose logging.
type Level int

const (
	// TraceLevel is the most verbose level, used for very detailed debugging.
	TraceLevel Level = iota
	// DebugLevel is used for diagnostic information useful during development.
	DebugLevel
	// InfoLevel is the default level for general informational messages.
	InfoLevel
	// WarnLevel is used for warning conditions that don't prevent operation.
	WarnLevel
	// ErrorLevel is used for error conditions that may affect functionality.
	ErrorLevel
	// CriticalLevel is the least verbose level for critical conditions requiring immediate attention.
	CriticalLevel
)

var levelNames = map[Level]string{
	TraceLevel:    "TRACE",
	DebugLevel:    "DEBUG",
	InfoLevel:     "INFO",
	WarnLevel:     "WARN",
	ErrorLevel:    "ERROR",
	CriticalLevel: "CRITICAL",
}

var nameLevels = map[string]Level{
	"TRACE":    TraceLevel,
	"DEBUG":    DebugLevel,
	"INFO":     InfoLevel,
	"WARN":     WarnLevel,
	"ERROR":    ErrorLevel,
	"CRITICAL": CriticalLevel,
}

// String returns the string representation of the log level.
// Returns "UNKNOWN" if the level is not recognized.
func (l Level) String() string {
	if name, ok := levelNames[l]; ok {
		return name
	}
	return "UNKNOWN"
}

// ParseLevel parses a string level name into a Level value.
// The parsing is case-insensitive. Returns the level and true if successful,
// or TraceLevel and false if the level name is not recognized.
//
// Example:
//
//	level, ok := logging.ParseLevel("INFO")
//	if !ok {
//		// handle invalid level
//	}
func ParseLevel(level string) (Level, bool) {
	l, ok := nameLevels[strings.ToUpper(level)]
	return l, ok
}
