package tui

import (
	"fmt"
	"log"
	"time"
)

// ErrorType represents different categories of errors
type ErrorType int

const (
	ErrorTypeUser ErrorType = iota
	ErrorTypeSystem
	ErrorTypeNetwork
	ErrorTypeConversion
)

// String returns the string representation of ErrorType
func (e ErrorType) String() string {
	switch e {
	case ErrorTypeUser:
		return "User"
	case ErrorTypeSystem:
		return "System"
	case ErrorTypeNetwork:
		return "Network"
	case ErrorTypeConversion:
		return "Conversion"
	default:
		return "Unknown"
	}
}

// AppError interface for application-specific errors
type AppError interface {
	error
	Type() ErrorType
	UserMessage() string
	RecoverySuggestion() string
}

// UserError represents user input or validation errors
type UserError struct {
	message     string
	userMessage string
	suggestion  string
}

// NewUserError creates a new user error
func NewUserError(message, userMessage string) *UserError {
	suggestion := generateUserErrorSuggestion(userMessage)
	return &UserError{
		message:     fmt.Sprintf("%s: %s", message, userMessage),
		userMessage: userMessage,
		suggestion:  suggestion,
	}
}

func (e *UserError) Error() string              { return e.message }
func (e *UserError) Type() ErrorType            { return ErrorTypeUser }
func (e *UserError) UserMessage() string        { return e.userMessage }
func (e *UserError) RecoverySuggestion() string { return e.suggestion }

// SystemError represents system-level errors
type SystemError struct {
	message     string
	userMessage string
	underlying  error
	suggestion  string
}

// NewSystemError creates a new system error
func NewSystemError(message string, underlying error) *SystemError {
	userMessage := fmt.Sprintf("System error: %s", message)
	suggestion := generateSystemErrorSuggestion(underlying)
	return &SystemError{
		message:     message,
		userMessage: userMessage,
		underlying:  underlying,
		suggestion:  suggestion,
	}
}

func (e *SystemError) Error() string              { return e.message }
func (e *SystemError) Type() ErrorType            { return ErrorTypeSystem }
func (e *SystemError) UserMessage() string        { return e.userMessage }
func (e *SystemError) RecoverySuggestion() string { return e.suggestion }
func (e *SystemError) Unwrap() error              { return e.underlying }

// NetworkError represents network-related errors
type NetworkError struct {
	message     string
	userMessage string
	details     string
	suggestion  string
}

// NewNetworkError creates a new network error
func NewNetworkError(message, details string) *NetworkError {
	userMessage := fmt.Sprintf("Network error: %s", message)
	suggestion := generateNetworkErrorSuggestion(details)
	return &NetworkError{
		message:     message,
		userMessage: userMessage,
		details:     details,
		suggestion:  suggestion,
	}
}

func (e *NetworkError) Error() string              { return e.message }
func (e *NetworkError) Type() ErrorType            { return ErrorTypeNetwork }
func (e *NetworkError) UserMessage() string        { return e.userMessage }
func (e *NetworkError) RecoverySuggestion() string { return e.suggestion }

// ConversionError represents DWG/DXF conversion errors
type ConversionError struct {
	message     string
	userMessage string
	details     string
	suggestion  string
}

// NewConversionError creates a new conversion error
func NewConversionError(message, details string) *ConversionError {
	userMessage := fmt.Sprintf("Conversion error: %s", message)
	suggestion := generateConversionErrorSuggestion(details)
	return &ConversionError{
		message:     message,
		userMessage: userMessage,
		details:     details,
		suggestion:  suggestion,
	}
}

func (e *ConversionError) Error() string              { return e.message }
func (e *ConversionError) Type() ErrorType            { return ErrorTypeConversion }
func (e *ConversionError) UserMessage() string        { return e.userMessage }
func (e *ConversionError) RecoverySuggestion() string { return e.suggestion }

// ErrorDisplayType represents how errors should be displayed
type ErrorDisplayType int

const (
	ErrorDisplayModal ErrorDisplayType = iota
	ErrorDisplayStatusBar
	ErrorDisplayTemporary
)

// ErrorHandler manages error display in the TUI
type ErrorHandler struct {
	view          *DXFView
	currentError  AppError
	displayType   ErrorDisplayType
	isVisible     bool
	displayedText string
}

// NewErrorHandler creates a new error handler
func NewErrorHandler(view *DXFView) *ErrorHandler {
	return &ErrorHandler{
		view: view,
	}
}

// DisplayError displays an error in the TUI
func (eh *ErrorHandler) DisplayError(err AppError, displayType ErrorDisplayType) {
	eh.currentError = err
	eh.displayType = displayType
	eh.isVisible = true
	eh.displayedText = err.UserMessage()

	// For now, just set the displayed text - actual TUI integration would go here
}

