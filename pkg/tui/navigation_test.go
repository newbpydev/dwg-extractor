package tui

import (
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/remym/go-dwg-extractor/pkg/data"
	"github.com/stretchr/testify/assert"
)

// TestAdvancedNavigation_BetweenPanes tests navigation between different panes
func TestAdvancedNavigation_BetweenPanes(t *testing.T) {
	tests := []struct {
		name              string
		initialFocus      string
		key               tcell.Key
		expectedFocus     string
		expectedDirection string
	}{
		{
			name:              "Tab from search to layers",
			initialFocus:      "search",
			key:               tcell.KeyTab,
			expectedFocus:     "layers",
			expectedDirection: "forward",
		},
		{
			name:              "Shift+Tab from layers to search",
			initialFocus:      "layers",
			key:               tcell.KeyBacktab,
			expectedFocus:     "search",
			expectedDirection: "backward",
		},
		{
			name:              "Tab from layers to entities when available",
			initialFocus:      "layers",
			key:               tcell.KeyTab,
			expectedFocus:     "entities",
			expectedDirection: "forward",
		},
		{
			name:              "F1 shows help view",
			initialFocus:      "layers",
			key:               tcell.KeyF1,
			expectedFocus:     "help",
			expectedDirection: "help",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := SetupTestApp(t)
			view := NewDXFView(app)

			// Create test data - use data with entities for entity navigation tests
			var testData *data.ExtractedData
			if tt.expectedFocus == "entities" {
				testData = createTestDataWithMultipleItems() // Has entities
				view.Update(testData)
				// Simulate selecting a layer to populate the entity list
				view.showLayerDetails(0) // Select the first layer which has entities
			} else {
				testData = createTestData() // Basic data
				view.Update(testData)
			}

			// This should fail initially - we need to implement advanced navigation
			navigator := view.GetNavigator()
			assert.NotNil(t, navigator, "Expected navigator to be available")

			// Set initial focus
			navigator.SetFocus(tt.initialFocus)

			// Simulate key press
			handled := navigator.HandleNavigation(tt.key, tcell.ModNone)
			assert.True(t, handled, "Expected navigation key to be handled")

			// Check expected focus
			currentFocus := navigator.GetCurrentFocus()
			assert.Equal(t, tt.expectedFocus, currentFocus, "Expected focus to change correctly")
		})
	}
}

// TestKeyboardNavigation_WithinLists tests keyboard navigation within lists
func TestKeyboardNavigation_WithinLists(t *testing.T) {
	tests := []struct {
		name          string
		listType      string
		key           tcell.Key
		initialIndex  int
		expectedIndex int
		shouldWrap    bool
	}{
		{
			name:          "Down arrow in layers list",
			listType:      "layers",
			key:           tcell.KeyDown,
			initialIndex:  0,
			expectedIndex: 1,
			shouldWrap:    false,
		},
		{
			name:          "Up arrow in layers list",
			listType:      "layers",
			key:           tcell.KeyUp,
			initialIndex:  1,
			expectedIndex: 0,
			shouldWrap:    false,
		},
		{
			name:          "Down arrow wraps at bottom",
			listType:      "layers",
			key:           tcell.KeyDown,
			initialIndex:  2, // last item
			expectedIndex: 0, // wraps to first
			shouldWrap:    true,
		},
		{
			name:          "Up arrow wraps at top",
			listType:      "layers",
			key:           tcell.KeyUp,
			initialIndex:  0, // first item
			expectedIndex: 2, // wraps to last
			shouldWrap:    true,
		},
		{
			name:          "Page down navigation",
			listType:      "entities",
			key:           tcell.KeyPgDn,
			initialIndex:  0,
			expectedIndex: 3, // should jump by page size, but limited to last item (3 entities + 1 back item = 4 total, so index 3)
			shouldWrap:    false,
		},
		{
			name:          "Home key goes to first item",
			listType:      "layers",
			key:           tcell.KeyHome,
			initialIndex:  2,
			expectedIndex: 0,
			shouldWrap:    false,
		},
		{
			name:          "End key goes to last item",
			listType:      "layers",
			key:           tcell.KeyEnd,
			initialIndex:  0,
			expectedIndex: 2, // last item
			shouldWrap:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := SetupTestApp(t)
			view := NewDXFView(app)

			// Create test data with multiple items
			testData := createTestDataWithMultipleItems()
			view.Update(testData)

			// If testing entity navigation, populate the entity list
			if tt.listType == "entities" {
				view.showLayerDetails(0) // Select the first layer to populate entity list
			}

			// This should fail initially - we need to implement enhanced keyboard navigation
			listNavigator := view.GetListNavigator(tt.listType)
			assert.NotNil(t, listNavigator, "Expected list navigator to be available")

			// Set initial position
			listNavigator.SetCurrentIndex(tt.initialIndex)
			listNavigator.SetWrapNavigation(tt.shouldWrap)

			// Simulate key press
			handled := listNavigator.HandleKeyPress(tt.key, tcell.ModNone)
			assert.True(t, handled, "Expected navigation key to be handled")

			// Check expected position
			currentIndex := listNavigator.GetCurrentIndex()
			assert.Equal(t, tt.expectedIndex, currentIndex, "Expected index to change correctly")
		})
	}
}

