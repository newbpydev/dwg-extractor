package tui

import (
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/remym/go-dwg-extractor/pkg/data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockClipboardManager is a mock for testing clipboard integration
type MockClipboardManager struct {
	mock.Mock
}

// CopyToClipboard mocks the clipboard copy operation
func (m *MockClipboardManager) CopyToClipboard(text string) error {
	args := m.Called(text)
	return args.Error(0)
}

// TestClipboardIntegration_KeyBindings tests clipboard key bindings
func TestClipboardIntegration_KeyBindings(t *testing.T) {
	tests := []struct {
		name         string
		key          tcell.Key
		modifiers    tcell.ModMask
		expectedCopy bool
		setupData    func() *data.ExtractedData
	}{
		{
			name:         "Ctrl+C copies selected items",
			key:          tcell.KeyCtrlC,
			modifiers:    tcell.ModCtrl,
			expectedCopy: true,
			setupData:    createTestDataWithMultipleItems,
		},
		{
			name:         "C key copies selected items (alternative)",
			key:          tcell.KeyRune,
			modifiers:    tcell.ModNone,
			expectedCopy: true,
			setupData:    createTestDataWithMultipleItems,
		},
		{
			name:         "Other keys don't trigger copy",
			key:          tcell.KeyEnter,
			modifiers:    tcell.ModNone,
			expectedCopy: false,
			setupData:    createTestDataWithMultipleItems,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := SetupTestApp(t)
			view := NewDXFView(app)

			testData := tt.setupData()
			view.Update(testData)

			// Setup mock clipboard
			mockClipboard := new(MockClipboardManager)
			if tt.expectedCopy {
				mockClipboard.On("CopyToClipboard", mock.AnythingOfType("string")).Return(nil)
			}

			// This should fail initially - we need to implement clipboard integration
			clipboardHandler := NewClipboardHandler(view, mockClipboard)
			assert.NotNil(t, clipboardHandler, "Expected clipboard handler to be created")

			// Simulate key press
			handled := clipboardHandler.HandleKeyPress(tt.key, tt.modifiers)

			if tt.expectedCopy {
				assert.True(t, handled, "Expected copy key to be handled")
				mockClipboard.AssertExpectations(t)
			} else {
				assert.False(t, handled, "Expected non-copy key to not be handled")
			}
		})
	}
}

// TestClipboardIntegration_CopySelectedItems tests copying selected items
func TestClipboardIntegration_CopySelectedItems(t *testing.T) {
	tests := []struct {
		name            string
		selectedItems   []int
		expectedContent []string
		expectError     bool
	}{
		{
			name:            "Copy single selected item",
			selectedItems:   []int{0},
			expectedContent: []string{"Line:"},
			expectError:     false,
		},
		{
			name:            "Copy multiple selected items",
			selectedItems:   []int{0, 1, 2},
			expectedContent: []string{"Line:", "Circle:"},
			expectError:     false,
		},
		{
			name:            "Copy with no selection",
			selectedItems:   []int{},
			expectedContent: []string{},
			expectError:     false,
		},
		{
			name:            "Copy with invalid selection index",
			selectedItems:   []int{999},
			expectedContent: []string{},
			expectError:     false, // Should handle gracefully
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := SetupTestApp(t)
			view := NewDXFView(app)

			testData := createTestDataWithMultipleItems()
			view.Update(testData)

			// Setup mock clipboard
			mockClipboard := new(MockClipboardManager)
			if len(tt.expectedContent) > 0 {
				mockClipboard.On("CopyToClipboard", mock.MatchedBy(func(content string) bool {
					// Check that content contains expected strings
					for _, expected := range tt.expectedContent {
						if !assert.Contains(t, content, expected) {
							return false
						}
					}
					return true
				})).Return(nil)
			}

			// This should fail initially - we need to implement clipboard integration
			clipboardHandler := NewClipboardHandler(view, mockClipboard)

			// Simulate selection
			for _, index := range tt.selectedItems {
				clipboardHandler.AddToSelection(index)
			}

			// Perform copy operation
			err := clipboardHandler.CopySelectedItems()

			if tt.expectError {
				assert.Error(t, err, "Expected error for invalid copy operation")
			} else {
				assert.NoError(t, err, "Expected no error for valid copy operation")
				if len(tt.expectedContent) > 0 {
					mockClipboard.AssertExpectations(t)
				}
			}
		})
	}
}

