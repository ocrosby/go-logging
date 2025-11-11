package logging

import "strings"

type Level int

const (
	TraceLevel Level = iota
	DebugLevel
	InfoLevel
	WarnLevel
	ErrorLevel
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

func (l Level) String() string {
	if name, ok := levelNames[l]; ok {
		return name
	}
	return "UNKNOWN"
}

func ParseLevel(level string) (Level, bool) {
	l, ok := nameLevels[strings.ToUpper(level)]
	return l, ok
}
