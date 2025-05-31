package tui

import (
	"testing"

	"github.com/remym/go-dwg-extractor/pkg/data"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
)

func TestDXFView_GetLayout(t *testing.T) {
	// Explicitly use tview to ensure import is used
	var _ *tview.Application
	app := SetupTestApp(t)
	view := NewDXFView(app)

	layout := view.GetLayout()
	assert.NotNil(t, layout, "Expected GetLayout to return a non-nil layout")
}

func TestDXFView_SetLayersChangedFunc(t *testing.T) {
	app := SetupTestApp(t)
	view := NewDXFView(app)

	// Set a no-op callback
	view.SetLayersChangedFunc(func(index int, name, secondaryText string, shortcut rune) {
		// No-op for test
	})

	// Verify the layers list is initialized
	assert.NotNil(t, view.layers, "Expected layers list to be initialized")
}

func TestDXFView_SetLayersSelectedFunc(t *testing.T) {
	app := SetupTestApp(t)
	view := NewDXFView(app)

	// Set a no-op callback
	view.SetLayersSelectedFunc(func(index int, name, secondaryText string, shortcut rune) {
		// No-op for test
	})

	// Verify the layers list is initialized
	assert.NotNil(t, view.layers, "Expected layers list to be initialized")
}

func TestDXFView_UpdateLayersList(t *testing.T) {
	app := SetupTestApp(t)
	view := NewDXFView(app)

	// Create test data with layers
	testData := &data.ExtractedData{
		DXFVersion: "R2020",
		Layers: []data.LayerInfo{
			{Name: "TestLayer1", IsOn: true, IsFrozen: false, Color: 1},
			{Name: "TestLayer2", IsOn: false, IsFrozen: true, Color: 2},
		},
	}

	// Update the view with test data
	view.Update(testData)

	// Verify the layers list was updated
	assert.Equal(t, 2, view.layers.GetItemCount(), "Expected 2 layers in the list")
}

func TestDXFView_ShowLayerDetails(t *testing.T) {
	app := SetupTestApp(t)
	view := NewDXFView(app)

	// Create test data with a layer and entities
	line1 := &data.LineInfo{
		StartPoint: data.Point{X: 0, Y: 0},
		EndPoint:   data.Point{X: 10, Y: 10},
		Layer:      "TestLayer",
		Color:      1,
	}

	testData := &data.ExtractedData{
		DXFVersion: "R2020",
		Layers: []data.LayerInfo{
			{
				Name:     "TestLayer",
				IsOn:     true,
				IsFrozen: false,
				Color:    1,
				Entities: []data.Entity{line1},
			},
		},
	}

	// Update the view and show layer details
	view.Update(testData)
	view.showLayerDetails(0)

	// Verify the text view contains layer information
	text := view.textView.GetText(true)
	assert.Contains(t, text, "Layer: TestLayer", "Expected layer details to be shown")
}

func TestDXFView_ShowEntitiesView(t *testing.T) {
	app := SetupTestApp(t)
	view := NewDXFView(app)

	// Create test data with a layer and entities
	line1 := &data.LineInfo{
		StartPoint: data.Point{X: 0, Y: 0},
		EndPoint:   data.Point{X: 10, Y: 10},
		Layer:      "TestLayer",
		Color:      1,
	}

	testData := &data.ExtractedData{
		DXFVersion: "R2020",
		Layers: []data.LayerInfo{
			{
				Name:     "TestLayer",
				IsOn:     true,
				IsFrozen: false,
				Color:    1,
				Entities: []data.Entity{line1},
			},
		},
	}

	// Update the view and show entities view
	view.Update(testData)
	view.showLayerDetails(0)
	view.showEntitiesView()

	// Verify the entities view is shown
	_, page := view.pages.GetFrontPage()
	assert.NotNil(t, page, "Expected entities view to be shown")
}
