package logging

import (
	"io"
	"log/slog"
	"regexp"
)

type OutputFormat int

const (
	CommonLogFormat OutputFormat = iota
	TextFormat
	JSONFormat
)

// Config provides backward compatibility with the old configuration system.
// Deprecated: Use LoggerConfig for new code.
type Config struct {
	Level          Level
	Output         io.Writer
	Format         OutputFormat
	IncludeFile    bool
	IncludeTime    bool
	UseShortFile   bool
	RedactPatterns []*regexp.Regexp
	StaticFields   map[string]interface{}
	Handler        slog.Handler
	UseSlog        bool
}

// ToLoggerConfig converts old Config to new LoggerConfig structure.
func (c *Config) ToLoggerConfig() *LoggerConfig {
	return &LoggerConfig{
		Core: &CoreConfig{
			Level:        c.Level,
			StaticFields: c.StaticFields,
		},
		Formatter: &FormatterConfig{
			Format:         c.Format,
			IncludeFile:    c.IncludeFile,
			IncludeTime:    c.IncludeTime,
			UseShortFile:   c.UseShortFile,
			RedactPatterns: c.RedactPatterns,
		},
		Output: &OutputConfig{
			Writer: c.Output,
		},
		Handler: c.Handler,
		UseSlog: c.UseSlog,
	}
}

type ConfigBuilder struct {
	builder *LoggerConfigBuilder
}

func NewConfig() *ConfigBuilder {
	return &ConfigBuilder{
		builder: NewLoggerConfig(),
	}
}

func (b *ConfigBuilder) WithLevel(level Level) *ConfigBuilder {
	b.builder.WithLevel(level)
	return b
}

func (b *ConfigBuilder) WithLevelString(level string) *ConfigBuilder {
	b.builder.WithLevelString(level)
	return b
}

func (b *ConfigBuilder) WithOutput(w io.Writer) *ConfigBuilder {
	b.builder.WithWriter(w)
	return b
}

func (b *ConfigBuilder) WithFormat(format OutputFormat) *ConfigBuilder {
	b.builder.WithFormat(format)
	return b
}

func (b *ConfigBuilder) WithJSONFormat() *ConfigBuilder {
	b.builder.WithJSONFormat()
	return b
}

func (b *ConfigBuilder) WithTextFormat() *ConfigBuilder {
	b.builder.WithTextFormat()
	return b
}

func (b *ConfigBuilder) WithCommonLogFormat() *ConfigBuilder {
	b.builder.WithCommonLogFormat()
	return b
}

func (b *ConfigBuilder) IncludeFile(include bool) *ConfigBuilder {
	b.builder.config.Formatter.IncludeFile = include
	return b
}

func (b *ConfigBuilder) IncludeTime(include bool) *ConfigBuilder {
	b.builder.config.Formatter.IncludeTime = include
	return b
}

func (b *ConfigBuilder) UseShortFile(useShort bool) *ConfigBuilder {
	b.builder.config.Formatter.UseShortFile = useShort
	return b
}

func (b *ConfigBuilder) AddRedactPattern(pattern string) *ConfigBuilder {
	if re, err := regexp.Compile(pattern); err == nil {
		b.builder.config.Formatter.RedactPatterns = append(b.builder.config.Formatter.RedactPatterns, re)
	}
	return b
}

func (b *ConfigBuilder) AddRedactRegex(re *regexp.Regexp) *ConfigBuilder {
	b.builder.config.Formatter.RedactPatterns = append(b.builder.config.Formatter.RedactPatterns, re)
	return b
}

func (b *ConfigBuilder) WithStaticField(key string, value interface{}) *ConfigBuilder {
	b.builder.config.Core.StaticFields[key] = value
	return b
}

func (b *ConfigBuilder) WithStaticFields(fields map[string]interface{}) *ConfigBuilder {
	for k, v := range fields {
		b.builder.config.Core.StaticFields[k] = v
	}
	return b
}

func (b *ConfigBuilder) WithHandler(handler slog.Handler) *ConfigBuilder {
	b.builder.WithHandler(handler)
	return b
}

func (b *ConfigBuilder) UseSlog(use bool) *ConfigBuilder {
	b.builder.UseSlog(use)
	return b
}

func (b *ConfigBuilder) FromEnvironment() *ConfigBuilder {
	b.builder.FromEnvironment()
	return b
}

func (b *ConfigBuilder) Build() *Config {
	loggerConfig := b.builder.Build()
	return &Config{
		Level:          loggerConfig.Core.Level,
		Output:         loggerConfig.Output.Writer,
		Format:         loggerConfig.Formatter.Format,
		IncludeFile:    loggerConfig.Formatter.IncludeFile,
		IncludeTime:    loggerConfig.Formatter.IncludeTime,
		UseShortFile:   loggerConfig.Formatter.UseShortFile,
		RedactPatterns: loggerConfig.Formatter.RedactPatterns,
		StaticFields:   loggerConfig.Core.StaticFields,
		Handler:        loggerConfig.Handler,
		UseSlog:        loggerConfig.UseSlog,
	}
}
