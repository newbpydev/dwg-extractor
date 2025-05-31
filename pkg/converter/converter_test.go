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
					// Verify the command and args are as expected for new format
					assert.Equal(t, "path/to/odaconverter", command)
					assert.Len(t, args, 7) // Should have 7 positional arguments

					// Check positional arguments: InputDir OutputDir Version FileType Recurse Audit Filter
					assert.Contains(t, args[0], tempDir)                       // Input directory should contain tempDir
					assert.Equal(t, filepath.Join(tempDir, "output"), args[1]) // Output directory
					assert.Equal(t, "ACAD2018", args[2])                       // Version
					assert.Equal(t, "DXF", args[3])                            // File type
					assert.Equal(t, "0", args[4])                              // Recurse
					assert.Equal(t, "0", args[5])                              // Audit
					assert.Equal(t, "*.DWG", args[6])                          // Filter

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
			name:      "output directory creation failure",
			dwgPath:   testDWGPath,
			outputDir: filepath.Join(tempDir, "fail-output"),
			setup: func() {
				// Make a file where the output dir should be, so MkdirAll fails
				failPath := filepath.Join(tempDir, "fail-output")
				_ = os.WriteFile(failPath, []byte("not a dir"), 0644)
			},
			expectError: true,
			errContains: "failed to create output directory",
		},
		{
			name:      "command execution failure",
			dwgPath:   testDWGPath,
			outputDir: filepath.Join(tempDir, "cmd-fail-output"),
			setup: func() {
				commandContext = func(ctx context.Context, command string, args ...string) *exec.Cmd {
					cmd := exec.CommandContext(ctx, "false") // always fails
					return cmd
				}
			},
			expectError: true,
			errContains: "failed to convert DWG to DXF",
		},
		{
			name:      "no DXF file generated, no alternate found",
			dwgPath:   testDWGPath,
			outputDir: filepath.Join(tempDir, "no-dxf-output"),
			setup: func() {
				commandContext = func(ctx context.Context, command string, args ...string) *exec.Cmd {
					// Don't create any DXF file
					return exec.CommandContext(ctx, "echo", "mock command")
				}
			},
			expectError: true,
			errContains: "conversion failed: no DXF file was generated",
		},
		{
			name:      "alternate DXF file found",
			dwgPath:   testDWGPath,
			outputDir: filepath.Join(tempDir, "alt-dxf-output"),
			setup: func() {
				commandContext = func(ctx context.Context, command string, args ...string) *exec.Cmd {
					// Create a DXF file with a different name
					altDXF := filepath.Join(filepath.Join(tempDir, "alt-dxf-output"), "altname.dxf")
					_ = os.MkdirAll(filepath.Dir(altDXF), 0755)
					_ = os.WriteFile(altDXF, []byte("DXF content"), 0644)
					return exec.CommandContext(ctx, "echo", "mock command")
				}
			},
			expectError: false,
		},
		{
			name:        "empty dwg path",
			dwgPath:     "",
			outputDir:   filepath.Join(tempDir, "output"),
			setup:       func() {},
			expectError: true,
			errContains: "DWG path cannot be empty",
		},
		{
			name:        "nonexistent dwg file",
			dwgPath:     filepath.Join(tempDir, "nonexistent.dwg"),
			outputDir:   filepath.Join(tempDir, "output"),
			setup:       func() {},
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

			// Create output directory only if we don't expect directory creation to fail
			if tt.name != "output directory creation failure" {
				err = os.MkdirAll(tt.outputDir, 0755)
				require.NoError(t, err)
			}

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
		name          string
		converterPath string
		expectError   bool
		errContains   string
	}{
		{
			name:          "valid converter path",
			converterPath: "path/to/odaconverter",
			expectError:   false,
		},
		{
			name:          "empty converter path",
			converterPath: "",
			expectError:   true,
			errContains:   "converter path cannot be empty",
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
