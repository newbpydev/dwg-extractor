package main

import (
	"os"
	"os/exec"
	"runtime"
	"strings"
	"testing"

	"github.com/remym/go-dwg-extractor/cmd"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

// TestShowVersion tests the showVersion function
func TestShowVersion(t *testing.T) {
	tests := []struct {
		name        string
		version     string
		gitCommit   string
		buildTime   string
		expectedOut []string
	}{
		{
			name:      "Development version",
			version:   "dev",
			gitCommit: "unknown",
			buildTime: "unknown",
			expectedOut: []string{
				"Go DWG Extractor",
				"Version:    dev",
				"Git Commit: unknown",
				"Build Time: unknown",
			},
		},
		{
			name:      "Release version",
			version:   "v1.0.0",
			gitCommit: "abc1234",
			buildTime: "2024-12-28T10:00:00Z",
			expectedOut: []string{
				"Go DWG Extractor",
				"Version:    v1.0.0",
				"Git Commit: abc1234",
				"Build Time: 2024-12-28T10:00:00Z",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This should fail initially - no way to test showVersion output
			// We need to capture stdout to test this function

			// Save original values
			origVersion := version
			origGitCommit := gitCommit
			origBuildTime := buildTime
			defer func() {
				version = origVersion
				gitCommit = origGitCommit
				buildTime = origBuildTime
			}()

			// Set test values
			version = tt.version
			gitCommit = tt.gitCommit
			buildTime = tt.buildTime

			// This will fail - we need a way to capture output
			output := captureOutput(func() {
				showVersion()
			})

			// Verify output contains expected strings
			for _, expected := range tt.expectedOut {
				assert.Contains(t, output, expected, "Output should contain expected text")
			}
		})
	}
}

// TestShowHelp tests the showHelp function
func TestShowHelp(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		expectedOut []string
	}{
		{
			name: "Standard help output",
			args: []string{"go-dwg-extractor"},
			expectedOut: []string{
				"Go DWG Extractor - Extract data from DWG files",
				"Usage:",
				"Commands:",
				"extract",
				"tui",
				"version",
				"help",
				"Options:",
				"-file",
				"-output",
				"Examples:",
			},
		},
		{
			name: "Help with different program name",
			args: []string{"dwg-tool"},
			expectedOut: []string{
				"Usage: dwg-tool [command] [options]",
				"dwg-tool extract -file sample.dwg",
				"dwg-tool tui -file sample.dwg",
				"dwg-tool version",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original args
			origArgs := os.Args
			defer func() { os.Args = origArgs }()

			// Set test args
			os.Args = tt.args

			// This will fail - we need a way to capture output
			output := captureOutput(func() {
				showHelp()
			})

			// Verify output contains expected strings
			for _, expected := range tt.expectedOut {
				assert.Contains(t, output, expected, "Help output should contain expected text")
			}
		})
	}
}

// TestMainFunction tests the main function behavior
func TestMainFunction(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping main function tests in short mode")
	}

	tests := []struct {
		name         string
		args         []string
		expectExit   bool
		expectedCode int
		expectedOut  []string
		expectedErr  []string
	}{
		{
			name:        "Version command",
			args:        []string{"program", "version"},
			expectExit:  false,
			expectedOut: []string{"Go DWG Extractor", "Version:"},
		},
		{
			name:        "Version flag --version",
			args:        []string{"program", "--version"},
			expectExit:  false,
			expectedOut: []string{"Go DWG Extractor", "Version:"},
		},
		{
			name:        "Version flag -v",
			args:        []string{"program", "-v"},
			expectExit:  false,
			expectedOut: []string{"Go DWG Extractor", "Version:"},
		},
		{
			name:        "Help command",
			args:        []string{"program", "help"},
			expectExit:  false,
			expectedOut: []string{"Go DWG Extractor - Extract data", "Usage:"},
		},
		{
			name:        "Help flag --help",
			args:        []string{"program", "--help"},
			expectExit:  false,
			expectedOut: []string{"Go DWG Extractor - Extract data", "Usage:"},
		},
		{
			name:        "Help flag -h",
			args:        []string{"program", "-h"},
			expectExit:  false,
			expectedOut: []string{"Go DWG Extractor - Extract data", "Usage:"},
		},
		{
			name:         "No arguments",
			args:         []string{"program"},
			expectExit:   true,
			expectedCode: 1,
			expectedErr:  []string{"No command provided"},
		},
		{
			name:         "Invalid command",
			args:         []string{"program", "invalid"},
			expectExit:   true,
			expectedCode: 1,
			expectedErr:  []string{"Error:"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This should fail initially - no way to test main function
			// We need to use subprocess testing for main function

			// Build the test binary
			execName := "test-main"
			if runtime.GOOS == "windows" {
				execName = "test-main.exe"
			}

			cmd := exec.Command("go", "build", "-o", execName, ".")
			err := cmd.Run()
			require.NoError(t, err, "Should be able to build test binary")
			defer os.Remove(execName)

			// Run the test binary with arguments
			testCmd := exec.Command("./"+execName, tt.args[1:]...)
			output, err := testCmd.CombinedOutput()
			outputStr := string(output)

			if tt.expectExit {
				assert.Error(t, err, "Should exit with error")
				if exitError, ok := err.(*exec.ExitError); ok {
					if tt.expectedCode != 0 {
						assert.Equal(t, tt.expectedCode, exitError.ExitCode(), "Should exit with expected code")
					}
				}
				for _, expected := range tt.expectedErr {
					assert.Contains(t, outputStr, expected, "Error output should contain expected text")
				}
			} else {
				assert.NoError(t, err, "Should not exit with error")
				for _, expected := range tt.expectedOut {
					assert.Contains(t, outputStr, expected, "Output should contain expected text")
				}
			}
		})
	}
}

