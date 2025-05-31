package config

import (
	"os"
	"path/filepath"
	"runtime"
)

// AppConfig holds the application configuration
type AppConfig struct {
	// ODAConverterPath is the path to the ODA File Converter executable
	ODAConverterPath string
}

// DefaultODAConverterPath is the default path where we expect to find the ODA File Converter
var DefaultODAConverterPath = getDefaultODAConverterPath()

// getDefaultODAConverterPath returns the default ODA converter path for the current platform
func getDefaultODAConverterPath() string {
	if runtime.GOOS == "windows" {
		return `C:\Program Files\ODA\ODAFileConverter 26.4.0\ODAFileConverter.exe`
	}
	return "/usr/local/bin/ODAFileConverter"
}

// LoadConfig loads the application configuration with priority: env var > bundled > default
func LoadConfig() (*AppConfig, error) {
	// Priority 1: Environment variable
	if envPath := os.Getenv("ODA_CONVERTER_PATH"); envPath != "" {
		return &AppConfig{
			ODAConverterPath: filepath.Clean(envPath),
		}, nil
	}

	// Priority 2: Bundled converter
	if bundledPath, exists := DetectBundledConverter(); exists {
		return &AppConfig{
			ODAConverterPath: bundledPath,
		}, nil
	}

	// Priority 3: Default path
	return &AppConfig{
		ODAConverterPath: DefaultODAConverterPath,
	}, nil
}

// Validate checks if the configuration is valid
func (c *AppConfig) Validate() error {
	if c.ODAConverterPath == "" {
		return ErrMissingODAConverterPath
	}

	// Check if the file exists and is executable
	info, err := os.Stat(c.ODAConverterPath)
	if os.IsNotExist(err) {
		return ErrODAConverterNotFound
	}
	if err != nil {
		return err
	}

	// On Windows, we can't directly check if a file is executable,
	// so we just check if it's a regular file
	if info.Mode().IsRegular() {
		return nil
	}

	return ErrInvalidODAConverterPath
}
