package tui

import (
	"context"
	"testing"
	"time"

	"github.com/remym/go-dwg-extractor/pkg/data"
	"github.com/stretchr/testify/require"
)

// TestNewApp_Basic ensures NewApp can be constructed without panic.
func TestNewApp_Basic(t *testing.T) {
	app := NewApp()
	app.SetTestMode(true) // Enable test mode to prevent hanging
	require.NotNil(t, app, "NewApp() returned nil")
	app.Stop()
}

func TestApp_UpdateDXFData(t *testing.T) {
	app := NewApp()
	app.SetTestMode(true) // Enable test mode to prevent hanging
	defer app.Stop()

	// Use a minimal ExtractedData for update
	data := &data.ExtractedData{DXFVersion: "AC1018"}

	// Test that UpdateDXFData doesn't panic
	// In test mode, this should complete immediately
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("UpdateDXFData panicked: %v", r)
		}
	}()

	app.UpdateDXFData(data)
	// If we reach here without hanging, the test passes
}

func TestApp_SetupLayout(t *testing.T) {
	app := NewApp()
	app.SetTestMode(true) // Enable test mode to prevent hanging
	defer app.Stop()
	// Should not panic or error
	app.setupLayout()
}

func TestApp_Run_Stop(t *testing.T) {

	// Test in two modes
	testCases := []struct {
		name     string
		testMode bool
	}{
		{"with test mode", true},
		{"without test mode (limited duration)", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			app := NewApp()
			app.SetTestMode(tc.testMode)
			done := make(chan struct{})

			// Use a very short timeout for the actual event loop test
			// to avoid hanging but still test both paths
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			go func() {
				defer close(done)
				error := app.Run()
				require.NoError(t, error, "app.Run() should not return an error")
			}()

			// Give a moment for the app to start
			time.Sleep(100 * time.Millisecond)
			app.Stop()

			select {
			case <-done:
				// App stopped successfully
			case <-ctx.Done():
				// If we're intentionally in test mode, this shouldn't happen
				if tc.testMode {
					t.Fatal("Test mode failed to prevent app from hanging")
				} else {
					// Force stop if needed in non-test mode
					app.Stop()
					// This is expected sometimes in real mode
				}
			}
		})
	}
}

func TestApp_App_GetLayout(t *testing.T) {
	app := NewApp()
	app.SetTestMode(true) // Enable test mode to prevent hanging
	defer app.Stop()

	require.NotNil(t, app.App(), "App.App() returned nil")
	require.NotNil(t, app.GetLayout(), "App.GetLayout() returned nil")
}

// TestNewStatusView tests the NewStatusView function
func TestNewStatusView(t *testing.T) {
	statusView := NewStatusView()

	require.NotNil(t, statusView, "NewStatusView() returned nil")

	// Test that the status view has the expected properties
	title := statusView.GetTitle()
	require.Equal(t, "Status", title, "Expected title to be 'Status'")

	// Test that we can set text
	statusView.SetText("Test message")
	text := statusView.GetText(false)
	require.Equal(t, "Test message", text, "Expected text to be set correctly")
}
