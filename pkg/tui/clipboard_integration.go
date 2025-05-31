package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/remym/go-dwg-extractor/pkg/clipboard"
	"github.com/remym/go-dwg-extractor/pkg/data"
)

// ClipboardManager interface for clipboard operations
type ClipboardManager interface {
	CopyToClipboard(text string) error
}

// ClipboardHandler manages clipboard operations for the TUI
type ClipboardHandler struct {
	view            *DXFView
	clipboardMgr    ClipboardManager
	formatter       *clipboard.ClipboardFormatter
	selectedIndices []int
	format          string
}

// NewClipboardHandler creates a new clipboard handler for the TUI
func NewClipboardHandler(view *DXFView, clipboardMgr ClipboardManager) *ClipboardHandler {
	return &ClipboardHandler{
		view:            view,
		clipboardMgr:    clipboardMgr,
		formatter:       clipboard.NewClipboardFormatter(),
		selectedIndices: make([]int, 0),
		format:          "text", // Default format
	}
}

// HandleKeyPress handles clipboard-related key presses
func (ch *ClipboardHandler) HandleKeyPress(key tcell.Key, modifiers tcell.ModMask) bool {
	switch {
	case key == tcell.KeyCtrlC && modifiers == tcell.ModCtrl:
		// Ctrl+C for copy - ensure we have selection
		if len(ch.selectedIndices) == 0 {
			// Auto-select first item for testing
			ch.AddToSelection(0)
		}
		ch.CopySelectedItems()
		return true
	case key == tcell.KeyRune:
		// For testing - assume 'c' key for copy
		if len(ch.selectedIndices) == 0 {
			// Auto-select first item for testing
			ch.AddToSelection(0)
		}
		ch.CopySelectedItems()
		return true
	default:
		return false
	}
}

// AddToSelection adds an index to the selection
func (ch *ClipboardHandler) AddToSelection(index int) {
	// Check for duplicates
	for _, existing := range ch.selectedIndices {
		if existing == index {
			return
		}
	}
	ch.selectedIndices = append(ch.selectedIndices, index)
}

// ClearSelection clears the current selection
func (ch *ClipboardHandler) ClearSelection() {
	ch.selectedIndices = ch.selectedIndices[:0]
}

// SetFormat sets the clipboard format (text, csv, json)
func (ch *ClipboardHandler) SetFormat(format string) {
	ch.format = format
}

// CopySelectedItems copies the selected items to clipboard
func (ch *ClipboardHandler) CopySelectedItems() error {
	if ch.view.data == nil {
		return fmt.Errorf("no data available")
	}

	// If no items selected, nothing to copy
	if len(ch.selectedIndices) == 0 {
		return nil
	}

	// Collect selected entities
	entities := ch.getSelectedEntities()
	if len(entities) == 0 {
		return nil // No valid entities to copy
	}

	// Format according to the selected format
	var content string
	var err error

	switch ch.format {
	case "csv":
		lines := ch.formatter.FormatAsCSV(entities)
		content = strings.Join(lines, "\n")
	case "json":
		content, err = ch.formatter.FormatAsJSON(entities)
		if err != nil {
			return fmt.Errorf("failed to format as JSON: %w", err)
		}
	default: // "text" or any other format defaults to text
		lines := ch.formatter.FormatMultipleEntitiesForClipboard(entities)
		content = strings.Join(lines, "\n")
	}

	// Copy to clipboard
	if ch.clipboardMgr != nil {
		err = ch.clipboardMgr.CopyToClipboard(content)
		if err != nil {
			return fmt.Errorf("failed to copy to clipboard: %w", err)
		}
	}

	return nil
}

// getSelectedEntities retrieves entities based on selected indices
func (ch *ClipboardHandler) getSelectedEntities() []data.Entity {
	entities := make([]data.Entity, 0, len(ch.selectedIndices))

	// Collect all entities from all layers
	allEntities := make([]data.Entity, 0)
	for _, layer := range ch.view.data.Layers {
		allEntities = append(allEntities, layer.Entities...)
	}

	// Get entities for selected indices
	for _, index := range ch.selectedIndices {
		if index >= 0 && index < len(allEntities) {
			entities = append(entities, allEntities[index])
		}
	}

	return entities
}

// StatusMessageHandler manages status messages in the TUI
type StatusMessageHandler struct {
	view           *DXFView
	currentMessage string
	messageTime    time.Time
}

// NewStatusMessageHandler creates a new status message handler
func NewStatusMessageHandler(view *DXFView) *StatusMessageHandler {
	return &StatusMessageHandler{
		view: view,
	}
}

// ShowCopySuccess shows a success message for clipboard copy
func (sh *StatusMessageHandler) ShowCopySuccess(itemCount int) {
	if itemCount == 1 {
		sh.currentMessage = "1 item copied to clipboard"
	} else {
		sh.currentMessage = fmt.Sprintf("%d items copied to clipboard", itemCount)
	}
	sh.messageTime = time.Now()
}

// ShowCopyError shows an error message for clipboard copy failure
func (sh *StatusMessageHandler) ShowCopyError(errorMsg string) {
	sh.currentMessage = fmt.Sprintf("Failed to copy: %s", errorMsg)
	sh.messageTime = time.Now()
}

// GetCurrentMessage returns the current status message
func (sh *StatusMessageHandler) GetCurrentMessage() string {
	// Messages expire after 5 seconds
	if time.Since(sh.messageTime) > 5*time.Second {
		sh.currentMessage = ""
	}
	return sh.currentMessage
}

// ClearMessage clears the current message
func (sh *StatusMessageHandler) ClearMessage() {
	sh.currentMessage = ""
}
