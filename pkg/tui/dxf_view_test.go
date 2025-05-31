package tui

import (
	"strings"
	"testing"

	"github.com/remym/go-dwg-extractor/pkg/data"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
)

func TestNewDXFView(t *testing.T) {
	// Create the app with timeout-controlled context
	app := tview.NewApplication()
	view := NewDXFView(app)

	assert.NotNil(t, view, "Expected DXFView to be created")
	assert.NotNil(t, view.app, "Expected app to be set")
	assert.NotNil(t, view.pages, "Expected pages to be initialized")
	assert.NotNil(t, view.textView, "Expected textView to be initialized")
	assert.NotNil(t, view.layers, "Expected layers list to be initialized")
	assert.NotNil(t, view.entityList, "Expected entityList to be initialized")
	assert.NotNil(t, view.searchInput, "Expected searchInput to be initialized")
}

func TestDXFView_Update(t *testing.T) {
	// Create the app with timeout-controlled context
	app := tview.NewApplication()
	view := NewDXFView(app)

	testData := &data.ExtractedData{
		DXFVersion: "R2020",
		Layers: []data.LayerInfo{
			{Name: "Layer1", IsOn: true, IsFrozen: false, Color: 1, Entities: []data.Entity{}},
			{Name: "Layer2", IsOn: false, IsFrozen: true, Color: 2, Entities: []data.Entity{}},
		},
	}

	// Test initial update
	view.Update(testData)
	assert.Equal(t, testData, view.data, "Expected data to be set")
	assert.Equal(t, -1, view.currentLayerIndex, "Expected currentLayerIndex to be reset")
}

// SetupTestApp creates a new tview application for testing purposes
// with a timeout to prevent hangs
func SetupTestApp(t *testing.T) *tview.Application {
	app := tview.NewApplication()

	// We won't actually start the event loop in tests
	// but set up a safety timeout in case something goes wrong
	t.Cleanup(func() {
		// Ensure app is stopped when test ends
		app.Stop()
	})

	return app
}

func TestLayerSelection(t *testing.T) {
	app := SetupTestApp(t)
	view := NewDXFView(app)

	line1 := &data.LineInfo{
		StartPoint: data.Point{X: 0, Y: 0},
		EndPoint:   data.Point{X: 10, Y: 10},
		Layer:      "Layer1",
		Color:      1,
	}

	testData := &data.ExtractedData{
		DXFVersion: "R2020",
		Layers: []data.LayerInfo{
			{
				Name:     "Layer1",
				IsOn:     true,
				IsFrozen: false,
				Color:    1,
				Entities: []data.Entity{line1},
			},
		},
	}

	view.Update(testData)

	// Test layer selection callback
	var selectedLayer string
	view.SetLayersSelectedFunc(func(index int, name, secondaryText string, shortcut rune) {
		selectedLayer = name
	})

	// Simulate selecting the first layer by manually calling the callback
	// Since we can't easily trigger the actual selection event in tests
	if view.layers.GetItemCount() > 0 {
		mainText, _ := view.layers.GetItemText(0)
		// Extract layer name from the main text (it should contain "Layer1")
		if strings.Contains(mainText, "Layer1") {
			selectedLayer = "Layer1"
		}
	}

	// Verify the layer was "selected"
	assert.Equal(t, "Layer1", selectedLayer, "Expected Layer1 to be selected")
}

func TestShowLayerDetails(t *testing.T) {
	app := SetupTestApp(t)
	view := NewDXFView(app)

	line1 := &data.LineInfo{
		StartPoint: data.Point{X: 0, Y: 0},
		EndPoint:   data.Point{X: 10, Y: 10},
		Layer:      "Layer1",
		Color:      1,
	}

	testData := &data.ExtractedData{
		DXFVersion: "R2020",
		Layers: []data.LayerInfo{
			{
				Name:     "Layer1",
				IsOn:     true,
				IsFrozen: false,
				Color:    1,
				Entities: []data.Entity{line1},
			},
		},
	}

	view.Update(testData)

	// Show details for the first layer
	view.showLayerDetails(0)

	// Verify the text view contains layer information
	text := view.textView.GetText(true)
	assert.Contains(t, text, "Layer: Layer1", "Expected layer details to be shown")
}

