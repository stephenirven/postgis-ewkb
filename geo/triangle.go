package geo

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"

	"github.com/google/go-cmp/cmp"
)

/*
	https://postgis.net/docs/using_postgis_dbmanagement.html#Triangle

A Triangle is a polygon defined by three distinct non-collinear vertices.
Because a Triangle is a polygon it is specified by four coordinates, with
the first and fourth being equal.
*/
type Triangle struct {
	Points     [4]Point
	Dimensions Dimensions
}

func (g Triangle) GetGISGeometryType() GISGeometryType {
	return TriangleType
}

// Stringer interface
func (t Triangle) String() string {
	var sb strings.Builder
	sb.WriteString("(Triangle ")
	sb.WriteString(t.Dimensions.String())
	sb.WriteString(" [")
	for _, p := range t.Points {
		sb.WriteString(p.String())
		sb.WriteString(" ")
	}

	sb.WriteString("])")
	return sb.String()
}

// Get the dimensions of the geometry
func (t Triangle) GetDimensions() Dimensions {
	return t.Dimensions
}

// Create a Triangle from an array of Points, length 4 of the same dimensions.
// First and last Points must be equal.
func NewTriangle(p [4]Point) (*Triangle, error) {
	if !cmp.Equal(p[0], p[3]) {
		return nil, fmt.Errorf("first and last point of triangle must be equal")
	}

	return &Triangle{
		Dimensions: p[0].GetDimensions(),
		Points:     p,
	}, nil

}

// Create a new Triangle from input byte buffer in EWKB format and dimensions.
func TriangleFromEWKB(b *bytes.Buffer, dimensions Dimensions) (*Triangle, error) {

	t := Triangle{}
	t.Dimensions = dimensions
	if b.Len() < 4 {
		return nil, fmt.Errorf("byte array for triangle must be at least length 4")
	}

	b.Next(4) // Move past geotype?

	pointCount := binary.LittleEndian.Uint32(b.Next(4))

	if pointCount != 4 {
		return nil, fmt.Errorf("triangle must contain 4 points (first & last must be the same)")
	}

	for i := 0; i < 4; i++ {
		point, err := PointFromEWKB(b, dimensions)
		if err != nil {
			return nil, fmt.Errorf(err.Error())
		}
		t.Points[i] = *point
	}
	if !cmp.Equal(t.Points[0], t.Points[3]) {
		return nil, fmt.Errorf("first and last point of triangle must be equal")
	}

	return &t, nil
}

// Get a byte slice containing the EKWB representation of the geometry
func (t Triangle) GetEWKB(includeGeoType bool) bytes.Buffer {
	buf := new(bytes.Buffer)

	// Include geotype encoding if requested
	if includeGeoType {
		buf.WriteByte(byte(LittleEndian)) // Add Byte Order Marker
		geoTypeBytes := encodeGeoType(TriangleType, false, t.Dimensions)
		buf.Write(geoTypeBytes)
	}

	ringLengthBytes := binary.LittleEndian.AppendUint32([]byte{}, uint32(1))
	buf.Write(ringLengthBytes)

	pointLengthBytes := binary.LittleEndian.AppendUint32([]byte{}, uint32(len(t.Points)))
	buf.Write(pointLengthBytes)

	for _, p := range t.Points {
		// no geotype stuff
		pb := p.GetEWKB(false)
		pb.WriteTo(buf)
	}
	return *buf
}
