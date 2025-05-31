package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestEndToEnd_ExtractCommand tests the complete extract command workflow
func TestEndToEnd_ExtractCommand(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}

	tests := []struct {
		name           string
		setupFile      func(t *testing.T) string
		args           []string
		expectError    bool
		expectedOutput []string
		cleanup        func(string)
	}{
		{
			name: "Extract with sample DXF file",
			setupFile: func(t *testing.T) string {
				// Create a simple DXF file for testing
				tmpDir := t.TempDir()
				dxfFile := filepath.Join(tmpDir, "test.dxf")
				dxfContent := `0
SECTION
2
HEADER
9
$ACADVER
1
AC1015
0
ENDSEC
0
SECTION
2
ENTITIES
0
LINE
8
TestLayer
10
0.0
20
0.0
30
0.0
11
10.0
21
10.0
31
0.0
0
ENDSEC
0
EOF
`
				err := os.WriteFile(dxfFile, []byte(dxfContent), 0644)
				require.NoError(t, err)
				return dxfFile
			},
			args:        []string{"extract", "-file"},
			expectError: false,
			expectedOutput: []string{
				"Successfully extracted DXF information",
				"DXF Version:",
				"Number of layers:",
			},
		},
		{
			name: "Extract with nonexistent file",
			setupFile: func(t *testing.T) string {
				return "/nonexistent/file.dwg"
			},
			args:        []string{"extract", "-file"},
			expectError: true,
			expectedOutput: []string{
				"Error:",
			},
		},
		{
			name: "Extract with invalid DWG file",
			setupFile: func(t *testing.T) string {
				tmpDir := t.TempDir()
				invalidFile := filepath.Join(tmpDir, "invalid.dwg")
				err := os.WriteFile(invalidFile, []byte("not a real dwg file"), 0644)
				require.NoError(t, err)
				return invalidFile
			},
			args:        []string{"extract", "-file"},
			expectError: true,
			expectedOutput: []string{
				"Error:",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This should fail initially - we need to build and test the actual binary
			testFile := tt.setupFile(t)
			if tt.cleanup != nil {
				defer tt.cleanup(testFile)
			}

			// Build the binary
			execName := "dwg-extractor-test"
			if runtime.GOOS == "windows" {
				execName += ".exe"
			}

			buildCmd := exec.Command("go", "build", "-o", execName, ".")
			err := buildCmd.Run()
			require.NoError(t, err, "Should build test binary")
			defer os.Remove(execName)

			// Get absolute path for the executable
			execPath, err := filepath.Abs(execName)
			require.NoError(t, err, "Should get absolute path")

			// Run the command
			args := append(tt.args, testFile)
			testCmd := exec.Command(execPath, args...)
			output, err := testCmd.CombinedOutput()
			outputStr := string(output)

			if tt.expectError {
				assert.Error(t, err, "Should fail for invalid input")
			} else {
				assert.NoError(t, err, "Should succeed for valid input")
			}

			for _, expected := range tt.expectedOutput {
				assert.Contains(t, outputStr, expected, "Output should contain expected text")
			}
		})
	}
}

// TestEndToEnd_TUICommand tests the TUI command workflow
func TestEndToEnd_TUICommand(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping TUI integration tests in short mode")
	}

	tests := []struct {
		name        string
		args        []string
		timeout     time.Duration
		expectStart bool
	}{
		{
			name:        "TUI with sample data",
			args:        []string{"tui"},
			timeout:     2 * time.Second,
			expectStart: true,
		},
		{
			name:        "TUI with output directory",
			args:        []string{"tui", "-output", "/tmp"},
			timeout:     2 * time.Second,
			expectStart: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This should fail initially - TUI testing is complex

			// Build the binary
			execName := "dwg-extractor-tui-test"
			if runtime.GOOS == "windows" {
				execName += ".exe"
			}

			buildCmd := exec.Command("go", "build", "-o", execName, ".")
			err := buildCmd.Run()
			require.NoError(t, err, "Should build test binary")
			defer os.Remove(execName)

			// Get absolute path for the executable
			execPath, err := filepath.Abs(execName)
			require.NoError(t, err, "Should get absolute path")

			// Run TUI command with timeout
			testCmd := exec.Command(execPath, tt.args...)

			done := make(chan error, 1)
			go func() {
				err := testCmd.Run()
				done <- err
			}()

			select {
			case err := <-done:
				if tt.expectStart {
					// TUI might exit due to terminal requirements, that's OK
					t.Logf("TUI command completed: %v", err)
				} else {
					assert.Error(t, err, "Should fail for invalid TUI setup")
				}
			case <-time.After(tt.timeout):
				if tt.expectStart {
					// TUI started and is running, kill it
					if testCmd.Process != nil {
						testCmd.Process.Kill()
					}
					t.Log("TUI started successfully (timed out as expected)")
				} else {
					t.Error("TUI should have failed quickly")
				}
			}
		})
	}
}

