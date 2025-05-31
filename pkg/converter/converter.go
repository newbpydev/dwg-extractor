// Package converter provides functionality to convert DWG files to DXF format
// using the ODA File Converter.
package converter

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// commandContext is a variable that holds the function to create commands
// This is used to allow mocking in tests
var commandContext = exec.CommandContext

// DWGConverter defines the interface for converting DWG files to DXF.
type DWGConverter interface {
	// ConvertToDXF converts a DWG file to DXF format.
	// It returns the path to the converted DXF file or an error if the conversion fails.
	ConvertToDXF(dwgPath, outputDir string) (string, error)
}

// odaconverter implements the DWGConverter interface.
type odaconverter struct {
	converterPath string // Path to the ODA File Converter executable
}

// NewDWGConverter creates a new instance of DWGConverter.
// It returns an error if the converter path is empty.
func NewDWGConverter(converterPath string) (DWGConverter, error) {
	if converterPath == "" {
		return nil, fmt.Errorf("converter path cannot be empty")
	}

	return &odaconverter{
		converterPath: converterPath,
	}, nil
}

// ConvertToDXF converts the specified DWG file to DXF format using the ODA File Converter.
// It returns the path to the converted DXF file or an error if the conversion fails.
func (c *odaconverter) ConvertToDXF(dwgPath, outputDir string) (string, error) {
	if dwgPath == "" {
		return "", fmt.Errorf("DWG path cannot be empty")
	}

	// Check if the input file exists
	if _, err := os.Stat(dwgPath); os.IsNotExist(err) {
		return "", fmt.Errorf("input file does not exist: %s", dwgPath)
	}

	// Ensure output directory exists
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create output directory: %w", err)
	}

	// Get the base name of the input file without extension
	baseName := filepath.Base(dwgPath)
	ext := filepath.Ext(baseName)
	if ext != "" {
		baseName = baseName[0 : len(baseName)-len(ext)]
	}

	// Generate DXF file path (same name as DWG but with .dxf extension)
	dxfPath := filepath.Join(outputDir, baseName+".dxf")

	// Create a context with timeout for the conversion
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Prepare the command to run the ODA File Converter
	// Command format: ODAFileConverter.exe <input> <output> version [input format] [output format] [recurse] [report] [audit]
	cmd := commandContext(
		ctx,
		c.converterPath,
		"-i", dwgPath,
		"-o", outputDir,
		"-f", "DXF",
		"-v", "ACAD2018",
	)

	// Set up output buffers
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Run the command
	if err := cmd.Run(); err != nil {
		// If the command failed, include stderr in the error message
		return "", fmt.Errorf("failed to convert DWG to DXF: %w\n%s", err, stderr.String())
	}

	// Verify the output file was created
	if _, err := os.Stat(dxfPath); os.IsNotExist(err) {
		// If the expected DXF file doesn't exist, check for other possible names
		// Sometimes the converter might use a different naming convention
		files, err := filepath.Glob(filepath.Join(outputDir, "*.dxf"))
		if err != nil || len(files) == 0 {
			return "", fmt.Errorf("conversion failed: no DXF file was generated")
		}

		// Filter files to only include those that were likely created by this conversion
		// We'll check modification time to avoid using pre-existing files
		var recentFiles []string
		conversionStartTime := time.Now().Add(-6 * time.Minute) // Allow 6 minutes for conversion

		for _, file := range files {
			info, err := os.Stat(file)
			if err != nil {
				continue
			}

			// Only consider files modified after we started the conversion
			if info.ModTime().After(conversionStartTime) {
				recentFiles = append(recentFiles, file)
			}
		}

		if len(recentFiles) == 0 {
			return "", fmt.Errorf("conversion failed: no recently created DXF file was found")
		}

		// Prefer files with the expected base name, but accept any recent file
		expectedBase := baseName + ".dxf"
		for _, file := range recentFiles {
			if filepath.Base(file) == expectedBase {
				dxfPath = file
				break
			}
		}

		// If no file with expected name found, use the first recent file
		if dxfPath == filepath.Join(outputDir, baseName+".dxf") {
			dxfPath = recentFiles[0]
		}
	}

	return dxfPath, nil
}
