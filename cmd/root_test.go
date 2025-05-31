package cmd

import (
	"flag"
	"io/ioutil"
	"os"
	"testing"
)

func TestRootCommand(t *testing.T) {
	// Save original command-line arguments
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	tests := []struct {
		name        string
		args        []string
		wantErr     bool
		errContains string
	}{
		{
			name:        "no arguments",
			args:        []string{"cmd"},
			wantErr:     true,
			errContains: "no DWG file specified",
		},
		{
			name:    "with file argument",
			args:    []string{"cmd", "-file"},
			wantErr: false,
		},
		{
			name:        "non-existent file",
			args:        []string{"cmd", "-file", "nonexistent.dwg"},
			wantErr:     true,
			errContains: "file does not exist",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary file for testing
			tempFile, err := ioutil.TempFile("", "testfile*.dwg")
			if err != nil {
				t.Fatalf("Failed to create temp file: %v", err)
			}
			tempFileName := tempFile.Name()
			tempFile.Close()
			defer os.Remove(tempFileName)

			// Update args with the temp file path
			args := make([]string, len(tt.args))
			copy(args, tt.args)
			if len(args) > 1 && args[1] == "-file" && len(args) == 2 {
				args = append(args, tempFileName)
			}

			// Set command-line arguments for the test
			os.Args = args
			// Reset the flag set to avoid flag redefinition errors
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

			execErr := Execute()
			if (execErr != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", execErr, tt.wantErr)
			}
			if tt.wantErr && execErr != nil && tt.errContains != "" {
				if execErr.Error() != tt.errContains && !contains(execErr.Error(), tt.errContains) {
					t.Errorf("Execute() error = %v, want error containing %q", execErr, tt.errContains)
				}
			}
		})
	}
}

// contains checks if a string contains another string
func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr
}
