package tui

import (
	"fmt"

	"github.com/remym/go-dwg-extractor/pkg/data"
	"github.com/rivo/tview"
)

// DXFView handles the display of DXF data
type DXFView struct {
	textView *tview.TextView
	layers   *tview.List
}

// NewDXFView creates a new DXF view
func NewDXFView() *DXFView {
	// Create the main text view
	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWrap(false)

	// Create the layers list
	layers := tview.NewList()
	layers.SetBorder(true).SetTitle("Layers")

	return &DXFView{
		textView: textView,
		layers:   layers,
	}
}

// Update updates the view with the given DXF data
func (v *DXFView) Update(data *data.ExtractedData) {
	// Clear the current content
	v.textView.Clear()

	// Display DXF version
	fmt.Fprintf(v.textView, "[green]DXF Version:[-] %s\n\n", data.DXFVersion)

	// Display number of layers
	fmt.Fprintf(v.textView, "[green]Layers:[-] %d\n\n", len(data.Layers))

	// Update layers list
	v.layers.Clear()
	for _, layer := range data.Layers {
		// Create a string representation of the layer
		onOff := "ON"
		if !layer.IsOn {
			onOff = "OFF"
		}
		frozen := ""
		if layer.IsFrozen {
			frozen = " (FROZEN)"
		}
		layerText := fmt.Sprintf("%s (Color: %d, %s%s)",
			layer.Name, layer.Color, onOff, frozen)

		v.layers.AddItem(layerText, "", 0, nil)
	}
}

// GetLayout returns a layout containing the DXF view and layers list
func (v *DXFView) GetLayout() *tview.Flex {
	// Create a flex layout with the layers list on the left and the main view on the right
	flex := tview.NewFlex()

	// Add the layers list (20% width)
	flex.AddItem(v.layers, 0, 1, false)

	// Add the main view (80% width)
	flex.AddItem(v.textView, 0, 4, true)

	return flex
}

// SetLayersChangedFunc sets the function to be called when a layer is selected
func (v *DXFView) SetLayersChangedFunc(handler func(index int, name string, secondaryText string, shortcut rune)) {
	v.layers.SetChangedFunc(handler)
}

// SetLayersSelectedFunc sets the function to be called when a layer is selected
func (v *DXFView) SetLayersSelectedFunc(handler func(index int, name string, secondaryText string, shortcut rune)) {
	v.layers.SetSelectedFunc(handler)
}
