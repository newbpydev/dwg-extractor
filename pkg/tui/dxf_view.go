package tui

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/remym/go-dwg-extractor/pkg/data"
	"github.com/rivo/tview"
)

// DXFView handles the display of DXF data
type DXFView struct {
	app               *tview.Application
	pages             *tview.Pages
	textView          *tview.TextView
	layers            *tview.List
	entityList        *tview.List
	searchInput       *tview.InputField
	data              *data.ExtractedData
	currentLayerIndex int

	// Navigation components
	navigator           Navigator
	layersNavigator     ListNavigator
	entitiesNavigator   ListNavigator
	categorySelector    CategorySelector
	itemSelector        ItemSelector
	breadcrumbNavigator BreadcrumbNavigator
}

// NewDXFView creates a new DXF view
func NewDXFView(app *tview.Application) *DXFView {
	// Create the main text view for layer details
	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWrap(false)

	// Create the layers list
	layers := tview.NewList()
	layers.SetBorder(true).SetTitle("Layers")

	// Create the entity list
	entityList := tview.NewList()
	entityList.SetBorder(true).SetTitle("Entities")

	// Create search input
	searchInput := tview.NewInputField().
		SetLabel("Search: ").
		SetFieldWidth(30).
		SetPlaceholder("Type to filter layers... (Space/t: toggle visibility)")

	// Create pages container
	pages := tview.NewPages()

	view := &DXFView{
		app:               app,
		pages:             pages,
		textView:          textView,
		layers:            layers,
		entityList:        entityList,
		searchInput:       searchInput,
		currentLayerIndex: -1,
	}

	// Set up search input handler
	searchInput.SetChangedFunc(func(text string) {
		view.FilterLayers(text)
	})

	// Set up keyboard navigation
	view.setupKeybindings()

	// Initialize navigation components
	view.navigator = NewTUINavigator(app, searchInput, layers, entityList)
	view.layersNavigator = NewTUIListNavigator(layers)
	view.entitiesNavigator = NewTUIListNavigator(entityList)
	view.categorySelector = NewEnhancedCategorySelector(view)
	view.itemSelector = NewEnhancedItemSelector(view)
	view.breadcrumbNavigator = NewTUIBreadcrumbNavigator()

	return view
}

// Update updates the view with the given DXF data
func (v *DXFView) Update(data *data.ExtractedData) {
	v.data = data
	v.currentLayerIndex = -1

	// Clear the current content
	v.textView.Clear()

	// Display DXF version
	fmt.Fprintf(v.textView, "[green]DXF Version:[-] %s\n\n", data.DXFVersion)

	// Display number of layers
	fmt.Fprintf(v.textView, "[green]Layers:[-] %d\n\n", len(data.Layers))

	// Update layers list
	v.updateLayersList()

	// Show the layers view
	v.showLayersView()
}

// updateLayersList updates the layers list with current data
func (v *DXFView) updateLayersList() {
	v.layers.Clear()
	for i, layer := range v.data.Layers {
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

		// Store the layer index as a reference
		index := i
		v.layers.AddItem(layerText, "", 0, func() {
			v.showLayerDetails(index)
		})
	}
}

// showLayersView shows the layers list view
func (v *DXFView) showLayersView() {
	// Create a flex container for the search input and layers list
	listFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	listFlex.AddItem(v.searchInput, 1, 0, false)
	listFlex.AddItem(v.layers, 0, 1, true)

	// Create the main flex container
	flex := tview.NewFlex().SetDirection(tview.FlexRow)
	flex.AddItem(v.textView, 0, 1, false)
	flex.AddItem(listFlex, 0, 3, true)

	v.pages.SwitchToPage("layers")
	v.app.SetFocus(v.searchInput)
}

