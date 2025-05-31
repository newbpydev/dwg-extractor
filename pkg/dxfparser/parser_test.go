package dxfparser

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewParser(t *testing.T) {
	p := NewParser()
	assert.NotNil(t, p, "NewParser should return a non-nil Parser")
}

func TestParseDXF_NonExistentFile(t *testing.T) {
	p := NewParser()
	_, err := p.ParseDXF("nonexistent.dxf")
	assert.Error(t, err, "Expected error for non-existent file")
}

func TestParseDXF_EmptyFile(t *testing.T) {
	// Create a temporary empty file
	tmpFile, err := os.CreateTemp("", "test-*.dxf")
	require.NoError(t, err, "Failed to create temp file")
	defer os.Remove(tmpFile.Name())

	p := NewParser()
	result, err := p.ParseDXF(tmpFile.Name())
	require.NoError(t, err, "Unexpected error parsing empty file")
	assert.Equal(t, "R12", result.DXFVersion, "Default DXF version should be R12")
	assert.Empty(t, result.Layers, "No layers should be extracted from empty file")
}

func TestParseDXF_WithLayers(t *testing.T) {
	// Create a simple DXF file with layers
	dxfContent := `0
SECTION
2
HEADER
9
$ACADVER
1
AC1015
0
ENDSEC
0
SECTION
2
TABLES
0
TABLE
2
LAYER
0
LAYER
2
0
70
0
62
7
6
CONTINUOUS
0
ENDTAB
0
ENDSEC
0
EOF`

	tmpFile, err := os.CreateTemp("", "test-*.dxf")
	require.NoError(t, err, "Failed to create temp file")
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(dxfContent)
	require.NoError(t, err, "Failed to write test DXF content")
	tmpFile.Close()

	p := NewParser()
	result, err := p.ParseDXF(tmpFile.Name())
	require.NoError(t, err, "Unexpected error parsing DXF with layers")

	// We expect at least the default layer to be present
	assert.GreaterOrEqual(t, len(result.Layers), 1, "Expected at least one layer")
	assert.Equal(t, "0", result.Layers[0].Name, "Expected default layer name")
	assert.Equal(t, 7, result.Layers[0].Color, "Expected default layer color")
}
