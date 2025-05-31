package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/remym/go-dwg-extractor/pkg/converter"
	"github.com/remym/go-dwg-extractor/pkg/data"
	"github.com/remym/go-dwg-extractor/pkg/dxfparser"
	"github.com/remym/go-dwg-extractor/pkg/tui"
	"github.com/remym/go-dwg-extractor/pkg/config"
)

// tuiCmd represents the tui command
var tuiCmd = flag.NewFlagSet("tui", flag.ExitOnError)
var tuiOutputDir string

func init() {
	tuiCmd.StringVar(&tuiOutputDir, "output", "", "Output directory for converted files (default: same as input file)")
}

// RunTUI runs the TUI command
func RunTUI(args []string) error {
	var dxfData *data.ExtractedData

	if args != nil && len(args) > 0 {
		// Load configuration
		cfg, err := config.LoadConfig()
		if err != nil {
			return fmt.Errorf("failed to load configuration: %w", err)
		}

		// Create a new DWG converter
		dwgConverter, err := converter.NewDWGConverter(cfg.ODAConverterPath)
		if err != nil {
			return fmt.Errorf("failed to create DWG converter: %w", err)
		}

		// Convert DWG to DXF
		dwgFile := args[0]
		dxfFile, err := dwgConverter.ConvertToDXF(dwgFile, "")
		if err != nil {
			return fmt.Errorf("conversion failed: %w", err)
		}

		// Parse the DXF file
		var dxfParser dxfparser.ParserInterface = dxfparser.NewParser()
		dxfData, err = dxfParser.ParseDXF(dxfFile)
		if err != nil {
			return fmt.Errorf("failed to parse DXF file: %w", err)
		}
	} else {
		// Use sample data if no file is provided
		dxfData = &data.ExtractedData{
			DXFVersion: "R2020 (Sample Data)",
			Layers: []data.LayerInfo{
				{Name: "0", IsOn: true, IsFrozen: false, Color: 7, LineType: "CONTINUOUS"},
				{Name: "Walls", IsOn: true, IsFrozen: false, Color: 1, LineType: "CONTINUOUS"},
				{Name: "Doors", IsOn: true, IsFrozen: false, Color: 2, LineType: "DASHED"},
				{Name: "Windows", IsOn: true, IsFrozen: true, Color: 3, LineType: "HIDDEN"},
			},
		}
	}

	// Create the TUI app
	app := tui.NewApp()
	
	// Start the application in a goroutine
	done := make(chan error, 1)
	go func() {
		done <- app.Run()
	}()
	
	// Wait for the application to be fully started
	app.App().Draw()
	
	// Update the UI with the DXF data
	app.UpdateDXFData(dxfData)
	
	// Wait for the application to finish
	return <-done
}

// ExecuteTUI executes the TUI command
func ExecuteTUI() error {
	if err := tuiCmd.Parse(os.Args[2:]); err != nil {
		return err
	}

	// If no DWG file is provided, use sample data
	return RunTUI(nil)
}