// showLayerDetails shows the details for a specific layer
func (v *DXFView) showLayerDetails(layerIndex int) {
	if layerIndex < 0 || layerIndex >= len(v.data.Layers) {
		return
	}

	v.currentLayerIndex = layerIndex
	layer := v.data.Layers[layerIndex]

	// Update the entity list
	v.entityList.Clear()

	// Add a back item
	v.entityList.AddItem("← Back to Layers", "", 'b', func() {
		v.showLayersView()
	})

	// Add entities for this layer
	entityCount := 0
	for _, entity := range layer.Entities {
		switch e := entity.(type) {
		case *data.LineInfo:
			v.entityList.AddItem(
				fmt.Sprintf("Line (%.1f,%.1f) to (%.1f,%.1f)",
					e.StartPoint.X, e.StartPoint.Y, e.EndPoint.X, e.EndPoint.Y),
				fmt.Sprintf("Layer: %s, Color: %d", e.Layer, e.Color),
				0, nil)
			entityCount++
		case *data.CircleInfo:
			v.entityList.AddItem(
				fmt.Sprintf("Circle center:(%.1f,%.1f) radius:%.1f",
					e.Center.X, e.Center.Y, e.Radius),
				fmt.Sprintf("Layer: %s, Color: %d", e.Layer, e.Color),
				0, nil)
			entityCount++
		case *data.TextInfo:
			v.entityList.AddItem(
				fmt.Sprintf("Text: %s at (%.1f,%.1f)",
					e.Value, e.InsertionPoint.X, e.InsertionPoint.Y),
				fmt.Sprintf("Layer: %s, Height: %.1f", e.Layer, e.Height),
				0, nil)
			entityCount++
		case *data.PolylineInfo:
			v.entityList.AddItem(
				fmt.Sprintf("Polyline with %d points", len(e.Points)),
				fmt.Sprintf("Layer: %s, Color: %d, Closed: %v", e.Layer, e.Color, e.IsClosed),
				0, nil)
			entityCount++
		case *data.BlockInfo:
			v.entityList.AddItem(
				fmt.Sprintf("Block: %s at (%.1f,%.1f)",
					e.Name, e.InsertionPoint.X, e.InsertionPoint.Y),
				fmt.Sprintf("Layer: %s, Rotation: %.1f", e.Layer, e.Rotation),
				0, nil)
			entityCount++
		default:
			// Handle any other entity types
			v.entityList.AddItem(
				fmt.Sprintf("Entity: %T", entity),
				fmt.Sprintf("Layer: %s", entity.GetLayer()),
				0, nil)
			entityCount++
		}
	}

	// Update the text view with layer details
	v.textView.Clear()
	fmt.Fprintf(v.textView, "[green]Layer:[-] %s\n", layer.Name)
	fmt.Fprintf(v.textView, "[green]Color:[-] %d\n", layer.Color)
	fmt.Fprintf(v.textView, "[green]Status:[-] %s\n", map[bool]string{true: "ON", false: "OFF"}[layer.IsOn])
	fmt.Fprintf(v.textView, "[green]Frozen:[-] %v\n", layer.IsFrozen)
	fmt.Fprintf(v.textView, "[green]Line Type:[-] %s\n", layer.LineType)
	fmt.Fprintf(v.textView, "[green]Entities:[-] %d\n\n", entityCount)

	// Show the entities view
	v.showEntitiesView()
}

// showEntitiesView shows the entities list view
func (v *DXFView) showEntitiesView() {
	// Create a flex layout with the entities list and details
	flex := tview.NewFlex().
		AddItem(v.entityList, 0, 1, true).
		AddItem(v.textView, 0, 1, false)

	// Add or update the entities page
	v.pages.AddAndSwitchToPage("entities", flex, true)
}

// setupKeybindings sets up keyboard shortcuts
func (v *DXFView) setupKeybindings() {
	// Handle key events for the search input
	v.searchInput.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyDown, tcell.KeyTab:
			// Move focus to the layers list
			v.app.SetFocus(v.layers)
			return nil
		case tcell.KeyEsc:
			// Clear search and reset focus
			v.searchInput.SetText("")
			v.app.SetFocus(v.layers)
			return nil
		}
		return event
	})

	// Handle key events for the layers list
	v.layers.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEnter:
			if v.layers.GetItemCount() > 0 {
				index := v.layers.GetCurrentItem()
				v.showLayerDetails(index)
			}
			return nil
		case tcell.KeyEsc, tcell.KeyBackspace, tcell.KeyBackspace2:
			// Move focus back to search
			v.app.SetFocus(v.searchInput)
			return nil
		case tcell.KeyRune:
			// Space or 't' toggles layer visibility
			if event.Rune() == ' ' || event.Rune() == 't' || event.Rune() == 'T' {
				idx := v.layers.GetCurrentItem()
				v.ToggleLayerVisibility(idx)
				return nil
			}
			// If a letter or number is pressed, focus on search and type
			if (event.Rune() >= 'a' && event.Rune() <= 'z') ||
				(event.Rune() >= 'A' && event.Rune() <= 'Z') ||
				(event.Rune() >= '0' && event.Rune() <= '9') ||
				event.Rune() == ' ' || event.Rune() == ':' {
				v.app.SetFocus(v.searchInput)
				// Append the pressed key to the search input
				v.searchInput.SetText(v.searchInput.GetText() + string(event.Rune()))
				return nil
			}
		}
		return event
	})

	// Handle key events for the entity list
	v.entityList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEsc, tcell.KeyBackspace, tcell.KeyBackspace2:
			v.showLayersView()
			return nil
		}
		return event
	})
}

// GetLayout returns the pages container for the DXF view
func (v *DXFView) GetLayout() *tview.Pages {
	return v.pages
}

// SetLayersChangedFunc sets the function to be called when a layer is selected
func (v *DXFView) SetLayersChangedFunc(handler func(index int, name string, secondaryText string, shortcut rune)) {
	v.layers.SetChangedFunc(handler)
}

