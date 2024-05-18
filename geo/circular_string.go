package geo

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
)

/*
	https://postgis.net/docs/using_postgis_dbmanagement.html#CircularString

CircularString is the basic curve type, similar to a LineString in the linear world.
A single arc segment is specified by three points: the start and end points (first and third)
and some other point on the arc. To specify a closed circle the start and end points are the
same and the middle point is the opposite point on the circle diameter (which is the center
of the arc). In a sequence of arcs the end point of the previous arc is the start point of the
next arc, just like the segments of a LineString. This means that a CircularString must have
an odd number of points greater than 1.
*/
type CircularString struct {
	Points     []Point
	Dimensions Dimensions
}

func (g CircularString) GetGISGeometryType() GISGeometryType {
	return CircularStringType
}

// Stringer interface
func (c CircularString) String() string {
	var sb strings.Builder
	sb.WriteString("(CircularString ")
	sb.WriteString(c.Dimensions.String())
	sb.WriteString(" [")
	for _, p := range c.Points {
		sb.WriteString(p.String())
		sb.WriteString(" ")
	}

	sb.WriteString("])")
	return sb.String()
}

// Get the dimensions of the geometry
func (c CircularString) GetDimensions() Dimensions {
	return c.Dimensions
}

// Create a new CircularString from input slice of Points
// A CircularString is specified by three points: the start and end points (first and third)
// and some other point on the arc.
// Point array must be an odd number of points greater than 1.
func NewCircularString(p []Point) (*CircularString, error) {

	if len(p)%2 != 1 || len(p) < 3 {
		return nil, fmt.Errorf("circularstring must contain an odd number of points greater than 1")
	}
	c := CircularString{}
	c.Dimensions = p[0].Dimensions
	c.Points = p
	return &c, nil
}

// Create Circular String from input byte buffer in EWKB format and dimensions.
// A CircularString is specified by three points: the start and end points (first and third)
// and some other point on the arc.
func CircularStringFromEWKB(b *bytes.Buffer, dimensions Dimensions) (*CircularString, error) {

	cs := CircularString{}
	cs.Dimensions = dimensions
	if b.Len() < 4 {
		return nil, fmt.Errorf("byte array for circularstring must be at least length 4")
	}
	count := binary.LittleEndian.Uint32(b.Next(4))
	if count < 3 || count%2 != 1 {
		return nil, fmt.Errorf("circularstring must contain an odd number of points greater than 1. %v provided", count)
	}

	// Iterate points in byte slice, adding to struct
	for i := 0; i < int(count); i++ {

		point, err := PointFromEWKB(b, dimensions)
		if err != nil {
			return nil, fmt.Errorf(err.Error())
		}
		cs.Points = append(cs.Points, *point)
	}
	if len(cs.Points) < 3 || len(cs.Points)%2 != 1 {
		return nil, fmt.Errorf("circular string must have an odd number of points, 3 or greater")
	}

	return &cs, nil
}

// Get a byte slice containing the EKWB representation of the geometry
func (c CircularString) GetEWKB(includeGeoType bool) bytes.Buffer {
	buf := new(bytes.Buffer)

	// Include geotype encoding if requested
	if includeGeoType {
		buf.WriteByte(byte(LittleEndian)) // Add Byte Order Marker
		geoTypeBytes := encodeGeoType(CircularStringType, false, c.Dimensions)
		buf.Write(geoTypeBytes)
	}

	// Add an encoded length to the slice
	lenBytes := binary.LittleEndian.AppendUint32([]byte{}, uint32(len(c.Points)))
	buf.Write(lenBytes)

	// iterate points to be added to EWKB
	for _, p := range c.Points {
		pb := p.GetEWKB(false)
		pb.WriteTo(buf)
	}
	return *buf
}
