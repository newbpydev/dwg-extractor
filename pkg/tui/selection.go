package tui

import (
	"fmt"

	"github.com/remym/go-dwg-extractor/pkg/data"
)

// SelectionState tracks the current selection state
type SelectionState struct {
	selectedItems   map[string]bool // key: itemID, value: selected
	currentCategory string          // current category filter
	currentView     string          // current view mode
}

// NewSelectionState creates a new selection state
func NewSelectionState() *SelectionState {
	return &SelectionState{
		selectedItems: make(map[string]bool),
		currentView:   "layers",
	}
}

// IsSelected checks if an item is selected
func (s *SelectionState) IsSelected(itemID string) bool {
	return s.selectedItems[itemID]
}

// ToggleSelection toggles the selection of an item
func (s *SelectionState) ToggleSelection(itemID string) {
	s.selectedItems[itemID] = !s.selectedItems[itemID]
}

// SelectAll selects all items in the given list
func (s *SelectionState) SelectAll(itemIDs []string) {
	for _, id := range itemIDs {
		s.selectedItems[id] = true
	}
}

// SelectNone deselects all items
func (s *SelectionState) SelectNone() {
	s.selectedItems = make(map[string]bool)
}

// GetSelectedCount returns the number of selected items
func (s *SelectionState) GetSelectedCount() int {
	count := 0
	for _, selected := range s.selectedItems {
		if selected {
			count++
		}
	}
	return count
}

// GetSelectedItemIDs returns a slice of selected item IDs
func (s *SelectionState) GetSelectedItemIDs() []string {
	var selected []string
	for id, isSelected := range s.selectedItems {
		if isSelected {
			selected = append(selected, id)
		}
	}
	return selected
}

// EnhancedCategorySelector implements CategorySelector with actual functionality
type EnhancedCategorySelector struct {
	view      *DXFView
	selection *SelectionState
}

// NewEnhancedCategorySelector creates a new enhanced category selector
func NewEnhancedCategorySelector(view *DXFView) *EnhancedCategorySelector {
	return &EnhancedCategorySelector{
		view:      view,
		selection: NewSelectionState(),
	}
}

// SelectCategory selects a category and updates the views accordingly
func (cs *EnhancedCategorySelector) SelectCategory(categoryType string, index int) error {
	if cs.view.data == nil {
		return fmt.Errorf("no data available")
	}

	cs.selection.currentCategory = categoryType

	switch categoryType {
	case "layer":
		return cs.selectLayer(index)
	case "block":
		return cs.selectBlockCategory(index)
	case "text":
		return cs.selectTextCategory(index)
	default:
		return fmt.Errorf("unknown category type: %s", categoryType)
	}
}

// selectLayer selects a layer and shows its entities
func (cs *EnhancedCategorySelector) selectLayer(index int) error {
	if index < 0 || index >= len(cs.view.data.Layers) {
		return fmt.Errorf("layer index %d out of range", index)
	}

	layer := cs.view.data.Layers[index]

	// Clear the current item list
	cs.view.layers.Clear()

	// Add entities from this layer to the list
	for i, entity := range layer.Entities {
		var itemText string
		switch e := entity.(type) {
		case *data.LineInfo:
			itemText = fmt.Sprintf("Line (%.1f,%.1f) to (%.1f,%.1f)",
				e.StartPoint.X, e.StartPoint.Y, e.EndPoint.X, e.EndPoint.Y)
		case *data.CircleInfo:
			itemText = fmt.Sprintf("Circle center:(%.1f,%.1f) radius:%.1f",
				e.Center.X, e.Center.Y, e.Radius)
		case *data.TextInfo:
			itemText = fmt.Sprintf("Text: %s at (%.1f,%.1f)",
				e.Value, e.InsertionPoint.X, e.InsertionPoint.Y)
		default:
			itemText = fmt.Sprintf("Entity: %T", entity)
		}

		// Store the entity index for later use
		entityIndex := i
		cs.view.layers.AddItem(itemText, "", 0, func() {
			// Update details when entity is selected
			cs.updateEntityDetails(layer.Entities[entityIndex])
		})
	}

	// Update the details view with layer information
	cs.view.textView.Clear()
	fmt.Fprintf(cs.view.textView, "[green]Layer:[-] %s\n", layer.Name)
	fmt.Fprintf(cs.view.textView, "[green]Color:[-] %d\n", layer.Color)
	fmt.Fprintf(cs.view.textView, "[green]Status:[-] %s\n", map[bool]string{true: "ON", false: "OFF"}[layer.IsOn])
	fmt.Fprintf(cs.view.textView, "[green]Frozen:[-] %v\n", layer.IsFrozen)
	fmt.Fprintf(cs.view.textView, "[green]Line Type:[-] %s\n", layer.LineType)
	fmt.Fprintf(cs.view.textView, "[green]Entities:[-] %d\n\n", len(layer.Entities))

	return nil
}

