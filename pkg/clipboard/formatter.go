package clipboard

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/remym/go-dwg-extractor/pkg/data"
)

// ClipboardFormatter handles formatting of DXF entities for clipboard operations
type ClipboardFormatter struct{}

// NewClipboardFormatter creates a new clipboard formatter
func NewClipboardFormatter() *ClipboardFormatter {
	return &ClipboardFormatter{}
}

// FormatEntityForClipboard formats a single entity for clipboard copying
func (f *ClipboardFormatter) FormatEntityForClipboard(entity data.Entity) string {
	if entity == nil {
		return "Unknown Entity"
	}

	switch e := entity.(type) {
	case *data.LineInfo:
		return fmt.Sprintf("Line: (%.1f, %.1f) to (%.1f, %.1f), Layer: %s, Color: %d",
			e.StartPoint.X, e.StartPoint.Y, e.EndPoint.X, e.EndPoint.Y, e.Layer, e.Color)

	case *data.CircleInfo:
		return fmt.Sprintf("Circle: Center (%.1f, %.1f), Radius: %.1f, Layer: %s, Color: %d",
			e.Center.X, e.Center.Y, e.Radius, e.Layer, e.Color)

	case *data.TextInfo:
		return fmt.Sprintf("Text: \"%s\", InsertionPoint: (%.1f, %.1f), Height: %.1f, Layer: %s",
			e.Value, e.InsertionPoint.X, e.InsertionPoint.Y, e.Height, e.Layer)

	case *data.BlockInfo:
		attributeStr := f.formatAttributes(e.Attributes)
		if attributeStr != "" {
			return fmt.Sprintf("Block: %s, InsertionPoint: (%.1f, %.1f), Rotation: %.1f, Scale: (%.1f, %.1f), Layer: %s, Attributes: [%s]",
				e.Name, e.InsertionPoint.X, e.InsertionPoint.Y, e.Rotation, e.Scale.X, e.Scale.Y, e.Layer, attributeStr)
		}
		return fmt.Sprintf("Block: %s, InsertionPoint: (%.1f, %.1f), Rotation: %.1f, Scale: (%.1f, %.1f), Layer: %s",
			e.Name, e.InsertionPoint.X, e.InsertionPoint.Y, e.Rotation, e.Scale.X, e.Scale.Y, e.Layer)

	case *data.PolylineInfo:
		return fmt.Sprintf("Polyline: %d points, Layer: %s, Color: %d, Closed: %v",
			len(e.Points), e.Layer, e.Color, e.IsClosed)

	default:
		return fmt.Sprintf("Entity: %T, Layer: %s", entity, entity.GetLayer())
	}
}

// FormatMultipleEntitiesForClipboard formats multiple entities as separate lines
func (f *ClipboardFormatter) FormatMultipleEntitiesForClipboard(entities []data.Entity) []string {
	if len(entities) == 0 {
		return []string{}
	}

	result := make([]string, len(entities))
	for i, entity := range entities {
		result[i] = f.FormatEntityForClipboard(entity)
	}

	return result
}

