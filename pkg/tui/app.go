package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/remym/go-dwg-extractor/pkg/data"
)

// App represents the main TUI application
type App struct {
	app     *tview.Application
	pages   *tview.Pages
	dxfView *DXFView
}

// NewApp creates a new TUI application
func NewApp() *App {
	app := tview.NewApplication()

	// Create the main pages container
	pages := tview.NewPages()

	// Set up the application
	tuiApp := &App{
		app:   app,
		pages: pages,
	}

	// Set up the main layout
	tuiApp.setupLayout()

	// Set the root
	app.SetRoot(pages, true).
		EnableMouse(true).
		SetFocus(pages)

	return tuiApp
}

// UpdateDXFData updates the DXF view with new data
func (a *App) UpdateDXFData(data *data.ExtractedData) {
	a.app.QueueUpdateDraw(func() {
		a.dxfView.Update(data)
	})
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

	// Add a status bar
	status := tview.NewTextView().
		SetTextAlign(tview.AlignLeft).
		SetText("Press Ctrl+C or Esc to exit")
	status.SetBorder(true)

	// Add components to the flex layout
	flex.AddItem(header, 3, 1, false).
		AddItem(a.dxfView.GetLayout(), 0, 1, true).
		AddItem(status, 3, 1, false)

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
