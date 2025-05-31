package main

import (
	"fmt"
	"log"
	"os"

	"github.com/remym/go-dwg-extractor/cmd"
)

// Build-time variables (injected via ldflags)
var (
	version   = "dev"
	gitCommit = "unknown"
	buildTime = "unknown"
)

func main() {
	// Check for version flag first
	if len(os.Args) >= 2 {
		switch os.Args[1] {
		case "version", "--version", "-v":
			showVersion()
			return
		case "help", "--help", "-h":
			showHelp()
			return
		}
	}

	// Ensure at least one command is provided
	if len(os.Args) < 2 {
		log.Fatalf("No command provided. Usage: %s [extract|tui] [options]", os.Args[0])
	}

	if err := cmd.Execute(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

// showVersion displays version information
func showVersion() {
	fmt.Printf("Go DWG Extractor\n")
	fmt.Printf("Version:    %s\n", version)
	fmt.Printf("Git Commit: %s\n", gitCommit)
	fmt.Printf("Build Time: %s\n", buildTime)
}

// showHelp displays help information
func showHelp() {
	fmt.Printf("Go DWG Extractor - Extract data from DWG files using a Terminal User Interface\n\n")
	fmt.Printf("Usage: %s [command] [options]\n\n", os.Args[0])
	fmt.Printf("Commands:\n")
	fmt.Printf("  extract    Extract data from DWG file and output to console\n")
	fmt.Printf("  tui        Launch Terminal User Interface\n")
	fmt.Printf("  version    Show version information\n")
	fmt.Printf("  help       Show this help message\n\n")
	fmt.Printf("Options:\n")
	fmt.Printf("  -file      Path to DWG file (required for extract command)\n")
	fmt.Printf("  -output    Output directory for conversion (optional)\n\n")
	fmt.Printf("Examples:\n")
	fmt.Printf("  %s extract -file sample.dwg\n", os.Args[0])
	fmt.Printf("  %s tui -file sample.dwg\n", os.Args[0])
	fmt.Printf("  %s tui  # Uses sample data if no file specified\n", os.Args[0])
	fmt.Printf("  %s version\n", os.Args[0])
}
