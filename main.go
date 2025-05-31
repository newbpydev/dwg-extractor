package main

import (
	"log"
	"os"

	"github.com/remym/go-dwg-extractor/cmd"
)

func main() {
	// Ensure at least one argument is provided
	if len(os.Args) < 2 {
		log.Fatalf("No command provided. Usage: %s [extract|tui] [options]", os.Args[0])
	}

	if err := cmd.Execute(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
