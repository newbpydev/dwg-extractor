package tui

import (
	"testing"

	"github.com/remym/go-dwg-extractor/pkg/data"
	"github.com/stretchr/testify/assert"
)

// TestSelectionState tests the selection state functionality
func TestSelectionState(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(t *testing.T, state *SelectionState)
	}{
		{
			name: "NewSelectionState creates empty state",
			testFunc: func(t *testing.T, state *SelectionState) {
				assert.Equal(t, 0, state.GetSelectedCount(), "Expected no selected items initially")
				assert.Equal(t, "layers", state.currentView, "Expected default view to be layers")
				assert.False(t, state.IsSelected("any_item"), "Expected no items to be selected initially")
			},
		},
		{
			name: "ToggleSelection works correctly",
			testFunc: func(t *testing.T, state *SelectionState) {
				itemID := "item1"

				// Select item
				state.ToggleSelection(itemID)
				assert.True(t, state.IsSelected(itemID), "Expected item to be selected")
				assert.Equal(t, 1, state.GetSelectedCount(), "Expected one selected item")

				// Deselect item
				state.ToggleSelection(itemID)
				assert.False(t, state.IsSelected(itemID), "Expected item to be deselected")
				assert.Equal(t, 0, state.GetSelectedCount(), "Expected no selected items")
			},
		},
		{
			name: "SelectAll works correctly",
			testFunc: func(t *testing.T, state *SelectionState) {
				items := []string{"item1", "item2", "item3"}

				state.SelectAll(items)
				assert.Equal(t, 3, state.GetSelectedCount(), "Expected all items to be selected")

				for _, item := range items {
					assert.True(t, state.IsSelected(item), "Expected item %s to be selected", item)
				}

				selectedIDs := state.GetSelectedItemIDs()
				assert.Len(t, selectedIDs, 3, "Expected 3 selected item IDs")
			},
		},
		{
			name: "SelectNone works correctly",
			testFunc: func(t *testing.T, state *SelectionState) {
				// First select some items
				items := []string{"item1", "item2", "item3"}
				state.SelectAll(items)
				assert.Equal(t, 3, state.GetSelectedCount(), "Expected items to be selected")

				// Then clear selection
				state.SelectNone()
				assert.Equal(t, 0, state.GetSelectedCount(), "Expected no selected items after SelectNone")

				for _, item := range items {
					assert.False(t, state.IsSelected(item), "Expected item %s to be deselected", item)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			state := NewSelectionState()
			tt.testFunc(t, state)
		})
	}
}

// TestEnhancedCategorySelector tests the enhanced category selector
func TestEnhancedCategorySelector(t *testing.T) {
	tests := []struct {
		name          string
		categoryType  string
		index         int
		expectError   bool
		expectedItems int
		testDataFunc  func() *data.ExtractedData
	}{
		{
			name:          "SelectCategory with layer type",
			categoryType:  "layer",
			index:         0,
			expectError:   false,
			expectedItems: 3, // entities in layer
			testDataFunc:  createTestDataWithMultipleItems,
		},
		{
			name:          "SelectCategory with block type",
			categoryType:  "block",
			index:         0,
			expectError:   false,
			expectedItems: 2, // blocks in test data
			testDataFunc:  createTestDataWithCategories,
		},
		{
			name:          "SelectCategory with text type",
			categoryType:  "text",
			index:         0,
			expectError:   false,
			expectedItems: 1, // text entities
			testDataFunc:  createTestDataWithCategories,
		},
		{
			name:          "SelectCategory with invalid type",
			categoryType:  "invalid",
			index:         0,
			expectError:   true,
			expectedItems: 0,
			testDataFunc:  createTestDataWithCategories,
		},
		{
			name:          "SelectCategory with invalid layer index",
			categoryType:  "layer",
			index:         999,
			expectError:   true,
			expectedItems: 0,
			testDataFunc:  createTestDataWithCategories,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := SetupTestApp(t)
			view := NewDXFView(app)

			testData := tt.testDataFunc()
			view.Update(testData)

			selector := NewEnhancedCategorySelector(view)

			err := selector.SelectCategory(tt.categoryType, tt.index)

			if tt.expectError {
				assert.Error(t, err, "Expected error for invalid category selection")
			} else {
				assert.NoError(t, err, "Expected no error for valid category selection")

				// Check that the list was updated with correct number of items
				if tt.expectedItems > 0 {
					assert.Equal(t, tt.expectedItems, view.layers.GetItemCount(),
						"Expected correct number of items in list")
				}
			}
		})
	}
}

// TestEnhancedCategorySelector_NilData tests error handling with nil data
func TestEnhancedCategorySelector_NilData(t *testing.T) {
	app := SetupTestApp(t)
	view := NewDXFView(app)
	// Don't set any data (view.data will be nil)

	selector := NewEnhancedCategorySelector(view)

	err := selector.SelectCategory("layer", 0)
	assert.Error(t, err, "Expected error when data is nil")
	assert.Contains(t, err.Error(), "no data available", "Expected specific error message")
}