// TestCategorySelection_UpdatesViews tests that selecting categories updates corresponding views
func TestCategorySelection_UpdatesViews(t *testing.T) {
	tests := []struct {
		name               string
		categoryType       string
		categoryIndex      int
		expectedItemCount  int
		expectedDetailText string
	}{
		{
			name:               "Select layer shows entities",
			categoryType:       "layer",
			categoryIndex:      0,
			expectedItemCount:  5, // lines + circles + texts + blocks (line1, circle1, text1, block1, block2)
			expectedDetailText: "Layer: Layer1",
		},
		{
			name:               "Select block category shows blocks",
			categoryType:       "block",
			categoryIndex:      0,
			expectedItemCount:  2, // test blocks
			expectedDetailText: "Block: TestBlock1",
		},
		{
			name:               "Select text category shows texts",
			categoryType:       "text",
			categoryIndex:      0,
			expectedItemCount:  1, // test texts
			expectedDetailText: "Text: Sample Text",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := SetupTestApp(t)
			view := NewDXFView(app)

			// Create test data
			testData := createTestDataWithCategories()
			view.Update(testData)

			// This should fail initially - we need to implement category selection
			categorySelector := view.GetCategorySelector()
			assert.NotNil(t, categorySelector, "Expected category selector to be available")

			// Select category
			err := categorySelector.SelectCategory(tt.categoryType, tt.categoryIndex)
			assert.NoError(t, err, "Expected category selection to succeed")

			// Check item list update
			itemList := view.GetCurrentItemList()
			assert.NotNil(t, itemList, "Expected item list to be available")
			assert.Equal(t, tt.expectedItemCount, itemList.GetItemCount(), "Expected correct item count")

			// Check details view update
			detailsView := view.GetDetailsView()
			assert.NotNil(t, detailsView, "Expected details view to be available")
			detailText := detailsView.GetText(true)
			assert.Contains(t, detailText, tt.expectedDetailText, "Expected correct detail text")
		})
	}
}

