package tui

import (
	"testing"

	"github.com/gdamore/tcell/v2"
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

// TestSetupKeybindings_Comprehensive tests the setupKeybindings function comprehensively
func TestSetupKeybindings_Comprehensive(t *testing.T) {
	app := SetupTestApp(t)
	view := NewDXFView(app)

	// Setup test data
	testData := &data.ExtractedData{
		DXFVersion: "R2020",
		Layers: []data.LayerInfo{
			{Name: "Layer1", IsOn: true, IsFrozen: false, Color: 1, LineType: "CONTINUOUS"},
			{Name: "Layer2", IsOn: false, IsFrozen: false, Color: 2, LineType: "DASHED"},
		},
	}
	view.Update(testData)

	// Test that setupKeybindings doesn't panic
	view.setupKeybindings()

	// Test search input key events
	searchEvent := tcell.NewEventKey(tcell.KeyDown, 0, tcell.ModNone)
	result := view.searchInput.GetInputCapture()(searchEvent)
	assert.Nil(t, result) // Should be handled and return nil

	tabEvent := tcell.NewEventKey(tcell.KeyTab, 0, tcell.ModNone)
	result = view.searchInput.GetInputCapture()(tabEvent)
	assert.Nil(t, result) // Should be handled and return nil

	escEvent := tcell.NewEventKey(tcell.KeyEsc, 0, tcell.ModNone)
	result = view.searchInput.GetInputCapture()(escEvent)
	assert.Nil(t, result) // Should be handled and return nil

	// Test layers list key events
	enterEvent := tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone)
	result = view.layers.GetInputCapture()(enterEvent)
	assert.Nil(t, result) // Should be handled and return nil

	backspaceEvent := tcell.NewEventKey(tcell.KeyBackspace, 0, tcell.ModNone)
	result = view.layers.GetInputCapture()(backspaceEvent)
	assert.Nil(t, result) // Should be handled and return nil

	backspace2Event := tcell.NewEventKey(tcell.KeyBackspace2, 0, tcell.ModNone)
	result = view.layers.GetInputCapture()(backspace2Event)
	assert.Nil(t, result) // Should be handled and return nil

	spaceEvent := tcell.NewEventKey(tcell.KeyRune, ' ', tcell.ModNone)
	result = view.layers.GetInputCapture()(spaceEvent)
	assert.Nil(t, result) // Should be handled and return nil

	tEvent := tcell.NewEventKey(tcell.KeyRune, 't', tcell.ModNone)
	result = view.layers.GetInputCapture()(tEvent)
	assert.Nil(t, result) // Should be handled and return nil

	TEvent := tcell.NewEventKey(tcell.KeyRune, 'T', tcell.ModNone)
	result = view.layers.GetInputCapture()(TEvent)
	assert.Nil(t, result) // Should be handled and return nil

	letterEvent := tcell.NewEventKey(tcell.KeyRune, 'a', tcell.ModNone)
	result = view.layers.GetInputCapture()(letterEvent)
	assert.Nil(t, result) // Should be handled and return nil

	numberEvent := tcell.NewEventKey(tcell.KeyRune, '5', tcell.ModNone)
	result = view.layers.GetInputCapture()(numberEvent)
	assert.Nil(t, result) // Should be handled and return nil

	colonEvent := tcell.NewEventKey(tcell.KeyRune, ':', tcell.ModNone)
	result = view.layers.GetInputCapture()(colonEvent)
	assert.Nil(t, result) // Should be handled and return nil

	// Test entity list key events
	entityEscEvent := tcell.NewEventKey(tcell.KeyEsc, 0, tcell.ModNone)
	result = view.entityList.GetInputCapture()(entityEscEvent)
	assert.Nil(t, result) // Should be handled and return nil

	entityBackspaceEvent := tcell.NewEventKey(tcell.KeyBackspace, 0, tcell.ModNone)
	result = view.entityList.GetInputCapture()(entityBackspaceEvent)
	assert.Nil(t, result) // Should be handled and return nil

	entityBackspace2Event := tcell.NewEventKey(tcell.KeyBackspace2, 0, tcell.ModNone)
	result = view.entityList.GetInputCapture()(entityBackspace2Event)
	assert.Nil(t, result) // Should be handled and return nil
}
