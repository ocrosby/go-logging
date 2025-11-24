package logging

import (
	"io"
	"log/slog"
	"os"
	"regexp"
)

const (
	jsonFormatString = "json"
	textFormatString = "text"
)

// CoreConfig contains the core logging configuration.
type CoreConfig struct {
	Level        Level
	StaticFields map[string]interface{}
}

// FormatterConfig contains formatting-related configuration.
type FormatterConfig struct {
	Format         OutputFormat
	IncludeFile    bool
	IncludeTime    bool
	UseShortFile   bool
	RedactPatterns []*regexp.Regexp
}

// OutputConfig contains output-related configuration.
type OutputConfig struct {
	Writer io.Writer
}

// LoggerConfig combines all configuration types.
type LoggerConfig struct {
	Core      *CoreConfig
	Formatter *FormatterConfig
	Output    *OutputConfig
	Handler   slog.Handler
	UseSlog   bool
}

// CoreConfigBuilder builds CoreConfig instances.
type CoreConfigBuilder struct {
	config *CoreConfig
}

// NewCoreConfig creates a new CoreConfigBuilder with defaults.
func NewCoreConfig() *CoreConfigBuilder {
	return &CoreConfigBuilder{
		config: &CoreConfig{
			Level:        InfoLevel,
			StaticFields: make(map[string]interface{}),
		},
	}
}

func (b *CoreConfigBuilder) WithLevel(level Level) *CoreConfigBuilder {
	b.config.Level = level
	return b
}

func (b *CoreConfigBuilder) WithLevelString(level string) *CoreConfigBuilder {
	if l, ok := ParseLevel(level); ok {
		b.config.Level = l
	}
	return b
}

func (b *CoreConfigBuilder) WithStaticField(key string, value interface{}) *CoreConfigBuilder {
	b.config.StaticFields[key] = value
	return b
}

func (b *CoreConfigBuilder) WithStaticFields(fields map[string]interface{}) *CoreConfigBuilder {
	for k, v := range fields {
		b.config.StaticFields[k] = v
	}
	return b
}

func (b *CoreConfigBuilder) Build() *CoreConfig {
	return b.config
}

// FormatterConfigBuilder builds FormatterConfig instances.
type FormatterConfigBuilder struct {
	config *FormatterConfig
}

// NewFormatterConfig creates a new FormatterConfigBuilder with defaults.
func NewFormatterConfig() *FormatterConfigBuilder {
	return &FormatterConfigBuilder{
		config: &FormatterConfig{
			Format:         CommonLogFormat,
			IncludeFile:    true,
			IncludeTime:    true,
			UseShortFile:   true,
			RedactPatterns: make([]*regexp.Regexp, 0),
		},
	}
}

func (b *FormatterConfigBuilder) WithFormat(format OutputFormat) *FormatterConfigBuilder {
	b.config.Format = format
	return b
}

func (b *FormatterConfigBuilder) WithJSONFormat() *FormatterConfigBuilder {
	b.config.Format = JSONFormat
	return b
}

func (b *FormatterConfigBuilder) WithTextFormat() *FormatterConfigBuilder {
	b.config.Format = TextFormat
	return b
}

func (b *FormatterConfigBuilder) WithCommonLogFormat() *FormatterConfigBuilder {
	b.config.Format = CommonLogFormat
	return b
}

func (b *FormatterConfigBuilder) IncludeFile(include bool) *FormatterConfigBuilder {
	b.config.IncludeFile = include
	return b
}

func (b *FormatterConfigBuilder) IncludeTime(include bool) *FormatterConfigBuilder {
	b.config.IncludeTime = include
	return b
}

func (b *FormatterConfigBuilder) UseShortFile(useShort bool) *FormatterConfigBuilder {
	b.config.UseShortFile = useShort
	return b
}

func (b *FormatterConfigBuilder) AddRedactPattern(pattern string) *FormatterConfigBuilder {
	if re, err := regexp.Compile(pattern); err == nil {
		b.config.RedactPatterns = append(b.config.RedactPatterns, re)
	}
	return b
}

