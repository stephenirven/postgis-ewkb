package geo

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"

	"github.com/google/go-cmp/cmp"
)

/* https://postgis.net/docs/using_postgis_dbmanagement.html#LinearRing
A LinearRing is a LineString which is both closed and simple. The first and
last points must be equal, and the line must not self-intersect.
*/

type LinearRing struct {
	Points     []Point
	Dimensions Dimensions
}

func (g LinearRing) GetGISGeometryType() GISGeometryType {
	return 0 // This doesn't exist as a standalone type
}

// Stringer interface
func (l LinearRing) String() string {
	var sb strings.Builder
	sb.WriteString("(LinearRing ")
	sb.WriteString(l.Dimensions.String())
	sb.WriteString(" [")
	for _, p := range l.Points {
		sb.WriteString(p.String())
		sb.WriteString(" ")
	}

	sb.WriteString("])")
	return sb.String()
}

// Get the dimensions of the geometry
func (l LinearRing) GetDimensions() Dimensions {
	return l.Dimensions
}

// Returns the expected byte length for a LinearRing of given dimensions and length
func LinearRingByteLength(dimensions Dimensions, length uint32) uint32 {
	return PointByteLength(dimensions) * length
}

// Create a LinearRing from a slice of Points of the same dimensions.
// First and last point must have the same coordinates, and length
// must be at least 3
func NewLinearRing(p []Point) (*LinearRing, error) {
	if len(p) < 3 {
		return nil, fmt.Errorf("linearring must have length of at least 3")
	}
	if !cmp.Equal(p[0], p[len(p)-1]) {
		return nil, fmt.Errorf("first and last point of linearring must be equal")
	}
	l := LinearRing{}
	l.Dimensions = p[0].GetDimensions()
	l.Points = p
	return &l, nil
}

// Create a new LinearRing from input byte buffer in EWKB format and dimensions.
func LinearRingFromEWKB(b *bytes.Buffer, dimensions Dimensions) (*LinearRing, error) {

	lr := LinearRing{}
	lr.Dimensions = dimensions
	if b.Len() < 4 {
		return nil, fmt.Errorf("byte array for linearring must be at least length 4")
	}
	count := binary.LittleEndian.Uint32(b.Next(4))

	if uint32(b.Len()) < LinearRingByteLength(dimensions, count) {
		return nil, fmt.Errorf("input for linearring is too short (%v) for requested length %v, dimensions %v", b.Len(), count, dimensions)
	}

	for i := 0; i < int(count); i++ {
		point, err := PointFromEWKB(b, dimensions)
		if err != nil {
			return nil, fmt.Errorf(err.Error())
		}

		lr.Points = append(lr.Points, *point)
	}
	if !cmp.Equal(lr.Points[0], lr.Points[len(lr.Points)-1]) {
		return nil, fmt.Errorf("first and last point of linearring must be equal")
	}

	return &lr, nil
}

// Get a byte slice containing the EKWB representation of the geometry
func (l LinearRing) GetEWKB(includeGeoType bool) bytes.Buffer {
	buf := new(bytes.Buffer)

	// Include geotype encoding if requested
	//if includeGeoType {
	// Not sure these can be directly encoded as geometries. Do they have a type?

	// buf.WriteByte(byte(LittleEndian)) // Add Byte Order Marker
	// geoTypeBytes := encodeGeoType(LinearRingType, false, l.Dimensions)
	// buf.Write(geoTypeBytes)
	//}

	// Encode the length
	lenBytes := binary.LittleEndian.AppendUint32([]byte{}, uint32(len(l.Points)))
	buf.Write(lenBytes)

	for _, p := range l.Points {
		// no BOM or geotype
		lr := p.GetEWKB(false)
		lr.WriteTo(buf)
	}
	return *buf
}
