package cmd

import "github.com/remym/go-dwg-extractor/pkg/converter"

// newDWGConverter is a variable that holds the function to create a new DWGConverter
// This is used to allow mocking in tests
var newDWGConverter = converter.NewDWGConverter
