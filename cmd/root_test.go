package cmd

import (
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/remym/go-dwg-extractor/pkg/converter"
	"github.com/remym/go-dwg-extractor/pkg/config"
	"github.com/remym/go-dwg-extractor/pkg/dxfparser"
	"github.com/remym/go-dwg-extractor/pkg/data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockDWGConverter is a mock implementation of the DWGConverter interface
type MockDWGConverter struct {
	ConvertToDXFFunc func(dwgPath, outputDir string) (string, error)
}

// MockParser is a mock implementation of the Parser interface
type MockParser struct {
	ParseDXFFunc func(dxfPath string) (*data.ExtractedData, error)
}

func (m *MockDWGConverter) ConvertToDXF(dwgPath, outputDir string) (string, error) {
	return m.ConvertToDXFFunc(dwgPath, outputDir)
}

func (m *MockParser) ParseDXF(dxfPath string) (*data.ExtractedData, error) {
	return m.ParseDXFFunc(dxfPath)
}

// Mock function to create a parser
var newParser = func() dxfparser.ParserInterface {
	return dxfparser.NewParser()
}

func TestRootCommand(t *testing.T) {
	// Save original command-line arguments and environment variables
	oldArgs := os.Args
	oldEnv := os.Getenv("ODA_CONVERTER_PATH")
	oldNewDWGConverter := newDWGConverter
	oldNewParser := newParser
	defer func() {
		os.Args = oldArgs
		os.Setenv("ODA_CONVERTER_PATH", oldEnv)
		newDWGConverter = oldNewDWGConverter
		newParser = oldNewParser
	}()

	// Create a temporary directory for test files
	tempDir, err := ioutil.TempDir("", "test-root-*")
	require.NoError(t, err, "Failed to create temp directory")
	defer os.RemoveAll(tempDir)

	// Create a test DWG file
	testDWGPath := filepath.Join(tempDir, "test.dwg")
	err = os.WriteFile(testDWGPath, []byte("test content"), 0644)
	require.NoError(t, err, "Failed to create test DWG file")

	// Create a test output directory
	outputDir := filepath.Join(tempDir, "output")
	err = os.MkdirAll(outputDir, 0755)
	require.NoError(t, err, "Failed to create output directory")

	// Set up test cases
	tests := []struct {
		name        string
		args        []string
		setup       func()
		wantErr     bool
		errContains string
	}{
		{
			name:        "no arguments",
			args:        []string{"cmd"},
			setup:       func() { newDWGConverter = converter.NewDWGConverter },
			wantErr:     true,
			errContains: "no DWG file specified",
		},
		{
			name: "successful conversion with default output",
			args: []string{"cmd", "-file", testDWGPath},
			setup: func() {
				newDWGConverter = func(path string) (converter.DWGConverter, error) {
					mock := &MockDWGConverter{
						ConvertToDXFFunc: func(dwgPath, outputDir string) (string, error) {
							dxfPath := filepath.Join(filepath.Dir(dwgPath), filepath.Base(dwgPath)+".dxf")
							// Create a dummy DXF file with layer information
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
EOF`
							_ = os.WriteFile(dxfPath, []byte(dxfContent), 0644)
							return dxfPath, nil
						},
					}
					return mock, nil
				}
			},
			wantErr: false,
		},
		{
			name: "successful conversion with custom output directory",
			args: []string{"cmd", "-file", testDWGPath, "-output", outputDir},
			setup: func() {
				newDWGConverter = func(path string) (converter.DWGConverter, error) {
					mock := &MockDWGConverter{
						ConvertToDXFFunc: func(dwgPath, outputDir string) (string, error) {
							dxfPath := filepath.Join(outputDir, filepath.Base(dwgPath)+".dxf")
							// Create a dummy DXF file with layer information
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
EOF`
							_ = os.WriteFile(dxfPath, []byte(dxfContent), 0644)
							return dxfPath, nil
						},
					}
					return mock, nil
				}
			},
			wantErr: false,
		},
		{
			name: "converter returns error",
			args: []string{"cmd", "-file", testDWGPath},
			setup: func() {
				newDWGConverter = func(path string) (converter.DWGConverter, error) {
					return nil, assert.AnError
				}
			},
			wantErr:     true,
			errContains: "failed to create DWG converter",
		},
		{
			name: "conversion fails",
			args: []string{"cmd", "-file", testDWGPath},
			setup: func() {
				newDWGConverter = func(path string) (converter.DWGConverter, error) {
					mock := &MockDWGConverter{
						ConvertToDXFFunc: func(dwgPath, outputDir string) (string, error) {
							return "", assert.AnError
						},
					}
					return mock, nil
				}
			},
			wantErr:     true,
			errContains: "conversion failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset the flag set to avoid flag redefinition errors
			flag.CommandLine = flag.NewFlagSet(tt.args[0], flag.ExitOnError)

			// Set up the test
			if tt.setup != nil {
				// Save the original newDWGConverter
				oldNewDWGConverter := newDWGConverter
				defer func() { newDWGConverter = oldNewDWGConverter }()

				tt.setup()
			}

			// Set command-line arguments for the test
			os.Args = tt.args

			// Execute the command
			execErr := Execute()

			// Verify the results
			if tt.wantErr {
				assert.Error(t, execErr)
				if tt.errContains != "" {
					assert.Contains(t, execErr.Error(), tt.errContains)
				}
			} else {
				assert.NoError(t, execErr)
			}
		})
	}
}

func TestConfigLoading(t *testing.T) {
	// Save original environment variable and function
	oldEnv := os.Getenv("ODA_CONVERTER_PATH")
	oldNewDWGConverter := newDWGConverter
	defer func() {
		os.Setenv("ODA_CONVERTER_PATH", oldEnv)
		newDWGConverter = oldNewDWGConverter
	}()

	// Mock the converter to avoid actual execution
	newDWGConverter = func(path string) (converter.DWGConverter, error) {
		return &MockDWGConverter{
			ConvertToDXFFunc: func(dwgPath, outputDir string) (string, error) {
				return "test.dxf", nil
			},
		}, nil
	}

	// Create a temporary file for testing
	tempFile, err := ioutil.TempFile("", "test-converter-*.exe")
	require.NoError(t, err, "Failed to create temp file")
	tempFilePath := tempFile.Name()
	tempFile.Close()
	defer os.Remove(tempFilePath)

	// Test with environment variable set
	os.Setenv("ODA_CONVERTER_PATH", tempFilePath)
	cfg, err := config.LoadConfig()
	require.NoError(t, err)
	require.NotNil(t, cfg)
	assert.Equal(t, tempFilePath, cfg.ODAConverterPath)

	// Test with default path (should use the default path since we can't test the actual default)
	os.Unsetenv("ODA_CONVERTER_PATH")
	cfg, err = config.LoadConfig()
	require.NoError(t, err)
	require.NotNil(t, cfg)
	assert.Equal(t, config.DefaultODAConverterPath, cfg.ODAConverterPath)
}

// contains checks if a string contains another string
func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr
}
