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
	wire.Bind(new(Logger), new(*unifiedLogger)),
)

// New sets using new config structure
var NewDefaultSet = wire.NewSet(
	ProvideLoggerConfig,
	ProvideRedactorChainFromLoggerConfig,
	ProvideLoggerFromConfig,
)

var NewJSONLoggerSet = wire.NewSet(
	ProvideLoggerConfig,
	ProvideRedactorChainFromLoggerConfig,
	ProvideLoggerFromConfig,
	wire.Bind(new(Logger), new(*unifiedLogger)),
)
