package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/remym/go-dwg-extractor/pkg/data"
	"github.com/rivo/tview"
)

// App represents the main TUI application
type App struct {
	app       *tview.Application
	pages     *tview.Pages
	dxfView   *DXFView
	statusBar *tview.TextView
	testMode  bool // Indicates if app is running in test mode
}

// NewApp creates a new TUI application
func NewApp() *App {
	app := tview.NewApplication()

	// Create the main pages container
	pages := tview.NewPages()

	// Set up the application
	tuiApp := &App{
		app:      app,
		pages:    pages,
		testMode: false,
	}

	// Set up the main layout
	tuiApp.setupLayout()

	// Set the root
	app.SetRoot(pages, true).
		EnableMouse(true).
		SetFocus(pages)

	return tuiApp
}

// SetTestMode enables or disables test mode
// When in test mode, Run() will not start the event loop
// This is useful for testing to prevent hanging
func (a *App) SetTestMode(enabled bool) {
	a.testMode = enabled
}

// UpdateDXFData updates the DXF view with new data
func (a *App) UpdateDXFData(data *data.ExtractedData) {
	if a.testMode {
		// In test mode, update directly without queuing
		a.dxfView.Update(data)
	} else {
		// In normal mode, queue the update for the event loop
		a.app.QueueUpdateDraw(func() {
			a.dxfView.Update(data)
		})
	}
}

// setupLayout sets up the main application layout
func (a *App) setupLayout() {
	// Create the DXF view with the application instance
	a.dxfView = NewDXFView(a.app)

	// Create a flex layout that will contain our main content
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow)

	// Add a header
	header := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetText("DWG Extractor")
	header.SetBorder(true)

	// Add a status bar and store reference for updates
	a.statusBar = tview.NewTextView().
		SetTextAlign(tview.AlignLeft).
		SetText("Press Ctrl+C or Esc to exit")
	a.statusBar.SetBorder(true)

	// Add components to the flex layout
	flex.AddItem(header, 3, 1, false).
		AddItem(a.dxfView.GetLayout(), 0, 1, true).
		AddItem(a.statusBar, 3, 1, false)

	// Add the layout to the pages
	a.pages.AddPage("main", flex, true, true)
}

// Run starts the TUI application
func (a *App) Run() error {
	// Set up keyboard shortcuts
	a.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlC, tcell.KeyEsc:
			a.Stop()
			return nil
		}
		return event
	})

	// If in test mode, don't actually run the event loop
	if a.testMode {
		return nil
	}

	// Run the application
	if err := a.app.Run(); err != nil {
		return err
	}
	return nil
}

// Stop gracefully shuts down the TUI application
func (a *App) Stop() {
	a.app.Stop()
}

// App returns the underlying tview.Application instance
func (a *App) App() *tview.Application {
	return a.app
}

// GetLayout returns the main pages layout for the TUI
func (a *App) GetLayout() *tview.Pages {
	return a.pages
}

// ShowStatus updates the status bar with a status message
func (a *App) ShowStatus(message string) {
	if a.testMode {
		// In test mode, update directly without queuing
		a.statusBar.SetText("[yellow]" + message + "[-]")
	} else {
		// In normal mode, queue the update for the event loop
		a.app.QueueUpdateDraw(func() {
			a.statusBar.SetText("[yellow]" + message + "[-]")
		})
	}
}

// ShowError updates the status bar with an error message
func (a *App) ShowError(message string) {
	if a.testMode {
		// In test mode, update directly without queuing
		a.statusBar.SetText("[red]Error: " + message + "[-]")
	} else {
		// In normal mode, queue the update for the event loop
		a.app.QueueUpdateDraw(func() {
			a.statusBar.SetText("[red]Error: " + message + "[-]")
		})
	}
}
