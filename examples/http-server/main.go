package main

import (
	"fmt"
	"net/http"

	"github.com/ocrosby/go-logging/pkg/logging"
)

func main() {
	config := logging.NewConfig().
		WithLevel(logging.InfoLevel).
		WithJSONFormat().
		Build()
	redactorChain := logging.ProvideRedactorChain(config)
	logger := logging.NewStandardLogger(config, redactorChain)

	mux := http.NewServeMux()

	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		logger.Fluent().Info().
			Ctx(ctx).
			Str("handler", "hello").
			Msg("Handling hello request")

		fmt.Fprintf(w, "Hello, World!")
	})

	mux.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		logger.Fluent().Info().
			Ctx(ctx).
			Str("handler", "user").
			Str("method", r.Method).
			Msg("Handling user request")

		fmt.Fprintf(w, "User endpoint")
	})

	handler := logging.TracingMiddleware(logger)(mux)

	logger.Info("Starting server on :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		logger.Critical("Server failed: %v", err)
	}
}
