package logging

import (
	"io"
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
	return NewStandardLogger(config, redactorChain)
}