// TestEnhancedItemSelector tests the enhanced item selector
func TestEnhancedItemSelector(t *testing.T) {
	tests := []struct {
		name        string
		itemType    string
		index       int
		expectError bool
	}{
		{
			name:        "SelectItem with line type",
			itemType:    "line",
			index:       0,
			expectError: false,
		},
		{
			name:        "SelectItem with circle type",
			itemType:    "circle",
			index:       0,
			expectError: false,
		},
		{
			name:        "SelectItem with text type",
			itemType:    "text",
			index:       0,
			expectError: false,
		},
		{
			name:        "SelectItem with block type",
			itemType:    "block",
			index:       0,
			expectError: false,
		},
		{
			name:        "SelectItem with invalid index",
			itemType:    "line",
			index:       999,
			expectError: true,
		},
		{
			name:        "SelectItem with invalid type",
			itemType:    "invalid",
			index:       0,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := SetupTestApp(t)
			view := NewDXFView(app)

			testData := createTestDataWithDetailedEntities()
			view.Update(testData)

			selector := NewEnhancedItemSelector(view)

			err := selector.SelectItem(tt.itemType, tt.index)

			if tt.expectError {
				assert.Error(t, err, "Expected error for invalid item selection")
			} else {
				assert.NoError(t, err, "Expected no error for valid item selection")

				// Check that the details view was updated
				detailText := view.textView.GetText(true)
				assert.NotEmpty(t, detailText, "Expected details view to be updated")

				// Check that the text contains entity-specific information
				switch tt.itemType {
				case "line":
					assert.Contains(t, detailText, "Line Entity", "Expected line entity details")
					assert.Contains(t, detailText, "Start Point", "Expected start point field")
					assert.Contains(t, detailText, "End Point", "Expected end point field")
				case "circle":
					assert.Contains(t, detailText, "Circle Entity", "Expected circle entity details")
					assert.Contains(t, detailText, "Center", "Expected center field")
					assert.Contains(t, detailText, "Radius", "Expected radius field")
				case "text":
					assert.Contains(t, detailText, "Text Entity", "Expected text entity details")
					assert.Contains(t, detailText, "Value", "Expected value field")
					assert.Contains(t, detailText, "Insertion Point", "Expected insertion point field")
				case "block":
					assert.Contains(t, detailText, "Block Entity", "Expected block entity details")
					assert.Contains(t, detailText, "Name", "Expected name field")
					assert.Contains(t, detailText, "Attributes", "Expected attributes field")
				}
			}
		})
	}
}

// TestEnhancedItemSelector_NilData tests error handling with nil data
func TestEnhancedItemSelector_NilData(t *testing.T) {
	app := SetupTestApp(t)
	view := NewDXFView(app)
	// Don't set any data (view.data will be nil)

	selector := NewEnhancedItemSelector(view)

	err := selector.SelectItem("line", 0)
	assert.Error(t, err, "Expected error when data is nil")
	assert.Contains(t, err.Error(), "no data available", "Expected specific error message")
}

// TestEnhancedItemSelector_GetSelectionState tests the selection state getter
func TestEnhancedItemSelector_GetSelectionState(t *testing.T) {
	app := SetupTestApp(t)
	view := NewDXFView(app)

	selector := NewEnhancedItemSelector(view)
	state := selector.GetSelectionState()

	assert.NotNil(t, state, "Expected selection state to be available")
	assert.Equal(t, 0, state.GetSelectedCount(), "Expected empty selection state")
}

// TestEntityDetailsFormatting tests the formatting of entity details
func TestEntityDetailsFormatting(t *testing.T) {
	tests := []struct {
		name           string
		entity         data.Entity
		expectedFields []string
	}{
		{
			name: "LineInfo formatting",
			entity: &data.LineInfo{
				StartPoint: data.Point{X: 1.5, Y: 2.5},
				EndPoint:   data.Point{X: 3.5, Y: 4.5},
				Layer:      "TestLayer",
				Color:      7,
			},
			expectedFields: []string{"Line Entity", "Start Point", "End Point", "Layer", "Color"},
		},
		{
			name: "CircleInfo formatting",
			entity: &data.CircleInfo{
				Center: data.Point{X: 10.0, Y: 15.0},
				Radius: 5.0,
				Layer:  "CircleLayer",
				Color:  3,
			},
			expectedFields: []string{"Circle Entity", "Center", "Radius", "Layer", "Color"},
		},
		{
			name: "TextInfo formatting",
			entity: &data.TextInfo{
				Value:          "Test Text Content",
				InsertionPoint: data.Point{X: 0.0, Y: 0.0},
				Height:         12.0,
				Layer:          "TextLayer",
			},
			expectedFields: []string{"Text Entity", "Value", "Insertion Point", "Height", "Layer"},
		},
		{
			name: "BlockInfo formatting with attributes",
			entity: &data.BlockInfo{
				Name:           "TestBlock",
				InsertionPoint: data.Point{X: 5.0, Y: 10.0},
				Rotation:       45.0,
				Scale:          data.Point{X: 2.0, Y: 1.5},
				Layer:          "BlockLayer",
				Attributes: []data.AttributeInfo{
					{Tag: "TAG1", Value: "Value1"},
					{Tag: "TAG2", Value: "Value2"},
				},
			},
			expectedFields: []string{"Block Entity", "Name", "Insertion Point", "Rotation", "Scale", "Attributes", "TAG1: Value1", "TAG2: Value2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := SetupTestApp(t)
			view := NewDXFView(app)

			// Create minimal test data
			testData := &data.ExtractedData{
				DXFVersion: "R2020",
				Layers: []data.LayerInfo{
					{
						Name:     "Layer1",
						IsOn:     true,
						IsFrozen: false,
						Color:    1,
						Entities: []data.Entity{tt.entity},
					},
				},
			}
			view.Update(testData)

			selector := NewEnhancedCategorySelector(view)

			// Use the updateEntityDetails method to format the entity
			selector.updateEntityDetails(tt.entity)

			// Check that all expected fields are present
			detailText := view.textView.GetText(true)
			for _, field := range tt.expectedFields {
				assert.Contains(t, detailText, field, "Expected field '%s' to be present in details", field)
			}
		})
	}
}
