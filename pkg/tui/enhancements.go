package tui

import (
	"github.com/gdamore/tcell/v2"
)

// HelpManager manages help view functionality
type HelpManager struct {
	view        *DXFView
	isVisible   bool
	helpContent string
	helpTitle   string
}

// NewHelpManager creates a new help manager
func NewHelpManager(view *DXFView) *HelpManager {
	content := generateHelpContent()
	return &HelpManager{
		view:        view,
		isVisible:   false,
		helpContent: content,
		helpTitle:   "Help",
	}
}

// HandleKeyPress handles help-related key presses
func (hm *HelpManager) HandleKeyPress(key tcell.Key, modifiers tcell.ModMask) bool {
	switch key {
	case tcell.KeyF1:
		hm.isVisible = !hm.isVisible
		return true
	case tcell.KeyEscape:
		// Escape always hides help
		hm.isVisible = false
		return true
	default:
		return false
	}
}

// IsHelpVisible returns whether help is currently visible
func (hm *HelpManager) IsHelpVisible() bool {
	return hm.isVisible
}

// GetHelpContent returns the help content
func (hm *HelpManager) GetHelpContent() string {
	return hm.helpContent
}

// GetHelpTitle returns the help title
func (hm *HelpManager) GetHelpTitle() string {
	return hm.helpTitle
}

// generateHelpContent creates comprehensive help content
func generateHelpContent() string {
	return `Key Bindings and Navigation:

Navigation:
  ↑/↓     - Navigate lists
  Tab     - Switch between panes
  Enter   - Select item
  
Search and Filter:
  Ctrl+F  - Focus search
  /       - Quick search
  Ctrl+R  - Refresh view
  
Selection and Copy:
  Space   - Toggle selection
  Ctrl+C  - Copy selected items
  Ctrl+A  - Select all
  
Help and Exit:
  F1      - Toggle this help
  Ctrl+H  - Show help
  Ctrl+Q  - Quit application
  Escape  - Clear selection/errors

View Controls:
  Ctrl+1  - Focus layers
  Ctrl+2  - Focus entities
  Ctrl+3  - Focus details
  
Accessibility:
  Ctrl++  - Increase text size
  Ctrl+-  - Decrease text size
  Ctrl+0  - Reset text size`
}

// StylingManager manages visual styling and colors
type StylingManager struct {
	view             *DXFView
	appliedStyles    map[string]string
	appliedColors    map[string][]string
	componentBorders map[string]bool
}

// NewStylingManager creates a new styling manager
func NewStylingManager(view *DXFView) *StylingManager {
	return &StylingManager{
		view:             view,
		appliedStyles:    make(map[string]string),
		appliedColors:    make(map[string][]string),
		componentBorders: make(map[string]bool),
	}
}

// ApplyModernStyling applies modern styling to all components
func (sm *StylingManager) ApplyModernStyling() {
	// Apply styling to each component
	components := []string{"layerList", "entityList", "detailView", "searchInput"}

	for _, component := range components {
		sm.appliedStyles[component] = "modern"
		sm.componentBorders[component] = true

		// Set component-specific colors
		switch component {
		case "layerList":
			sm.appliedColors[component] = []string{"blue", "green", "yellow"}
		case "entityList":
			sm.appliedColors[component] = []string{"blue", "green"}
		case "detailView":
			sm.appliedColors[component] = []string{"white", "gray"}
		case "searchInput":
			sm.appliedColors[component] = []string{"cyan", "white"}
		}
	}
}

// HasBorder returns whether a component has a border
func (sm *StylingManager) HasBorder(component string) bool {
	return sm.componentBorders[component]
}

// GetAppliedColors returns the colors applied to a component
func (sm *StylingManager) GetAppliedColors(component string) []string {
	return sm.appliedColors[component]
}

// GetComponentStyle returns the style applied to a component
func (sm *StylingManager) GetComponentStyle(component string) string {
	return sm.appliedStyles[component]
}

// ShortcutManager manages keyboard shortcuts
type ShortcutManager struct {
	view *DXFView
}

// NewShortcutManager creates a new shortcut manager
func NewShortcutManager(view *DXFView) *ShortcutManager {
	return &ShortcutManager{
		view: view,
	}
}

// HandleKeyPress handles keyboard shortcuts and returns action and handled status
func (sm *ShortcutManager) HandleKeyPress(key tcell.Key, modifiers tcell.ModMask) (string, bool) {
	// Handle Ctrl+key combinations
	if modifiers == tcell.ModCtrl {
		switch key {
		case tcell.KeyCtrlQ:
			return "quit", true
		case tcell.KeyCtrlH:
			return "help", true
		case tcell.KeyCtrlR:
			return "refresh", true
		case tcell.KeyCtrlF:
			return "focus_search", true
		}
	}

	// Handle other keys
	switch key {
	case tcell.KeyEscape:
		return "clear", true
	}

	return "", false
}

