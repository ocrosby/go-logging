package logging

import (
	"io"
	"log/slog"
	"os"
	"regexp"
)

func ProvideConfig() *Config {
	return NewConfig().
		FromEnvironment().
		Build()
}

func ProvideConfigWithLevel(level Level) *Config {
	return NewConfig().
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

func ProvideLogger(config *Config, redactorChain RedactorChainInterface) Logger {
	if config.UseSlog {
		return NewSlogLoggerFromConfig(config, redactorChain)
	}
	return NewStandardLogger(config, redactorChain)
}

func NewSlogLoggerFromConfig(config *Config, redactorChain RedactorChainInterface) Logger {
	var handler slog.Handler

	if config.Handler != nil {
		handler = config.Handler
	} else {
		opts := &slog.HandlerOptions{
			Level:     levelToSlogLevel(config.Level),
			AddSource: config.IncludeFile,
		}

		if config.Format == JSONFormat {
			handler = slog.NewJSONHandler(config.Output, opts)
		} else {
			handler = slog.NewTextHandler(config.Output, opts)
		}

		if len(config.StaticFields) > 0 {
			attrs := make([]slog.Attr, 0, len(config.StaticFields))
			for k, v := range config.StaticFields {
				attrs = append(attrs, slog.Any(k, v))
			}
			handler = handler.WithAttrs(attrs)
		}
	}

	logger := NewSlogLogger(handler, redactorChain).(*slogLogger)
	logger.level = config.Level
	return logger
}

func levelToSlogLevel(level Level) slog.Level {
	switch level {
	case TraceLevel:
		return LevelTrace
	case DebugLevel:
		return slog.LevelDebug
	case InfoLevel:
		return slog.LevelInfo
	case WarnLevel:
		return slog.LevelWarn
	case ErrorLevel:
		return slog.LevelError
	case CriticalLevel:
		return LevelCritical
	default:
		return slog.LevelInfo
	}
}
