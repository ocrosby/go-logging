package logging

import (
	"io"
	"log/slog"
	"os"
	"regexp"
)

type OutputFormat int

const (
	TextFormat OutputFormat = iota
	JSONFormat
)

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

type ConfigBuilder struct {
	config *Config
}

func NewConfig() *ConfigBuilder {
	return &ConfigBuilder{
		config: &Config{
			Level:          InfoLevel,
			Output:         os.Stdout,
			Format:         TextFormat,
			IncludeFile:    true,
			IncludeTime:    true,
			UseShortFile:   true,
			RedactPatterns: make([]*regexp.Regexp, 0),
			StaticFields:   make(map[string]interface{}),
		},
	}
}

func (b *ConfigBuilder) WithLevel(level Level) *ConfigBuilder {
	b.config.Level = level
	return b
}

func (b *ConfigBuilder) WithLevelString(level string) *ConfigBuilder {
	if l, ok := ParseLevel(level); ok {
		b.config.Level = l
	}
	return b
}

func (b *ConfigBuilder) WithOutput(w io.Writer) *ConfigBuilder {
	b.config.Output = w
	return b
}

func (b *ConfigBuilder) WithFormat(format OutputFormat) *ConfigBuilder {
	b.config.Format = format
	return b
}

func (b *ConfigBuilder) WithJSONFormat() *ConfigBuilder {
	b.config.Format = JSONFormat
	return b
}

func (b *ConfigBuilder) WithTextFormat() *ConfigBuilder {
	b.config.Format = TextFormat
	return b
}

func (b *ConfigBuilder) IncludeFile(include bool) *ConfigBuilder {
	b.config.IncludeFile = include
	return b
}

func (b *ConfigBuilder) IncludeTime(include bool) *ConfigBuilder {
	b.config.IncludeTime = include
	return b
}

func (b *ConfigBuilder) UseShortFile(useShort bool) *ConfigBuilder {
	b.config.UseShortFile = useShort
	return b
}

func (b *ConfigBuilder) AddRedactPattern(pattern string) *ConfigBuilder {
	if re, err := regexp.Compile(pattern); err == nil {
		b.config.RedactPatterns = append(b.config.RedactPatterns, re)
	}
	return b
}

func (b *ConfigBuilder) AddRedactRegex(re *regexp.Regexp) *ConfigBuilder {
	b.config.RedactPatterns = append(b.config.RedactPatterns, re)
	return b
}

func (b *ConfigBuilder) WithStaticField(key string, value interface{}) *ConfigBuilder {
	b.config.StaticFields[key] = value
	return b
}

func (b *ConfigBuilder) WithStaticFields(fields map[string]interface{}) *ConfigBuilder {
	for k, v := range fields {
		b.config.StaticFields[k] = v
	}
	return b
}

func (b *ConfigBuilder) WithHandler(handler slog.Handler) *ConfigBuilder {
	b.config.Handler = handler
	b.config.UseSlog = true
	return b
}

func (b *ConfigBuilder) UseSlog(use bool) *ConfigBuilder {
	b.config.UseSlog = use
	return b
}

func (b *ConfigBuilder) FromEnvironment() *ConfigBuilder {
	if level := os.Getenv("LOG_LEVEL"); level != "" {
		b.WithLevelString(level)
	}
	if format := os.Getenv("LOG_FORMAT"); format == "json" {
		b.WithJSONFormat()
	}
	return b
}

func (b *ConfigBuilder) Build() *Config {
	return b.config
}