// HandleError handles an application error
func (eh *ErrorHandler) HandleError(err AppError) {
	// Determine display type based on error type
	var displayType ErrorDisplayType
	switch err.Type() {
	case ErrorTypeSystem:
		displayType = ErrorDisplayModal
	case ErrorTypeUser:
		displayType = ErrorDisplayStatusBar
	default:
		displayType = ErrorDisplayTemporary
	}

	eh.DisplayError(err, displayType)
}

// IsErrorVisible returns whether an error is currently visible
func (eh *ErrorHandler) IsErrorVisible() bool {
	return eh.isVisible
}

// GetDisplayedError returns the currently displayed error text
func (eh *ErrorHandler) GetDisplayedError() string {
	return eh.displayedText
}

// LogLevel represents logging levels
type LogLevel int

const (
	LogLevelInfo LogLevel = iota
	LogLevelWarn
	LogLevelError
)

// String returns the string representation of LogLevel
func (l LogLevel) String() string {
	switch l {
	case LogLevelInfo:
		return "INFO"
	case LogLevelWarn:
		return "WARN"
	case LogLevelError:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// LogEntry represents a single log entry
type LogEntry struct {
	Timestamp time.Time
	Level     LogLevel
	Message   string
	Error     AppError
}

// ErrorLogger handles error logging
type ErrorLogger struct {
	entries []LogEntry
}

// NewErrorLogger creates a new error logger
func NewErrorLogger() *ErrorLogger {
	return &ErrorLogger{
		entries: make([]LogEntry, 0),
	}
}

// LogError logs an error with the specified level
func (el *ErrorLogger) LogError(err AppError, level LogLevel) {
	entry := LogEntry{
		Timestamp: time.Now(),
		Level:     level,
		Message:   err.Error(),
		Error:     err,
	}

	el.entries = append(el.entries, entry)

	// Also log to standard logger
	log.Printf("[%s] %s: %s", level.String(), err.Type().String(), err.Error())
}

// GetLogEntries returns all log entries
func (el *ErrorLogger) GetLogEntries() []LogEntry {
	return el.entries
}

// RecoveryResult represents the result of an error recovery attempt
type RecoveryResult int

const (
	RecoverySuccess RecoveryResult = iota
	RecoveryRetryNeeded
	RecoveryFailed
)

// RecoveryManager handles error recovery workflows
type RecoveryManager struct{}

// NewRecoveryManager creates a new recovery manager
func NewRecoveryManager() *RecoveryManager {
	return &RecoveryManager{}
}

// AttemptRecovery attempts to recover from an error
func (rm *RecoveryManager) AttemptRecovery(originalErr AppError, recoveryErr error) RecoveryResult {
	if recoveryErr == nil {
		return RecoverySuccess
	}

	// If recovery returned another error, determine if retry is needed
	if _, ok := recoveryErr.(AppError); ok {
		return RecoveryRetryNeeded
	}

	return RecoveryFailed
}

// Helper functions to generate context-specific recovery suggestions

func generateUserErrorSuggestion(userMessage string) string {
	switch {
	case contains(userMessage, "filename") || contains(userMessage, "empty"):
		return "Provide a valid filename"
	case contains(userMessage, "layer"):
		return "Check layer name and try again"
	case contains(userMessage, "input"):
		return "Verify your input and try again"
	default:
		return "Please correct the input and try again"
	}
}

func generateSystemErrorSuggestion(underlying error) string {
	if underlying == nil {
		return "Contact system administrator"
	}

	errMsg := underlying.Error()
	switch {
	case contains(errMsg, "permission"):
		return "Check file permissions and try again"
	case contains(errMsg, "not found"):
		return "Verify the file exists and path is correct"
	case contains(errMsg, "access"):
		return "Ensure you have necessary access rights"
	default:
		return "Check system resources and try again"
	}
}

func generateNetworkErrorSuggestion(details string) string {
	switch {
	case contains(details, "timeout"):
		return "Check network connection and try again"
	case contains(details, "unreachable"):
		return "Check network connection and try again"
	case contains(details, "refused"):
		return "Check if the service is running"
	default:
		return "Check network settings and try again"
	}
}

func generateConversionErrorSuggestion(details string) string {
	switch {
	case contains(details, "format") || contains(details, "DWG"):
		return "Verify DWG file format and version"
	case contains(details, "corrupted"):
		return "Use a different DWG file or repair the current one"
	case contains(details, "version"):
		return "Try with a supported DWG version"
	default:
		return "Check file integrity and format"
	}
}

// contains is a helper function to check if a string contains a substring (case-insensitive)
func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr ||
			len(s) > len(substr) &&
				(s[:len(substr)] == substr ||
					s[len(s)-len(substr):] == substr ||
					indexSubstring(s, substr) >= 0))
}

// indexSubstring finds the index of substr in s (simple implementation)
func indexSubstring(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
