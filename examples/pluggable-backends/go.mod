module github.com/ocrosby/go-logging/examples/pluggable-backends

go 1.24

require (
	github.com/ocrosby/go-logging v0.0.0
	github.com/rs/zerolog v1.34.0
	go.uber.org/zap v1.27.1
	go.uber.org/zap/exp v0.3.0
)

require (
	github.com/google/wire v0.7.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	golang.org/x/sys v0.12.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/ocrosby/go-logging => ../..