// selectBlockCategory filters and shows only block entities
func (cs *EnhancedCategorySelector) selectBlockCategory(index int) error {
	// Collect all block entities from all layers
	var blocks []*data.BlockInfo
	for _, layer := range cs.view.data.Layers {
		for _, entity := range layer.Entities {
			if block, ok := entity.(*data.BlockInfo); ok {
				blocks = append(blocks, block)
			}
		}
	}

	// Clear the current list and add blocks
	cs.view.layers.Clear()
	for i, block := range blocks {
		itemText := fmt.Sprintf("Block: %s at (%.1f,%.1f)",
			block.Name, block.InsertionPoint.X, block.InsertionPoint.Y)

		blockIndex := i
		cs.view.layers.AddItem(itemText, "", 0, func() {
			cs.updateEntityDetails(blocks[blockIndex])
		})
	}

	// Update details view with the specific block if index is valid
	cs.view.textView.Clear()
	if len(blocks) > 0 && index >= 0 && index < len(blocks) {
		block := blocks[index]
		fmt.Fprintf(cs.view.textView, "[green]Block:[-] %s\n", block.Name)
		fmt.Fprintf(cs.view.textView, "[green]Layer:[-] %s\n", block.Layer)
		fmt.Fprintf(cs.view.textView, "[green]Insertion Point:[-] (%.1f, %.1f)\n",
			block.InsertionPoint.X, block.InsertionPoint.Y)
		fmt.Fprintf(cs.view.textView, "[green]Rotation:[-] %.1f\n", block.Rotation)
		fmt.Fprintf(cs.view.textView, "[green]Scale:[-] (%.1f, %.1f)\n", block.Scale.X, block.Scale.Y)

		if len(block.Attributes) > 0 {
			fmt.Fprintf(cs.view.textView, "[green]Attributes:[-]\n")
			for _, attr := range block.Attributes {
				fmt.Fprintf(cs.view.textView, "  %s: %s\n", attr.Tag, attr.Value)
			}
		}
	} else {
		fmt.Fprintf(cs.view.textView, "[green]Blocks Found:[-] %d\n", len(blocks))
	}

	return nil
}

// selectTextCategory filters and shows only text entities
func (cs *EnhancedCategorySelector) selectTextCategory(index int) error {
	// Collect all text entities from all layers
	var texts []*data.TextInfo
	for _, layer := range cs.view.data.Layers {
		for _, entity := range layer.Entities {
			if text, ok := entity.(*data.TextInfo); ok {
				texts = append(texts, text)
			}
		}
	}

	// Clear the current list and add texts
	cs.view.layers.Clear()
	for i, text := range texts {
		itemText := fmt.Sprintf("Text: %s at (%.1f,%.1f)",
			text.Value, text.InsertionPoint.X, text.InsertionPoint.Y)

		textIndex := i
		cs.view.layers.AddItem(itemText, "", 0, func() {
			cs.updateEntityDetails(texts[textIndex])
		})
	}

	// Update details view
	cs.view.textView.Clear()
	if len(texts) > 0 && index < len(texts) {
		text := texts[index]
		fmt.Fprintf(cs.view.textView, "[green]Text:[-] %s\n", text.Value)
		fmt.Fprintf(cs.view.textView, "[green]Layer:[-] %s\n", text.Layer)
		fmt.Fprintf(cs.view.textView, "[green]Insertion Point:[-] (%.1f, %.1f)\n",
			text.InsertionPoint.X, text.InsertionPoint.Y)
		fmt.Fprintf(cs.view.textView, "[green]Height:[-] %.1f\n", text.Height)
	} else {
		fmt.Fprintf(cs.view.textView, "[green]Texts Found:[-] %d\n", len(texts))
	}

	return nil
}

