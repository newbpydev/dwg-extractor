package main

import (
	"log"


	"github.com/remym/go-dwg-extractor/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatalf("Error executing command: %v", err)
	}
}