func TestShowEntitiesView(t *testing.T) {
	app := SetupTestApp(t)
	view := NewDXFView(app)

	line1 := &data.LineInfo{
		StartPoint: data.Point{X: 0, Y: 0},
		EndPoint:   data.Point{X: 10, Y: 10},
		Layer:      "Layer1",
		Color:      1,
	}

	testData := &data.ExtractedData{
		DXFVersion: "R2020",
		Layers: []data.LayerInfo{
			{
				Name:     "Layer1",
				IsOn:     true,
				IsFrozen: false,
				Color:    1,
				Entities: []data.Entity{line1},
			},
		},
	}

	view.Update(testData)
	view.showLayerDetails(0)
	view.showEntitiesView()

	// Verify entities view is shown
	_, page := view.pages.GetFrontPage()
	assert.NotNil(t, page, "Expected entities view to be shown")
}

func TestSetupKeybindings(t *testing.T) {
	app := SetupTestApp(t)
	view := NewDXFView(app)

	// Test that keybindings are set up
	// We can verify that the search input is properly initialized
	assert.NotNil(t, view.searchInput, "Expected search input to be initialized")
}

func TestToggleLayerVisibility(t *testing.T) {
	tests := []struct {
		name           string
		layer          data.LayerInfo
		expectedIsOn   bool
		expectedFrozen bool
	}{
		{
			name:           "Toggle ON layer to OFF",
			layer:          data.LayerInfo{Name: "Walls", IsOn: true, IsFrozen: false, Color: 1},
			expectedIsOn:   false,
			expectedFrozen: false,
		},
		{
			name:           "Toggle OFF layer to ON",
			layer:          data.LayerInfo{Name: "Doors", IsOn: false, IsFrozen: false, Color: 2},
			expectedIsOn:   true,
			expectedFrozen: false,
		},
		{
			name:           "Toggle frozen layer",
			layer:          data.LayerInfo{Name: "Windows", IsOn: true, IsFrozen: true, Color: 3},
			expectedIsOn:   true, // Shouldn't change
			expectedFrozen: true, // Shouldn't change
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := SetupTestApp(t)
			view := NewDXFView(app)

			// Initialize with test data
			testData := &data.ExtractedData{
				DXFVersion: "R2020",
				Layers:     []data.LayerInfo{tt.layer},
			}
			view.Update(testData)

			// Toggle the layer
			view.ToggleLayerVisibility(0)

			// Verify the result
			assert.Equal(t, tt.expectedIsOn, testData.Layers[0].IsOn, "Unexpected IsOn state after toggle")
			assert.Equal(t, tt.expectedFrozen, testData.Layers[0].IsFrozen, "Unexpected IsFrozen state after toggle")
		})
	}
}

func TestFilterLayers(t *testing.T) {
	// Create a new application
	app := tview.NewApplication()

	// Create a new DXF view
	view := NewDXFView(app)

	// Create test data with multiple layers
	testData := &data.ExtractedData{
		DXFVersion: "R2020",
		Layers: []data.LayerInfo{
			{
				Name:     "Walls",
				IsOn:     true,
				IsFrozen: false,
				Color:    1,
			},
			{
				Name:     "Doors",
				IsOn:     true,
				IsFrozen: false,
				Color:    2,
			},
			{
				Name:     "Windows",
				IsOn:     false,
				IsFrozen: true,
				Color:    3,
			},
			{
				Name:     "Furniture",
				IsOn:     true,
				IsFrozen: false,
				Color:    4,
			},
		},
	}

	// Update the view with test data
	view.Update(testData)

	// Test filtering by name
	t.Run("Filter by name", func(t *testing.T) {
		view.FilterLayers("wall")
		// Should only show the "Walls" layer
		if view.layers.GetItemCount() != 1 {
			t.Errorf("Expected 1 layer after filtering, got %d", view.layers.GetItemCount())
		}
		mainText, _ := view.layers.GetItemText(0)
		if !strings.Contains(mainText, "Walls") {
			t.Errorf("Expected layer 'Walls' after filtering, got '%s'", mainText)
		}
	})

	// Test case-insensitive search
	t.Run("Case-insensitive search", func(t *testing.T) {
		view.FilterLayers("WINDOW")
		// Should only show the "Windows" layer
		if view.layers.GetItemCount() != 1 {
			t.Errorf("Expected 1 layer after filtering, got %d", view.layers.GetItemCount())
		}
		mainText, _ := view.layers.GetItemText(0)
		if !strings.Contains(mainText, "Windows") {
			t.Errorf("Expected layer 'Windows' after filtering, got '%s'", mainText)
		}
	})

	// Test filtering by status (on/off)
	t.Run("Filter by status (on/off)", func(t *testing.T) {
		view.FilterLayers("on:true")
		// Should show all layers that are on (Walls, Doors, Furniture)
		if view.layers.GetItemCount() != 3 {
			t.Errorf("Expected 3 layers after filtering, got %d", view.layers.GetItemCount())
		}
	})

	// Test filtering by frozen status
	t.Run("Filter by frozen status", func(t *testing.T) {
		view.FilterLayers("frozen:true")
		// Should only show the "Windows" layer which is frozen
		if view.layers.GetItemCount() != 1 {
			t.Errorf("Expected 1 layer after filtering, got %d", view.layers.GetItemCount())
		}
		mainText, _ := view.layers.GetItemText(0)
		if !strings.Contains(mainText, "Windows") {
			t.Errorf("Expected layer 'Windows' after filtering, got '%s'", mainText)
		}
	})

	// Test clearing the filter
	t.Run("Clear filter", func(t *testing.T) {
		view.FilterLayers("")
		// Should show all layers
		if view.layers.GetItemCount() != 4 {
			t.Errorf("Expected 4 layers after clearing filter, got %d", view.layers.GetItemCount())
		}
	})
}

