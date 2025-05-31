package clipboard

import (
	"testing"

	"github.com/remym/go-dwg-extractor/pkg/data"
	"github.com/stretchr/testify/assert"
)

// TestFormatEntityForClipboard tests formatting individual entities
func TestFormatEntityForClipboard(t *testing.T) {
	tests := []struct {
		name           string
		entity         data.Entity
		expectedFormat string
		expectedFields []string
	}{
		{
			name: "LineInfo formatting",
			entity: &data.LineInfo{
				StartPoint: data.Point{X: 10.5, Y: 20.5},
				EndPoint:   data.Point{X: 30.5, Y: 40.5},
				Layer:      "Layer1",
				Color:      7,
			},
			expectedFormat: "Line: (10.5, 20.5) to (30.5, 40.5), Layer: Layer1, Color: 7",
			expectedFields: []string{"Line:", "10.5", "20.5", "30.5", "40.5", "Layer1", "Color: 7"},
		},
		{
			name: "CircleInfo formatting",
			entity: &data.CircleInfo{
				Center: data.Point{X: 15.0, Y: 25.0},
				Radius: 5.5,
				Layer:  "CircleLayer",
				Color:  3,
			},
			expectedFormat: "Circle: Center (15.0, 25.0), Radius: 5.5, Layer: CircleLayer, Color: 3",
			expectedFields: []string{"Circle:", "Center", "15.0", "25.0", "Radius: 5.5", "CircleLayer", "Color: 3"},
		},
		{
			name: "TextInfo formatting",
			entity: &data.TextInfo{
				Value:          "Sample Text Content",
				InsertionPoint: data.Point{X: 5.0, Y: 10.0},
				Height:         12.0,
				Layer:          "TextLayer",
			},
			expectedFormat: "Text: \"Sample Text Content\", InsertionPoint: (5.0, 10.0), Height: 12.0, Layer: TextLayer",
			expectedFields: []string{"Text:", "Sample Text Content", "InsertionPoint:", "5.0", "10.0", "Height: 12.0", "TextLayer"},
		},
		{
			name: "BlockInfo formatting with attributes",
			entity: &data.BlockInfo{
				Name:           "TestBlock",
				InsertionPoint: data.Point{X: 0.0, Y: 0.0},
				Rotation:       45.0,
				Scale:          data.Point{X: 2.0, Y: 1.5},
				Layer:          "BlockLayer",
				Attributes: []data.AttributeInfo{
					{Tag: "TAG1", Value: "Value1"},
					{Tag: "TAG2", Value: "Value2"},
				},
			},
			expectedFormat: "Block: TestBlock, InsertionPoint: (0.0, 0.0), Rotation: 45.0, Scale: (2.0, 1.5), Layer: BlockLayer, Attributes: [TAG1:Value1, TAG2:Value2]",
			expectedFields: []string{"Block:", "TestBlock", "InsertionPoint:", "0.0", "Rotation: 45.0", "Scale:", "2.0", "1.5", "BlockLayer", "TAG1:Value1", "TAG2:Value2"},
		},
		{
			name: "BlockInfo formatting without attributes",
			entity: &data.BlockInfo{
				Name:           "SimpleBlock",
				InsertionPoint: data.Point{X: 10.0, Y: 20.0},
				Rotation:       0.0,
				Scale:          data.Point{X: 1.0, Y: 1.0},
				Layer:          "Layer1",
				Attributes:     []data.AttributeInfo{},
			},
			expectedFormat: "Block: SimpleBlock, InsertionPoint: (10.0, 20.0), Rotation: 0.0, Scale: (1.0, 1.0), Layer: Layer1",
			expectedFields: []string{"Block:", "SimpleBlock", "InsertionPoint:", "10.0", "20.0", "Rotation: 0.0", "Scale:", "1.0", "Layer1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This should fail initially - we need to implement FormatEntityForClipboard
			formatter := NewClipboardFormatter()
			result := formatter.FormatEntityForClipboard(tt.entity)

			assert.Equal(t, tt.expectedFormat, result, "Expected exact format match")

			// Check that all expected fields are present
			for _, field := range tt.expectedFields {
				assert.Contains(t, result, field, "Expected field '%s' to be present", field)
			}
		})
	}
}

