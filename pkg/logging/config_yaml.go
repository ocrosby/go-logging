package logging

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

// Constants for YAML configuration
const (
	stdoutString = "stdout"
	stderrString = "stderr"
	fileString   = "file"
	infoString   = "info"
)

// YAMLConfig represents the complete YAML configuration structure.
type YAMLConfig struct {
	// Core logging configuration
	Level        string                 `yaml:"level"`
	StaticFields map[string]interface{} `yaml:"static_fields,omitempty"`

	// Formatting configuration
	Format       string   `yaml:"format"`
	IncludeFile  bool     `yaml:"include_file"`
	IncludeTime  bool     `yaml:"include_time"`
	UseShortFile bool     `yaml:"use_short_file"`
	RedactList   []string `yaml:"redact_patterns,omitempty"`

	// Output configuration
	Output YAMLOutputConfig `yaml:"output"`

	// Slog configuration
	UseSlog bool            `yaml:"use_slog"`
	Slog    *YAMLSlogConfig `yaml:"slog,omitempty"`

	// Presets for common configurations
	Preset string `yaml:"preset,omitempty"`
}

// YAMLOutputConfig represents output configuration in YAML.
type YAMLOutputConfig struct {
	Type   string `yaml:"type"`             // "stdout", "stderr", "file"
	Target string `yaml:"target,omitempty"` // file path for type "file"
}

// YAMLSlogConfig represents slog-specific configuration in YAML.
type YAMLSlogConfig struct {
	HandlerType string                 `yaml:"handler_type"` // "text", "json"
	Options     map[string]interface{} `yaml:"options,omitempty"`
}

// LoadFromYAML loads configuration from a YAML file.
func LoadFromYAML(filename string) (Logger, error) {
	// Expand user home directory if needed
	if strings.HasPrefix(filename, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get user home directory: %w", err)
		}
		filename = filepath.Join(home, filename[2:])
	}

	// Make relative paths relative to current directory
	if !filepath.IsAbs(filename) {
		wd, err := os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("failed to get working directory: %w", err)
		}
		filename = filepath.Join(wd, filename)
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read YAML file %s: %w", filename, err)
	}

	return LoadFromYAMLData(data)
}

// LoadFromYAMLData loads configuration from YAML data bytes.
func LoadFromYAMLData(data []byte) (Logger, error) {
	var yamlConfig YAMLConfig
	if err := yaml.Unmarshal(data, &yamlConfig); err != nil {
		return nil, fmt.Errorf("failed to parse YAML configuration: %w", err)
	}

	return buildLoggerFromYAML(&yamlConfig)
}

// LoadFromYAMLString loads configuration from a YAML string.
func LoadFromYAMLString(yamlStr string) (Logger, error) {
	return LoadFromYAMLData([]byte(yamlStr))
}

// buildLoggerFromYAML builds a logger from the parsed YAML configuration.
func buildLoggerFromYAML(yamlConfig *YAMLConfig) (Logger, error) {
	// Apply preset if specified
	if yamlConfig.Preset != "" {
		if err := applyPreset(yamlConfig, yamlConfig.Preset); err != nil {
			return nil, fmt.Errorf("failed to apply preset '%s': %w", yamlConfig.Preset, err)
		}
	}

	// Build LoggerConfig
	builder := NewLoggerConfig()

	// Core configuration
	if err := configureCoreFromYAML(builder, yamlConfig); err != nil {
		return nil, fmt.Errorf("failed to configure core: %w", err)
	}

	// Formatter configuration
	if err := configureFormatterFromYAML(builder, yamlConfig); err != nil {
		return nil, fmt.Errorf("failed to configure formatter: %w", err)
	}

	// Output configuration
	if err := configureOutputFromYAML(builder, yamlConfig); err != nil {
		return nil, fmt.Errorf("failed to configure output: %w", err)
	}

	// Slog configuration
	if yamlConfig.UseSlog {
		builder.UseSlog(true)
		// TODO: Add custom slog handler configuration support
	}

	config := builder.Build()
	redactorChain := ProvideRedactorChainFromLoggerConfig(config)
	return NewUnifiedLogger(config, redactorChain), nil
}

