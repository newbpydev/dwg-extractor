package tui

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Navigator handles navigation between different panes/views
type Navigator interface {
	SetFocus(focusTarget string) error
	HandleNavigation(key tcell.Key, mod tcell.ModMask) bool
	GetCurrentFocus() string
}

// ListNavigator handles navigation within lists
type ListNavigator interface {
	SetCurrentIndex(index int) error
	SetWrapNavigation(wrap bool)
	HandleKeyPress(key tcell.Key, mod tcell.ModMask) bool
	GetCurrentIndex() int
}

// CategorySelector handles category selection and view updates
type CategorySelector interface {
	SelectCategory(categoryType string, index int) error
}

// ItemSelector handles item selection and details pane updates
type ItemSelector interface {
	SelectItem(itemType string, index int) error
}

// BreadcrumbNavigator handles breadcrumb navigation
type BreadcrumbNavigator interface {
	NavigateTo(step string) error
	GetBreadcrumb() string
	CanGoBack() bool
}

// TUINavigator is the concrete implementation of Navigator
type TUINavigator struct {
	app            *tview.Application
	searchInput    *tview.InputField
	layers         *tview.List
	entityList     *tview.List
	helpView       *tview.TextView
	currentFocus   string
	focusableViews map[string]tview.Primitive
}

// NewTUINavigator creates a new navigator
func NewTUINavigator(app *tview.Application, searchInput *tview.InputField, layers *tview.List, entityList *tview.List) *TUINavigator {
	helpView := tview.NewTextView()
	helpView.SetText("Help View\n\nNavigation:\nTab/Shift+Tab: Switch panes\nArrows: Navigate lists\nEnter: Select item\nEsc: Go back\nF1: Show/hide help")
	helpView.SetBorder(true)
	helpView.SetTitle("Help")

	navigator := &TUINavigator{
		app:          app,
		searchInput:  searchInput,
		layers:       layers,
		entityList:   entityList,
		helpView:     helpView,
		currentFocus: "search",
		focusableViews: map[string]tview.Primitive{
			"search":   searchInput,
			"layers":   layers,
			"entities": entityList,
			"help":     helpView,
		},
	}

	return navigator
}

// SetFocus sets the focus to the specified target
func (n *TUINavigator) SetFocus(focusTarget string) error {
	if view, exists := n.focusableViews[focusTarget]; exists {
		n.currentFocus = focusTarget
		n.app.SetFocus(view)
		return nil
	}
	return fmt.Errorf("invalid focus target: %s", focusTarget)
}

// HandleNavigation handles navigation key presses
func (n *TUINavigator) HandleNavigation(key tcell.Key, mod tcell.ModMask) bool {
	switch key {
	case tcell.KeyTab:
		if mod == tcell.ModShift {
			return n.navigateBackward()
		}
		return n.navigateForward()
	case tcell.KeyBacktab:
		return n.navigateBackward()
	case tcell.KeyF1:
		return n.toggleHelp()
	}
	return false
}

// GetCurrentFocus returns the currently focused view
func (n *TUINavigator) GetCurrentFocus() string {
	return n.currentFocus
}

// navigateForward navigates to the next view
func (n *TUINavigator) navigateForward() bool {
	switch n.currentFocus {
	case "search":
		n.SetFocus("layers")
	case "layers":
		// Check if entities are available
		if n.entityList.GetItemCount() > 0 {
			n.SetFocus("entities")
		} else {
			n.SetFocus("search")
		}
	case "entities":
		n.SetFocus("search")
	case "help":
		n.SetFocus("search")
	default:
		n.SetFocus("search")
	}
	return true
}

// navigateBackward navigates to the previous view
func (n *TUINavigator) navigateBackward() bool {
	switch n.currentFocus {
	case "search":
		if n.entityList.GetItemCount() > 0 {
			n.SetFocus("entities")
		} else {
			n.SetFocus("layers")
		}
	case "layers":
		n.SetFocus("search")
	case "entities":
		n.SetFocus("layers")
	case "help":
		n.SetFocus("search")
	default:
		n.SetFocus("search")
	}
	return true
}

// toggleHelp toggles the help view
func (n *TUINavigator) toggleHelp() bool {
	if n.currentFocus == "help" {
		n.SetFocus("search")
	} else {
		n.SetFocus("help")
	}
	return true
}

// TUIListNavigator is the concrete implementation of ListNavigator
type TUIListNavigator struct {
	list         *tview.List
	wrapEnabled  bool
	currentIndex int
}

// NewTUIListNavigator creates a new list navigator
func NewTUIListNavigator(list *tview.List) *TUIListNavigator {
	return &TUIListNavigator{
		list:         list,
		wrapEnabled:  false,
		currentIndex: 0,
	}
}

// SetCurrentIndex sets the current index
func (ln *TUIListNavigator) SetCurrentIndex(index int) error {
	if ln.list == nil {
		return fmt.Errorf("list is nil")
	}

	itemCount := ln.list.GetItemCount()
	if itemCount == 0 {
		ln.currentIndex = 0
		return nil
	}

	if index < 0 || index >= itemCount {
		return fmt.Errorf("index %d out of range [0, %d)", index, itemCount)
	}

	ln.currentIndex = index
	ln.list.SetCurrentItem(index)
	return nil
}

