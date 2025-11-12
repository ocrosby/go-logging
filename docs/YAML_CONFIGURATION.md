# YAML Configuration Guide

This guide covers the comprehensive YAML-based configuration system that simplifies complex logging setups through declarative configuration files.

## Table of Contents

- [Quick Start](#quick-start)
- [Configuration Schema](#configuration-schema)
- [Presets](#presets)
- [Examples](#examples)
- [Advanced Features](#advanced-features)
- [Best Practices](#best-practices)
- [Troubleshooting](#troubleshooting)

## Quick Start

### Basic Usage

```go
import "github.com/ocrosby/go-logging/pkg/logging"

// Load from file
logger, err := logging.NewFromYAMLFile("config/logging.yaml")
if err != nil {
    log.Fatal(err)
}

// Load from environment variable (with fallback)
logger := logging.NewFromYAMLEnv("LOG_CONFIG_FILE")

// Load from string
yamlConfig := `
preset: production
static_fields:
  service: my-app
`
logger, err := logging.LoadFromYAMLString(yamlConfig)
```

### Minimal Configuration

```yaml
# config/logging.yaml
preset: development
```

This single line gives you a fully configured logger optimized for development!

## Configuration Schema

### Complete Structure

```yaml
# Optional preset for common configurations
preset: development | production | debug | minimal | structured

# Core logging settings
level: trace | debug | info | warn | error | critical
static_fields:
  key: value
  service: "my-app"
  version: "1.0.0"

# Output formatting
format: text | json
include_file: true | false      # Include file and line info
include_time: true | false      # Include timestamps
use_short_file: true | false    # Use short file paths

# Output destination
output:
  type: stdout | stderr | file
  target: "/path/to/logfile"    # Required for type: file

# Backend selection
use_slog: true | false          # Use Go's slog backend

# Security redaction
redact_patterns:
  - "password=[^\\s]*"
  - "token=[^\\s]*"
  - "apikey=[^\\s]*"

# Advanced slog configuration (optional)
slog:
  handler_type: text | json
  options:
    level: debug
    add_source: true
```

### Field Descriptions

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `preset` | string | none | Apply predefined configuration |
| `level` | string | info | Minimum logging level |
| `format` | string | text | Output format |
| `include_file` | bool | false | Include file/line information |
| `include_time` | bool | true | Include timestamps |
| `use_short_file` | bool | true | Use short file paths |
| `use_slog` | bool | false | Use slog backend for performance |
| `static_fields` | map | {} | Fields included in every log entry |
| `redact_patterns` | []string | [] | Regex patterns for sensitive data |

## Presets

Presets provide optimized configurations for common scenarios:

### Available Presets

| Preset | Level | Format | File | Time | slog | Description |
|--------|-------|--------|------|------|------|-------------|
| `development` | debug | text | ✅ | ✅ | ❌ | Human-readable for local dev |
| `production` | info | json | ❌ | ✅ | ✅ | Optimized for production |
| `debug` | trace | text | ✅ | ✅ | ❌ | Maximum verbosity |
| `minimal` | info | text | ❌ | ❌ | ❌ | Bare minimum output |
| `structured` | info | json | ✅ | ✅ | ✅ | Rich structured logging |

### Using Presets

```yaml
# Start with a preset
preset: production

# Override specific settings
level: debug
static_fields:
  custom_field: "value"
```

Presets are applied first, then individual settings override preset values.

## Examples

### Development Configuration

```yaml
# configs/development.yaml
preset: development
static_fields:
  environment: development
  service: my-app
  version: "1.0.0-dev"
redact_patterns:
  - "password=\\w+"
  - "token=\\w+"
```

### Production Configuration

```yaml
# configs/production.yaml
preset: production
static_fields:
  environment: production
  service: my-app
  version: "1.0.0"
  region: "us-east-1"
redact_patterns:
  - "password=[^\\s]*"
  - "token=[^\\s]*"
  - "apikey=[^\\s]*"
  - "authorization:[^\\r\\n]*"
```

### File Logging

```yaml
# configs/file-logging.yaml
level: info
format: json
include_file: true
include_time: true
static_fields:
  service: my-app
  component: file-logger
output:
  type: file
  target: ~/logs/application.log
```

### Microservice Configuration

```yaml
# configs/microservice.yaml
preset: structured
static_fields:
  service: user-service
  version: "2.1.0"
  environment: production
  region: us-west-2
  team: platform
  component: api
redact_patterns:
  - "password=[^\\s&]*"
  - "token=[^\\s&]*"
  - "bearer [^\\s]*"
  - "jwt [^\\s]*"
  - "x-api-key:[^\\r\\n]*"
```

## Advanced Features

### Environment Variable Integration

YAML configuration works seamlessly with environment variables:

```bash
# Set config file via environment
export LOG_CONFIG_FILE=configs/production.yaml

# Use in application
logger := logging.NewFromYAMLEnv("LOG_CONFIG_FILE")
```

The logger falls back to `NewSimple()` if the environment variable is not set or the file cannot be loaded.

### File Output Features

When using `output.type: file`:

#### Automatic Directory Creation
```yaml
output:
  type: file
  target: ~/logs/deep/nested/app.log  # Creates all parent directories
```

#### Home Directory Expansion
```yaml
output:
  type: file
  target: ~/logs/app.log              # Expands to /home/user/logs/app.log
```

#### Relative Paths
```yaml
output:
  type: file
  target: logs/app.log                # Relative to current directory
```

### Security and Redaction

#### Built-in Security Patterns

Common security patterns are included in example configs:

```yaml
redact_patterns:
  - "password=[^\\s]*"               # password=secret123
  - "token=[^\\s]*"                  # token=abc123
  - "apikey=[^\\s]*"                 # apikey=xyz789
  - "authorization:[^\\r\\n]*"       # Authorization: Bearer token
  - "bearer [^\\s]*"                 # bearer token123
  - "jwt [^\\s]*"                    # jwt eyJ0eXAi...
```

#### Custom Patterns

Add your own redaction patterns:

```yaml
redact_patterns:
  - "secret_key=[^\\s]*"
  - "\\b\\d{4}[-\\s]?\\d{4}[-\\s]?\\d{4}[-\\s]?\\d{4}\\b"  # Credit cards
  - "\\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\\.[A-Z|a-z]{2,}\\b"  # Email addresses
```

### Static Fields

Static fields are included in every log entry:

```yaml
static_fields:
  service: "user-service"
  version: "1.2.3"
  environment: "production"
  region: "us-east-1"
  team: "platform"
  component: "api"
  deployment_id: "abc123"
```

These fields appear in all log messages automatically.

## Best Practices

### 1. Use Environment-Specific Configs

Create separate configuration files for each environment:

```
configs/
├── development.yaml
├── staging.yaml
├── production.yaml
└── testing.yaml
```

### 2. Start with Presets

Always start with a preset and customize:

```yaml
# Good
preset: production
static_fields:
  custom_field: "value"

# Less optimal
level: info
format: json
include_time: true
use_slog: true
# ... manually specify everything
```

### 3. Version Your Configurations

Keep configuration files in version control alongside your application code.

### 4. Security First

Always include redaction patterns:

```yaml
redact_patterns:
  - "password=[^\\s]*"
  - "token=[^\\s]*"
  - "apikey=[^\\s]*"
  # Add more as needed
```

### 5. Use Static Fields for Context

Include consistent fields across all logs:

```yaml
static_fields:
  service: "my-service"
  version: "1.0.0"
  environment: "production"
```

### 6. Test Configurations

Always test configuration files:

```go
func TestConfig(t *testing.T) {
    logger, err := logging.NewFromYAMLFile("configs/production.yaml")
    assert.NoError(t, err)
    assert.NotNil(t, logger)
    
    logger.Info("Test message")
}
```

### 7. Document Your Patterns

Comment your redaction patterns:

```yaml
redact_patterns:
  - "password=[^\\s]*"      # Basic password fields
  - "token=[^\\s]*"         # API tokens
  - "jwt [^\\s]*"           # JWT tokens
  - "ssn=\\d{3}-\\d{2}-\\d{4}"  # Social Security Numbers
```

## Troubleshooting

### Common Issues

#### 1. Configuration File Not Found

```
Error: failed to read YAML file config/logging.yaml: no such file or directory
```

**Solutions:**
- Check file path is correct
- Use absolute paths for clarity
- Verify file exists and is readable

#### 2. Invalid YAML Syntax

```
Error: failed to parse YAML configuration: yaml: line 5: mapping values are not allowed in this context
```

**Solutions:**
- Validate YAML syntax with online validator
- Check indentation (use spaces, not tabs)
- Ensure proper key-value formatting

#### 3. Invalid Configuration Values

```
Error: failed to configure core: invalid log level: invalid_level
```

**Solutions:**
- Check valid values in schema documentation
- Use one of: trace, debug, info, warn, error, critical
- Check preset names are correct

#### 4. File Permission Issues

```
Error: failed to create log directory /var/log/myapp: permission denied
```

**Solutions:**
- Ensure process has write permissions
- Create directories manually with correct permissions
- Use user-writable locations like `~/logs/`

### Debugging Configuration

#### Verify Configuration Loading

```go
logger, err := logging.NewFromYAMLFile("config/logging.yaml")
if err != nil {
    fmt.Printf("Configuration error: %v\n", err)
    return
}

logger.Info("Configuration loaded successfully")
```

#### Test with Simple String

```go
testYAML := `
preset: minimal
static_fields:
  test: true
`

logger, err := logging.LoadFromYAMLString(testYAML)
if err != nil {
    fmt.Printf("String config error: %v\n", err)
} else {
    logger.Info("String config works")
}
```

#### Enable Debug Logging

```yaml
# Temporary debug configuration
level: debug
format: text
include_file: true
```

### Validation Tools

#### YAML Syntax Validation

Use online tools:
- [YAML Lint](https://www.yamllint.com/)
- [Online YAML Parser](https://yaml-online-parser.appspot.com/)

#### Regex Pattern Testing

Test redaction patterns:
- [Regex101](https://regex101.com/)
- [RegExr](https://regexr.com/)

## Migration Guide

### From Environment Variables

**Before:**
```bash
export LOG_LEVEL=info
export LOG_FORMAT=json
```

**After:**
```yaml
# config/logging.yaml
level: info
format: json
```

```go
// Replace
logger := logging.NewFromEnvironment()

// With
logger, _ := logging.NewFromYAMLFile("config/logging.yaml")
```

### From Builder Pattern

**Before:**
```go
logger := logging.NewEasyBuilder().
    Level(logging.InfoLevel).
    JSON().
    WithFile().
    Field("service", "my-app").
    Build()
```

**After:**
```yaml
# config/logging.yaml
level: info
format: json
include_file: true
static_fields:
  service: my-app
```

```go
logger, _ := logging.NewFromYAMLFile("config/logging.yaml")
```

## Integration Examples

### Docker

```dockerfile
# Dockerfile
COPY configs/production.yaml /app/config/logging.yaml
ENV LOG_CONFIG_FILE=/app/config/logging.yaml
```

```go
// Application
logger := logging.NewFromYAMLEnv("LOG_CONFIG_FILE")
```

### Kubernetes

```yaml
# ConfigMap
apiVersion: v1
kind: ConfigMap
metadata:
  name: logging-config
data:
  logging.yaml: |
    preset: production
    static_fields:
      service: my-app
      cluster: prod-cluster
    redact_patterns:
      - "password=[^\\s]*"
```

```yaml
# Deployment
spec:
  containers:
  - name: app
    env:
    - name: LOG_CONFIG_FILE
      value: /etc/config/logging.yaml
    volumeMounts:
    - name: config
      mountPath: /etc/config
  volumes:
  - name: config
    configMap:
      name: logging-config
```

### Helm Charts

```yaml
# values.yaml
logging:
  level: info
  format: json
  staticFields:
    service: my-app
    version: "1.0.0"

# templates/configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: {{ .Release.Name }}-logging
data:
  logging.yaml: |
    level: {{ .Values.logging.level }}
    format: {{ .Values.logging.format }}
    static_fields:
      {{- toYaml .Values.logging.staticFields | nindent 6 }}
```

This comprehensive YAML configuration system makes complex logging setups simple and maintainable!