// AccessibilityManager manages accessibility features
type AccessibilityManager struct {
	view            *DXFView
	enabledFeatures map[string]bool
	featureValues   map[string]interface{}
}

// NewAccessibilityManager creates a new accessibility manager
func NewAccessibilityManager(view *DXFView) *AccessibilityManager {
	manager := &AccessibilityManager{
		view:            view,
		enabledFeatures: make(map[string]bool),
		featureValues:   make(map[string]interface{}),
	}

	// Initialize accessibility features
	manager.initializeFeatures()
	return manager
}

// initializeFeatures sets up default accessibility features
func (am *AccessibilityManager) initializeFeatures() {
	// Enable all accessibility features by default
	features := map[string]interface{}{
		"high_contrast":       true,
		"screen_reader":       "aria-labels",
		"focus_indicators":    "prominent",
		"colorblind_friendly": "viridis",
		"large_text":          1.5,
	}

	for feature, value := range features {
		am.enabledFeatures[feature] = true
		am.featureValues[feature] = value
	}
}

// IsFeatureEnabled returns whether an accessibility feature is enabled
func (am *AccessibilityManager) IsFeatureEnabled(feature string) bool {
	return am.enabledFeatures[feature]
}

// GetFeatureValue returns the value of an accessibility feature
func (am *AccessibilityManager) GetFeatureValue(feature string) interface{} {
	return am.featureValues[feature]
}

// QuitManager manages graceful application exit
type QuitManager struct {
	view              *DXFView
	hasUnsavedChanges bool
	shouldExit        bool
	promptShown       bool
}

// NewQuitManager creates a new quit manager
func NewQuitManager(view *DXFView) *QuitManager {
	return &QuitManager{
		view: view,
	}
}

// SetUnsavedChanges sets whether there are unsaved changes
func (qm *QuitManager) SetUnsavedChanges(hasChanges bool) {
	qm.hasUnsavedChanges = hasChanges
}

// AttemptQuit attempts to quit the application and returns whether a prompt was shown
func (qm *QuitManager) AttemptQuit(userChoice string) bool {
	// Force quit bypasses all prompts
	if userChoice == "force" {
		qm.shouldExit = true
		qm.promptShown = false
		return false
	}

	// No unsaved changes - quit immediately
	if !qm.hasUnsavedChanges {
		qm.shouldExit = true
		qm.promptShown = false
		return false
	}

	// Unsaved changes - show prompt
	qm.promptShown = true
	switch userChoice {
	case "yes":
		qm.shouldExit = true
	case "no":
		qm.shouldExit = false
	default:
		qm.shouldExit = false
	}

	return true
}

// ShouldExit returns whether the application should exit
func (qm *QuitManager) ShouldExit() bool {
	return qm.shouldExit
}

// LayoutManager manages responsive UI layout
type LayoutManager struct {
	view              *DXFView
	terminalWidth     int
	terminalHeight    int
	currentLayout     string
	visibleComponents []string
}

// NewLayoutManager creates a new layout manager
func NewLayoutManager(view *DXFView) *LayoutManager {
	return &LayoutManager{
		view:              view,
		currentLayout:     "four_pane",
		visibleComponents: []string{},
	}
}

// SetTerminalSize sets the terminal dimensions
func (lm *LayoutManager) SetTerminalSize(width, height int) {
	lm.terminalWidth = width
	lm.terminalHeight = height
}

// ApplyResponsiveLayout applies responsive layout based on terminal size
func (lm *LayoutManager) ApplyResponsiveLayout() {
	if lm.terminalWidth >= 120 && lm.terminalHeight >= 40 {
		lm.currentLayout = "four_pane"
		lm.visibleComponents = []string{"search", "layers", "entities", "details", "help"}
	} else if lm.terminalWidth >= 80 && lm.terminalHeight >= 25 {
		lm.currentLayout = "three_pane"
		lm.visibleComponents = []string{"search", "layers", "entities"}
	} else if lm.terminalWidth >= 60 && lm.terminalHeight >= 20 {
		lm.currentLayout = "two_pane"
		lm.visibleComponents = []string{"layers", "entities"}
	} else {
		lm.currentLayout = "single_pane"
		lm.visibleComponents = []string{"layers"}
	}
}

// GetCurrentLayout returns the current layout name
func (lm *LayoutManager) GetCurrentLayout() string {
	return lm.currentLayout
}

// GetVisibleComponents returns the list of visible components
func (lm *LayoutManager) GetVisibleComponents() []string {
	return lm.visibleComponents
}
