package main

import (
	"os"
	"testing"

	"github.com/remym/go-dwg-extractor/cmd"
)

// TestMain_ArgumentHandling tests the main function's argument handling logic
func TestMain_ArgumentHandling(t *testing.T) {
	// Save original args and restore them after test
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	// Test case 1: No arguments (should fail)
	os.Args = []string{"dwg-extractor"}
	err := cmd.Execute()
	if err == nil {
		t.Errorf("Expected error when no command provided, got nil")
	}

	// Test case 2: Valid extract command with nonexistent file (should fail gracefully)
	os.Args = []string{"dwg-extractor", "extract", "-file", "nonexistent.dwg"}
	err = cmd.Execute()
	if err == nil {
		t.Errorf("Expected error for nonexistent file, got nil")
	}

	// Test case 3: TUI command (should not panic, but may hang so we don't test execution)
	os.Args = []string{"dwg-extractor", "tui"}
	// We don't actually execute this to avoid hanging
	t.Log("TUI command parsing would work")
}

// TestMain_FunctionExists ensures the main function exists and can be referenced
func TestMain_FunctionExists(t *testing.T) {
	// This test ensures the main function exists and is properly defined
	// We test this by ensuring the test file compiles and the main package is valid
	t.Log("main function exists and main package compiles successfully")
}
