package tui

import "github.com/rivo/tview"

// NewStatusView creates a simple TextView for status/progress/error display
func NewStatusView() *tview.TextView {
	view := tview.NewTextView()
	view.SetDynamicColors(true)
	view.SetTextAlign(tview.AlignCenter)
	view.SetBorder(true)
	view.SetTitle("Status")
	return view
}
