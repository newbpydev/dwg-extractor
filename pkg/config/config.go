package config

import (
	"os"
	"path/filepath"
)

// AppConfig holds the application configuration
type AppConfig struct {
	// ODAConverterPath is the path to the ODA File Converter executable
	ODAConverterPath string
}

// DefaultODAConverterPath is the default path where we expect to find the ODA File Converter
var DefaultODAConverterPath = `C:\Program Files\ODA\ODAFileConverter 26.4.0\ODAFileConverter.exe`

// LoadConfig loads the application configuration from environment variables
func LoadConfig() (*AppConfig, error) {
	// First, try to get the path from environment variable
	path := os.Getenv("ODA_CONVERTER_PATH")
	
	// If not set, use the default path
	if path == "" {
		path = DefaultODAConverterPath
	}
	
	// Clean the path to handle any path separators
	path = filepath.Clean(path)
	
	// Create and return the config
	return &AppConfig{
		ODAConverterPath: path,
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