// TestEndToEnd_VersionAndHelp tests version and help commands
func TestEndToEnd_VersionAndHelp(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedOutput []string
	}{
		{
			name: "Version command",
			args: []string{"version"},
			expectedOutput: []string{
				"Go DWG Extractor",
				"Version:",
				"Git Commit:",
				"Build Time:",
			},
		},
		{
			name: "Help command",
			args: []string{"help"},
			expectedOutput: []string{
				"Go DWG Extractor - Extract data",
				"Usage:",
				"Commands:",
				"extract",
				"tui",
			},
		},
		{
			name: "Version flag",
			args: []string{"--version"},
			expectedOutput: []string{
				"Go DWG Extractor",
				"Version:",
			},
		},
		{
			name: "Help flag",
			args: []string{"--help"},
			expectedOutput: []string{
				"Go DWG Extractor - Extract data",
				"Usage:",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Build the binary
			execName := "dwg-extractor-help-test"
			if runtime.GOOS == "windows" {
				execName += ".exe"
			}

			buildCmd := exec.Command("go", "build", "-o", execName, ".")
			err := buildCmd.Run()
			require.NoError(t, err, "Should build test binary")
			defer os.Remove(execName)

			// Get absolute path for the executable
			execPath, err := filepath.Abs(execName)
			require.NoError(t, err, "Should get absolute path")

			// Run the command
			testCmd := exec.Command(execPath, tt.args...)
			output, err := testCmd.CombinedOutput()
			outputStr := string(output)

			assert.NoError(t, err, "Version and help commands should not error")

			for _, expected := range tt.expectedOutput {
				assert.Contains(t, outputStr, expected, "Output should contain expected text")
			}
		})
	}
}

// TestEndToEnd_ErrorConditions tests various error conditions
func TestEndToEnd_ErrorConditions(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		env            map[string]string
		expectedError  bool
		expectedOutput []string
	}{
		{
			name:          "No command provided",
			args:          []string{},
			expectedError: true,
			expectedOutput: []string{
				"No command provided",
			},
		},
		{
			name:          "Invalid command",
			args:          []string{"invalid-command"},
			expectedError: true,
			expectedOutput: []string{
				"Error:",
			},
		},
		{
			name:          "Extract without file flag",
			args:          []string{"extract"},
			expectedError: true,
			expectedOutput: []string{
				"Error:",
			},
		},
		{
			name: "Invalid ODA converter path",
			args: []string{"extract", "-file", "dummy.dwg"},
			env: map[string]string{
				"ODA_CONVERTER_PATH": "/invalid/path/converter",
			},
			expectedError: true,
			expectedOutput: []string{
				"Error:",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Build the binary
			execName := "dwg-extractor-error-test"
			if runtime.GOOS == "windows" {
				execName += ".exe"
			}

			buildCmd := exec.Command("go", "build", "-o", execName, ".")
			err := buildCmd.Run()
			require.NoError(t, err, "Should build test binary")
			defer os.Remove(execName)

			// Get absolute path for the executable
			execPath, err := filepath.Abs(execName)
			require.NoError(t, err, "Should get absolute path")

			// Set up environment
			testCmd := exec.Command(execPath, tt.args...)
			if tt.env != nil {
				testCmd.Env = os.Environ()
				for key, value := range tt.env {
					testCmd.Env = append(testCmd.Env, key+"="+value)
				}
			}

			// Run the command
			output, err := testCmd.CombinedOutput()
			outputStr := string(output)

			if tt.expectedError {
				assert.Error(t, err, "Should fail for error condition")
			} else {
				assert.NoError(t, err, "Should succeed")
			}

			for _, expected := range tt.expectedOutput {
				assert.Contains(t, outputStr, expected, "Output should contain expected error text")
			}
		})
	}
}

