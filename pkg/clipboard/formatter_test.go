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

// TestFormatAsCSV_ComprehensiveEdgeCases tests all edge cases for CSV formatting
func TestFormatAsCSV_ComprehensiveEdgeCases(t *testing.T) {
	formatter := NewClipboardFormatter()

	tests := []struct {
		name           string
		entities       []data.Entity
		expectedRows   int
		expectedHeader string
		checkContent   func(t *testing.T, result []string)
	}{
		{
			name:           "Empty entities list",
			entities:       []data.Entity{},
			expectedRows:   1,
			expectedHeader: "Type,Layer,Details",
			checkContent: func(t *testing.T, result []string) {
				assert.Equal(t, "Type,Layer,Details", result[0])
			},
		},
		{
			name: "Nil entity in list",
			entities: []data.Entity{
				&data.LineInfo{StartPoint: data.Point{X: 0, Y: 0}, EndPoint: data.Point{X: 1, Y: 1}, Layer: "Layer1", Color: 1},
				nil,
				&data.CircleInfo{Center: data.Point{X: 0, Y: 0}, Radius: 1, Layer: "Layer2", Color: 2},
			},
			expectedRows:   3, // Header + 2 valid entities (nil skipped)
			expectedHeader: "Type,Layer,Details",
			checkContent: func(t *testing.T, result []string) {
				assert.Contains(t, result[1], "Line")
				assert.Contains(t, result[2], "Circle")
			},
		},
		{
			name: "Text with quotes in CSV",
			entities: []data.Entity{
				&data.TextInfo{
					Value:          "Text with \"embedded quotes\"",
					InsertionPoint: data.Point{X: 1, Y: 2},
					Height:         10,
					Layer:          "TextLayer",
				},
			},
			expectedRows:   2,
			expectedHeader: "Type,Layer,Details",
			checkContent: func(t *testing.T, result []string) {
				// Should escape quotes properly
				assert.Contains(t, result[1], "\"\"embedded quotes\"\"")
			},
		},
		{
			name: "Block with attributes in CSV",
			entities: []data.Entity{
				&data.BlockInfo{
					Name:           "TestBlock",
					InsertionPoint: data.Point{X: 0, Y: 0},
					Rotation:       45,
					Scale:          data.Point{X: 1, Y: 1},
					Layer:          "BlockLayer",
					Attributes: []data.AttributeInfo{
						{Tag: "TAG1", Value: "Value1"},
						{Tag: "TAG2", Value: "Value2"},
					},
				},
			},
			expectedRows:   2,
			expectedHeader: "Type,Layer,Details",
			checkContent: func(t *testing.T, result []string) {
				assert.Contains(t, result[1], "Block")
				assert.Contains(t, result[1], "TAG1:Value1")
				assert.Contains(t, result[1], "TAG2:Value2")
			},
		},
		{
			name: "Block without attributes in CSV",
			entities: []data.Entity{
				&data.BlockInfo{
					Name:           "SimpleBlock",
					InsertionPoint: data.Point{X: 5, Y: 10},
					Rotation:       0,
					Scale:          data.Point{X: 1, Y: 1},
					Layer:          "SimpleLayer",
					Attributes:     []data.AttributeInfo{},
				},
			},
			expectedRows:   2,
			expectedHeader: "Type,Layer,Details",
			checkContent: func(t *testing.T, result []string) {
				assert.Contains(t, result[1], "Block")
				assert.Contains(t, result[1], "SimpleBlock")
				assert.NotContains(t, result[1], "Attributes:")
			},
		},
		{
			name: "Polyline in CSV",
			entities: []data.Entity{
				&data.PolylineInfo{
					Points:   []data.Point{{X: 0, Y: 0}, {X: 1, Y: 1}, {X: 2, Y: 0}},
					Layer:    "PolyLayer",
					Color:    3,
					IsClosed: true,
				},
			},
			expectedRows:   2,
			expectedHeader: "Type,Layer,Details",
			checkContent: func(t *testing.T, result []string) {
				assert.Contains(t, result[1], "Polyline")
				assert.Contains(t, result[1], "3 points")
				assert.Contains(t, result[1], "true")
			},
		},
		{
			name: "Unknown entity type in CSV",
			entities: []data.Entity{
				&unknownEntity{layer: "UnknownLayer"},
			},
			expectedRows:   2,
			expectedHeader: "Type,Layer,Details",
			checkContent: func(t *testing.T, result []string) {
				assert.Contains(t, result[1], "Unknown")
				assert.Contains(t, result[1], "UnknownLayer")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatter.FormatAsCSV(tt.entities)

			assert.Equal(t, tt.expectedRows, len(result), "Expected correct number of CSV rows")

			if len(result) > 0 {
				assert.Equal(t, tt.expectedHeader, result[0], "Expected correct CSV header")
			}

			if tt.checkContent != nil {
				tt.checkContent(t, result)
			}
		})
	}
}