// updateEntityDetails updates the details view with entity-specific information
func (cs *EnhancedCategorySelector) updateEntityDetails(entity data.Entity) {
	cs.view.textView.Clear()

	switch e := entity.(type) {
	case *data.LineInfo:
		fmt.Fprintf(cs.view.textView, "[green]Line Entity[-]\n\n")
		fmt.Fprintf(cs.view.textView, "[green]Start Point:[-] (%.1f, %.1f)\n", e.StartPoint.X, e.StartPoint.Y)
		fmt.Fprintf(cs.view.textView, "[green]End Point:[-] (%.1f, %.1f)\n", e.EndPoint.X, e.EndPoint.Y)
		fmt.Fprintf(cs.view.textView, "[green]Layer:[-] %s\n", e.Layer)
		fmt.Fprintf(cs.view.textView, "[green]Color:[-] %d\n", e.Color)

	case *data.CircleInfo:
		fmt.Fprintf(cs.view.textView, "[green]Circle Entity[-]\n\n")
		fmt.Fprintf(cs.view.textView, "[green]Center:[-] (%.1f, %.1f)\n", e.Center.X, e.Center.Y)
		fmt.Fprintf(cs.view.textView, "[green]Radius:[-] %.1f\n", e.Radius)
		fmt.Fprintf(cs.view.textView, "[green]Layer:[-] %s\n", e.Layer)
		fmt.Fprintf(cs.view.textView, "[green]Color:[-] %d\n", e.Color)

	case *data.TextInfo:
		fmt.Fprintf(cs.view.textView, "[green]Text Entity[-]\n\n")
		fmt.Fprintf(cs.view.textView, "[green]Value:[-] %s\n", e.Value)
		fmt.Fprintf(cs.view.textView, "[green]Insertion Point:[-] (%.1f, %.1f)\n", e.InsertionPoint.X, e.InsertionPoint.Y)
		fmt.Fprintf(cs.view.textView, "[green]Height:[-] %.1f\n", e.Height)
		fmt.Fprintf(cs.view.textView, "[green]Layer:[-] %s\n", e.Layer)

	case *data.BlockInfo:
		fmt.Fprintf(cs.view.textView, "[green]Block Entity[-]\n\n")
		fmt.Fprintf(cs.view.textView, "[green]Name:[-] %s\n", e.Name)
		fmt.Fprintf(cs.view.textView, "[green]Insertion Point:[-] (%.1f, %.1f)\n", e.InsertionPoint.X, e.InsertionPoint.Y)
		fmt.Fprintf(cs.view.textView, "[green]Rotation:[-] %.1f\n", e.Rotation)
		fmt.Fprintf(cs.view.textView, "[green]Scale:[-] (%.1f, %.1f)\n", e.Scale.X, e.Scale.Y)
		fmt.Fprintf(cs.view.textView, "[green]Layer:[-] %s\n", e.Layer)

		if len(e.Attributes) > 0 {
			fmt.Fprintf(cs.view.textView, "[green]Attributes:[-]\n")
			for _, attr := range e.Attributes {
				fmt.Fprintf(cs.view.textView, "  %s: %s\n", attr.Tag, attr.Value)
			}
		}

	default:
		fmt.Fprintf(cs.view.textView, "[green]Entity:[-] %T\n", entity)
		fmt.Fprintf(cs.view.textView, "[green]Layer:[-] %s\n", entity.GetLayer())
	}
}

// EnhancedItemSelector implements ItemSelector with actual functionality
type EnhancedItemSelector struct {
	view      *DXFView
	selection *SelectionState
}

// NewEnhancedItemSelector creates a new enhanced item selector
func NewEnhancedItemSelector(view *DXFView) *EnhancedItemSelector {
	return &EnhancedItemSelector{
		view:      view,
		selection: NewSelectionState(),
	}
}

