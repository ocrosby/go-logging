# Quick Start Example

The fastest way to get started with go-logging.

## Run the Example

```bash
go run main.go
```

## What This Example Shows

- Create a logger with one line: `logging.NewSimple()`
- Log at different levels (Debug, Info, Warn, Error)
- Switch to JSON format with `logging.NewEasyJSON()`
- Change log levels with `logging.NewEasyJSONWithLevel()`

## Key Takeaways

1. **Default behavior**: `NewSimple()` creates a text logger at INFO level
2. **Debug messages**: Won't show unless you set DEBUG level
3. **JSON format**: Perfect for production and log aggregation systems
4. **No configuration needed**: Just import and start logging

## Next Steps

- Try [configuration examples](../configuration/) for more advanced setup
- See [structured logging](../structured/) to add rich data to your logs