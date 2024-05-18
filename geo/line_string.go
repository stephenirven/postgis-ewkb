package geo

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
)

/*
	https://postgis.net/docs/using_postgis_dbmanagement.html#LineString

A LineString is a 1-dimensional line formed by a contiguous sequence of line
segments. Each line segment is defined by two points, with the end point of
one segment forming the start point of the next segment. An OGC-valid
LineString has either zero or two or more points, but PostGIS also allows
single-point LineStrings. LineStrings may cross themselves (self-intersect).
A LineString is closed if the start and end points are the same. A LineString
is simple if it does not self-intersect.
*/
type LineString struct {
	Points     []Point
	Dimensions Dimensions
}

func (g LineString) GetGISGeometryType() GISGeometryType {
	return LineStringType
}

// Stringer interface
func (l LineString) String() string {
	var sb strings.Builder
	sb.WriteString("(LineString ")
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
func (l LineString) GetDimensions() Dimensions {
	return l.Dimensions
}

// Create a LineString from a slice of Points of the same dimensions.
// Length must be at least 1
func NewLineString(p []Point) (*LineString, error) {
	l := LineString{}
	l.Points = p
	if len(p) == 0 {
		return nil, fmt.Errorf("error creating linestring, no points provided")
	}
	l.Dimensions = p[0].GetDimensions()
	return &l, nil
}

// Returns the expected byte length for a LineString of given dimensions and length
func LineStringByteLength(dimensions Dimensions, length uint32) uint32 {
	return PointByteLength(dimensions) * length
}

// Create a new Linestring from input byte buffer in EWKB format and dimensions.
func LineStringFromEWKB(b *bytes.Buffer, dimensions Dimensions) (*LineString, error) {

	ls := LineString{}
	ls.Dimensions = dimensions

	// Get the length of the LineString
	count := binary.LittleEndian.Uint32(b.Next(4))

	if uint32(b.Len()) < LineStringByteLength(dimensions, count) {
		return nil, fmt.Errorf("input for linestring is too short (%v) for requested length %v, dimensions %v", b.Len(), count, dimensions)
	}

	// Add point data for requested line
	for i := 0; i < int(count); i++ {
		point, err := PointFromEWKB(b, dimensions)
		if err != nil {
			return nil, fmt.Errorf(err.Error())
		}
		ls.Points = append(ls.Points, *point)
	}

	return &ls, nil
}

// Get a byte slice containing the EKWB representation of the geometry
func (l LineString) GetEWKB(includeGeoType bool) bytes.Buffer {
	buf := new(bytes.Buffer)

	// Include geotype encoding if requested
	if includeGeoType {
		buf.WriteByte(byte(LittleEndian)) // Add Byte Order Marker
		geoTypeBytes := encodeGeoType(LineStringType, false, l.Dimensions)
		buf.Write(geoTypeBytes)
	}

	// Encode the length
	lenBytes := binary.LittleEndian.AppendUint32([]byte{}, uint32(len(l.Points)))

	buf.Write(lenBytes)

	for _, p := range l.Points {
		//no BOM or geotype
		qb := p.GetEWKB(false)
		qb.WriteTo(buf)
	}
	return *buf
}