// SetLayersSelectedFunc sets the function to be called when a layer is selected
func (v *DXFView) SetLayersSelectedFunc(handler func(index int, name string, secondaryText string, shortcut rune)) {
	v.layers.SetSelectedFunc(handler)
}

// ToggleLayerVisibility toggles the visibility (IsOn) of the layer at the given visible index.
func (v *DXFView) ToggleLayerVisibility(visibleIndex int) {
	if v.data == nil || v.layers.GetItemCount() == 0 || visibleIndex < 0 || visibleIndex >= v.layers.GetItemCount() {
		return
	}
	// Find the actual layer index in v.data.Layers by matching name
	mainText, _ := v.layers.GetItemText(visibleIndex)
	// Extract the layer name from the display string (before first ' (')
	name := mainText
	if idx := strings.Index(mainText, " ("); idx > 0 {
		name = mainText[:idx]
	}
	for i := range v.data.Layers {
		if v.data.Layers[i].Name == name {
			// Don't toggle frozen layers
			if !v.data.Layers[i].IsFrozen {
				v.data.Layers[i].IsOn = !v.data.Layers[i].IsOn
			}
			break
		}
	}
	// Re-filter or update the list to reflect the change
	if v.searchInput != nil && v.searchInput.GetText() != "" {
		v.FilterLayers(v.searchInput.GetText())
	} else {
		v.updateLayersList()
	}
}

// FilterLayers filters the layers list based on the provided query string.
// The query can be:
// - A simple string to filter by layer name (case-insensitive)
// - "on:true" or "on:false" to filter by layer on/off status
// - "frozen:true" or "frozen:false" to filter by frozen status
// - An empty string to clear all filters
func (v *DXFView) FilterLayers(query string) {
	if v.data == nil {
		return
	}

	// Store the current scroll position
	_, currentOffset := v.layers.GetOffset()

	// Clear the current list
	v.layers.Clear()

	// If query is empty, show all layers
	if query == "" {
		v.updateLayersList()
		v.layers.SetOffset(0, currentOffset)
		return
	}

	// Convert query to lowercase for case-insensitive comparison
	query = strings.ToLower(query)

	// Check for special filter types
	var filterFunc func(layer data.LayerInfo) bool

	switch {
	case strings.HasPrefix(query, "on:true"):
		filterFunc = func(layer data.LayerInfo) bool {
			return layer.IsOn
		}
	case strings.HasPrefix(query, "on:false"):
		filterFunc = func(layer data.LayerInfo) bool {
			return !layer.IsOn
		}
	case strings.HasPrefix(query, "frozen:true"):
		filterFunc = func(layer data.LayerInfo) bool {
			return layer.IsFrozen
		}
	case strings.HasPrefix(query, "frozen:false"):
		filterFunc = func(layer data.LayerInfo) bool {
			return !layer.IsFrozen
		}
	default:
		// Filter by name (case-insensitive)
		filterFunc = func(layer data.LayerInfo) bool {
			return strings.Contains(strings.ToLower(layer.Name), query)
		}
	}

	// Filter and add layers
	for i, layer := range v.data.Layers {
		if filterFunc(layer) {
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

			// Store the layer index as a reference
			index := i
			v.layers.AddItem(layerText, "", 0, func() {
				v.showLayerDetails(index)
			})
		}
	}

	// Restore scroll position if possible
	v.layers.SetOffset(0, currentOffset)
}

// Navigation Getter Methods

// GetNavigator returns the main navigator for pane switching
func (v *DXFView) GetNavigator() Navigator {
	return v.navigator
}

// GetListNavigator returns the list navigator for the specified list type
func (v *DXFView) GetListNavigator(listType string) ListNavigator {
	switch listType {
	case "layers":
		return v.layersNavigator
	case "entities":
		return v.entitiesNavigator
	default:
		return v.layersNavigator // Default to layers
	}
}

// GetCategorySelector returns the category selector
func (v *DXFView) GetCategorySelector() CategorySelector {
	return v.categorySelector
}

// GetCurrentItemList returns the currently active item list
func (v *DXFView) GetCurrentItemList() *tview.List {
	// For now, return the layers list as the default
	// TODO: Implement logic to determine the current active list
	return v.layers
}

// GetDetailsView returns the details text view
func (v *DXFView) GetDetailsView() *tview.TextView {
	return v.textView
}

// GetItemSelector returns the item selector
func (v *DXFView) GetItemSelector() ItemSelector {
	return v.itemSelector
}

// GetDetailsPane returns the details pane (same as details view for now)
func (v *DXFView) GetDetailsPane() *tview.TextView {
	return v.textView
}

// GetBreadcrumbNavigator returns the breadcrumb navigator
func (v *DXFView) GetBreadcrumbNavigator() BreadcrumbNavigator {
	return v.breadcrumbNavigator
}
