package config

import (
	"os"
	"path/filepath"
	"runtime"
)

// DetectBundledConverter checks if a bundled ODA converter exists
func DetectBundledConverter() (string, bool) {
	execPath, err := GetExecutablePath()
	if err != nil {
		return "", false
	}

	// Get platform-specific bundled converter path
	bundledPath := GetBundledConverterPath(runtime.GOOS)
	fullPath := ResolveBundledPath(execPath, bundledPath)

	// Validate the bundled converter exists and is executable
	if ValidateBundledConverter(fullPath) {
		return fullPath, true
	}

	return "", false
}

// GetBundledConverterPath returns the platform-specific bundled converter path
func GetBundledConverterPath(platform string) string {
	var executableName string

	switch platform {
	case "windows":
		executableName = "ODAFileConverter.exe"
	case "linux", "darwin":
		executableName = "ODAFileConverter"
	default:
		executableName = "ODAFileConverter"
	}

	return filepath.Join("assets", "oda_converter", platform, executableName)
}

// GetExecutablePath returns the path of the current executable
func GetExecutablePath() (string, error) {
	execPath, err := os.Executable()
	if err != nil {
		return "", err
	}

	// Resolve any symlinks
	execPath, err = filepath.EvalSymlinks(execPath)
	if err != nil {
		return "", err
	}

	return filepath.Dir(execPath), nil
}

// ResolveBundledPath resolves a bundled asset path relative to the executable
func ResolveBundledPath(execPath, bundledPath string) string {
	return filepath.Join(execPath, bundledPath)
}

// ValidateBundledConverter validates that a bundled converter exists and is executable
func ValidateBundledConverter(path string) bool {
	// Check if file exists
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	// Check if it's a file (not directory)
	if info.IsDir() {
		return false
	}

	// Check if it's executable (on Unix systems)
	if runtime.GOOS != "windows" {
		mode := info.Mode()
		if mode&0111 == 0 { // Check if any execute bit is set
			return false
		}
	}

	return true
}

// GetBundledConverterForPlatform returns bundled converter path for specific platform
func GetBundledConverterForPlatform(goos, goarch string) string {
	var executableName string

	switch goos {
	case "windows":
		executableName = "ODAFileConverter.exe"
	case "linux", "darwin":
		executableName = "ODAFileConverter"
	default:
		executableName = "ODAFileConverter"
	}

	// For now, we don't differentiate by architecture, only OS
	return filepath.Join("assets", "oda_converter", goos, executableName)
}

// LoadConfigWithPriority loads configuration with priority order: env var > bundled > default
func LoadConfigWithPriority() (*AppConfig, error) {
	config := &AppConfig{}

	// Priority 1: Environment variable
	if envPath := os.Getenv("ODA_CONVERTER_PATH"); envPath != "" {
		config.ODAConverterPath = envPath
		return config, nil
	}

	// Priority 2: Bundled converter
	if bundledPath, exists := DetectBundledConverter(); exists {
		config.ODAConverterPath = bundledPath
		return config, nil
	}

	// Priority 3: Default path
	config.ODAConverterPath = getDefaultODAConverterPath()
	return config, nil
}
