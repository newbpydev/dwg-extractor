package cmd

import (
	"flag"
	"os"

	"github.com/gdamore/tcell/v2"
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
	var statusMsg string

	app := tui.NewApp()
	statusView := tui.NewStatusView()
	app.App().SetRoot(statusView, true)

	if args != nil && len(args) > 0 {
		statusView.SetText("[yellow]Converting DWG to DXF, please wait...[-]")
		app.App().Draw()

		// Load configuration
		cfg, err := config.LoadConfig()
		if err != nil {
			statusView.SetText("[red]Failed to load configuration:[-] " + err.Error() + "\nPress any key to exit.")
			app.App().Run()
			return err
		}

		// Create a new DWG converter
		dwgConverter, err := converter.NewDWGConverter(cfg.ODAConverterPath)
		if err != nil {
			statusView.SetText("[red]Failed to create DWG converter:[-] " + err.Error() + "\nPress any key to exit.")
			app.App().Run()
			return err
		}

		// Convert DWG to DXF
		dwgFile := args[0]
		statusView.SetText("[yellow]Converting: [-]" + dwgFile)
		app.App().Draw()
		dxfFile, err := dwgConverter.ConvertToDXF(dwgFile, "")
		if err != nil {
			statusView.SetText("[red]Conversion failed:[-] " + err.Error() + "\nPress any key to exit.")
			app.App().Run()
			return err
		}

		// Parse the DXF file
		statusView.SetText("[yellow]Parsing DXF file...[-]")
		app.App().Draw()
		var dxfParser dxfparser.ParserInterface = dxfparser.NewParser()
		dxfData, err = dxfParser.ParseDXF(dxfFile)
		if err != nil {
			statusView.SetText("[red]Failed to parse DXF file:[-] " + err.Error() + "\nPress any key to exit.")
			app.App().Run()
			return err
		}
		statusMsg = "[green]Conversion and parsing successful![-]"
	} else {
		// Use sample data if no file is provided
		statusMsg = "[yellow]No DWG file provided. Using sample data.[-]"
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

	statusView.SetText(statusMsg + "\n[gray]Press any key to continue...[-]")
	app.App().Draw()
	app.App().SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// Any key continues to main TUI
		app.App().SetRoot(app.GetLayout(), true)
		app.UpdateDXFData(dxfData)
		app.App().SetInputCapture(nil)
		return nil
	})

	return app.Run()
}


// ExecuteTUI executes the TUI command
func ExecuteTUI() error {
	if err := tuiCmd.Parse(os.Args[2:]); err != nil {
		return err
	}

	// If no DWG file is provided, use sample data
	return RunTUI(nil)
}
