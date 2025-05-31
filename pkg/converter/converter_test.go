package converter

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDWGConverter_ConvertToDXF(t *testing.T) {
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
			dwgPath:   "test.dwg",
			outputDir: "output",
			setup: func() {
				// This will be implemented with proper mocking in Task 2.2
			},
			expectError: false,
		},
		{
			name:        "empty dwg path",
			dwgPath:     "",
			outputDir:   "output",
			setup:       func() {},
			expectError: true,
			errContains: "DWG path cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test environment if needed
			if tt.setup != nil {
				tt.setup()
			}

			// Create a new converter instance
			converter, err := NewDWGConverter("path/to/odaconverter")
			assert.NoError(t, err)
			assert.NotNil(t, converter)

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
				assert.Equal(t, filepath.Join(tt.outputDir, "test.dxf"), dxfPath)
			}
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