// TestFormatAsJSON_ComprehensiveEdgeCases tests all edge cases for JSON formatting
func TestFormatAsJSON_ComprehensiveEdgeCases(t *testing.T) {
	formatter := NewClipboardFormatter()

	tests := []struct {
		name         string
		entities     []data.Entity
		wantErr      bool
		checkContent func(t *testing.T, result string)
	}{
		{
			name:     "Empty entities list",
			entities: []data.Entity{},
			wantErr:  false,
			checkContent: func(t *testing.T, result string) {
				assert.Equal(t, "[]", result)
			},
		},
		{
			name: "Nil entity in list",
			entities: []data.Entity{
				&data.LineInfo{StartPoint: data.Point{X: 0, Y: 0}, EndPoint: data.Point{X: 1, Y: 1}, Layer: "Layer1", Color: 1},
				nil,
				&data.CircleInfo{Center: data.Point{X: 0, Y: 0}, Radius: 1, Layer: "Layer2", Color: 2},
			},
			wantErr: false,
			checkContent: func(t *testing.T, result string) {
				assert.Contains(t, result, "Line")
				assert.Contains(t, result, "Circle")
				// Should have 2 entities, not 3 (nil skipped)
				assert.Contains(t, result, "\"type\": \"Line\"")
				assert.Contains(t, result, "\"type\": \"Circle\"")
			},
		},
		{
			name: "Line entity JSON",
			entities: []data.Entity{
				&data.LineInfo{
					StartPoint: data.Point{X: 1.5, Y: 2.5},
					EndPoint:   data.Point{X: 3.5, Y: 4.5},
					Layer:      "LineLayer",
					Color:      7,
				},
			},
			wantErr: false,
			checkContent: func(t *testing.T, result string) {
				assert.Contains(t, result, "\"type\": \"Line\"")
				assert.Contains(t, result, "\"startPoint\"")
				assert.Contains(t, result, "\"endPoint\"")
				assert.Contains(t, result, "\"x\": 1.5")
				assert.Contains(t, result, "\"y\": 2.5")
				assert.Contains(t, result, "\"color\": 7")
				assert.Contains(t, result, "\"layer\": \"LineLayer\"")
			},
		},
		{
			name: "Circle entity JSON",
			entities: []data.Entity{
				&data.CircleInfo{
					Center: data.Point{X: 10.0, Y: 20.0},
					Radius: 5.5,
					Layer:  "CircleLayer",
					Color:  3,
				},
			},
			wantErr: false,
			checkContent: func(t *testing.T, result string) {
				assert.Contains(t, result, "\"type\": \"Circle\"")
				assert.Contains(t, result, "\"center\"")
				assert.Contains(t, result, "\"radius\": 5.5")
				assert.Contains(t, result, "\"x\": 10")
				assert.Contains(t, result, "\"y\": 20")
				assert.Contains(t, result, "\"color\": 3")
			},
		},
		{
			name: "Text entity JSON",
			entities: []data.Entity{
				&data.TextInfo{
					Value:          "JSON Test Text",
					InsertionPoint: data.Point{X: 7.5, Y: 8.5},
					Height:         12.5,
					Layer:          "TextLayer",
				},
			},
			wantErr: false,
			checkContent: func(t *testing.T, result string) {
				assert.Contains(t, result, "\"type\": \"Text\"")
				assert.Contains(t, result, "\"value\": \"JSON Test Text\"")
				assert.Contains(t, result, "\"insertionPoint\"")
				assert.Contains(t, result, "\"height\": 12.5")
			},
		},
		{
			name: "Block with attributes JSON",
			entities: []data.Entity{
				&data.BlockInfo{
					Name:           "JSONBlock",
					InsertionPoint: data.Point{X: 1, Y: 2},
					Rotation:       90,
					Scale:          data.Point{X: 2, Y: 3},
					Layer:          "BlockLayer",
					Attributes: []data.AttributeInfo{
						{Tag: "ATTR1", Value: "Value1"},
						{Tag: "ATTR2", Value: "Value2"},
					},
				},
			},
			wantErr: false,
			checkContent: func(t *testing.T, result string) {
				assert.Contains(t, result, "\"type\": \"Block\"")
				assert.Contains(t, result, "\"name\": \"JSONBlock\"")
				assert.Contains(t, result, "\"rotation\": 90")
				assert.Contains(t, result, "\"attributes\"")
				assert.Contains(t, result, "\"tag\": \"ATTR1\"")
				assert.Contains(t, result, "\"value\": \"Value1\"")
			},
		},
		{
			name: "Block without attributes JSON",
			entities: []data.Entity{
				&data.BlockInfo{
					Name:           "SimpleJSONBlock",
					InsertionPoint: data.Point{X: 5, Y: 6},
					Rotation:       0,
					Scale:          data.Point{X: 1, Y: 1},
					Layer:          "SimpleLayer",
					Attributes:     []data.AttributeInfo{},
				},
			},
			wantErr: false,
			checkContent: func(t *testing.T, result string) {
				assert.Contains(t, result, "\"type\": \"Block\"")
				assert.Contains(t, result, "\"name\": \"SimpleJSONBlock\"")
				assert.NotContains(t, result, "\"attributes\"")
			},
		},
		{
			name: "Polyline entity JSON",
			entities: []data.Entity{
				&data.PolylineInfo{
					Points:   []data.Point{{X: 0, Y: 0}, {X: 1, Y: 1}},
					Layer:    "PolyLayer",
					Color:    5,
					IsClosed: false,
				},
			},
			wantErr: false,
			checkContent: func(t *testing.T, result string) {
				assert.Contains(t, result, "\"type\": \"Polyline\"")
				assert.Contains(t, result, "\"pointCount\": 2")
				assert.Contains(t, result, "\"closed\": false")
				assert.Contains(t, result, "\"color\": 5")
			},
		},
		{
			name: "Unknown entity type JSON",
			entities: []data.Entity{
				&unknownEntity{layer: "UnknownLayer"},
			},
			wantErr: false,
			checkContent: func(t *testing.T, result string) {
				assert.Contains(t, result, "\"type\": \"Unknown\"")
				assert.Contains(t, result, "\"layer\": \"UnknownLayer\"")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := formatter.FormatAsJSON(tt.entities)

			if tt.wantErr {
				assert.Error(t, err, "Expected error for invalid JSON formatting")
			} else {
				assert.NoError(t, err, "Expected no error for valid JSON formatting")

				if tt.checkContent != nil {
					tt.checkContent(t, result)
				}
			}
		})
	}
}

