package cmd

import (
	"github.com/remym/go-dwg-extractor/pkg/converter"
	"github.com/remym/go-dwg-extractor/pkg/dxfparser"
)

// newDWGConverter is a variable that holds the function to create a new DWGConverter
// This is used to allow mocking in tests
var newDWGConverter = converter.NewDWGConverter

// newParser is a variable that holds the function to create a new Parser
// This is used to allow mocking in tests
var newParser = func() dxfparser.ParserInterface {
	return dxfparser.NewParser()
}