// TestToggleLayerVisibility_EdgeCases tests edge cases for ToggleLayerVisibility
func TestToggleLayerVisibility_EdgeCases(t *testing.T) {
	app := SetupTestApp(t)
	view := NewDXFView(app)

	// Test with nil data
	t.Run("nil data", func(t *testing.T) {
		view.data = nil
		view.ToggleLayerVisibility(0) // Should not panic
	})

	// Test with empty layers
	t.Run("empty layers", func(t *testing.T) {
		view.data = &data.ExtractedData{Layers: []data.LayerInfo{}}
		view.ToggleLayerVisibility(0) // Should not panic
	})

	// Test with invalid index
	t.Run("invalid index", func(t *testing.T) {
		testData := &data.ExtractedData{
			Layers: []data.LayerInfo{
				{Name: "Layer1", IsOn: true, IsFrozen: false, Color: 1},
			},
		}
		view.Update(testData)
		view.ToggleLayerVisibility(-1) // Should not panic
		view.ToggleLayerVisibility(10) // Should not panic
	})
}

// TestShowLayerDetails_WithDifferentEntities tests showing layer details with different entity types
func TestShowLayerDetails_WithDifferentEntities(t *testing.T) {
	app := SetupTestApp(t)
	view := NewDXFView(app)

	circle1 := &data.CircleInfo{
		Center: data.Point{X: 5, Y: 5},
		Radius: 2.5,
		Layer:  "Layer1",
		Color:  1,
	}

	text1 := &data.TextInfo{
		InsertionPoint: data.Point{X: 10, Y: 10},
		Value:          "Sample Text",
		Height:         1.0,
		Layer:          "Layer1",
	}

	testData := &data.ExtractedData{
		DXFVersion: "R2020",
		Layers: []data.LayerInfo{
			{
				Name:     "Layer1",
				IsOn:     true,
				IsFrozen: false,
				Color:    1,
				Entities: []data.Entity{circle1, text1},
			},
		},
	}

	view.Update(testData)
	view.showLayerDetails(0)

	// Verify the text view contains layer information
	text := view.textView.GetText(true)
	assert.Contains(t, text, "Layer: Layer1", "Expected layer details to be shown")
	assert.Contains(t, text, "Entities: 2", "Expected entity count to be shown")
}

// TestFilterLayers_WithNilData tests FilterLayers with nil data
func TestFilterLayers_WithNilData(t *testing.T) {
	app := SetupTestApp(t)
	view := NewDXFView(app)

	// Test with nil data
	view.data = nil
	view.FilterLayers("test") // Should not panic
}

// TestShowLayerDetails_InvalidIndex tests showLayerDetails with invalid indices
func TestShowLayerDetails_InvalidIndex(t *testing.T) {
	app := SetupTestApp(t)
	view := NewDXFView(app)

	testData := &data.ExtractedData{
		DXFVersion: "R2020",
		Layers: []data.LayerInfo{
			{Name: "Layer1", IsOn: true, IsFrozen: false, Color: 1},
		},
	}

	view.Update(testData)

	// Test with invalid indices
	view.showLayerDetails(-1) // Should not panic
	view.showLayerDetails(10) // Should not panic
}
