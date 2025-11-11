//go:build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/ocrosby/go-logging/pkg/logging"
)

func InitializeLogger() logging.Logger {
	wire.Build(logging.DefaultSet)
	return nil
}
