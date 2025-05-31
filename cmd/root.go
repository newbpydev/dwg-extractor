package cmd

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/remym/go-dwg-extractor/pkg/config"
)

var (
	dwgFile   string
	outputDir string
	cfg       *config.AppConfig
)

// Execute runs the root command
func Execute() error {
	// Parse command line flags
	rootCmd := flag.String("file", "", "Path to the DWG file to process")
	flag.StringVar(&outputDir, "output", "", "Output directory for converted files (default: same as input file)")
	flag.Parse()

	// Load configuration
	var err error
	cfg, err = config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	// Check if file is provided
	if *rootCmd == "" {
		return fmt.Errorf("no DWG file specified. Please provide a file using the -file flag")
	}

	// Check if file exists
	if _, err := os.Stat(*rootCmd); os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %s", *rootCmd)
	}

	dwgFile = *rootCmd

	// Set default output directory if not provided
	if outputDir == "" {
		outputDir = filepath.Dir(dwgFile)
	}

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	fmt.Printf("Processing DWG file: %s\n", dwgFile)
	fmt.Printf("Output directory: %s\n", outputDir)
	fmt.Printf("Using ODA Converter: %s\n", cfg.ODAConverterPath)

	// Create converter instance
	converter, err := newDWGConverter(cfg.ODAConverterPath)
	if err != nil {
		return fmt.Errorf("failed to create DWG converter: %w", err)
	}

	// Convert DWG to DXF
	dxfPath, err := converter.ConvertToDXF(dwgFile, outputDir)
	if err != nil {
		return fmt.Errorf("conversion failed: %w", err)
	}

	fmt.Printf("Successfully converted to DXF: %s\n", dxfPath)
	return nil
}