// TestFormatEntityForClipboard_AllEntityTypes tests all entity types including edge cases
func TestFormatEntityForClipboard_AllEntityTypes(t *testing.T) {
	formatter := NewClipboardFormatter()

	tests := []struct {
		name           string
		entity         data.Entity
		expectedFormat string
	}{
		{
			name:           "Nil entity",
			entity:         nil,
			expectedFormat: "Unknown Entity",
		},
		{
			name: "PolylineInfo formatting",
			entity: &data.PolylineInfo{
				Points:   []data.Point{{X: 0, Y: 0}, {X: 1, Y: 1}, {X: 2, Y: 0}},
				Layer:    "PolyLayer",
				Color:    4,
				IsClosed: true,
			},
			expectedFormat: "Polyline: 3 points, Layer: PolyLayer, Color: 4, Closed: true",
		},
		{
			name:           "Unknown entity type",
			entity:         &unknownEntity{layer: "TestLayer"},
			expectedFormat: "Entity: *clipboard.unknownEntity, Layer: TestLayer",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatter.FormatEntityForClipboard(tt.entity)
			assert.Equal(t, tt.expectedFormat, result)
		})
	}
}

// TestFormatAttributes_EdgeCases tests the formatAttributes helper function
func TestFormatAttributes_EdgeCases(t *testing.T) {
	formatter := NewClipboardFormatter()

	tests := []struct {
		name       string
		attributes []data.AttributeInfo
		expected   string
	}{
		{
			name:       "Empty attributes",
			attributes: []data.AttributeInfo{},
			expected:   "",
		},
		{
			name:       "Nil attributes",
			attributes: nil,
			expected:   "",
		},
		{
			name: "Single attribute",
			attributes: []data.AttributeInfo{
				{Tag: "TAG1", Value: "Value1"},
			},
			expected: "TAG1:Value1",
		},
		{
			name: "Multiple attributes",
			attributes: []data.AttributeInfo{
				{Tag: "TAG1", Value: "Value1"},
				{Tag: "TAG2", Value: "Value2"},
				{Tag: "TAG3", Value: "Value3"},
			},
			expected: "TAG1:Value1, TAG2:Value2, TAG3:Value3",
		},
		{
			name: "Attributes with special characters",
			attributes: []data.AttributeInfo{
				{Tag: "TAG:1", Value: "Value,1"},
				{Tag: "TAG\"2", Value: "Value\"2"},
			},
			expected: "TAG:1:Value,1, TAG\"2:Value\"2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatter.formatAttributes(tt.attributes)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// unknownEntity is a test helper for testing unknown entity types
type unknownEntity struct {
	layer string
}

func (u *unknownEntity) GetLayer() string {
	return u.layer
}
