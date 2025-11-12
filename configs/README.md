# YAML Configuration Examples

This directory contains example YAML configuration files for various logging scenarios. These files demonstrate how to use the YAML-based configuration system to simplify complex logging setups.

## Available Configurations

### 🚀 [`development.yaml`](development.yaml)
Perfect for local development with detailed debugging information:
- **Level**: Debug
- **Format**: Text (human-readable)
- **Includes**: File info, timestamps
- **Use case**: Local development, debugging

### 📦 [`production.yaml`](production.yaml)
Optimized for production environments:
- **Level**: Info
- **Format**: JSON (structured)
- **Backend**: slog for performance
- **Features**: Security redaction patterns
- **Use case**: Production deployments

### 📁 [`file-logging.yaml`](file-logging.yaml)
Logs to files with rotation-friendly settings:
- **Output**: File (`~/logs/application.log`)
- **Format**: JSON for parsing
- **Features**: Automatic directory creation
- **Use case**: File-based logging systems

### 🎯 [`minimal.yaml`](minimal.yaml)
Absolute minimum configuration:
- **Level**: Info
- **Format**: Text
- **Features**: No timestamps, no file info
- **Use case**: Simple applications, containers

### 🔧 [`microservice.yaml`](microservice.yaml)
Comprehensive configuration for microservices:
- **Features**: Rich static fields, security patterns
- **Backend**: slog for performance
- **Format**: JSON for log aggregation
- **Use case**: Microservice architectures

## Usage Examples

### Basic Usage

```go
package main

import (
    "log"
    "github.com/ocrosby/go-logging/pkg/logging"
)

func main() {
    // Load from specific config file
    logger, err := logging.NewFromYAMLFile("configs/production.yaml")
    if err != nil {
        log.Fatal(err)
    }
    
    logger.Info("Application started with YAML config")
}
```

### Environment-Based Configuration

```go
// Use LOG_CONFIG environment variable
logger := logging.NewFromYAMLEnv("LOG_CONFIG")

// Falls back to simple logger if env var not set or file load fails
logger.Info("This always works")
```

### Embedded Configuration

```go
yamlConfig := `
level: debug
format: json
static_fields:
  service: my-app
output:
  type: stdout
`

logger, err := logging.NewFromYAMLString(yamlConfig)
if err != nil {
    log.Fatal(err)
}
```

## Configuration Schema

### Complete YAML Structure

```yaml
# Optional preset to apply common configurations
preset: "development" | "production" | "debug" | "minimal" | "structured"

# Core configuration
level: trace | debug | info | warn | error | critical
static_fields:
  service: "my-app"
  version: "1.0.0"
  custom_field: "value"

# Formatting
format: text | json
include_file: true | false
include_time: true | false
use_short_file: true | false

# Output configuration
output:
  type: stdout | stderr | file
  target: "/path/to/logfile"  # required for type: file

# Slog backend (optional)
use_slog: true | false

# Security (optional)
redact_patterns:
  - "password=\\w+"
  - "token=[^\\s]*"
```

### Available Presets

| Preset | Level | Format | File Info | slog | Use Case |
|--------|-------|--------|-----------|------|----------|
| `development` | debug | text | ✅ | ❌ | Local development |
| `production` | info | json | ❌ | ✅ | Production systems |
| `debug` | trace | text | ✅ (full) | ❌ | Debugging issues |
| `minimal` | info | text | ❌ | ❌ | Simple applications |
| `structured` | info | json | ✅ | ✅ | Microservices |

## File Output Features

When using `output.type: file`:

- **Automatic directory creation**: Parent directories are created if they don't exist
- **Home directory expansion**: `~/logs/app.log` expands to user home
- **Append mode**: Files are opened in append mode for log rotation compatibility

## Security Features

### Built-in Redaction Patterns

All configurations include patterns to redact sensitive data:

- Passwords: `password=***`
- API Keys: `apikey=***`
- Tokens: `token=***`
- Authorization headers: `authorization: ***`

### Custom Patterns

Add your own redaction patterns:

```yaml
redact_patterns:
  - "secret_key=[^\\s]*"
  - "user_token:[^\\r\\n]*"
  - "\\b\\d{4}[-\\s]?\\d{4}[-\\s]?\\d{4}[-\\s]?\\d{4}\\b"  # Credit cards
```

## Environment Variables

You can override any YAML setting with environment variables:

```bash
export LOG_LEVEL=debug
export LOG_FORMAT=json
export LOG_CONFIG=configs/production.yaml

# Then use
logger := logging.NewFromYAMLEnv("LOG_CONFIG")
```

## Best Practices

1. **Use presets**: Start with a preset and customize as needed
2. **Version your configs**: Keep different configs for different environments
3. **Test locally**: Use `development.yaml` for local testing
4. **Security first**: Always include redaction patterns for sensitive data
5. **File permissions**: Ensure log directories have correct permissions
6. **Rotation friendly**: Use JSON format for log rotation and parsing tools

## Creating Custom Configurations

1. Copy an existing configuration file
2. Modify the settings for your use case
3. Test with a simple application
4. Deploy with your application

## Troubleshooting

### Common Issues

- **File not found**: Check file paths and permissions
- **Invalid YAML**: Use a YAML validator
- **Pattern errors**: Test regex patterns separately
- **Permission denied**: Check directory permissions for file output

### Validation

```go
// Test your configuration
logger, err := logging.NewFromYAMLFile("your-config.yaml")
if err != nil {
    fmt.Printf("Configuration error: %v\n", err)
    return
}

logger.Info("Configuration loaded successfully")
```