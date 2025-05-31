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

// Entity is the interface that all DXF entities must implement.
type Entity interface {
	GetLayer() string
}

// LayerInfo holds information about a DXF layer.
type LayerInfo struct {
	Name     string
	Color    int
	IsOn     bool
	IsFrozen bool
	LineType string
	LineWeight float64
	Entities []Entity // Entities that belong to this layer
}

// AttributeInfo holds information about a block attribute.
type AttributeInfo struct {
	Tag      string
	Value    string
	Position Point
	Layer    string
}

// GetLayer implements the Entity interface for BlockInfo.
func (b BlockInfo) GetLayer() string {
	return b.Layer
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

// GetLayer implements the Entity interface for TextInfo.
func (t TextInfo) GetLayer() string {
	return t.Layer
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

// GetLayer implements the Entity interface for LineInfo.
func (l LineInfo) GetLayer() string {
	return l.Layer
}

// LineInfo holds information about a Line entity.
type LineInfo struct {
	StartPoint Point
	EndPoint   Point
	Layer      string
	Color      int
}

// GetLayer implements the Entity interface for CircleInfo.
func (c CircleInfo) GetLayer() string {
	return c.Layer
}

// CircleInfo holds information about a Circle entity.
type CircleInfo struct {
	Center     Point
	Radius     float64
	Layer      string
	Color      int
}

// GetLayer implements the Entity interface for PolylineInfo.
func (p PolylineInfo) GetLayer() string {
	return p.Layer
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
