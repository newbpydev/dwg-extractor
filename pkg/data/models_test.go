package data

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPoint_String(t *testing.T) {
	tests := []struct {
		name string
		p    Point
		want string
	}{
		{
			name: "3D point",
			p:    Point{X: 1.0, Y: 2.0, Z: 3.0},
			want: "(1.00, 2.00, 3.00)",
		},
		{
			name: "2D point",
			p:    Point{X: 1.5, Y: 2.5},
			want: "(1.50, 2.50, 0.00)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.p.String())
		})
	}
}

func TestExtractedData_AddLayer(t *testing.T) {
	data := &ExtractedData{}
	layer := LayerInfo{
		Name:     "TestLayer",
		Color:    1,
		IsOn:     true,
		IsFrozen: false,
	}

	data.Layers = append(data.Layers, layer)

	assert.Len(t, data.Layers, 1)
	assert.Equal(t, "TestLayer", data.Layers[0].Name)
	assert.True(t, data.Layers[0].IsOn)
}

func TestBlockInfo_AddAttribute(t *testing.T) {
	block := BlockInfo{
		Name:  "TestBlock",
		Layer: "0",
	}

	attr := AttributeInfo{
		Tag:   "TAG1",
		Value: "Value1",
	}

	block.Attributes = append(block.Attributes, attr)

	assert.Len(t, block.Attributes, 1)
	assert.Equal(t, "TAG1", block.Attributes[0].Tag)
	assert.Equal(t, "Value1", block.Attributes[0].Value)
}

func TestNewExtractedData(t *testing.T) {
	data := &ExtractedData{
		DXFVersion: "AC1018",
	}

	assert.Equal(t, "AC1018", data.DXFVersion)
	assert.Empty(t, data.Layers)
	assert.Empty(t, data.Blocks)
	assert.Empty(t, data.Texts)
}
