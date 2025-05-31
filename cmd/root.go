package cmd

import (
	"flag"
	"fmt"
	"os"
)

var (
	dwgFile string
)

// Execute runs the root command
func Execute() error {
	rootCmd := flag.String("file", "", "Path to the DWG file to process")
	flag.Parse()

	if *rootCmd == "" {
		return fmt.Errorf("no DWG file specified. Please provide a file using the -file flag")
	}

	// Check if file exists
	if _, err := os.Stat(*rootCmd); os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %s", *rootCmd)
	}

	dwgFile = *rootCmd
	fmt.Printf("Processing DWG file: %s\n", dwgFile)

	return nil
}
