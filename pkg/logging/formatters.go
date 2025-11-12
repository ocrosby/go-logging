package logging

import (
	"encoding/json"
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/ocrosby/go-logging/pkg/logging/internal"
)

// JSONFormatter formats log entries as JSON.
type JSONFormatter struct {
	config *FormatterConfig
}

// NewJSONFormatter creates a new JSON formatter.
func NewJSONFormatter(config *FormatterConfig) *JSONFormatter {
	if config == nil {
		config = NewFormatterConfig().WithJSONFormat().Build()
	}
	return &JSONFormatter{config: config}
}

// Format formats a log entry as JSON bytes.
func (f *JSONFormatter) Format(entry LogEntry) ([]byte, error) {
	data := make(map[string]interface{})

	if f.config.IncludeTime {
		data["timestamp"] = entry.Timestamp.UTC().Format(time.RFC3339)
	}

	data["level"] = entry.Level.String()
	data["message"] = f.applyRedaction(entry.Message)

	// Add fields
	for k, v := range entry.Fields {
		data[k] = v
	}

	// Add file information if configured
	if f.config.IncludeFile && (entry.File != "" || entry.Line != 0) {
		if entry.File == "" && entry.Line == 0 {
			// Try to get caller info
			if _, file, line, ok := runtime.Caller(4); ok {
				entry.File = file
				entry.Line = line
			}
		}
		if entry.File != "" {
			data["file"] = f.formatFilename(entry.File, entry.Line)
		}
	}

	// Add context information
	contextFields := internal.ExtractContextFields(entry.Context)
	contextFields.AddToMap(data)

	return json.Marshal(data)
}

func (f *JSONFormatter) applyRedaction(message string) string {
	return internal.ApplyRedactionPatterns(message, f.config.RedactPatterns)
}

func (f *JSONFormatter) formatFilename(file string, line int) string {
	return internal.FormatFilename(file, line, f.config.UseShortFile)
}

// TextFormatter formats log entries as human-readable text.
type TextFormatter struct {
	config *FormatterConfig
}

// NewTextFormatter creates a new text formatter.
func NewTextFormatter(config *FormatterConfig) *TextFormatter {
	if config == nil {
		config = NewFormatterConfig().WithTextFormat().Build()
	}
	return &TextFormatter{config: config}
}

// Format formats a log entry as text bytes.
func (f *TextFormatter) Format(entry LogEntry) ([]byte, error) {
	var parts []string

	// Add timestamp if configured
	if f.config.IncludeTime {
		parts = append(parts, entry.Timestamp.Format("2006/01/02 15:04:05"))
	}

	// Add level
	parts = append(parts, fmt.Sprintf("[%s]", strings.ToUpper(entry.Level.String())))

	// Add file information if configured
	if f.config.IncludeFile {
		if entry.File == "" && entry.Line == 0 {
			// Try to get caller info
			if _, file, line, ok := runtime.Caller(4); ok {
				entry.File = file
				entry.Line = line
			}
		}
		if entry.File != "" {
			parts = append(parts, f.formatFilename(entry.File, entry.Line))
		}
	}

	// Add the main message
	message := f.applyRedaction(entry.Message)
	parts = append(parts, message)

	// Add fields if present
	if len(entry.Fields) > 0 {
		var fieldParts []string
		for k, v := range entry.Fields {
			fieldParts = append(fieldParts, fmt.Sprintf("%s=%v", k, v))
		}
		parts = append(parts, fmt.Sprintf("{%s}", strings.Join(fieldParts, " ")))
	}

	// Add context information if present
	contextFields := internal.ExtractContextFields(entry.Context)
	if contextText := contextFields.FormatForText(); contextText != "" {
		parts = append(parts, contextText)
	}

	result := strings.Join(parts, " ") + "\n"
	return []byte(result), nil
}

func (f *TextFormatter) applyRedaction(message string) string {
	return internal.ApplyRedactionPatterns(message, f.config.RedactPatterns)
}

func (f *TextFormatter) formatFilename(file string, line int) string {
	return internal.FormatFilename(file, line, f.config.UseShortFile)
}

// ConsoleFormatter provides colored output for console/terminal usage.
type ConsoleFormatter struct {
	config      *FormatterConfig
	useColors   bool
	levelColors map[Level]string
}

// NewConsoleFormatter creates a new console formatter with color support.
func NewConsoleFormatter(config *FormatterConfig, useColors bool) *ConsoleFormatter {
	if config == nil {
		config = NewFormatterConfig().WithTextFormat().Build()
	}

	return &ConsoleFormatter{
		config:    config,
		useColors: useColors,
		levelColors: map[Level]string{
			TraceLevel:    "\033[36m", // Cyan
			DebugLevel:    "\033[37m", // White
			InfoLevel:     "\033[32m", // Green
			WarnLevel:     "\033[33m", // Yellow
			ErrorLevel:    "\033[31m", // Red
			CriticalLevel: "\033[35m", // Magenta
		},
	}
}

// Format formats a log entry with optional colors for console output.
func (f *ConsoleFormatter) Format(entry LogEntry) ([]byte, error) {
	var parts []string

	// Add timestamp if configured
	if f.config.IncludeTime {
		timestamp := entry.Timestamp.Format("15:04:05")
		if f.useColors {
			timestamp = "\033[90m" + timestamp + "\033[0m" // Dark gray
		}
		parts = append(parts, timestamp)
	}

	// Add level with color
	levelStr := strings.ToUpper(entry.Level.String())
	if f.useColors {
		if color, exists := f.levelColors[entry.Level]; exists {
			levelStr = color + levelStr + "\033[0m"
		}
	}
	parts = append(parts, fmt.Sprintf("[%s]", levelStr))

	// Add the main message
	message := f.applyRedaction(entry.Message)
	parts = append(parts, message)

	// Add fields in a compact format
	if len(entry.Fields) > 0 {
		var fieldParts []string
		for k, v := range entry.Fields {
			fieldStr := fmt.Sprintf("%s=%v", k, v)
			if f.useColors {
				fieldStr = "\033[90m" + fieldStr + "\033[0m" // Dark gray
			}
			fieldParts = append(fieldParts, fieldStr)
		}
		parts = append(parts, strings.Join(fieldParts, " "))
	}

	result := strings.Join(parts, " ") + "\n"
	return []byte(result), nil
}

func (f *ConsoleFormatter) applyRedaction(message string) string {
	return internal.ApplyRedactionPatterns(message, f.config.RedactPatterns)
}