// TestFormatMultipleEntitiesForClipboard tests formatting multiple entities
func TestFormatMultipleEntitiesForClipboard(t *testing.T) {
	tests := []struct {
		name             string
		entities         []data.Entity
		expectedLines    int
		expectedContains []string
	}{
		{
			name: "Multiple different entities",
			entities: []data.Entity{
				&data.LineInfo{
					StartPoint: data.Point{X: 0, Y: 0},
					EndPoint:   data.Point{X: 10, Y: 10},
					Layer:      "Layer1",
					Color:      1,
				},
				&data.CircleInfo{
					Center: data.Point{X: 5, Y: 5},
					Radius: 2.5,
					Layer:  "Layer1",
					Color:  2,
				},
				&data.TextInfo{
					Value:          "Test Text",
					InsertionPoint: data.Point{X: 1, Y: 1},
					Height:         12.0,
					Layer:          "Layer1",
				},
			},
			expectedLines:    3,
			expectedContains: []string{"Line:", "Circle:", "Text:", "Layer1"},
		},
		{
			name:             "Empty entities list",
			entities:         []data.Entity{},
			expectedLines:    0,
			expectedContains: []string{},
		},
		{
			name: "Single entity",
			entities: []data.Entity{
				&data.LineInfo{
					StartPoint: data.Point{X: 1, Y: 2},
					EndPoint:   data.Point{X: 3, Y: 4},
					Layer:      "SingleLayer",
					Color:      5,
				},
			},
			expectedLines:    1,
			expectedContains: []string{"Line:", "SingleLayer"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This should fail initially - we need to implement FormatMultipleEntitiesForClipboard
			formatter := NewClipboardFormatter()
			result := formatter.FormatMultipleEntitiesForClipboard(tt.entities)

			if tt.expectedLines == 0 {
				assert.Empty(t, result, "Expected empty result for no entities")
			} else {
				lines := len(result)
				assert.Equal(t, tt.expectedLines, lines, "Expected correct number of formatted lines")

				// Join all lines to check for content
				allContent := ""
				for _, line := range result {
					allContent += line + " "
				}

				for _, expected := range tt.expectedContains {
					assert.Contains(t, allContent, expected, "Expected content '%s' to be present", expected)
				}
			}
		})
	}
}

// TestFormatForCSV tests CSV formatting for spreadsheet compatibility
func TestFormatForCSV(t *testing.T) {
	tests := []struct {
		name         string
		entities     []data.Entity
		expectedCols []string
		expectedRows int
	}{
		{
			name: "Mixed entities to CSV",
			entities: []data.Entity{
				&data.LineInfo{
					StartPoint: data.Point{X: 0, Y: 0},
					EndPoint:   data.Point{X: 10, Y: 10},
					Layer:      "Layer1",
					Color:      1,
				},
				&data.TextInfo{
					Value:          "CSV Test",
					InsertionPoint: data.Point{X: 5, Y: 5},
					Height:         10.0,
					Layer:          "Layer2",
				},
			},
			expectedCols: []string{"Type", "Layer", "Details"},
			expectedRows: 3, // Header + 2 data rows
		},
		{
			name:         "Empty entities to CSV",
			entities:     []data.Entity{},
			expectedCols: []string{"Type", "Layer", "Details"},
			expectedRows: 1, // Header only
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This should fail initially - we need to implement FormatAsCSV
			formatter := NewClipboardFormatter()
			result := formatter.FormatAsCSV(tt.entities)

			lines := len(result)
			assert.Equal(t, tt.expectedRows, lines, "Expected correct number of CSV rows")

			if len(result) > 0 {
				// Check header row
				headerRow := result[0]
				for _, col := range tt.expectedCols {
					assert.Contains(t, headerRow, col, "Expected column '%s' in header", col)
				}
			}
		})
	}
}

// TestFormatForJSON tests JSON formatting
func TestFormatForJSON(t *testing.T) {
	tests := []struct {
		name     string
		entities []data.Entity
		wantErr  bool
	}{
		{
			name: "Valid entities to JSON",
			entities: []data.Entity{
				&data.LineInfo{
					StartPoint: data.Point{X: 0, Y: 0},
					EndPoint:   data.Point{X: 5, Y: 5},
					Layer:      "JSONLayer",
					Color:      1,
				},
			},
			wantErr: false,
		},
		{
			name:     "Empty entities to JSON",
			entities: []data.Entity{},
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This should fail initially - we need to implement FormatAsJSON
			formatter := NewClipboardFormatter()
			result, err := formatter.FormatAsJSON(tt.entities)

			if tt.wantErr {
				assert.Error(t, err, "Expected error for invalid JSON formatting")
			} else {
				assert.NoError(t, err, "Expected no error for valid JSON formatting")
				assert.NotEmpty(t, result, "Expected non-empty JSON result")

				// Basic JSON validation - should start with [ or {
				if len(tt.entities) > 0 {
					assert.True(t, result[0] == '[' || result[0] == '{', "Expected valid JSON format")
				}
			}
		})
	}
}

// TestClipboardFormatterEdgeCases tests edge cases and error conditions
func TestClipboardFormatterEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(t *testing.T)
	}{
		{
			name: "NewClipboardFormatter creates valid formatter",
			testFunc: func(t *testing.T) {
				formatter := NewClipboardFormatter()
				assert.NotNil(t, formatter, "Expected formatter to be created")
			},
		},
		{
			name: "Nil entity handling",
			testFunc: func(t *testing.T) {
				formatter := NewClipboardFormatter()
				// This should not panic and should handle gracefully
				result := formatter.FormatEntityForClipboard(nil)
				assert.NotEmpty(t, result, "Expected some result for nil entity")
				assert.Contains(t, result, "Unknown", "Expected unknown entity indication")
			},
		},
		{
			name: "Special characters in text entities",
			testFunc: func(t *testing.T) {
				formatter := NewClipboardFormatter()
				entity := &data.TextInfo{
					Value:          "Text with \"quotes\" and, commas",
					InsertionPoint: data.Point{X: 0, Y: 0},
					Height:         10.0,
					Layer:          "SpecialLayer",
				}

				result := formatter.FormatEntityForClipboard(entity)
				assert.Contains(t, result, "quotes", "Expected quotes to be preserved")
				assert.Contains(t, result, "commas", "Expected commas to be preserved")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunc(t)
		})
	}
}