// TestItemSelection_UpdatesDetailsPane tests that selecting items updates the details pane
func TestItemSelection_UpdatesDetailsPane(t *testing.T) {
	tests := []struct {
		name           string
		itemType       string
		itemIndex      int
		expectedDetail string
		expectedFields []string
	}{
		{
			name:           "Select line shows line details",
			itemType:       "line",
			itemIndex:      0,
			expectedDetail: "Line Entity",
			expectedFields: []string{"Start Point", "End Point", "Layer", "Color"},
		},
		{
			name:           "Select circle shows circle details",
			itemType:       "circle",
			itemIndex:      0,
			expectedDetail: "Circle Entity",
			expectedFields: []string{"Center", "Radius", "Layer", "Color"},
		},
		{
			name:           "Select text shows text details",
			itemType:       "text",
			itemIndex:      0,
			expectedDetail: "Text Entity",
			expectedFields: []string{"Value", "Insertion Point", "Height", "Layer"},
		},
		{
			name:           "Select block shows block details",
			itemType:       "block",
			itemIndex:      0,
			expectedDetail: "Block Entity",
			expectedFields: []string{"Name", "Insertion Point", "Rotation", "Scale", "Attributes"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := SetupTestApp(t)
			view := NewDXFView(app)

			// Create test data
			testData := createTestDataWithDetailedEntities()
			view.Update(testData)

			// This should fail initially - we need to implement item selection
			itemSelector := view.GetItemSelector()
			assert.NotNil(t, itemSelector, "Expected item selector to be available")

			// Select item
			err := itemSelector.SelectItem(tt.itemType, tt.itemIndex)
			assert.NoError(t, err, "Expected item selection to succeed")

			// Check details pane update
			detailsPane := view.GetDetailsPane()
			assert.NotNil(t, detailsPane, "Expected details pane to be available")

			detailText := detailsPane.GetText(true)
			assert.Contains(t, detailText, tt.expectedDetail, "Expected correct detail header")

			// Check that all expected fields are present
			for _, field := range tt.expectedFields {
				assert.Contains(t, detailText, field, "Expected field %s to be present", field)
			}
		})
	}
}

// TestBreadcrumbNavigation tests breadcrumb navigation functionality
func TestBreadcrumbNavigation(t *testing.T) {
	tests := []struct {
		name               string
		navigationPath     []string
		expectedBreadcrumb string
		canGoBack          bool
	}{
		{
			name:               "Root view shows no breadcrumb",
			navigationPath:     []string{},
			expectedBreadcrumb: "",
			canGoBack:          false,
		},
		{
			name:               "Layer view shows layer breadcrumb",
			navigationPath:     []string{"layers"},
			expectedBreadcrumb: "Layers",
			canGoBack:          true,
		},
		{
			name:               "Entity view shows layer > entities breadcrumb",
			navigationPath:     []string{"layers", "Layer1", "entities"},
			expectedBreadcrumb: "Layers > Layer1 > Entities",
			canGoBack:          true,
		},
		{
			name:               "Detail view shows full breadcrumb",
			navigationPath:     []string{"layers", "Layer1", "entities", "Line_001"},
			expectedBreadcrumb: "Layers > Layer1 > Entities > Line_001",
			canGoBack:          true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := SetupTestApp(t)
			view := NewDXFView(app)

			// Create test data
			testData := createTestData()
			view.Update(testData)

			// This should fail initially - we need to implement breadcrumb navigation
			navigator := view.GetBreadcrumbNavigator()
			assert.NotNil(t, navigator, "Expected breadcrumb navigator to be available")

			// Navigate through path
			for _, step := range tt.navigationPath {
				err := navigator.NavigateTo(step)
				assert.NoError(t, err, "Expected navigation step to succeed")
			}

			// Check breadcrumb
			breadcrumb := navigator.GetBreadcrumb()
			assert.Equal(t, tt.expectedBreadcrumb, breadcrumb, "Expected correct breadcrumb")

			// Check back navigation availability
			canGoBack := navigator.CanGoBack()
			assert.Equal(t, tt.canGoBack, canGoBack, "Expected correct back navigation availability")
		})
	}
}

// Helper functions to create test data

func createTestData() *data.ExtractedData {
	return &data.ExtractedData{
		DXFVersion: "R2020",
		Layers: []data.LayerInfo{
			{Name: "Layer1", IsOn: true, IsFrozen: false, Color: 1},
			{Name: "Layer2", IsOn: true, IsFrozen: false, Color: 2},
			{Name: "Layer3", IsOn: false, IsFrozen: false, Color: 3},
		},
	}
}

