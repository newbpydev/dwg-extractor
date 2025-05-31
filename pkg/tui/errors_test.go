package tui

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestErrorCategories tests different error categories
func TestErrorCategories(t *testing.T) {
	tests := []struct {
		name         string
		err          error
		expectedType ErrorType
		expectedMsg  string
	}{
		{
			name:         "User error with validation failure",
			err:          NewUserError("invalid layer name", "Layer names cannot be empty"),
			expectedType: ErrorTypeUser,
			expectedMsg:  "Layer names cannot be empty",
		},
		{
			name:         "System error with file access",
			err:          NewSystemError("file read failed", errors.New("permission denied")),
			expectedType: ErrorTypeSystem,
			expectedMsg:  "file read failed",
		},
		{
			name:         "Network error with connection timeout",
			err:          NewNetworkError("connection failed", "timeout after 30s"),
			expectedType: ErrorTypeNetwork,
			expectedMsg:  "connection failed",
		},
		{
			name:         "Conversion error with ODA failure",
			err:          NewConversionError("DWG conversion failed", "invalid DWG format"),
			expectedType: ErrorTypeConversion,
			expectedMsg:  "DWG conversion failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This should fail initially - we need to implement AppError interface
			appErr, ok := tt.err.(AppError)
			assert.True(t, ok, "Expected error to implement AppError interface")

			assert.Equal(t, tt.expectedType, appErr.Type(), "Expected correct error type")
			assert.Contains(t, appErr.Error(), tt.expectedMsg, "Expected error message to contain expected text")
			assert.NotEmpty(t, appErr.UserMessage(), "Expected user-friendly message")
		})
	}
}

// TestErrorRecoverySuggestions tests error recovery suggestions
func TestErrorRecoverySuggestions(t *testing.T) {
	tests := []struct {
		name               string
		err                AppError
		expectedSuggestion string
	}{
		{
			name:               "File permission error suggests chmod",
			err:                NewSystemError("file access denied", errors.New("permission denied")),
			expectedSuggestion: "Check file permissions",
		},
		{
			name:               "Network error suggests retry",
			err:                NewNetworkError("connection timeout", "network unreachable"),
			expectedSuggestion: "Check network connection",
		},
		{
			name:               "Conversion error suggests file format check",
			err:                NewConversionError("unsupported format", "not a valid DWG file"),
			expectedSuggestion: "Verify DWG file format",
		},
		{
			name:               "User error suggests input validation",
			err:                NewUserError("invalid input", "empty filename"),
			expectedSuggestion: "Provide a valid filename",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This should fail initially - we need to implement RecoverySuggestion method
			suggestion := tt.err.RecoverySuggestion()
			assert.Contains(t, suggestion, tt.expectedSuggestion, "Expected appropriate recovery suggestion")
		})
	}
}

// TestErrorDisplayInTUI tests error display in the TUI
func TestErrorDisplayInTUI(t *testing.T) {
	tests := []struct {
		name            string
		err             AppError
		displayType     ErrorDisplayType
		expectedVisible bool
		expectedText    string
	}{
		{
			name:            "Critical error shows modal dialog",
			err:             NewSystemError("critical failure", errors.New("system crash")),
			displayType:     ErrorDisplayModal,
			expectedVisible: true,
			expectedText:    "critical failure",
		},
		{
			name:            "Warning shows status bar",
			err:             NewUserError("minor issue", "layer not found"),
			displayType:     ErrorDisplayStatusBar,
			expectedVisible: true,
			expectedText:    "layer not found",
		},
		{
			name:            "Info message shows temporarily",
			err:             NewUserError("info", "operation completed with warnings"),
			displayType:     ErrorDisplayTemporary,
			expectedVisible: true,
			expectedText:    "operation completed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := SetupTestApp(t)
			view := NewDXFView(app)

			// This should fail initially - we need to implement ErrorHandler
			errorHandler := NewErrorHandler(view)
			assert.NotNil(t, errorHandler, "Expected error handler to be created")

			// Display the error
			errorHandler.DisplayError(tt.err, tt.displayType)

			// Check if error is visible
			isVisible := errorHandler.IsErrorVisible()
			assert.Equal(t, tt.expectedVisible, isVisible, "Expected correct error visibility")

			if tt.expectedVisible {
				displayedText := errorHandler.GetDisplayedError()
				assert.Contains(t, displayedText, tt.expectedText, "Expected error text to be displayed")
			}
		})
	}
}