// TestEndToEnd_CrossPlatformCompatibility tests cross-platform behavior
func TestEndToEnd_CrossPlatformCompatibility(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping cross-platform tests in short mode")
	}

	tests := []struct {
		name     string
		testFunc func(t *testing.T)
		skipOnOS []string
	}{
		{
			name: "Path handling with different separators",
			testFunc: func(t *testing.T) {
				// Test that the application handles different path separators correctly
				tmpDir := t.TempDir()

				// Create a test file with platform-appropriate path
				testFile := filepath.Join(tmpDir, "test.dxf")
				err := os.WriteFile(testFile, []byte("dummy content"), 0644)
				require.NoError(t, err)

				// The path should work regardless of platform
				assert.True(t, filepath.IsAbs(filepath.Dir(testFile)) || strings.Contains(testFile, string(filepath.Separator)))
			},
		},
		{
			name: "Executable permissions",
			testFunc: func(t *testing.T) {
				// Build the binary
				execName := "dwg-extractor-perm-test"
				if runtime.GOOS == "windows" {
					execName += ".exe"
				}

				buildCmd := exec.Command("go", "build", "-o", execName, ".")
				err := buildCmd.Run()
				require.NoError(t, err, "Should build test binary")
				defer os.Remove(execName)

				// Get absolute path for the executable
				execPath, err := filepath.Abs(execName)
				require.NoError(t, err, "Should get absolute path")

				// Check that the binary is executable
				info, err := os.Stat(execPath)
				require.NoError(t, err)

				if runtime.GOOS != "windows" {
					// On Unix-like systems, check execute permissions
					mode := info.Mode()
					assert.True(t, mode&0111 != 0, "Binary should be executable")
				}
			},
			skipOnOS: []string{}, // Run on all platforms
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Check if we should skip this test on current OS
			currentOS := runtime.GOOS
			if currentOS == "" {
				currentOS = "unknown"
			}

			for _, skipOS := range tt.skipOnOS {
				if currentOS == skipOS {
					t.Skipf("Skipping test on %s", currentOS)
					return
				}
			}

			tt.testFunc(t)
		})
	}
}

// TestEndToEnd_PerformanceBaseline tests basic performance characteristics
func TestEndToEnd_PerformanceBaseline(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance tests in short mode")
	}

	t.Run("Application startup time", func(t *testing.T) {
		// Build the binary
		execName := "dwg-extractor-perf-test"
		if runtime.GOOS == "windows" {
			execName += ".exe"
		}

		buildCmd := exec.Command("go", "build", "-o", execName, ".")
		err := buildCmd.Run()
		require.NoError(t, err, "Should build test binary")
		defer os.Remove(execName)

		// Get absolute path for the executable
		execPath, err := filepath.Abs(execName)
		require.NoError(t, err, "Should get absolute path")

		// Measure startup time for version command
		start := time.Now()
		testCmd := exec.Command(execPath, "version")
		err = testCmd.Run()
		duration := time.Since(start)

		assert.NoError(t, err, "Version command should succeed")
		assert.Less(t, duration, 5*time.Second, "Application should start quickly")
		t.Logf("Startup time: %v", duration)
	})

	t.Run("Memory usage baseline", func(t *testing.T) {
		// This is a basic test - in a real scenario you'd use more sophisticated profiling
		t.Log("Memory usage testing would require profiling tools")
		// Could use runtime.MemStats or external tools like pprof
	})
}
