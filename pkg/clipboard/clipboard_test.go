package clipboard

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockClipboard is a mock implementation for testing
type MockClipboard struct {
	mock.Mock
}

// WriteAll mocks the clipboard write operation
func (m *MockClipboard) WriteAll(text string) error {
	args := m.Called(text)
	return args.Error(0)
}

// TestCopyToClipboard tests the core clipboard functionality
func TestCopyToClipboard(t *testing.T) {
	tests := []struct {
		name          string
		text          string
		mockSetup     func(*MockClipboard)
		expectedError bool
		errorMessage  string
	}{
		{
			name: "Successful copy to clipboard",
			text: "Test content to copy",
			mockSetup: func(m *MockClipboard) {
				m.On("WriteAll", "Test content to copy").Return(nil)
			},
			expectedError: false,
		},
		{
			name: "Empty string copy",
			text: "",
			mockSetup: func(m *MockClipboard) {
				m.On("WriteAll", "").Return(nil)
			},
			expectedError: false,
		},
		{
			name: "Large text copy",
			text: string(make([]byte, 10000)), // Large text
			mockSetup: func(m *MockClipboard) {
				m.On("WriteAll", mock.AnythingOfType("string")).Return(nil)
			},
			expectedError: false,
		},
		{
			name: "Clipboard write failure",
			text: "Test content",
			mockSetup: func(m *MockClipboard) {
				m.On("WriteAll", "Test content").Return(errors.New("clipboard access denied"))
			},
			expectedError: true,
			errorMessage:  "clipboard access denied",
		},
		{
			name: "Unicode text copy",
			text: "Test with unicode: 你好世界 🌍",
			mockSetup: func(m *MockClipboard) {
				m.On("WriteAll", "Test with unicode: 你好世界 🌍").Return(nil)
			},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			mockClipboard := new(MockClipboard)
			tt.mockSetup(mockClipboard)

			// Create clipboard manager with mock
			manager := NewClipboardManager(mockClipboard)

			// Test the copy operation - this should fail initially
			err := manager.CopyToClipboard(tt.text)

			if tt.expectedError {
				assert.Error(t, err, "Expected error for test case")
				if tt.errorMessage != "" {
					assert.Contains(t, err.Error(), tt.errorMessage, "Expected specific error message")
				}
			} else {
				assert.NoError(t, err, "Expected no error for successful copy")
			}

			// Verify mock expectations
			mockClipboard.AssertExpectations(t)
		})
	}
}

// TestClipboardManager tests the clipboard manager structure
func TestClipboardManager(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(t *testing.T)
	}{
		{
			name: "NewClipboardManager creates valid manager",
			testFunc: func(t *testing.T) {
				mockClipboard := new(MockClipboard)
				manager := NewClipboardManager(mockClipboard)

				assert.NotNil(t, manager, "Expected manager to be created")
				// This should fail initially - we need to implement ClipboardManager
			},
		},
		{
			name: "NewClipboardManager with nil writer",
			testFunc: func(t *testing.T) {
				manager := NewClipboardManager(nil)

				err := manager.CopyToClipboard("test")
				assert.Error(t, err, "Expected error when clipboard writer is nil")
				assert.Contains(t, err.Error(), "clipboard writer is nil", "Expected specific error message")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunc(t)
		})
	}
}

// TestRealClipboardIntegration tests with the actual clipboard library
func TestRealClipboardIntegration(t *testing.T) {
	// Skip in CI/headless environments
	if testing.Short() {
		t.Skip("Skipping real clipboard test in short mode")
	}

	tests := []struct {
		name     string
		text     string
		testFunc func(t *testing.T, text string)
	}{
		{
			name: "Real clipboard write and verify",
			text: "Integration test content",
			testFunc: func(t *testing.T, text string) {
				// Create real clipboard manager - this should fail initially
				manager := NewRealClipboardManager()

				err := manager.CopyToClipboard(text)
				assert.NoError(t, err, "Expected successful copy to real clipboard")

				// Verify content was written (if possible in test environment)
				// Note: Reading from clipboard might not work in all test environments
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunc(t, tt.text)
		})
	}
}

// TestClipboardErrorHandling tests comprehensive error scenarios
func TestClipboardErrorHandling(t *testing.T) {
	tests := []struct {
		name        string
		setupError  error
		expectedErr string
	}{
		{
			name:        "Permission denied error",
			setupError:  errors.New("permission denied: clipboard access restricted"),
			expectedErr: "permission denied",
		},
		{
			name:        "System clipboard unavailable",
			setupError:  errors.New("no clipboard available"),
			expectedErr: "no clipboard available",
		},
		{
			name:        "Generic system error",
			setupError:  errors.New("system error occurred"),
			expectedErr: "system error occurred",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClipboard := new(MockClipboard)
			mockClipboard.On("WriteAll", mock.AnythingOfType("string")).Return(tt.setupError)

			manager := NewClipboardManager(mockClipboard)
			err := manager.CopyToClipboard("test content")

			assert.Error(t, err, "Expected error from clipboard operation")
			assert.Contains(t, err.Error(), tt.expectedErr, "Expected specific error message")

			mockClipboard.AssertExpectations(t)
		})
	}
}

// TestClipboardSizeHandling tests handling of various text sizes
func TestClipboardSizeHandling(t *testing.T) {
	tests := []struct {
		name     string
		textSize int
		expected bool
	}{
		{
			name:     "Small text (100 bytes)",
			textSize: 100,
			expected: true,
		},
		{
			name:     "Medium text (10KB)",
			textSize: 10000,
			expected: true,
		},
		{
			name:     "Large text (100KB)",
			textSize: 100000,
			expected: true,
		},
		{
			name:     "Very large text (1MB)",
			textSize: 1000000,
			expected: true, // Should handle large text gracefully
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClipboard := new(MockClipboard)
			text := string(make([]byte, tt.textSize))

			if tt.expected {
				mockClipboard.On("WriteAll", mock.AnythingOfType("string")).Return(nil)
			} else {
				mockClipboard.On("WriteAll", mock.AnythingOfType("string")).Return(errors.New("text too large"))
			}

			manager := NewClipboardManager(mockClipboard)
			err := manager.CopyToClipboard(text)

			if tt.expected {
				assert.NoError(t, err, "Expected successful copy for size %d", tt.textSize)
			} else {
				assert.Error(t, err, "Expected error for size %d", tt.textSize)
			}

			mockClipboard.AssertExpectations(t)
		})
	}
}
