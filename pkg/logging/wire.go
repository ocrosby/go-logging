package logging

import "github.com/google/wire"

var DefaultSet = wire.NewSet(
	ProvideConfig,
	ProvideRedactorChain,
	ProvideLogger,
)

var JSONLoggerSet = wire.NewSet(
	ProvideConfig,
	ProvideRedactorChain,
	ProvideLogger,
	wire.Bind(new(Logger), new(*standardLogger)),
)
