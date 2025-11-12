package logging

import (
	"io"
	"os"
	"regexp"
)

// Legacy providers for backward compatibility with old Config type
func ProvideConfig() *Config {
	// Convert new LoggerConfig to old Config format for backward compatibility
	newConfig := NewLoggerConfig().FromEnvironment().Build()
	return &Config{
		Level:          newConfig.Core.Level,
		Output:         newConfig.Output.Writer,
		Format:         newConfig.Formatter.Format,
		IncludeFile:    newConfig.Formatter.IncludeFile,
		IncludeTime:    newConfig.Formatter.IncludeTime,
		UseShortFile:   newConfig.Formatter.UseShortFile,
		RedactPatterns: newConfig.Formatter.RedactPatterns,
		StaticFields:   newConfig.Core.StaticFields,
		Handler:        newConfig.Handler,
		UseSlog:        newConfig.UseSlog,
	}
}

func ProvideConfigWithLevel(level Level) *Config {
	// Convert new LoggerConfig to old Config format for backward compatibility
	newConfig := NewLoggerConfig().WithLevel(level).FromEnvironment().Build()
	return &Config{
		Level:          newConfig.Core.Level,
		Output:         newConfig.Output.Writer,
		Format:         newConfig.Formatter.Format,
		IncludeFile:    newConfig.Formatter.IncludeFile,
		IncludeTime:    newConfig.Formatter.IncludeTime,
		UseShortFile:   newConfig.Formatter.UseShortFile,
		RedactPatterns: newConfig.Formatter.RedactPatterns,
		StaticFields:   newConfig.Core.StaticFields,
		Handler:        newConfig.Handler,
		UseSlog:        newConfig.UseSlog,
	}
}

// New providers using new config structure
func ProvideLoggerConfig() *LoggerConfig {
	return NewLoggerConfig().
		FromEnvironment().
		Build()
}

func ProvideLoggerConfigWithLevel(level Level) *LoggerConfig {
	return NewLoggerConfig().
		WithLevel(level).
		FromEnvironment().
		Build()
}

func ProvideOutput() io.Writer {
	return os.Stdout
}

func ProvideRedactorChain(config *Config) RedactorChainInterface {
	return NewRedactorChain(config.RedactPatterns...)
}

func ProvideRedactorChainWithPatterns(patterns ...*regexp.Regexp) RedactorChainInterface {
	return NewRedactorChain(patterns...)
}

func ProvideRedactorChainFromLoggerConfig(config *LoggerConfig) RedactorChainInterface {
	return NewRedactorChain(config.Formatter.RedactPatterns...)
}

// Updated provider using unified logger
func ProvideLogger(config *Config, redactorChain RedactorChainInterface) Logger {
	// Convert old config to new config format
	loggerConfig := &LoggerConfig{
		Core: &CoreConfig{
			Level:        config.Level,
			StaticFields: config.StaticFields,
		},
		Formatter: &FormatterConfig{
			Format:         config.Format,
			IncludeFile:    config.IncludeFile,
			IncludeTime:    config.IncludeTime,
			UseShortFile:   config.UseShortFile,
			RedactPatterns: config.RedactPatterns,
		},
		Output: &OutputConfig{
			Writer: config.Output,
		},
		Handler: config.Handler,
		UseSlog: config.UseSlog,
	}

	return NewUnifiedLogger(loggerConfig, redactorChain)
}

// New provider using new config structure
func ProvideLoggerFromConfig(config *LoggerConfig, redactorChain RedactorChainInterface) Logger {
	return NewUnifiedLogger(config, redactorChain)
}