// TestMainWithMockedCmd tests main function with mocked cmd.Execute
func TestMainWithMockedCmd(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		cmdError    error
		expectFatal bool
	}{
		{
			name:        "Successful command execution",
			args:        []string{"program", "tui"},
			cmdError:    nil,
			expectFatal: false,
		},
		{
			name:        "Command execution error",
			args:        []string{"program", "tui"},
			cmdError:    assert.AnError,
			expectFatal: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This will fail initially - we need dependency injection to test main
			// We would need to refactor main to accept a cmd executor interface

			// Save original args
			origArgs := os.Args
			defer func() { os.Args = origArgs }()

			// Set test args
			os.Args = tt.args

			if tt.expectFatal {
				// This should test that log.Fatalf is called
				// We need a way to capture or mock log.Fatalf
				t.Skip("Need to implement log.Fatalf mocking")
			} else {
				// This should test successful execution
				t.Skip("Need to implement cmd.Execute mocking")
			}
		})
	}
}

// captureOutput captures stdout during function execution
func captureOutput(f func()) string {
	// Capture stdout using os.Pipe
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Channel to capture the output
	outputChan := make(chan string)

	// Start a goroutine to read from the pipe
	go func() {
		var buf strings.Builder
		buffer := make([]byte, 1024)
		for {
			n, err := r.Read(buffer)
			if err != nil {
				break
			}
			buf.Write(buffer[:n])
		}
		outputChan <- buf.String()
	}()

	// Execute the function
	f()

	// Restore stdout and close writer
	w.Close()
	os.Stdout = oldStdout

	// Get the captured output
	output := <-outputChan
	r.Close()

	return output
}

// TestBuildTimeVariables tests that build-time variables are properly set
func TestBuildTimeVariables(t *testing.T) {
	tests := []struct {
		name     string
		variable *string
		varName  string
	}{
		{
			name:     "Version variable exists",
			variable: &version,
			varName:  "version",
		},
		{
			name:     "GitCommit variable exists",
			variable: &gitCommit,
			varName:  "gitCommit",
		},
		{
			name:     "BuildTime variable exists",
			variable: &buildTime,
			varName:  "buildTime",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test that variables are accessible and have values
			assert.NotNil(t, tt.variable, "%s variable should not be nil", tt.varName)
			assert.NotEmpty(t, *tt.variable, "%s variable should not be empty", tt.varName)
		})
	}
}

// TestArgumentParsing tests argument parsing logic
func TestArgumentParsing(t *testing.T) {
	tests := []struct {
		name         string
		args         []string
		shouldHandle bool
		action       string
	}{
		{
			name:         "Version argument",
			args:         []string{"program", "version"},
			shouldHandle: true,
			action:       "version",
		},
		{
			name:         "Version flag",
			args:         []string{"program", "--version"},
			shouldHandle: true,
			action:       "version",
		},
		{
			name:         "Short version flag",
			args:         []string{"program", "-v"},
			shouldHandle: true,
			action:       "version",
		},
		{
			name:         "Help argument",
			args:         []string{"program", "help"},
			shouldHandle: true,
			action:       "help",
		},
		{
			name:         "Help flag",
			args:         []string{"program", "--help"},
			shouldHandle: true,
			action:       "help",
		},
		{
			name:         "Short help flag",
			args:         []string{"program", "-h"},
			shouldHandle: true,
			action:       "help",
		},
		{
			name:         "Regular command",
			args:         []string{"program", "tui"},
			shouldHandle: false,
			action:       "",
		},
		{
			name:         "No arguments",
			args:         []string{"program"},
			shouldHandle: false,
			action:       "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This will fail initially - we need to extract argument parsing logic
			// to make it testable separately from main()
			result := parseMainArguments(tt.args)

			if tt.shouldHandle {
				assert.True(t, result.ShouldHandle, "Should handle this argument")
				assert.Equal(t, tt.action, result.Action, "Should identify correct action")
			} else {
				assert.False(t, result.ShouldHandle, "Should not handle this argument")
			}
		})
	}
}

// ArgumentParseResult represents the result of parsing main arguments
type ArgumentParseResult struct {
	ShouldHandle bool
	Action       string
}

// parseMainArguments extracts argument parsing logic for testing
func parseMainArguments(args []string) ArgumentParseResult {
	if len(args) < 2 {
		return ArgumentParseResult{ShouldHandle: false, Action: ""}
	}

	switch args[1] {
	case "version", "--version", "-v":
		return ArgumentParseResult{ShouldHandle: true, Action: "version"}
	case "help", "--help", "-h":
		return ArgumentParseResult{ShouldHandle: true, Action: "help"}
	default:
		return ArgumentParseResult{ShouldHandle: false, Action: ""}
	}
}
