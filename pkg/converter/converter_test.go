package converter

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockCmd is a mock implementation of the command runner
type mockCmd struct {
	name string
	args []string
}

// Run simulates command execution for testing
func (m *mockCmd) Run() error {
	// Simulate successful command execution
	return nil
}

// mockCommandContext is a test helper to mock exec.CommandContext
func mockCommandContext(ctx context.Context, command string, args ...string) *exec.Cmd {
	// Create a command that will be handled by our test
	cmd := exec.CommandContext(ctx, "echo", "mock command")
	return cmd
}

func TestDWGConverter_ConvertToDXF(t *testing.T) {
	// Mock the command runner
	originalCommand := commandContext
	defer func() { commandContext = originalCommand }()

	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "converter-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a test DWG file
	testDWGPath := filepath.Join(tempDir, "test.dwg")
	err = os.WriteFile(testDWGPath, []byte("test content"), 0644)
	require.NoError(t, err)

	tests := []struct {
		name        string
		dwgPath     string
		outputDir   string
		setup       func()
		expectError bool
		errContains string
	}{
		{
			name:      "successful conversion",
			dwgPath:   testDWGPath,
			outputDir: filepath.Join(tempDir, "output"),
			setup: func() {
				commandContext = func(ctx context.Context, command string, args ...string) *exec.Cmd {
					// Verify the command and args are as expected
					assert.Equal(t, "path/to/odaconverter", command)
					assert.Contains(t, args, "-i")
					assert.Contains(t, args, testDWGPath)
					assert.Contains(t, args, "-o")
					assert.Contains(t, args, filepath.Join(tempDir, "output"))
					assert.Contains(t, args, "-f")
					assert.Contains(t, args, "DXF")
					assert.Contains(t, args, "-v")
					assert.Contains(t, args, "ACAD2018")
					
					// Create a mock DXF file to simulate conversion
					dxfPath := filepath.Join(filepath.Join(tempDir, "output"), "test.dxf")
					_ = os.MkdirAll(filepath.Dir(dxfPath), 0755)
					_ = os.WriteFile(dxfPath, []byte("DXF content"), 0644)
					
					return exec.CommandContext(ctx, "echo", "mock command")
				}
			},
			expectError: false,
		},
		{
			name:        "empty dwg path",
			dwgPath:     "",
			outputDir:   filepath.Join(tempDir, "output"),
			setup:       func(){},
			expectError: true,
			errContains: "DWG path cannot be empty",
		},
		{
			name:        "nonexistent dwg file",
			dwgPath:     filepath.Join(tempDir, "nonexistent.dwg"),
			outputDir:   filepath.Join(tempDir, "output"),
			setup:       func(){},
			expectError: true,
			errContains: "input file does not exist",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset commandContext to default
			commandContext = originalCommand

			// Setup test environment if needed
			if tt.setup != nil {
				tt.setup()
			}

			// Create a new converter instance
			converter, err := NewDWGConverter("path/to/odaconverter")
			require.NoError(t, err)
			require.NotNil(t, converter)

			// Create output directory
			err = os.MkdirAll(tt.outputDir, 0755)
			require.NoError(t, err)

			// Call the method under test
			dxfPath, err := converter.ConvertToDXF(tt.dwgPath, tt.outputDir)

			// Assert the results
			if tt.expectError {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, dxfPath)
				// Verify the DXF file has the correct extension
				ext := filepath.Ext(dxfPath)
				assert.Equal(t, ".dxf", ext)
			}

			// Cleanup
			os.RemoveAll(tt.outputDir)
		})
	}
}

func TestNewDWGConverter(t *testing.T) {
	tests := []struct {
		name        string
		converterPath string
		expectError bool
		errContains string
	}{
		{
			name:         "valid converter path",
			converterPath: "path/to/odaconverter",
			expectError:  false,
		},
		{
			name:         "empty converter path",
			converterPath: "",
			expectError:  true,
			errContains:  "converter path cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			converter, err := NewDWGConverter(tt.converterPath)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, converter)
			}
		})
	}
}
