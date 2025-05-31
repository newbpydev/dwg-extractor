package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	// Create a temporary file for testing
	tempFile, err := os.CreateTemp("", "test-converter-*.exe")
	require.NoError(t, err, "Failed to create temp file")
	defer os.Remove(tempFile.Name())

	tests := []struct {
		name          string
		envVars       map[string]string
		wantPath      string
		wantErr       bool
		expectedError error
	}{
		{
			name:          "no env var, use default",
			envVars:       map[string]string{},
			wantPath:      DefaultODAConverterPath,
			wantErr:       false,
			expectedError: nil,
		},
		{
			name: "custom path from env var",
			envVars: map[string]string{
				"ODA_CONVERTER_PATH": "C:\\custom\\path\\to\\converter.exe",
			},
			wantPath:      "C:\\custom\\path\\to\\converter.exe",
			wantErr:       false,
			expectedError: nil,
		},
		{
			name: "empty path from env var",
			envVars: map[string]string{
				"ODA_CONVERTER_PATH": "",
			},
			wantPath:      DefaultODAConverterPath,
			wantErr:       false,
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up environment variables
			for k, v := range tt.envVars {
				t.Setenv(k, v)
			}

			// Load the config
			cfg, err := LoadConfig()

			// Check for expected errors
			if tt.wantErr {
				assert.Error(t, err)
				if tt.expectedError != nil {
					assert.ErrorIs(t, err, tt.expectedError)
				}
				return
			}

			// Check for no errors
			assert.NoError(t, err)
			require.NotNil(t, cfg)
			assert.Equal(t, tt.wantPath, cfg.ODAConverterPath)
		})
	}
}

func TestValidate(t *testing.T) {
	// Create a temporary file for testing
	tempFile, err := os.CreateTemp("", "test-converter-*.exe")
	require.NoError(t, err, "Failed to create temp file")
	tempFilePath := tempFile.Name()
	tempFile.Close()
	defer os.Remove(tempFilePath)

	tests := []struct {
		name          string
		setup         func() *AppConfig
		wantErr       bool
		expectedError error
	}{
		{
			name: "valid config",
			setup: func() *AppConfig {
				return &AppConfig{
					ODAConverterPath: tempFilePath,
				}
			},
			wantErr:       false,
			expectedError: nil,
		},
		{
			name: "empty path",
			setup: func() *AppConfig {
				return &AppConfig{
					ODAConverterPath: "",
				}
			},
			wantErr:       true,
			expectedError: ErrMissingODAConverterPath,
		},
		{
			name: "non-existent file",
			setup: func() *AppConfig {
				return &AppConfig{
					ODAConverterPath: filepath.Join(os.TempDir(), "non-existent-file.exe"),
				}
			},
			wantErr:       true,
			expectedError: ErrODAConverterNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := tt.setup()
			err := cfg.Validate()

			if tt.wantErr {
				assert.Error(t, err)
				if tt.expectedError != nil {
					assert.ErrorIs(t, err, tt.expectedError)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
