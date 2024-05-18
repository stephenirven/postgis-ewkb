package geo

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"strings"
)

/*
	https://postgis.net/docs/using_postgis_dbmanagement.html#Point

A Point is a 0-dimensional geometry that represents a single location in coordinate space.
Points may have dimensions XY, XYZ, XYM, or XYZM where X,Y,Z are coordinate values for axes,
and M may be a "Measure" for additional data, time etc.
*/
type Point struct {
	Coords     []float64
	Dimensions Dimensions
}

func (g Point) GetGISGeometryType() GISGeometryType {
	return PointType
}

// Get the dimensions of the geometry
func (p Point) GetDimensions() Dimensions {
	return p.Dimensions
}

// Stringer interface
func (p Point) String() string {
	return "(" + p.Dimensions.String() + " point: [" + strings.Join(strings.Fields(fmt.Sprint(p.Coords)), ",") + "])"
}

// Returns the expected byte length for a point of given dimensions
func PointByteLength(dimensions Dimensions) uint32 {
	switch dimensions {
	case XY:
		return 16 // 2*8 byte float64s
	case XYZ:
		return 24 // 3*8 byte float64s
	case XYM:
		return 24 // 3*8 byte float64s
	case XYZM:
		return 32 // 4*8 byte float64s
	}
	// This can't happen under current spec	(...)
	return 0
}

// Create a Point from a slice of float64 and Dimensions. Length of slice must correspond
// with dimensions.
func NewPoint(c []float64, dimensions Dimensions) (*Point, error) {
	if len(c) < len(fmt.Sprintf("%v", dimensions)) {
		return nil, fmt.Errorf("point of %v dimensions needs length %v, %v provided", dimensions, PointByteLength(dimensions), len(c))
	}

	return &Point{
		Dimensions: dimensions,
		Coords:     c,
	}, nil
}

// Create a new Point from input byte buffer in EWKB format and dimensions.
func PointFromEWKB(buffer *bytes.Buffer, dimensions Dimensions) (*Point, error) {

	if buffer.Len() < int(PointByteLength(dimensions)) {
		return nil, fmt.Errorf("input for point is too short (%v) for requested dimensions %v", buffer.Len(), dimensions)
	}
	p := Point{}
	p.Dimensions = dimensions
	for i := 0; i < len(string(fmt.Sprintf("%v", p.Dimensions))); i++ {
		p.Coords = append(p.Coords, math.Float64frombits(binary.LittleEndian.Uint64(buffer.Next(8))))
	}

	return &p, nil
}

// Get a byte slice containing the EKWB representation of the geometry
func (p Point) GetEWKB(includeGeoType bool) bytes.Buffer {
	buf := new(bytes.Buffer)

	// Include geotype encoding if requested
	if includeGeoType {
		buf.WriteByte(byte(LittleEndian)) // Add Byte Order Marker
		geoTypeBytes := encodeGeoType(PointType, false, p.Dimensions)
		buf.Write(geoTypeBytes)
	}

	// Point encoding is a simple concatenation of float64 bits
	for _, c := range p.Coords {
		bits := binary.LittleEndian.AppendUint64([]byte{}, math.Float64bits(c))
		buf.Write(bits)
	}
	return *buf
}
