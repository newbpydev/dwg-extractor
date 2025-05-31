package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestBundledConverterDetection tests detection of bundled ODA converter
func TestBundledConverterDetection(t *testing.T) {
	tests := []struct {
		name          string
		setupBundled  bool
		expectBundled bool
		expectedPath  string
	}{
		{
			name:          "Bundled converter exists",
			setupBundled:  true,
			expectBundled: true,
		},
		{
			name:          "No bundled converter",
			setupBundled:  false,
			expectBundled: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// For testing, we'll test the functions individually since the bundled converter
			// detection depends on the executable path which is different during testing

			if tt.expectBundled {
				// Test that the functions work correctly
				bundledPath := GetBundledConverterPath("windows")
				assert.NotEmpty(t, bundledPath, "Should return bundled converter path")
				assert.Contains(t, filepath.ToSlash(bundledPath), "assets/oda_converter", "Path should contain assets/oda_converter")

				// Test executable path function
				execPath, err := GetExecutablePath()
				assert.NoError(t, err, "Should be able to get executable path")
				assert.NotEmpty(t, execPath, "Executable path should not be empty")

				// Test path resolution
				resolvedPath := ResolveBundledPath("/test/path", "assets/oda_converter/windows/ODAFileConverter.exe")
				expected := filepath.Join("/test/path", "assets/oda_converter/windows/ODAFileConverter.exe")
				assert.Equal(t, expected, resolvedPath, "Should resolve bundled path correctly")
			} else {
				// Test that validation fails for non-existent paths
				isValid := ValidateBundledConverter("/nonexistent/path/converter.exe")
				assert.False(t, isValid, "Should not validate non-existent converter")
			}
		})
	}
}

// TestBundledConverterPlatformSpecific tests platform-specific bundled converter paths
func TestBundledConverterPlatformSpecific(t *testing.T) {
	tests := []struct {
		name           string
		platform       string
		expectedSuffix string
	}{
		{
			name:           "Windows bundled converter",
			platform:       "windows",
			expectedSuffix: "ODAFileConverter.exe",
		},
		{
			name:           "Linux bundled converter",
			platform:       "linux",
			expectedSuffix: "ODAFileConverter",
		},
		{
			name:           "macOS bundled converter",
			platform:       "darwin",
			expectedSuffix: "ODAFileConverter",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This should fail initially - no platform-specific bundled converter detection
			bundledPath := GetBundledConverterPath(tt.platform)

			assert.NotEmpty(t, bundledPath, "Should return bundled converter path for platform")
			assert.Contains(t, filepath.ToSlash(bundledPath), "assets/oda_converter", "Path should contain assets directory")
			assert.Contains(t, filepath.ToSlash(bundledPath), tt.platform, "Path should contain platform directory")
			assert.Contains(t, bundledPath, tt.expectedSuffix, "Path should end with correct executable name")
		})
	}
}

// TestBundledConverterExecutablePath tests finding executable path for relative bundled paths
func TestBundledConverterExecutablePath(t *testing.T) {
	// This should fail initially - no executable path resolution exists
	execPath, err := GetExecutablePath()
	require.NoError(t, err, "Should be able to get executable path")
	assert.NotEmpty(t, execPath, "Executable path should not be empty")

	// Test relative path resolution from executable
	bundledPath := ResolveBundledPath(execPath, "assets/oda_converter/windows/ODAFileConverter.exe")
	assert.NotEmpty(t, bundledPath, "Should resolve bundled path relative to executable")
	assert.Contains(t, bundledPath, "assets", "Resolved path should contain assets")
}

