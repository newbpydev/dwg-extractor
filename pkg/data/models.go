package data

import "fmt"

// Point defines a 2D or 3D point.
type Point struct {
	X, Y, Z float64
}

// String returns a string representation of the Point.
func (p Point) String() string {
	return fmt.Sprintf("(%.2f, %.2f, %.2f)", p.X, p.Y, p.Z)
}

// LayerInfo holds information about a DXF layer.
type LayerInfo struct {
	Name     string
	Color    int
	IsOn     bool
	IsFrozen bool
	LineType string
	LineWeight float64
}

// AttributeInfo holds information about a block attribute.
type AttributeInfo struct {
	Tag      string
	Value    string
	Position Point
	Layer    string
}

// BlockInfo holds information about a block instance (Insert entity).
type BlockInfo struct {
	Name           string
	Layer          string
	InsertionPoint Point
	Rotation       float64
	Scale          Point
	Attributes     []AttributeInfo
}

// TextInfo holds information about a Text entity.
type TextInfo struct {
	Value          string
	Layer          string
	InsertionPoint Point
	Height         float64
	Rotation       float64
	Style          string
}

// LineInfo holds information about a Line entity.
type LineInfo struct {
	StartPoint Point
	EndPoint   Point
	Layer      string
	Color      int
}

// CircleInfo holds information about a Circle entity.
type CircleInfo struct {
	Center     Point
	Radius     float64
	Layer      string
	Color      int
}

// PolylineInfo holds information about a Polyline entity.
type PolylineInfo struct {
	Points     []Point
	Layer      string
	Color      int
	IsClosed   bool
}

// ExtractedData holds all data parsed from the DXF.
type ExtractedData struct {
	DXFVersion string
	Layers     []LayerInfo
	Blocks     []BlockInfo
	Texts      []TextInfo
	Lines      []LineInfo
	Circles    []CircleInfo
	Polylines  []PolylineInfo
}
