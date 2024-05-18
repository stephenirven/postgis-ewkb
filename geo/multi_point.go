package geo

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
)

/*
	https://postgis.net/docs/using_postgis_dbmanagement.html#MultiPoint

A MultiPoint is a collection of Points.
*/
type MultiPoint struct {
	Points     []Point
	Dimensions Dimensions
}

func (g MultiPoint) GetGISGeometryType() GISGeometryType {
	return MultiPointType
}

// Stringer interface
func (mp MultiPoint) String() string {
	var sb strings.Builder
	sb.WriteString("(MultiPoint ")
	sb.WriteString(mp.Dimensions.String())
	sb.WriteString(" [")
	for _, p := range mp.Points {
		sb.WriteString(p.String())
		sb.WriteString(" ")
	}

	sb.WriteString("])")
	return sb.String()
}

// Get the dimensions of the geometry
func (mp MultiPoint) GetDimensions() Dimensions {
	return mp.Dimensions
}

// Create a MultiPoint from a slice of Points of the same dimensions.
// Length must be at least 1
func NewMultiPoint(p []Point) (*MultiPoint, error) {
	mp := MultiPoint{}
	mp.Points = p
	if len(p) == 0 {
		return nil, fmt.Errorf("error creating linestring, no points provided")
	}
	mp.Dimensions = p[0].GetDimensions()
	return &mp, nil
}

// Create a new MultiPoint from input byte buffer in EWKB format and dimensions.
func MultiPointFromEWKB(buffer *bytes.Buffer, dimensions Dimensions) (*MultiPoint, error) {

	mp := MultiPoint{}
	mp.Dimensions = dimensions
	if buffer.Len() < 4 {
		return nil, fmt.Errorf("byte array for multipoint must be at least length 4")
	}
	count := binary.LittleEndian.Uint32(buffer.Next(4))

	for i := 0; i < int(count); i++ {
		buffer.Next(1) // Move past Byte Order Marker
		buffer.Next(4) // Move past GeoType

		point, err := PointFromEWKB(buffer, dimensions)
		if err != nil {
			return nil, fmt.Errorf(err.Error())
		}
		mp.Points = append(mp.Points, *point)
	}

	return &mp, nil
}

// Get a byte slice containing the EKWB representation of the geometry
func (mp MultiPoint) GetEWKB(includeGeoType bool) bytes.Buffer {
	buf := new(bytes.Buffer)

	// Include geotype encoding if requested
	if includeGeoType {
		buf.WriteByte(byte(LittleEndian)) // Add Byte Order Marker
		geoTypeBytes := encodeGeoType(MultiPointType, false, mp.Dimensions)
		buf.Write(geoTypeBytes)
	}

	lenBytes := binary.LittleEndian.AppendUint32([]byte{}, uint32(len(mp.Points)))
	buf.Write(lenBytes)

	for _, p := range mp.Points {
		buf.WriteByte(byte(LittleEndian))
		geoTypeBytes := encodeGeoType(PointType, false, mp.Dimensions)
		buf.Write(geoTypeBytes)
		pb := p.GetEWKB(false)
		pb.WriteTo(buf)
	}
	return *buf
}