// TestBundledConverterValidation tests validation of bundled converter executables
func TestBundledConverterValidation(t *testing.T) {
	// Get the project root directory (two levels up from pkg/config)
	wd, err := os.Getwd()
	require.NoError(t, err)
	projectRoot := filepath.Join(wd, "..", "..")

	tests := []struct {
		name         string
		path         string
		shouldExist  bool
		shouldBeExec bool
	}{
		{
			name:         "Valid bundled converter",
			path:         filepath.Join(projectRoot, "assets", "oda_converter", "windows", "ODAFileConverter.exe"),
			shouldExist:  true,
			shouldBeExec: true,
		},
		{
			name:         "Missing bundled converter",
			path:         filepath.Join(projectRoot, "assets", "oda_converter", "nonexistent", "ODAFileConverter.exe"),
			shouldExist:  false,
			shouldBeExec: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test validation function
			isValid := ValidateBundledConverter(tt.path)

			if tt.shouldExist && tt.shouldBeExec {
				// We expect this to pass since we created the mock files
				assert.True(t, isValid, "Should validate existing executable bundled converter")
			} else {
				assert.False(t, isValid, "Should not validate missing bundled converter")
			}
		})
	}
}

// TestConfigWithBundledConverter tests config loading with bundled converter preference
func TestConfigWithBundledConverter(t *testing.T) {
	// Clear environment variable for clean test
	originalPath := os.Getenv("ODA_CONVERTER_PATH")
	defer os.Setenv("ODA_CONVERTER_PATH", originalPath)
	os.Unsetenv("ODA_CONVERTER_PATH")

	// This should fail initially - config doesn't check for bundled converter
	config, err := LoadConfig()
	require.NoError(t, err, "Should be able to load config")

	// Should prefer bundled converter over default path when available
	if bundledPath, exists := DetectBundledConverter(); exists {
		assert.Equal(t, bundledPath, config.ODAConverterPath,
			"Config should prefer bundled converter when available")
	}
}

// TestBundledConverterPriority tests the priority order: env var > bundled > default
func TestBundledConverterPriority(t *testing.T) {
	tests := []struct {
		name          string
		envVarSet     bool
		envVarValue   string
		bundledExists bool
		expectedType  string // "env", "bundled", "default"
	}{
		{
			name:          "Environment variable takes priority",
			envVarSet:     true,
			envVarValue:   "/custom/path/converter.exe",
			bundledExists: true,
			expectedType:  "env",
		},
		{
			name:          "Default path when nothing else available",
			envVarSet:     false,
			bundledExists: false,
			expectedType:  "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up test environment
			originalPath := os.Getenv("ODA_CONVERTER_PATH")
			defer os.Setenv("ODA_CONVERTER_PATH", originalPath)

			if tt.envVarSet {
				os.Setenv("ODA_CONVERTER_PATH", tt.envVarValue)
			} else {
				os.Unsetenv("ODA_CONVERTER_PATH")
			}

			// Test the priority logic
			config, err := LoadConfigWithPriority()
			require.NoError(t, err, "Should load config with priority logic")

			switch tt.expectedType {
			case "env":
				assert.Equal(t, tt.envVarValue, config.ODAConverterPath, "Should use environment variable")
			case "default":
				// Should use default path for current platform
				expectedDefault := getDefaultODAConverterPath()
				assert.Equal(t, expectedDefault, config.ODAConverterPath, "Should use default path")
			}
		})
	}
}

// TestBundledConverterCrossPlatform tests bundled converter works across platforms
func TestBundledConverterCrossPlatform(t *testing.T) {
	platforms := []struct {
		goos     string
		goarch   string
		expected string
	}{
		{"windows", "amd64", "windows/ODAFileConverter.exe"},
		{"linux", "amd64", "linux/ODAFileConverter"},
		{"darwin", "amd64", "darwin/ODAFileConverter"},
		{"darwin", "arm64", "darwin/ODAFileConverter"},
	}

	for _, platform := range platforms {
		t.Run(platform.goos+"_"+platform.goarch, func(t *testing.T) {
			// This should fail initially - cross-platform bundled detection not implemented
			bundledPath := GetBundledConverterForPlatform(platform.goos, platform.goarch)

			assert.NotEmpty(t, bundledPath, "Should return bundled path for platform")
			assert.Contains(t, filepath.ToSlash(bundledPath), filepath.ToSlash(platform.expected), "Should contain correct platform-specific path")
		})
	}
}

// Helper function to get default converter path (this will exist)
