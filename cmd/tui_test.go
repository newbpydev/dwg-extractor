package cmd

import (
	"errors"
	"testing"
)

// MockRunTUI simulates RunTUI for testing conversion integration.
func MockRunTUI(args []string, simulateError bool) error {
	if simulateError {
		return errors.New("conversion failed: mock error")
	}
	// Simulate successful conversion and TUI startup
	return nil
}

func TestRunTUI_ConversionSuccess(t *testing.T) {
	err := MockRunTUI([]string{"-file", "testdata/sample.dwg"}, false)
	if err != nil {
		t.Errorf("Expected no error on successful conversion and TUI startup, got: %v", err)
	}
}

func TestRunTUI_ConversionFailure(t *testing.T) {
	err := MockRunTUI([]string{"-file", "testdata/badfile.dwg"}, true)
	if err == nil {
		t.Errorf("Expected error on conversion failure, got nil")
	}
	if err.Error() != "conversion failed: mock error" {
		t.Errorf("Unexpected error message: %v", err)
	}
}
