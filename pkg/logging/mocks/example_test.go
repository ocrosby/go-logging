package mocks_test

import (
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/ocrosby/go-logging/pkg/logging"
	"github.com/ocrosby/go-logging/pkg/logging/mocks"
)

func TestLoggerWithMock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mocks.NewMockLogger(ctrl)

	mockLogger.EXPECT().
		Info("test message").
		Times(1)

	mockLogger.EXPECT().
		WithField("user_id", 123).
		Return(mockLogger)

	mockLogger.EXPECT().
		IsLevelEnabled(logging.DebugLevel).
		Return(true)

	mockLogger.Info("test message")
	logger := mockLogger.WithField("user_id", 123)
	enabled := mockLogger.IsLevelEnabled(logging.DebugLevel)

	if logger == nil {
		t.Error("Expected logger to be returned")
	}

	if !enabled {
		t.Error("Expected debug level to be enabled")
	}
}

func TestRedactorWithMock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRedactor := mocks.NewMockRedactor(ctrl)

	input := "secret data: password123"
	expected := "secret data: ***REDACTED***"

	mockRedactor.EXPECT().
		Redact(input).
		Return(expected).
		Times(1)

	result := mockRedactor.Redact(input)

	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func TestRedactorChainWithMock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockChain := mocks.NewMockRedactorChainInterface(ctrl)
	mockRedactor := mocks.NewMockRedactor(ctrl)

	mockChain.EXPECT().
		AddRedactor(mockRedactor).
		Times(1)

	input := "sensitive data"
	expected := "***REDACTED***"

	mockChain.EXPECT().
		Redact(input).
		Return(expected).
		Times(1)

	mockChain.AddRedactor(mockRedactor)
	result := mockChain.Redact(input)

	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}