func (b *FormatterConfigBuilder) AddRedactRegex(re *regexp.Regexp) *FormatterConfigBuilder {
	b.config.RedactPatterns = append(b.config.RedactPatterns, re)
	return b
}

func (b *FormatterConfigBuilder) Build() *FormatterConfig {
	return b.config
}

// OutputConfigBuilder builds OutputConfig instances.
type OutputConfigBuilder struct {
	config *OutputConfig
}

// NewOutputConfig creates a new OutputConfigBuilder with defaults.
func NewOutputConfig() *OutputConfigBuilder {
	return &OutputConfigBuilder{
		config: &OutputConfig{
			Writer: os.Stdout,
		},
	}
}

func (b *OutputConfigBuilder) WithWriter(w io.Writer) *OutputConfigBuilder {
	b.config.Writer = w
	return b
}

func (b *OutputConfigBuilder) Build() *OutputConfig {
	return b.config
}

// LoggerConfigBuilder builds complete LoggerConfig instances.
type LoggerConfigBuilder struct {
	config *LoggerConfig
}

// NewLoggerConfig creates a new LoggerConfigBuilder with defaults.
func NewLoggerConfig() *LoggerConfigBuilder {
	return &LoggerConfigBuilder{
		config: &LoggerConfig{
			Core:      NewCoreConfig().Build(),
			Formatter: NewFormatterConfig().Build(),
			Output:    NewOutputConfig().Build(),
			UseSlog:   false,
		},
	}
}

func (b *LoggerConfigBuilder) WithCore(core *CoreConfig) *LoggerConfigBuilder {
	b.config.Core = core
	return b
}

func (b *LoggerConfigBuilder) WithFormatter(formatter *FormatterConfig) *LoggerConfigBuilder {
	b.config.Formatter = formatter
	return b
}

func (b *LoggerConfigBuilder) WithOutput(output *OutputConfig) *LoggerConfigBuilder {
	b.config.Output = output
	return b
}

func (b *LoggerConfigBuilder) WithHandler(handler slog.Handler) *LoggerConfigBuilder {
	b.config.Handler = handler
	b.config.UseSlog = true
	return b
}

func (b *LoggerConfigBuilder) UseSlog(use bool) *LoggerConfigBuilder {
	b.config.UseSlog = use
	return b
}

func (b *LoggerConfigBuilder) FromEnvironment() *LoggerConfigBuilder {
	if level := os.Getenv("LOG_LEVEL"); level != "" {
		if l, ok := ParseLevel(level); ok {
			b.config.Core.Level = l
		}
	}
	if format := os.Getenv("LOG_FORMAT"); format == jsonFormatString {
		b.config.Formatter.Format = JSONFormat
	}
	return b
}

func (b *LoggerConfigBuilder) Build() *LoggerConfig {
	return b.config
}

// Convenience methods for backward compatibility
func (b *LoggerConfigBuilder) WithLevel(level Level) *LoggerConfigBuilder {
	b.config.Core.Level = level
	return b
}

func (b *LoggerConfigBuilder) WithLevelString(level string) *LoggerConfigBuilder {
	if l, ok := ParseLevel(level); ok {
		b.config.Core.Level = l
	}
	return b
}

func (b *LoggerConfigBuilder) WithWriter(w io.Writer) *LoggerConfigBuilder {
	b.config.Output.Writer = w
	return b
}

func (b *LoggerConfigBuilder) WithFormat(format OutputFormat) *LoggerConfigBuilder {
	b.config.Formatter.Format = format
	return b
}

func (b *LoggerConfigBuilder) WithJSONFormat() *LoggerConfigBuilder {
	b.config.Formatter.Format = JSONFormat
	return b
}

func (b *LoggerConfigBuilder) WithTextFormat() *LoggerConfigBuilder {
	b.config.Formatter.Format = TextFormat
	return b
}

func (b *LoggerConfigBuilder) WithCommonLogFormat() *LoggerConfigBuilder {
	b.config.Formatter.Format = CommonLogFormat
	return b
}