// SelectItem selects an item and updates the details pane
func (is *EnhancedItemSelector) SelectItem(itemType string, index int) error {
	if is.view.data == nil {
		return fmt.Errorf("no data available")
	}

	// Find entities of the specified type
	var entities []data.Entity
	for _, layer := range is.view.data.Layers {
		for _, entity := range layer.Entities {
			switch itemType {
			case "line":
				if _, ok := entity.(*data.LineInfo); ok {
					entities = append(entities, entity)
				}
			case "circle":
				if _, ok := entity.(*data.CircleInfo); ok {
					entities = append(entities, entity)
				}
			case "text":
				if _, ok := entity.(*data.TextInfo); ok {
					entities = append(entities, entity)
				}
			case "block":
				if _, ok := entity.(*data.BlockInfo); ok {
					entities = append(entities, entity)
				}
			}
		}
	}

	if index < 0 || index >= len(entities) {
		return fmt.Errorf("entity index %d out of range for type %s", index, itemType)
	}

	// Update the details pane with the selected entity
	entity := entities[index]
	is.updateDetailsPane(entity)

	return nil
}

// updateDetailsPane updates the details pane with entity information
func (is *EnhancedItemSelector) updateDetailsPane(entity data.Entity) {
	is.view.textView.Clear()

	switch e := entity.(type) {
	case *data.LineInfo:
		fmt.Fprintf(is.view.textView, "[green]Line Entity[-]\n\n")
		fmt.Fprintf(is.view.textView, "[green]Start Point:[-] (%.1f, %.1f)\n", e.StartPoint.X, e.StartPoint.Y)
		fmt.Fprintf(is.view.textView, "[green]End Point:[-] (%.1f, %.1f)\n", e.EndPoint.X, e.EndPoint.Y)
		fmt.Fprintf(is.view.textView, "[green]Layer:[-] %s\n", e.Layer)
		fmt.Fprintf(is.view.textView, "[green]Color:[-] %d\n", e.Color)

	case *data.CircleInfo:
		fmt.Fprintf(is.view.textView, "[green]Circle Entity[-]\n\n")
		fmt.Fprintf(is.view.textView, "[green]Center:[-] (%.1f, %.1f)\n", e.Center.X, e.Center.Y)
		fmt.Fprintf(is.view.textView, "[green]Radius:[-] %.1f\n", e.Radius)
		fmt.Fprintf(is.view.textView, "[green]Layer:[-] %s\n", e.Layer)
		fmt.Fprintf(is.view.textView, "[green]Color:[-] %d\n", e.Color)

	case *data.TextInfo:
		fmt.Fprintf(is.view.textView, "[green]Text Entity[-]\n\n")
		fmt.Fprintf(is.view.textView, "[green]Value:[-] %s\n", e.Value)
		fmt.Fprintf(is.view.textView, "[green]Insertion Point:[-] (%.1f, %.1f)\n", e.InsertionPoint.X, e.InsertionPoint.Y)
		fmt.Fprintf(is.view.textView, "[green]Height:[-] %.1f\n", e.Height)
		fmt.Fprintf(is.view.textView, "[green]Layer:[-] %s\n", e.Layer)

	case *data.BlockInfo:
		fmt.Fprintf(is.view.textView, "[green]Block Entity[-]\n\n")
		fmt.Fprintf(is.view.textView, "[green]Name:[-] %s\n", e.Name)
		fmt.Fprintf(is.view.textView, "[green]Insertion Point:[-] (%.1f, %.1f)\n", e.InsertionPoint.X, e.InsertionPoint.Y)
		fmt.Fprintf(is.view.textView, "[green]Rotation:[-] %.1f\n", e.Rotation)
		fmt.Fprintf(is.view.textView, "[green]Scale:[-] (%.1f, %.1f)\n", e.Scale.X, e.Scale.Y)
		fmt.Fprintf(is.view.textView, "[green]Layer:[-] %s\n", e.Layer)

		if len(e.Attributes) > 0 {
			fmt.Fprintf(is.view.textView, "[green]Attributes:[-]\n")
			for _, attr := range e.Attributes {
				fmt.Fprintf(is.view.textView, "  %s: %s\n", attr.Tag, attr.Value)
			}
		}
	}
}

// GetSelectionState returns the current selection state
func (is *EnhancedItemSelector) GetSelectionState() *SelectionState {
	return is.selection
}