// configureCoreFromYAML configures core settings from YAML.
func configureCoreFromYAML(builder *LoggerConfigBuilder, yamlConfig *YAMLConfig) error {
	// Set log level
	if yamlConfig.Level != "" {
		level, ok := ParseLevel(yamlConfig.Level)
		if !ok {
			return fmt.Errorf("invalid log level: %s", yamlConfig.Level)
		}
		builder.WithLevel(level)
	}

	// Set static fields
	if len(yamlConfig.StaticFields) > 0 {
		for k, v := range yamlConfig.StaticFields {
			builder.config.Core.StaticFields[k] = v
		}
	}

	return nil
}

// configureFormatterFromYAML configures formatter settings from YAML.
func configureFormatterFromYAML(builder *LoggerConfigBuilder, yamlConfig *YAMLConfig) error {
	// Set format
	switch strings.ToLower(yamlConfig.Format) {
	case jsonFormatString:
		builder.WithJSONFormat()
	case textFormatString, "":
		builder.WithTextFormat()
	default:
		return fmt.Errorf("invalid format: %s (must be 'json' or 'text')", yamlConfig.Format)
	}

	// Set formatter options
	builder.config.Formatter.IncludeFile = yamlConfig.IncludeFile
	builder.config.Formatter.IncludeTime = yamlConfig.IncludeTime
	builder.config.Formatter.UseShortFile = yamlConfig.UseShortFile

	// Add redact patterns
	for _, pattern := range yamlConfig.RedactList {
		if re, err := regexp.Compile(pattern); err == nil {
			builder.config.Formatter.RedactPatterns = append(builder.config.Formatter.RedactPatterns, re)
		} else {
			return fmt.Errorf("invalid redact pattern '%s': %w", pattern, err)
		}
	}

	return nil
}

// configureOutputFromYAML configures output settings from YAML.
func configureOutputFromYAML(builder *LoggerConfigBuilder, yamlConfig *YAMLConfig) error {
	outputType := strings.ToLower(yamlConfig.Output.Type)

	switch outputType {
	case stdoutString, "":
		builder.WithWriter(os.Stdout)
	case stderrString:
		builder.WithWriter(os.Stderr)
	case fileString:
		writer, err := createFileWriter(yamlConfig.Output.Target)
		if err != nil {
			return err
		}
		builder.WithWriter(writer)
	default:
		return fmt.Errorf("invalid output type: %s (must be '%s', '%s', or '%s')", yamlConfig.Output.Type, stdoutString, stderrString, fileString)
	}

	return nil
}

// createFileWriter creates a file writer with proper path handling.
func createFileWriter(target string) (io.Writer, error) {
	if target == "" {
		return nil, fmt.Errorf("file output requires target path")
	}

	// Expand user home directory if needed
	expandedTarget, err := expandHomePath(target)
	if err != nil {
		return nil, err
	}

	// Create directory if it doesn't exist
	if dirErr := createLogDirectory(expandedTarget); dirErr != nil {
		return nil, dirErr
	}

	// Open file for writing
	file, err := os.OpenFile(expandedTarget, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file %s: %w", expandedTarget, err)
	}

	return file, nil
}

// expandHomePath expands ~ to user home directory.
func expandHomePath(target string) (string, error) {
	if strings.HasPrefix(target, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get user home directory: %w", err)
		}
		return filepath.Join(home, target[2:]), nil
	}
	return target, nil
}

// createLogDirectory creates the log directory if it doesn't exist.
func createLogDirectory(target string) error {
	dir := filepath.Dir(target)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory %s: %w", dir, err)
	}
	return nil
}

// applyPreset applies a predefined configuration preset.
func applyPreset(yamlConfig *YAMLConfig, preset string) error {
	switch strings.ToLower(preset) {
	case "development", "dev":
		applyDevelopmentPreset(yamlConfig)
	case "production", "prod":
		applyProductionPreset(yamlConfig)
	case "debug":
		applyDebugPreset(yamlConfig)
	case "minimal":
		applyMinimalPreset(yamlConfig)
	case "structured":
		applyStructuredPreset(yamlConfig)
	default:
		return fmt.Errorf("unknown preset: %s", preset)
	}
	return nil
}

// applyDevelopmentPreset applies the development preset.
func applyDevelopmentPreset(yamlConfig *YAMLConfig) {
	if yamlConfig.Level == "" {
		yamlConfig.Level = "debug"
	}
	if yamlConfig.Format == "" {
		yamlConfig.Format = textFormatString
	}
	if yamlConfig.Output.Type == "" {
		yamlConfig.Output.Type = stdoutString
	}
	yamlConfig.IncludeFile = true
	yamlConfig.IncludeTime = true
	yamlConfig.UseShortFile = true
}