// FormatAsCSV formats entities as CSV for spreadsheet compatibility
func (f *ClipboardFormatter) FormatAsCSV(entities []data.Entity) []string {
	result := []string{"Type,Layer,Details"}

	if len(entities) == 0 {
		return result
	}

	for _, entity := range entities {
		if entity == nil {
			continue
		}

		var entityType, details string
		layer := entity.GetLayer()

		switch e := entity.(type) {
		case *data.LineInfo:
			entityType = "Line"
			details = fmt.Sprintf("\"(%.1f,%.1f) to (%.1f,%.1f), Color: %d\"",
				e.StartPoint.X, e.StartPoint.Y, e.EndPoint.X, e.EndPoint.Y, e.Color)

		case *data.CircleInfo:
			entityType = "Circle"
			details = fmt.Sprintf("\"Center (%.1f,%.1f), Radius: %.1f, Color: %d\"",
				e.Center.X, e.Center.Y, e.Radius, e.Color)

		case *data.TextInfo:
			entityType = "Text"
			details = fmt.Sprintf("\"%s at (%.1f,%.1f), Height: %.1f\"",
				strings.ReplaceAll(e.Value, "\"", "\"\""), e.InsertionPoint.X, e.InsertionPoint.Y, e.Height)

		case *data.BlockInfo:
			entityType = "Block"
			attributeStr := f.formatAttributes(e.Attributes)
			if attributeStr != "" {
				details = fmt.Sprintf("\"%s at (%.1f,%.1f), Rotation: %.1f, Attributes: %s\"",
					e.Name, e.InsertionPoint.X, e.InsertionPoint.Y, e.Rotation, attributeStr)
			} else {
				details = fmt.Sprintf("\"%s at (%.1f,%.1f), Rotation: %.1f\"",
					e.Name, e.InsertionPoint.X, e.InsertionPoint.Y, e.Rotation)
			}

		case *data.PolylineInfo:
			entityType = "Polyline"
			details = fmt.Sprintf("\"%d points, Color: %d, Closed: %v\"",
				len(e.Points), e.Color, e.IsClosed)

		default:
			entityType = "Unknown"
			details = fmt.Sprintf("\"%T\"", entity)
		}

		csvLine := fmt.Sprintf("%s,%s,%s", entityType, layer, details)
		result = append(result, csvLine)
	}

	return result
}

// FormatAsJSON formats entities as JSON
func (f *ClipboardFormatter) FormatAsJSON(entities []data.Entity) (string, error) {
	// Create a simplified structure for JSON serialization
	jsonEntities := make([]map[string]interface{}, 0, len(entities))

	for _, entity := range entities {
		if entity == nil {
			continue
		}

		entityMap := map[string]interface{}{
			"layer": entity.GetLayer(),
		}

		switch e := entity.(type) {
		case *data.LineInfo:
			entityMap["type"] = "Line"
			entityMap["startPoint"] = map[string]float64{"x": e.StartPoint.X, "y": e.StartPoint.Y}
			entityMap["endPoint"] = map[string]float64{"x": e.EndPoint.X, "y": e.EndPoint.Y}
			entityMap["color"] = e.Color

		case *data.CircleInfo:
			entityMap["type"] = "Circle"
			entityMap["center"] = map[string]float64{"x": e.Center.X, "y": e.Center.Y}
			entityMap["radius"] = e.Radius
			entityMap["color"] = e.Color

		case *data.TextInfo:
			entityMap["type"] = "Text"
			entityMap["value"] = e.Value
			entityMap["insertionPoint"] = map[string]float64{"x": e.InsertionPoint.X, "y": e.InsertionPoint.Y}
			entityMap["height"] = e.Height

		case *data.BlockInfo:
			entityMap["type"] = "Block"
			entityMap["name"] = e.Name
			entityMap["insertionPoint"] = map[string]float64{"x": e.InsertionPoint.X, "y": e.InsertionPoint.Y}
			entityMap["rotation"] = e.Rotation
			entityMap["scale"] = map[string]float64{"x": e.Scale.X, "y": e.Scale.Y}

			if len(e.Attributes) > 0 {
				attributes := make([]map[string]string, len(e.Attributes))
				for i, attr := range e.Attributes {
					attributes[i] = map[string]string{"tag": attr.Tag, "value": attr.Value}
				}
				entityMap["attributes"] = attributes
			}

		case *data.PolylineInfo:
			entityMap["type"] = "Polyline"
			entityMap["pointCount"] = len(e.Points)
			entityMap["color"] = e.Color
			entityMap["closed"] = e.IsClosed

		default:
			entityMap["type"] = "Unknown"
		}

		jsonEntities = append(jsonEntities, entityMap)
	}

	jsonBytes, err := json.MarshalIndent(jsonEntities, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal entities to JSON: %w", err)
	}

	return string(jsonBytes), nil
}

// formatAttributes formats a list of attributes as "Tag:Value" pairs
func (f *ClipboardFormatter) formatAttributes(attributes []data.AttributeInfo) string {
	if len(attributes) == 0 {
		return ""
	}

	attributeStrings := make([]string, len(attributes))
	for i, attr := range attributes {
		attributeStrings[i] = fmt.Sprintf("%s:%s", attr.Tag, attr.Value)
	}

	return strings.Join(attributeStrings, ", ")
}
