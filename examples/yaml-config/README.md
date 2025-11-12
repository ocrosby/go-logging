# YAML Configuration Example

This example demonstrates how to use YAML configuration files to control logging setup declaratively, providing the most powerful and flexible configuration method available in the logging library.

## Overview

YAML configuration allows you to:
- Define complex logging setups declaratively
- Use presets for common configurations
- Configure multiple environments consistently
- Include security redaction patterns
- Set up file output with automatic directory creation

## Running the Example

```bash
cd examples/yaml-config
go run main.go
```

## What This Example Demonstrates

### 1. Loading from Configuration Files
Shows how to load configurations from the predefined YAML files in the `configs/` directory.

### 2. Embedded YAML Configuration
Demonstrates loading configuration from a YAML string embedded in code.

### 3. Environment-Based Configuration
Shows how to use environment variables to specify configuration files with automatic fallback.

### 4. Structured Logging with YAML
Demonstrates how YAML-configured static fields automatically appear in all log entries.

### 5. Preset Usage
Shows how different presets (`minimal`, `debug`) change logging behavior.

## Example Output

The example produces different output formats based on the configuration:

**Development Config** (text format):
```
2025/11/11 21:44:43 unified_logger.go:105: [DEBUG] Development logger loaded successfully
```

**Production Config** (JSON format with static fields):
```json
{"time":"2025-11-11T21:44:43.386947-05:00","level":"INFO","msg":"Production logger loaded successfully","environment":"production","service":"my-app","version":"1.0.0","region":"us-east-1"}
```

**Embedded Config** (JSON with file info and custom static fields):
```json
{"application":"yaml-example","environment":"embedded","file":"/path/to/file.go:195","level":"INFO","message":"Embedded logger created successfully","timestamp":"2025-11-12T02:44:43Z"}
```

## Configuration Files Used

The example references configuration files from the `configs/` directory:

### `configs/development.yaml`
- **Preset**: development  
- **Level**: debug
- **Format**: text
- **Features**: File info, human-readable output

### `configs/production.yaml`
- **Preset**: production
- **Level**: info
- **Format**: JSON
- **Features**: slog backend, structured output, security patterns

### `configs/microservice.yaml`
- **Preset**: structured
- **Features**: Rich static fields, comprehensive security patterns

## Key Learning Points

### 1. Presets Simplify Configuration
```yaml
preset: production  # Instantly configures for production use
```

### 2. Static Fields Are Powerful
```yaml
static_fields:
  service: my-app
  version: "1.0.0"
  environment: production
```
These fields appear in every log entry automatically.

### 3. Environment Integration
```go
logger := logging.NewFromYAMLEnv("LOG_CONFIG_FILE")
```
Falls back gracefully if environment variable not set.

### 4. Security by Default
All configuration files include redaction patterns for sensitive data.

## Try It Yourself

### 1. Modify Configurations
Edit the files in `configs/` and re-run the example to see different output.

### 2. Set Environment Variables
```bash
export LOG_CONFIG_FILE=$PWD/../../configs/development.yaml
go run main.go
```

### 3. Create Custom Configuration
Create your own YAML file:

```yaml
# my-config.yaml
preset: debug
static_fields:
  experiment: custom-config
  user: $(whoami)
redact_patterns:
  - "secret=[^\\s]*"
```

## Best Practices Demonstrated

1. **Use presets as starting points**
2. **Include static fields for context**
3. **Always include security redaction patterns**
4. **Use environment variables for deployment flexibility**
5. **Test different configurations during development**

## Related Files

- `/configs/` - Example YAML configuration files
- `/docs/YAML_CONFIGURATION.md` - Comprehensive YAML configuration guide
- `/pkg/logging/config_yaml.go` - YAML configuration implementation
- `/pkg/logging/config_yaml_test.go` - YAML configuration tests