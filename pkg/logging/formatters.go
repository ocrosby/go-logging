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

	f.addBaseFields(entry, data)
	f.addUserFields(entry, data)
	f.addFileInfo(entry, data)
	f.addContextFields(entry, data)

	return json.Marshal(data)
}

func (f *JSONFormatter) addBaseFields(entry LogEntry, data map[string]interface{}) {
	if f.config.IncludeTime {
		data["timestamp"] = entry.Timestamp.UTC().Format(time.RFC3339)
	}
	data["level"] = entry.Level.String()
	data["message"] = f.applyRedaction(entry.Message)
}

func (f *JSONFormatter) addUserFields(entry LogEntry, data map[string]interface{}) {
	for k, v := range entry.Fields {
		data[k] = v
	}
}

func (f *JSONFormatter) addFileInfo(entry LogEntry, data map[string]interface{}) {
	if !f.config.IncludeFile {
		return
	}

	file, line := f.getFileInfo(entry)
	if file != "" {
		data["file"] = f.formatFilename(file, line)
	}
}

func (f *JSONFormatter) getFileInfo(entry LogEntry) (string, int) {
	if entry.File != "" || entry.Line != 0 {
		return entry.File, entry.Line
	}

	if _, file, line, ok := runtime.Caller(4); ok {
		return file, line
	}

	return "", 0
}

func (f *JSONFormatter) addContextFields(entry LogEntry, data map[string]interface{}) {
	contextFields := internal.ExtractContextFields(entry.Context)
	contextFields.AddToMap(data)
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

	f.addTimestamp(&parts, entry)
	f.addLevel(&parts, entry)
	f.addFileInfoText(&parts, entry)
	f.addMessage(&parts, entry)
	f.addFieldsText(&parts, entry)
	f.addContextText(&parts, entry)

	result := strings.Join(parts, " ") + "\n"
	return []byte(result), nil
}

func (f *TextFormatter) addTimestamp(parts *[]string, entry LogEntry) {
	if f.config.IncludeTime {
		*parts = append(*parts, entry.Timestamp.Format("2006/01/02 15:04:05"))
	}
}

func (f *TextFormatter) addLevel(parts *[]string, entry LogEntry) {
	*parts = append(*parts, fmt.Sprintf("[%s]", strings.ToUpper(entry.Level.String())))
}

func (f *TextFormatter) addFileInfoText(parts *[]string, entry LogEntry) {
	if !f.config.IncludeFile {
		return
	}

	file, line := f.getFileInfoText(entry)
	if file != "" {
		*parts = append(*parts, f.formatFilename(file, line))
	}
}

func (f *TextFormatter) getFileInfoText(entry LogEntry) (string, int) {
	if entry.File != "" || entry.Line != 0 {
		return entry.File, entry.Line
	}

	if _, file, line, ok := runtime.Caller(4); ok {
		return file, line
	}

	return "", 0
}

func (f *TextFormatter) addMessage(parts *[]string, entry LogEntry) {
	message := f.applyRedaction(entry.Message)
	*parts = append(*parts, message)
}

func (f *TextFormatter) addFieldsText(parts *[]string, entry LogEntry) {
	if len(entry.Fields) == 0 {
		return
	}

	var fieldParts []string
	for k, v := range entry.Fields {
		fieldParts = append(fieldParts, fmt.Sprintf("%s=%v", k, v))
	}
	*parts = append(*parts, fmt.Sprintf("{%s}", strings.Join(fieldParts, " ")))
}

func (f *TextFormatter) addContextText(parts *[]string, entry LogEntry) {
	contextFields := internal.ExtractContextFields(entry.Context)
	if contextText := contextFields.FormatForText(); contextText != "" {
		*parts = append(*parts, contextText)
	}
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

	f.addTimestampConsole(&parts, entry)
	f.addLevelConsole(&parts, entry)
	f.addMessageConsole(&parts, entry)
	f.addFieldsConsole(&parts, entry)

	result := strings.Join(parts, " ") + "\n"
	return []byte(result), nil
}

func (f *ConsoleFormatter) addTimestampConsole(parts *[]string, entry LogEntry) {
	if !f.config.IncludeTime {
		return
	}

	timestamp := entry.Timestamp.Format("15:04:05")
	if f.useColors {
		timestamp = "\033[90m" + timestamp + "\033[0m" // Dark gray
	}
	*parts = append(*parts, timestamp)
}

func (f *ConsoleFormatter) addLevelConsole(parts *[]string, entry LogEntry) {
	levelStr := strings.ToUpper(entry.Level.String())
	if f.useColors {
		if color, exists := f.levelColors[entry.Level]; exists {
			levelStr = color + levelStr + "\033[0m"
		}
	}
	*parts = append(*parts, fmt.Sprintf("[%s]", levelStr))
}

func (f *ConsoleFormatter) addMessageConsole(parts *[]string, entry LogEntry) {
	message := f.applyRedaction(entry.Message)
	*parts = append(*parts, message)
}

func (f *ConsoleFormatter) addFieldsConsole(parts *[]string, entry LogEntry) {
	if len(entry.Fields) == 0 {
		return
	}

	var fieldParts []string
	for k, v := range entry.Fields {
		fieldStr := fmt.Sprintf("%s=%v", k, v)
		if f.useColors {
			fieldStr = "\033[90m" + fieldStr + "\033[0m" // Dark gray
		}
		fieldParts = append(fieldParts, fieldStr)
	}
	*parts = append(*parts, strings.Join(fieldParts, " "))
}

func (f *ConsoleFormatter) applyRedaction(message string) string {
	return internal.ApplyRedactionPatterns(message, f.config.RedactPatterns)
}

// CommonLogFormatter formats log entries in the NCSA Common Log Format.
type CommonLogFormatter struct {
	config *FormatterConfig
}

// NewCommonLogFormatter creates a new Common Log Format formatter.
func NewCommonLogFormatter(config *FormatterConfig) *CommonLogFormatter {
	if config == nil {
		config = NewFormatterConfig().Build()
	}
	return &CommonLogFormatter{config: config}
}

// Format formats a log entry according to Common Log Format.
// Format: host ident authuser [timestamp] "request-line" status bytes
func (f *CommonLogFormatter) Format(entry LogEntry) ([]byte, error) {
	host := f.getField(entry, "host", "-")
	ident := f.getField(entry, "ident", "-")
	authuser := f.getField(entry, "authuser", "-")
	timestamp := entry.Timestamp.Format("02/Jan/2006:15:04:05 -0700")
	requestLine := f.getField(entry, "request", entry.Message)
	status := f.getField(entry, "status", "-")
	bytes := f.getField(entry, "bytes", "-")

	result := fmt.Sprintf("%s %s %s [%s] \"%s\" %s %s\n",
		host, ident, authuser, timestamp, requestLine, status, bytes)

	return []byte(result), nil
}

func (f *CommonLogFormatter) getField(entry LogEntry, key string, defaultValue string) string {
	if value, ok := entry.Fields[key]; ok {
		return fmt.Sprintf("%v", value)
	}
	return defaultValue
}
