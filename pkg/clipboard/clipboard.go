package clipboard

import (
	"fmt"

	"github.com/atotto/clipboard"
)

// ClipboardWriter interface for dependency injection
type ClipboardWriter interface {
	WriteAll(text string) error
}

// ClipboardManager manages clipboard operations
type ClipboardManager struct {
	writer ClipboardWriter
}

// RealClipboardWriter is the production implementation using atotto/clipboard
type RealClipboardWriter struct{}

// WriteAll implements the clipboard write operation using the real clipboard
func (r *RealClipboardWriter) WriteAll(text string) error {
	return clipboard.WriteAll(text)
}

// NewClipboardManager creates a new clipboard manager with the given writer
func NewClipboardManager(writer ClipboardWriter) *ClipboardManager {
	return &ClipboardManager{
		writer: writer,
	}
}

// NewRealClipboardManager creates a clipboard manager with the real clipboard implementation
func NewRealClipboardManager() *ClipboardManager {
	return &ClipboardManager{
		writer: &RealClipboardWriter{},
	}
}

// CopyToClipboard copies the given text to the clipboard
func (c *ClipboardManager) CopyToClipboard(text string) error {
	if c.writer == nil {
		return fmt.Errorf("clipboard writer is nil")
	}

	err := c.writer.WriteAll(text)
	if err != nil {
		return fmt.Errorf("failed to copy to clipboard: %w", err)
	}

	return nil
}
