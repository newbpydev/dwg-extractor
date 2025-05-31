package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/remym/go-dwg-extractor/pkg/config"
)

var (
	dwgFile string
	cfg     *config.AppConfig
)

// Execute runs the root command
func Execute() error {
	// Parse command line flags
	rootCmd := flag.String("file", "", "Path to the DWG file to process")
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
	fmt.Printf("Processing DWG file: %s\n", dwgFile)
	fmt.Printf("Using ODA Converter: %s\n", cfg.ODAConverterPath)

	return nil
}
