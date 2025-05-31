package tui

import (
	"strings"
	"testing"

	"github.com/remym/go-dwg-extractor/pkg/data"
	"github.com/rivo/tview"
)

func TestDXFView_Update(t *testing.T) {
	// Create a new application
	app := tview.NewApplication()
	
	// Create a new DXF view
	view := NewDXFView(app)

	// Create test data
	testData := &data.ExtractedData{
		DXFVersion: "R2020",
		Layers: []data.LayerInfo{
			{Name: "Layer1", IsOn: true, IsFrozen: false, Color: 1},
			{Name: "Layer2", IsOn: false, IsFrozen: true, Color: 2},
		},
	}

	// Update the view with test data
	view.Update(testData)
	
	// Basic test to verify the view was updated
	// We'll add more specific tests once we have getters for the view
	if view == nil {
		t.Error("Expected DXFView to be initialized")
	}
}

func TestLayerSelection(t *testing.T) {
	// Create a new application
	app := tview.NewApplication()
	
	// Create a new DXF view
	view := NewDXFView(app)

	// Create test data with entities
	testData := &data.ExtractedData{
		DXFVersion: "R2020",
		Layers: []data.LayerInfo{
			{
				Name:     "Layer1",
				IsOn:     true,
				IsFrozen: false,
				Color:    1,
				Entities: []data.Entity{
					&data.LineInfo{
						StartPoint: data.Point{X: 0, Y: 0},
						EndPoint:   data.Point{X: 10, Y: 10},
						Layer:      "Layer1",
						Color:      1,
					},
				},
			},
		},
	}

	// Update the view with test data
	view.Update(testData)

	// Simulate selecting the first layer
	// We'll need to add functionality to handle layer selection
}

// Add more test cases as needed

func TestToggleLayerVisibility(t *testing.T) {
	app := tview.NewApplication()
	view := NewDXFView(app)

	// Test data
	layers := []data.LayerInfo{
		{Name: "Walls", IsOn: true, IsFrozen: false, Color: 1},
		{Name: "Doors", IsOn: false, IsFrozen: false, Color: 2},
		{Name: "Windows", IsOn: true, IsFrozen: true, Color: 3},
	}
	testData := &data.ExtractedData{DXFVersion: "R2020", Layers: layers}
	view.Update(testData)

	// Toggle the first layer OFF
	view.ToggleLayerVisibility(0)
	if testData.Layers[0].IsOn {
		t.Errorf("Expected layer 'Walls' to be OFF after toggle, got ON")
	}
	// Toggle the second layer ON
	view.ToggleLayerVisibility(1)
	if !testData.Layers[1].IsOn {
		t.Errorf("Expected layer 'Doors' to be ON after toggle, got OFF")
	}
	// Toggle a frozen layer
	view.ToggleLayerVisibility(2)
	if testData.Layers[2].IsOn {
		t.Errorf("Expected layer 'Windows' to be OFF after toggle, got ON")
	}

	// Toggle after filtering
	view.FilterLayers("Doors")
	if view.layers.GetItemCount() != 1 {
		t.Fatalf("Expected 1 layer after filter, got %d", view.layers.GetItemCount())
	}
	// Should toggle the only visible (filtered) layer
	view.ToggleLayerVisibility(view.layers.GetCurrentItem())
	if testData.Layers[1].IsOn {
		t.Errorf("Expected filtered 'Doors' layer to be OFF after toggle, got ON")
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
