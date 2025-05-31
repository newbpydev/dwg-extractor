package cmd

import (
	"flag"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/remym/go-dwg-extractor/pkg/config"
	"github.com/remym/go-dwg-extractor/pkg/converter"
	"github.com/remym/go-dwg-extractor/pkg/data"
	"github.com/remym/go-dwg-extractor/pkg/dxfparser"
	"github.com/remym/go-dwg-extractor/pkg/tui"
)

// tuiCmd represents the tui command
var tuiCmd = flag.NewFlagSet("tui", flag.ExitOnError)
var tuiOutputDir string
var tuiFileFlag string

func init() {
	tuiCmd.StringVar(&tuiOutputDir, "output", "", "Output directory for converted files (default: same as input file)")
	tuiCmd.StringVar(&tuiFileFlag, "file", "", "Path to the DWG file to process")
}

// RunTUI runs the TUI command
func RunTUI(args []string) error {
	app := tui.NewApp()

	// Start the app and handle initialization after event loop starts
	go func() {
		// Wait a moment for the app to start
		time.Sleep(100 * time.Millisecond)

		if args != nil && len(args) > 0 {
			dwgFile := args[0]

			// Check if the file is a DXF file (for testing)
			if strings.ToLower(filepath.Ext(dwgFile)) == ".dxf" {
				app.ShowStatus("Parsing DXF file: " + dwgFile)
				var dxfParser dxfparser.ParserInterface = dxfparser.NewParser()
				dxfData, err := dxfParser.ParseDXF(dwgFile)
				if err != nil {
					app.ShowError("Failed to parse DXF file: " + err.Error())
					return
				}
				app.ShowStatus("DXF parsing successful!")
				app.UpdateDXFData(dxfData)
				return
			}

			// Process DWG file
			// Load configuration
			cfg, err := config.LoadConfig()
			if err != nil {
				app.ShowError("Failed to load configuration: " + err.Error())
				return
			}

			// Create a new DWG converter
			dwgConverter, err := converter.NewDWGConverter(cfg.ODAConverterPath)
			if err != nil {
				app.ShowError("Failed to create DWG converter: " + err.Error())
				return
			}

			// Convert DWG to DXF
			app.ShowStatus("Converting: " + dwgFile)

			// Determine output directory
			outputDir := tuiOutputDir
			if outputDir == "" {
				// If no output directory specified, use a temp directory
				tempDir, err := os.MkdirTemp("", "dwg-extractor-*")
				if err != nil {
					app.ShowError("Failed to create temp directory: " + err.Error())
					return
				}
				outputDir = tempDir
			}

			dxfFile, err := dwgConverter.ConvertToDXF(dwgFile, outputDir)
			if err != nil {
				app.ShowError("Conversion failed: " + err.Error())
				return
			}

			// Parse the DXF file
			app.ShowStatus("Parsing DXF file...")
			var dxfParser dxfparser.ParserInterface = dxfparser.NewParser()
			dxfData, err := dxfParser.ParseDXF(dxfFile)
			if err != nil {
				app.ShowError("Failed to parse DXF file: " + err.Error())
				return
			}

			// Update the UI with the parsed data
			app.ShowStatus("Conversion and parsing successful!")
			app.UpdateDXFData(dxfData)
		} else {
			// Use sample data if no file is provided
			dxfData := &data.ExtractedData{
				DXFVersion: "R2020 (Sample Data)",
				Layers: []data.LayerInfo{
					{Name: "0", IsOn: true, IsFrozen: false, Color: 7, LineType: "CONTINUOUS"},
					{Name: "Walls", IsOn: true, IsFrozen: false, Color: 1, LineType: "CONTINUOUS"},
					{Name: "Doors", IsOn: true, IsFrozen: false, Color: 2, LineType: "DASHED"},
					{Name: "Windows", IsOn: true, IsFrozen: true, Color: 3, LineType: "HIDDEN"},
				},
			}
			app.ShowStatus("No DWG file provided. Using sample data.")
			app.UpdateDXFData(dxfData)
		}
	}()

	return app.Run()
}

// ExecuteTUI executes the TUI command
func ExecuteTUI() error {
	if err := tuiCmd.Parse(os.Args[2:]); err != nil {
		return err
	}

	// Check if a file argument was provided via flag or positional argument
	var args []string
	if tuiFileFlag != "" {
		args = []string{tuiFileFlag}
	} else if len(tuiCmd.Args()) > 0 {
		args = tuiCmd.Args()
	}

	if len(args) > 0 {
		return RunTUI(args)
	}

	// If no DWG file is provided, use sample data
	return RunTUI(nil)
}
