package cmd

import (
	"flag"
	"io/ioutil"
	"os"
	"testing"

	"github.com/remym/go-dwg-extractor/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRootCommand(t *testing.T) {
	// Save original command-line arguments and environment variables
	oldArgs := os.Args
	oldEnv := os.Getenv("ODA_CONVERTER_PATH")
	defer func() { 
		os.Args = oldArgs 
		os.Setenv("ODA_CONVERTER_PATH", oldEnv)
	}()

	// Create a temporary file for testing ODA converter
	tempConverter, err := ioutil.TempFile("", "test-converter-*.exe")
	require.NoError(t, err, "Failed to create temp converter file")
	tempConverterPath := tempConverter.Name()
	tempConverter.Close()
	defer os.Remove(tempConverterPath)

	// Set the ODA_CONVERTER_PATH environment variable for testing
	os.Setenv("ODA_CONVERTER_PATH", tempConverterPath)

	tests := []struct {
		name        string
		args        []string
		envVars     map[string]string
		wantErr     bool
		errContains string
	}{
		{
			name:        "no arguments",
			args:        []string{"cmd"},
			wantErr:     true,
			errContains: "no DWG file specified",
		},
		{
			name:    "with file argument",
			args:    []string{"cmd", "-file"},
			wantErr: false,
		},
		{
			name:        "non-existent file",
			args:        []string{"cmd", "-file", "nonexistent.dwg"},
			wantErr:     true,
			errContains: "file does not exist",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up environment variables
			for k, v := range tt.envVars {
				os.Setenv(k, v)
			}

			// Create a temporary file for testing
			tempFile, err := ioutil.TempFile("", "testfile*.dwg")
			require.NoError(t, err, "Failed to create temp file")
			tempFileName := tempFile.Name()
			tempFile.Close()
			defer os.Remove(tempFileName)

			// Update args with the temp file path
			args := make([]string, len(tt.args))
			copy(args, tt.args)
			if len(args) > 1 && args[1] == "-file" && len(args) == 2 {
				args = append(args, tempFileName)
			}

			// Set command-line arguments for the test
			os.Args = args
			// Reset the flag set to avoid flag redefinition errors
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

			execErr := Execute()
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
	// Save original environment variable
	oldEnv := os.Getenv("ODA_CONVERTER_PATH")
	defer os.Setenv("ODA_CONVERTER_PATH", oldEnv)

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
