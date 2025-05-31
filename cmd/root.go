package cmd

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/remym/go-dwg-extractor/pkg/converter"
	"github.com/remym/go-dwg-extractor/pkg/config"
	"github.com/remym/go-dwg-extractor/pkg/dxfparser"
)

var (
	rootCmd  string
	outputDir string
	cfg       *config.AppConfig
)

// Execute runs the root command
func Execute() error {
	// Check if no command is provided
	if len(os.Args) < 2 {
		return fmt.Errorf("no command provided. Use 'extract' or 'tui'")
	}

	// Handle the command
	command := os.Args[1]
	if command == "tui" {
		// For TUI, just run it without any file requirements
		return ExecuteTUI()
	} else if command == "extract" {
		// For extract, a DWG file is required
		if len(os.Args) < 3 {
			return fmt.Errorf("no DWG file specified. Usage: %s extract [DWG file]", os.Args[0])
		}
		// Remove the "extract" command from args
		os.Args = append(os.Args[:1], os.Args[2:]...)

		// Parse command line flags for extract command
		fileFlag := flag.String("file", "", "Path to the DWG file to process")
		flag.StringVar(&outputDir, "output", "", "Output directory for converted files (default: same as input file)")
		flag.Parse()

		// Set the root command from the flag
		rootCmd = *fileFlag

		// Check if file is provided for extract command
		if rootCmd == "" {
			return fmt.Errorf("no DWG file specified. Please provide a file using the -file flag")
		}

		// Load configuration
		var err error
		cfg, err = config.LoadConfig()
		if err != nil {
			return fmt.Errorf("failed to load configuration: %w", err)
		}

		// If output directory is not specified, use the same directory as the input file
		if outputDir == "" {
			outputDir = filepath.Dir(rootCmd)
		}

		// Create a new DWG converter
		dwgConverter, err := converter.NewDWGConverter(cfg.ODAConverterPath)
		if err != nil {
			return fmt.Errorf("failed to create DWG converter: %w", err)
		}

		// Convert DWG to DXF
		dxfFile, err := dwgConverter.ConvertToDXF(rootCmd, outputDir)
		if err != nil {
			return fmt.Errorf("conversion failed: %w", err)
		}

		// Parse the DXF file
		var dxfParser dxfparser.ParserInterface = dxfparser.NewParser()
		dxfData, err := dxfParser.ParseDXF(dxfFile)
		if err != nil {
			return fmt.Errorf("failed to parse DXF file: %w", err)
		}

		// Display the extracted information
		fmt.Println("Successfully extracted DXF information:")
		fmt.Printf("DXF Version: %s\n", dxfData.DXFVersion)
		fmt.Printf("Number of layers: %d\n", len(dxfData.Layers))
		for _, layer := range dxfData.Layers {
			onOff := "ON"
			if !layer.IsOn {
				onOff = "OFF"
			}
			frozen := ""
			if layer.IsFrozen {
				frozen = " (FROZEN)"
			}

			fmt.Printf("\nLayer: %s\n", layer.Name)
			fmt.Printf("  Color: %d, Line Type: %s, %s%s\n", layer.Color, layer.LineType, onOff, frozen)
		}

		return nil
	}

	return fmt.Errorf("unknown command: %s. Use 'extract' or 'tui'", command)

	// No duplicate code here

	return nil
}
