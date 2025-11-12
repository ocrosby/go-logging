package logging

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadFromYAMLString(t *testing.T) {
	tests := []struct {
		name     string
		yaml     string
		wantErr  bool
		checkLog bool
	}{
		{
			name: "minimal config",
			yaml: `
level: info
format: text
output:
  type: stdout
`,
			wantErr:  false,
			checkLog: true,
		},
		{
			name: "json config with static fields",
			yaml: `
level: debug
format: json
include_file: true
include_time: true
static_fields:
  service: test-service
  version: 1.0.0
output:
  type: stdout
`,
			wantErr:  false,
			checkLog: true,
		},
		{
			name: "development preset",
			yaml: `
preset: development
static_fields:
  app: test-app
`,
			wantErr:  false,
			checkLog: true,
		},
		{
			name: "production preset",
			yaml: `
preset: production
static_fields:
  environment: test
`,
			wantErr:  false,
			checkLog: true,
		},
		{
			name: "with redaction patterns",
			yaml: `
level: info
format: json
redact_patterns:
  - "password=[^\\s]*"
  - "token=[^\\s]*"
output:
  type: stdout
`,
			wantErr:  false,
			checkLog: true,
		},
		{
			name: "invalid level",
			yaml: `
level: invalid
format: text
output:
  type: stdout
`,
			wantErr:  true,
			checkLog: false,
		},
		{
			name: "invalid format",
			yaml: `
level: info
format: invalid
output:
  type: stdout
`,
			wantErr:  true,
			checkLog: false,
		},
		{
			name: "invalid output type",
			yaml: `
level: info
format: text
output:
  type: invalid
`,
			wantErr:  true,
			checkLog: false,
		},
		{
			name: "file output without target",
			yaml: `
level: info
format: text
output:
  type: file
`,
			wantErr:  true,
			checkLog: false,
		},
		{
			name: "invalid regex pattern",
			yaml: `
level: info
format: text
redact_patterns:
  - "[invalid regex"
output:
  type: stdout
`,
			wantErr:  true,
			checkLog: false,
		},
		{
			name: "unknown preset",
			yaml: `
preset: unknown
`,
			wantErr:  true,
			checkLog: false,
		},
		{
			name: "malformed yaml",
			yaml: `
level: info
format: [invalid yaml structure
`,
			wantErr:  true,
			checkLog: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger, err := LoadFromYAMLString(tt.yaml)

			if tt.wantErr {
				if err == nil {
					t.Errorf("LoadFromYAMLString() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("LoadFromYAMLString() unexpected error = %v", err)
				return
			}

			if logger == nil {
				t.Errorf("LoadFromYAMLString() returned nil logger")
				return
			}

			if tt.checkLog {
				// Test that the logger can actually log
				logger.Info("Test message")
			}
		})
	}
}

func TestLoadFromYAMLFile(t *testing.T) {
	// Create temporary directory for test files
	tmpDir := t.TempDir()

	// Test cases
	tests := []struct {
		name     string
		filename string
		yaml     string
		wantErr  bool
	}{
		{
			name:     "valid config file",
			filename: "valid.yaml",
			yaml: `
level: info
format: json
output:
  type: stdout
static_fields:
  test: true
`,
			wantErr: false,
		},
		{
			name:     "file not found",
			filename: "nonexistent.yaml",
			yaml:     "",
			wantErr:  true,
		},
		{
			name:     "invalid yaml content",
			filename: "invalid.yaml",
			yaml:     "invalid: [yaml content",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var filePath string

			if tt.name != "file not found" {
				// Create test file
				filePath = filepath.Join(tmpDir, tt.filename)
				err := os.WriteFile(filePath, []byte(tt.yaml), 0644)
				if err != nil {
					t.Fatalf("Failed to create test file: %v", err)
				}
			} else {
				// Non-existent file
				filePath = filepath.Join(tmpDir, tt.filename)
			}

			logger, err := LoadFromYAML(filePath)

			if tt.wantErr {
				if err == nil {
					t.Errorf("LoadFromYAML() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("LoadFromYAML() unexpected error = %v", err)
				return
			}

			if logger == nil {
				t.Errorf("LoadFromYAML() returned nil logger")
				return
			}

			// Test that the logger can actually log
			logger.Info("Test message from file config")
		})
	}
}

func TestYAMLFileOutput(t *testing.T) {
	// Create temporary directory for test log file
	tmpDir := t.TempDir()
	logFile := filepath.Join(tmpDir, "test.log")

	yaml := `
level: info
format: json
include_time: true
static_fields:
  test: file-output
output:
  type: file
  target: ` + logFile

	logger, err := LoadFromYAMLString(yaml)
	if err != nil {
		t.Fatalf("LoadFromYAMLString() error = %v", err)
	}

	// Log a message
	logger.Info("Test message to file")

	// Verify file was created and contains content
	if _, statErr := os.Stat(logFile); os.IsNotExist(statErr) {
		t.Errorf("Log file was not created")
		return
	}

	content, err := os.ReadFile(logFile)
	if err != nil {
		t.Errorf("Failed to read log file: %v", err)
		return
	}

	if len(content) == 0 {
		t.Errorf("Log file is empty")
		return
	}

	// Check that the content looks like JSON
	contentStr := string(content)
	if !strings.Contains(contentStr, "Test message to file") {
		t.Errorf("Log file does not contain expected message")
	}
}

func TestNewFromYAMLEnv(t *testing.T) {
	// Test with non-existent environment variable
	t.Run("env var not set", func(t *testing.T) {
		logger := NewFromYAMLEnv("NON_EXISTENT_CONFIG_VAR")
		if logger == nil {
			t.Errorf("NewFromYAMLEnv() returned nil logger")
		}

		// Should fall back to simple logger and work
		logger.Info("Test message with fallback")
	})

	// Test with valid config file
	t.Run("env var with valid file", func(t *testing.T) {
		tmpDir := t.TempDir()
		configFile := filepath.Join(tmpDir, "test-config.yaml")

		yaml := `
level: debug
format: text
static_fields:
  env_test: true
output:
  type: stdout
`
		err := os.WriteFile(configFile, []byte(yaml), 0644)
		if err != nil {
			t.Fatalf("Failed to create test config file: %v", err)
		}

		// Set environment variable
		envVar := "TEST_LOG_CONFIG"
		t.Setenv(envVar, configFile)

		logger := NewFromYAMLEnv(envVar)
		if logger == nil {
			t.Errorf("NewFromYAMLEnv() returned nil logger")
			return
		}

		logger.Debug("Test debug message from env config")
		logger.Info("Test info message from env config")
	})

	// Test with invalid config file (should fall back to simple)
	t.Run("env var with invalid file", func(t *testing.T) {
		tmpDir := t.TempDir()
		configFile := filepath.Join(tmpDir, "invalid-config.yaml")

		// Create invalid YAML file
		err := os.WriteFile(configFile, []byte("invalid: [yaml"), 0644)
		if err != nil {
			t.Fatalf("Failed to create invalid config file: %v", err)
		}

		envVar := "TEST_LOG_CONFIG_INVALID"
		t.Setenv(envVar, configFile)

		logger := NewFromYAMLEnv(envVar)
		if logger == nil {
			t.Errorf("NewFromYAMLEnv() returned nil logger")
			return
		}

		// Should work with fallback
		logger.Info("Test message with invalid config fallback")
	})
}

func TestPresets(t *testing.T) {
	presets := []string{"development", "production", "debug", "minimal", "structured"}

	for _, preset := range presets {
		t.Run(preset, func(t *testing.T) {
			yaml := `preset: ` + preset

			logger, err := LoadFromYAMLString(yaml)
			if err != nil {
				t.Errorf("LoadFromYAMLString() with preset %s error = %v", preset, err)
				return
			}

			if logger == nil {
				t.Errorf("LoadFromYAMLString() with preset %s returned nil logger", preset)
				return
			}

			// Test logging at different levels to verify preset behavior
			logger.Info("Info message with preset " + preset)
			logger.Debug("Debug message with preset " + preset)
			logger.Error("Error message with preset " + preset)
		})
	}
}

func TestYAMLConfigWithStructuredLogging(t *testing.T) {
	yaml := `
level: info
format: json
include_time: true
static_fields:
  service: test-service
  version: 1.0.0
output:
  type: stdout
`

	logger, err := LoadFromYAMLString(yaml)
	if err != nil {
		t.Fatalf("LoadFromYAMLString() error = %v", err)
	}

	// Test structured logging
	logger.Info("User action",
		"user_id", 12345,
		"action", "login",
		"success", true,
	)

	logger.Error("Database error",
		"error", "connection timeout",
		"database", "postgres",
		"retry_count", 3,
	)
}

func TestRedactionPatterns(t *testing.T) {
	yaml := `
level: info
format: text
redact_patterns:
  - "password=[^\\s]*"
  - "token=[^\\s]*"
output:
  type: stdout
`

	logger, err := LoadFromYAMLString(yaml)
	if err != nil {
		t.Fatalf("LoadFromYAMLString() error = %v", err)
	}

	// These should be redacted
	logger.Info("Login attempt with password=secret123 and token=abc456")
	logger.Error("API call failed with apikey=xyz789")
}