// TestClipboardIntegration_StatusMessages tests status message display
func TestClipboardIntegration_StatusMessages(t *testing.T) {
	tests := []struct {
		name            string
		operation       string
		expectedMessage string
		shouldDisplay   bool
	}{
		{
			name:            "Successful copy shows success message",
			operation:       "copy",
			expectedMessage: "copied to clipboard",
			shouldDisplay:   true,
		},
		{
			name:            "Failed copy shows error message",
			operation:       "copy_error",
			expectedMessage: "Failed to copy",
			shouldDisplay:   true,
		},
		{
			name:            "No operation shows no message",
			operation:       "none",
			expectedMessage: "",
			shouldDisplay:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := SetupTestApp(t)
			view := NewDXFView(app)

			testData := createTestDataWithMultipleItems()
			view.Update(testData)

			// This should fail initially - we need to implement status message display
			statusHandler := NewStatusMessageHandler(view)
			assert.NotNil(t, statusHandler, "Expected status handler to be created")

			// Simulate operation
			switch tt.operation {
			case "copy":
				statusHandler.ShowCopySuccess(3) // 3 items copied
			case "copy_error":
				statusHandler.ShowCopyError("clipboard access denied")
			case "none":
				// No operation
			}

			// Check if message is displayed
			if tt.shouldDisplay {
				message := statusHandler.GetCurrentMessage()
				assert.Contains(t, message, tt.expectedMessage, "Expected status message to contain expected text")
			} else {
				message := statusHandler.GetCurrentMessage()
				assert.Empty(t, message, "Expected no status message")
			}
		})
	}
}

// TestClipboardIntegration_ErrorHandling tests clipboard error scenarios
func TestClipboardIntegration_ErrorHandling(t *testing.T) {
	tests := []struct {
		name           string
		clipboardError string
		expectedError  bool
		expectedStatus string
	}{
		{
			name:           "Clipboard access denied",
			clipboardError: "permission denied",
			expectedError:  true,
			expectedStatus: "permission denied",
		},
		{
			name:           "Clipboard unavailable",
			clipboardError: "no clipboard available",
			expectedError:  true,
			expectedStatus: "clipboard unavailable",
		},
		{
			name:           "Generic clipboard error",
			clipboardError: "system error",
			expectedError:  true,
			expectedStatus: "system error",
		},
		{
			name:           "Successful operation",
			clipboardError: "",
			expectedError:  false,
			expectedStatus: "success",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := SetupTestApp(t)
			view := NewDXFView(app)

			testData := createTestDataWithMultipleItems()
			view.Update(testData)

			// Setup mock clipboard with error
			mockClipboard := new(MockClipboardManager)
			if tt.clipboardError != "" {
				mockClipboard.On("CopyToClipboard", mock.AnythingOfType("string")).
					Return(assert.AnError)
			} else {
				mockClipboard.On("CopyToClipboard", mock.AnythingOfType("string")).
					Return(nil)
			}

			// This should fail initially - we need to implement error handling
			clipboardHandler := NewClipboardHandler(view, mockClipboard)
			clipboardHandler.AddToSelection(0) // Select first item

			err := clipboardHandler.CopySelectedItems()

			if tt.expectedError {
				assert.Error(t, err, "Expected error for clipboard failure")
			} else {
				assert.NoError(t, err, "Expected no error for successful operation")
			}

			mockClipboard.AssertExpectations(t)
		})
	}
}

// TestClipboardIntegration_FormatOptions tests different clipboard formats
func TestClipboardIntegration_FormatOptions(t *testing.T) {
	tests := []struct {
		name           string
		format         string
		expectedFormat string
	}{
		{
			name:           "Plain text format",
			format:         "text",
			expectedFormat: "Line:",
		},
		{
			name:           "CSV format",
			format:         "csv",
			expectedFormat: "Type,Layer,Details",
		},
		{
			name:           "JSON format",
			format:         "json",
			expectedFormat: "[",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := SetupTestApp(t)
			view := NewDXFView(app)

			testData := createTestDataWithMultipleItems()
			view.Update(testData)

			// Setup mock clipboard
			mockClipboard := new(MockClipboardManager)
			mockClipboard.On("CopyToClipboard", mock.MatchedBy(func(content string) bool {
				return assert.Contains(t, content, tt.expectedFormat)
			})).Return(nil)

			// This should fail initially - we need to implement format options
			clipboardHandler := NewClipboardHandler(view, mockClipboard)
			clipboardHandler.AddToSelection(0)
			clipboardHandler.SetFormat(tt.format)

			err := clipboardHandler.CopySelectedItems()
			assert.NoError(t, err, "Expected no error for formatted copy")

			mockClipboard.AssertExpectations(t)
		})
	}
}
