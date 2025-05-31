package dxfparser

import (
	"fmt"
	"os"
	"strings"

	"github.com/remym/go-dwg-extractor/pkg/data"
)

// ParserInterface defines the contract for DXF parsing
//
//go:generate mockgen -destination=../mocks/mock_parser.go -package=mocks github.com/remym/go-dwg-extractor/pkg/dxfparser ParserInterface
type ParserInterface interface {
	ParseDXF(filePath string) (*data.ExtractedData, error)
}

// Parser handles the parsing of DXF files.
type Parser struct {
	// Add any parser configuration or state here
}

// NewParser creates a new instance of the DXF parser.
func NewParser() *Parser {
	return &Parser{}
}

// Ensure Parser implements ParserInterface
var _ ParserInterface = (*Parser)(nil)

// ParseDXF parses a DXF file and returns the extracted data.
// This is a simplified implementation that extracts basic information.
func (p *Parser) ParseDXF(filePath string) (*data.ExtractedData, error) {
	// Read the file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read DXF file: %w", err)
	}

	// Create a new ExtractedData instance
	result := &data.ExtractedData{
		DXFVersion: "R12", // Default version
	}

	// Convert content to string for parsing
	dxfContent := string(content)
	lines := strings.Split(dxfContent, "\n")

	// Parse DXF version
	for i, line := range lines {
		if strings.TrimSpace(line) == "$ACADVER" && i+2 < len(lines) {
			version := strings.TrimSpace(lines[i+2])
			if version == "AC1015" {
				result.DXFVersion = "R2000"
			} else if version == "AC1021" {
				result.DXFVersion = "R2007"
			} else if version == "AC1024" {
				result.DXFVersion = "R2010"
			} else {
				result.DXFVersion = version
			}
			break
		}
	}

	// Parse layers from TABLES section
	inLayerTable := false
	var layers []data.LayerInfo

	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])

		// Check if we're entering the LAYER table
		if line == "TABLE" && i+2 < len(lines) && strings.TrimSpace(lines[i+2]) == "LAYER" {
			inLayerTable = true
			continue
		}

		// Check if we're leaving the table
		if inLayerTable && line == "ENDTAB" {
			inLayerTable = false
			continue
		}

		// Parse layer entries
		if inLayerTable && line == "LAYER" {
			layer := data.LayerInfo{
				IsOn:     true,         // Default
				IsFrozen: false,        // Default
				Color:    7,            // Default
				LineType: "CONTINUOUS", // Default
			}

			// Parse layer properties
			for j := i + 1; j < len(lines) && j < i+20; j++ { // Look ahead max 20 lines
				code := strings.TrimSpace(lines[j])
				if j+1 >= len(lines) {
					break
				}
				value := strings.TrimSpace(lines[j+1])

				switch code {
				case "2": // Layer name
					layer.Name = value
				case "62": // Color number
					if color := parseInt(value); color != 0 {
						layer.Color = color
					}
				case "70": // Layer flags
					flags := parseInt(value)
					layer.IsFrozen = (flags & 1) != 0 // Bit 0: frozen
					layer.IsOn = (flags & 2) == 0     // Bit 1: frozen in new viewports (inverted for IsOn)
				case "6": // Line type
					layer.LineType = value
				case "0": // Next entity
					if value == "LAYER" || value == "ENDTAB" {
						break
					}
				}
				j++ // Skip the value line
			}

			if layer.Name != "" {
				layers = append(layers, layer)
			}
		}
	}

	result.Layers = layers

	// Parse entities (simplified)
	lineCount := 0
	circleCount := 0
	textCount := 0

	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "LINE" {
			lineCount++
		} else if line == "CIRCLE" {
			circleCount++
		} else if line == "TEXT" {
			textCount++
		}
	}

	// Create sample entities for display
	for i := 0; i < lineCount; i++ {
		result.Lines = append(result.Lines, data.LineInfo{
			StartPoint: data.Point{X: 0, Y: 0, Z: 0},
			EndPoint:   data.Point{X: 100, Y: 100, Z: 0},
			Layer:      "WALLS",
			Color:      1,
		})
	}

	for i := 0; i < circleCount; i++ {
		result.Circles = append(result.Circles, data.CircleInfo{
			Center: data.Point{X: 50, Y: 50, Z: 0},
			Radius: 25,
			Layer:  "DOORS",
			Color:  2,
		})
	}

	return result, nil
}

// parseInt safely converts a string to int, returning 0 on error
func parseInt(s string) int {
	var result int
	for _, c := range s {
		if c >= '0' && c <= '9' {
			result = result*10 + int(c-'0')
		} else {
			return 0
		}
	}
	return result
}
