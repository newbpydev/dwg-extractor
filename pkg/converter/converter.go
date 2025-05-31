// Package converter provides functionality to convert DWG files to DXF format
// using the ODA File Converter.
package converter

import (
	"fmt"
	"os"
	"path/filepath"
)

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

// ConvertToDXF converts the specified DWG file to DXF format.
// It returns the path to the converted DXF file or an error if the conversion fails.
func (c *odaconverter) ConvertToDXF(dwgPath, outputDir string) (string, error) {
	if dwgPath == "" {
		return "", fmt.Errorf("DWG path cannot be empty")
	}

	// Ensure output directory exists
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create output directory: %w", err)
	}

	// Generate DXF file path (same name as DWG but with .dxf extension)
	ext := filepath.Ext(dwgPath)
	dxfPath := filepath.Join(outputDir, filepath.Base(dwgPath)[0:len(dwgPath)-len(ext)]+".dxf")

	// TODO: Implement actual conversion using ODA File Converter in Task 2.2
	// This is a placeholder that will be replaced with the actual implementation
	// that calls the ODA File Converter CLI.

	return dxfPath, nil
}