// TestErrorLogging tests error logging functionality
func TestErrorLogging(t *testing.T) {
	tests := []struct {
		name         string
		err          AppError
		logLevel     LogLevel
		expectLogged bool
	}{
		{
			name:         "Critical error logs at error level",
			err:          NewSystemError("system failure", errors.New("critical error")),
			logLevel:     LogLevelError,
			expectLogged: true,
		},
		{
			name:         "Warning logs at warn level",
			err:          NewUserError("validation failed", "invalid input"),
			logLevel:     LogLevelWarn,
			expectLogged: true,
		},
		{
			name:         "Info logs at info level",
			err:          NewUserError("operation info", "process completed"),
			logLevel:     LogLevelInfo,
			expectLogged: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This should fail initially - we need to implement ErrorLogger
			logger := NewErrorLogger()
			assert.NotNil(t, logger, "Expected error logger to be created")

			// Log the error
			logger.LogError(tt.err, tt.logLevel)

			// Check if error was logged
			if tt.expectLogged {
				entries := logger.GetLogEntries()
				assert.NotEmpty(t, entries, "Expected log entries to be recorded")

				lastEntry := entries[len(entries)-1]
				assert.Equal(t, tt.logLevel, lastEntry.Level, "Expected correct log level")
				assert.Contains(t, lastEntry.Message, tt.err.Error(), "Expected error message in log")
			}
		})
	}
}

// TestErrorHandlerIntegration tests integration between error handler and TUI
func TestErrorHandlerIntegration(t *testing.T) {
	tests := []struct {
		name            string
		operation       func(*DXFView) error
		expectedError   bool
		expectedDisplay bool
	}{
		{
			name: "File load error displays appropriately",
			operation: func(view *DXFView) error {
				// Simulate file load error
				return NewSystemError("file load failed", errors.New("file not found"))
			},
			expectedError:   true,
			expectedDisplay: true,
		},
		{
			name: "Successful operation shows no error",
			operation: func(view *DXFView) error {
				// Simulate successful operation
				return nil
			},
			expectedError:   false,
			expectedDisplay: false,
		},
		{
			name: "Conversion error shows with recovery suggestion",
			operation: func(view *DXFView) error {
				return NewConversionError("conversion failed", "invalid DWG version")
			},
			expectedError:   true,
			expectedDisplay: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := SetupTestApp(t)
			view := NewDXFView(app)

			// This should fail initially - we need to implement error integration
			errorHandler := NewErrorHandler(view)

			// Execute operation
			err := tt.operation(view)

			if tt.expectedError {
				assert.Error(t, err, "Expected operation to return error")

				// Handle the error
				if appErr, ok := err.(AppError); ok {
					errorHandler.HandleError(appErr)

					if tt.expectedDisplay {
						assert.True(t, errorHandler.IsErrorVisible(), "Expected error to be displayed")
						displayText := errorHandler.GetDisplayedError()
						assert.NotEmpty(t, displayText, "Expected error display text")
					}
				}
			} else {
				assert.NoError(t, err, "Expected operation to succeed")
				assert.False(t, errorHandler.IsErrorVisible(), "Expected no error display")
			}
		})
	}
}

// TestErrorRecoveryFlow tests error recovery workflows
func TestErrorRecoveryFlow(t *testing.T) {
	tests := []struct {
		name           string
		initialError   AppError
		recoveryAction func() error
		expectedResult string
	}{
		{
			name:         "File permission error recovery",
			initialError: NewSystemError("permission denied", errors.New("access forbidden")),
			recoveryAction: func() error {
				// Simulate permission fix
				return nil
			},
			expectedResult: "success",
		},
		{
			name:         "Network error with retry",
			initialError: NewNetworkError("connection failed", "timeout"),
			recoveryAction: func() error {
				// Simulate successful retry
				return nil
			},
			expectedResult: "success",
		},
		{
			name:         "Conversion error with format fix",
			initialError: NewConversionError("invalid format", "corrupted file"),
			recoveryAction: func() error {
				// Simulate format correction
				return NewUserError("file still invalid", "manual intervention required")
			},
			expectedResult: "retry_needed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := SetupTestApp(t)
			view := NewDXFView(app)

			// This should fail initially - we need to implement recovery workflows
			errorHandler := NewErrorHandler(view)
			recoveryManager := NewRecoveryManager()

			// Initial error
			errorHandler.HandleError(tt.initialError)
			assert.True(t, errorHandler.IsErrorVisible(), "Expected initial error to be visible")

			// Attempt recovery
			recoveryErr := tt.recoveryAction()
			result := recoveryManager.AttemptRecovery(tt.initialError, recoveryErr)

			switch tt.expectedResult {
			case "success":
				assert.Equal(t, RecoverySuccess, result, "Expected successful recovery")
			case "retry_needed":
				assert.Equal(t, RecoveryRetryNeeded, result, "Expected retry needed")
			}
		})
	}
}