// applyProductionPreset applies the production preset.
func applyProductionPreset(yamlConfig *YAMLConfig) {
	if yamlConfig.Level == "" {
		yamlConfig.Level = infoString
	}
	if yamlConfig.Format == "" {
		yamlConfig.Format = jsonFormatString
	}
	if yamlConfig.Output.Type == "" {
		yamlConfig.Output.Type = stdoutString
	}
	yamlConfig.IncludeFile = false
	yamlConfig.IncludeTime = true
	yamlConfig.UseShortFile = true
	yamlConfig.UseSlog = true
}

// applyDebugPreset applies the debug preset.
func applyDebugPreset(yamlConfig *YAMLConfig) {
	if yamlConfig.Level == "" {
		yamlConfig.Level = "trace"
	}
	if yamlConfig.Format == "" {
		yamlConfig.Format = textFormatString
	}
	if yamlConfig.Output.Type == "" {
		yamlConfig.Output.Type = stdoutString
	}
	yamlConfig.IncludeFile = true
	yamlConfig.IncludeTime = true
	yamlConfig.UseShortFile = false
}

// applyMinimalPreset applies the minimal preset.
func applyMinimalPreset(yamlConfig *YAMLConfig) {
	if yamlConfig.Level == "" {
		yamlConfig.Level = infoString
	}
	if yamlConfig.Format == "" {
		yamlConfig.Format = textFormatString
	}
	if yamlConfig.Output.Type == "" {
		yamlConfig.Output.Type = stdoutString
	}
	yamlConfig.IncludeFile = false
	yamlConfig.IncludeTime = false
	yamlConfig.UseShortFile = true
}

// applyStructuredPreset applies the structured preset.
func applyStructuredPreset(yamlConfig *YAMLConfig) {
	if yamlConfig.Level == "" {
		yamlConfig.Level = infoString
	}
	if yamlConfig.Format == "" {
		yamlConfig.Format = jsonFormatString
	}
	if yamlConfig.Output.Type == "" {
		yamlConfig.Output.Type = stdoutString
	}
	yamlConfig.IncludeFile = true
	yamlConfig.IncludeTime = true
	yamlConfig.UseShortFile = true
	yamlConfig.UseSlog = true
}

// NewFromYAMLFile is a convenience function to create a logger from a YAML file.
// This provides a simple factory function similar to the other New* functions.
func NewFromYAMLFile(filename string) (Logger, error) {
	return LoadFromYAML(filename)
}

// NewFromYAMLEnv creates a logger from a YAML file specified by an environment variable.
// If the environment variable is not set, returns a simple logger with default settings.
func NewFromYAMLEnv(envVar string) Logger {
	filename := os.Getenv(envVar)
	if filename == "" {
		return NewSimple()
	}

	logger, err := LoadFromYAML(filename)
	if err != nil {
		// Fall back to simple logger if YAML loading fails
		return NewSimple()
	}

	return logger
}

// SaveToYAML saves the current logger configuration to a YAML file.
// This is useful for generating configuration templates.
func SaveToYAML(config *LoggerConfig, filename string) error {
	yamlConfig := &YAMLConfig{
		Level:        config.Core.Level.String(),
		StaticFields: config.Core.StaticFields,
		IncludeFile:  config.Formatter.IncludeFile,
		IncludeTime:  config.Formatter.IncludeTime,
		UseShortFile: config.Formatter.UseShortFile,
		UseSlog:      config.UseSlog,
	}

	// Set format
	if config.Formatter.Format == JSONFormat {
		yamlConfig.Format = jsonFormatString
	} else {
		yamlConfig.Format = textFormatString
	}

	// Set output (simplified - only supports stdout/stderr detection)
	yamlConfig.Output.Type = stdoutString // Default assumption

	// Convert redact patterns to strings (simplified)
	for _, pattern := range config.Formatter.RedactPatterns {
		yamlConfig.RedactList = append(yamlConfig.RedactList, pattern.String())
	}

	data, err := yaml.Marshal(yamlConfig)
	if err != nil {
		return fmt.Errorf("failed to marshal YAML: %w", err)
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	if err := os.WriteFile(filename, data, 0600); err != nil {
		return fmt.Errorf("failed to write YAML file %s: %w", filename, err)
	}

	return nil
}