// SetWrapNavigation enables or disables wrap navigation
func (ln *TUIListNavigator) SetWrapNavigation(wrap bool) {
	ln.wrapEnabled = wrap
}

// HandleKeyPress handles key presses for list navigation
func (ln *TUIListNavigator) HandleKeyPress(key tcell.Key, mod tcell.ModMask) bool {
	if ln.list == nil {
		return false
	}

	itemCount := ln.list.GetItemCount()
	if itemCount == 0 {
		return false
	}

	switch key {
	case tcell.KeyDown:
		return ln.moveDown()
	case tcell.KeyUp:
		return ln.moveUp()
	case tcell.KeyPgDn:
		return ln.pageDown()
	case tcell.KeyPgUp:
		return ln.pageUp()
	case tcell.KeyHome:
		return ln.moveToFirst()
	case tcell.KeyEnd:
		return ln.moveToLast()
	}

	return false
}

// GetCurrentIndex returns the current index
func (ln *TUIListNavigator) GetCurrentIndex() int {
	return ln.currentIndex
}

// moveDown moves selection down
func (ln *TUIListNavigator) moveDown() bool {
	itemCount := ln.list.GetItemCount()
	if itemCount == 0 {
		return false
	}

	newIndex := ln.currentIndex + 1
	if newIndex >= itemCount {
		if ln.wrapEnabled {
			newIndex = 0
		} else {
			newIndex = itemCount - 1
		}
	}

	ln.SetCurrentIndex(newIndex)
	return true
}

// moveUp moves selection up
func (ln *TUIListNavigator) moveUp() bool {
	itemCount := ln.list.GetItemCount()
	if itemCount == 0 {
		return false
	}

	newIndex := ln.currentIndex - 1
	if newIndex < 0 {
		if ln.wrapEnabled {
			newIndex = itemCount - 1
		} else {
			newIndex = 0
		}
	}

	ln.SetCurrentIndex(newIndex)
	return true
}

// pageDown moves selection down by page size
func (ln *TUIListNavigator) pageDown() bool {
	itemCount := ln.list.GetItemCount()
	if itemCount == 0 {
		return false
	}

	pageSize := 5 // Default page size
	newIndex := ln.currentIndex + pageSize
	if newIndex >= itemCount {
		newIndex = itemCount - 1
	}

	ln.SetCurrentIndex(newIndex)
	return true
}

// pageUp moves selection up by page size
func (ln *TUIListNavigator) pageUp() bool {
	itemCount := ln.list.GetItemCount()
	if itemCount == 0 {
		return false
	}

	pageSize := 5 // Default page size
	newIndex := ln.currentIndex - pageSize
	if newIndex < 0 {
		newIndex = 0
	}

	ln.SetCurrentIndex(newIndex)
	return true
}

// moveToFirst moves to the first item
func (ln *TUIListNavigator) moveToFirst() bool {
	return ln.SetCurrentIndex(0) == nil
}

// moveToLast moves to the last item
func (ln *TUIListNavigator) moveToLast() bool {
	itemCount := ln.list.GetItemCount()
	if itemCount == 0 {
		return false
	}
	return ln.SetCurrentIndex(itemCount-1) == nil
}

// TUIBreadcrumbNavigator is the concrete implementation of BreadcrumbNavigator
type TUIBreadcrumbNavigator struct {
	path       []string
	breadcrumb string
}

// NewTUIBreadcrumbNavigator creates a new breadcrumb navigator
func NewTUIBreadcrumbNavigator() *TUIBreadcrumbNavigator {
	return &TUIBreadcrumbNavigator{
		path:       make([]string, 0),
		breadcrumb: "",
	}
}

// NavigateTo navigates to a specific step
func (bn *TUIBreadcrumbNavigator) NavigateTo(step string) error {
	bn.path = append(bn.path, step)
	bn.updateBreadcrumb()
	return nil
}

// GetBreadcrumb returns the current breadcrumb
func (bn *TUIBreadcrumbNavigator) GetBreadcrumb() string {
	return bn.breadcrumb
}

// CanGoBack returns whether back navigation is possible
func (bn *TUIBreadcrumbNavigator) CanGoBack() bool {
	return len(bn.path) > 0
}

// updateBreadcrumb updates the breadcrumb string
func (bn *TUIBreadcrumbNavigator) updateBreadcrumb() {
	if len(bn.path) == 0 {
		bn.breadcrumb = ""
		return
	}

	// Capitalize first letter of each path element
	capitalizedPath := make([]string, len(bn.path))
	for i, p := range bn.path {
		if len(p) > 0 {
			capitalizedPath[i] = strings.ToUpper(p[:1]) + p[1:]
		}
	}

	bn.breadcrumb = strings.Join(capitalizedPath, " > ")
}