func createTestDataWithMultipleItems() *data.ExtractedData {
	line1 := &data.LineInfo{StartPoint: data.Point{X: 0, Y: 0}, EndPoint: data.Point{X: 10, Y: 10}, Layer: "Layer1", Color: 1}
	line2 := &data.LineInfo{StartPoint: data.Point{X: 10, Y: 10}, EndPoint: data.Point{X: 20, Y: 20}, Layer: "Layer1", Color: 1}
	circle1 := &data.CircleInfo{Center: data.Point{X: 5, Y: 5}, Radius: 2.5, Layer: "Layer1", Color: 1}

	return &data.ExtractedData{
		DXFVersion: "R2020",
		Layers: []data.LayerInfo{
			{
				Name: "Layer1", IsOn: true, IsFrozen: false, Color: 1,
				Entities: []data.Entity{line1, line2, circle1},
			},
			{Name: "Layer2", IsOn: true, IsFrozen: false, Color: 2},
			{Name: "Layer3", IsOn: false, IsFrozen: false, Color: 3},
		},
	}
}

func createTestDataWithCategories() *data.ExtractedData {
	line1 := &data.LineInfo{StartPoint: data.Point{X: 0, Y: 0}, EndPoint: data.Point{X: 10, Y: 10}, Layer: "Layer1", Color: 1}
	circle1 := &data.CircleInfo{Center: data.Point{X: 5, Y: 5}, Radius: 2.5, Layer: "Layer1", Color: 1}
	text1 := &data.TextInfo{Value: "Sample Text", InsertionPoint: data.Point{X: 1, Y: 1}, Height: 12.0, Layer: "Layer1"}

	// Add blocks to match the test expectations (expected 2 blocks)
	block1 := &data.BlockInfo{
		Name:           "TestBlock1",
		InsertionPoint: data.Point{X: 0, Y: 0},
		Rotation:       0.0,
		Scale:          data.Point{X: 1, Y: 1},
		Layer:          "Layer1",
		Attributes: []data.AttributeInfo{
			{Tag: "ATTR1", Value: "Value1"},
		},
	}
	block2 := &data.BlockInfo{
		Name:           "TestBlock2",
		InsertionPoint: data.Point{X: 5, Y: 5},
		Rotation:       45.0,
		Scale:          data.Point{X: 2, Y: 2},
		Layer:          "Layer1",
		Attributes: []data.AttributeInfo{
			{Tag: "ATTR2", Value: "Value2"},
		},
	}

	return &data.ExtractedData{
		DXFVersion: "R2020",
		Layers: []data.LayerInfo{
			{
				Name: "Layer1", IsOn: true, IsFrozen: false, Color: 1,
				Entities: []data.Entity{line1, circle1, text1, block1, block2},
			},
		},
	}
}

func createTestDataWithDetailedEntities() *data.ExtractedData {
	line1 := &data.LineInfo{
		StartPoint: data.Point{X: 0, Y: 0},
		EndPoint:   data.Point{X: 10, Y: 10},
		Layer:      "Layer1",
		Color:      1,
	}
	circle1 := &data.CircleInfo{
		Center: data.Point{X: 5, Y: 5},
		Radius: 2.5,
		Layer:  "Layer1",
		Color:  1,
	}
	text1 := &data.TextInfo{
		Value:          "Test Text",
		InsertionPoint: data.Point{X: 1, Y: 1},
		Height:         12.0,
		Layer:          "Layer1",
	}
	block1 := &data.BlockInfo{
		Name:           "TestBlock",
		InsertionPoint: data.Point{X: 0, Y: 0},
		Rotation:       0.0,
		Scale:          data.Point{X: 1, Y: 1},
		Layer:          "Layer1",
		Attributes: []data.AttributeInfo{
			{Tag: "ATTR1", Value: "Value1"},
		},
	}

	return &data.ExtractedData{
		DXFVersion: "R2020",
		Layers: []data.LayerInfo{
			{
				Name: "Layer1", IsOn: true, IsFrozen: false, Color: 1,
				Entities: []data.Entity{line1, circle1, text1, block1},
			},
		},
	}
}
