package geo

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
)

/*
	https://postgis.net/docs/using_postgis_dbmanagement.html#MultiPolygon

A MultiPolygon is a collection of non-overlapping, non-adjacent Polygons.
Polygons in the collection may touch only at a finite number of points.
*/
type MultiPolygon struct {
	Polygons   []Polygon
	Dimensions Dimensions
}

func (g MultiPolygon) GetGISGeometryType() GISGeometryType {
	return MultiPolygonType
}

// Stringer interface
func (mp MultiPolygon) String() string {
	var sb strings.Builder
	sb.WriteString("(MultiPolygon ")
	sb.WriteString(mp.Dimensions.String())
	sb.WriteString(" [")
	for _, g := range mp.Polygons {
		sb.WriteString(g.String())
		sb.WriteString(" ")
	}

	sb.WriteString("])")
	return sb.String()
}

// Get the dimensions of the geometry
func (mp MultiPolygon) GetDimensions() Dimensions {
	return mp.Dimensions
}

// Create a MultiPolygon from a slice of Polygons of the same dimensions.
// Length must be at least 1
func NewMultiPolygon(p []Polygon) (*MultiPolygon, error) {
	mp := MultiPolygon{}
	mp.Polygons = p
	if len(p) == 0 {
		return nil, fmt.Errorf("error creating multipolygon, no polygons provided")
	}
	mp.Dimensions = p[0].GetDimensions()
	return &mp, nil
}

// Create a new MultiPoint from input byte buffer in EWKB format and dimensions.
func MultiPolygonFromEWKB(buffer *bytes.Buffer, dimensions Dimensions) (*MultiPolygon, error) {

	multiPoly := MultiPolygon{}
	multiPoly.Dimensions = dimensions
	if buffer.Len() < 4 {
		return nil, fmt.Errorf("byte array for multipolygon must be at least length 4")
	}

	count := binary.LittleEndian.Uint32(buffer.Next(4))

	for i := 0; i < int(count); i++ {

		// we don't need these as they are passed in from parent
		//bom := data[0]
		buffer.Next(1)
		//geometryType := data[1:5]
		buffer.Next(4)

		polygon, err := PolygonFromEWKB(buffer, dimensions)
		if err != nil {
			return nil, fmt.Errorf(err.Error())
		}

		multiPoly.Polygons = append(multiPoly.Polygons, *polygon)
	}

	return &multiPoly, nil
}

// Get a byte slice containing the EKWB representation of the geometry
func (mp MultiPolygon) GetEWKB(includeGeoType bool) bytes.Buffer {
	buf := new(bytes.Buffer)

	// Include geotype encoding if requested
	if includeGeoType {
		buf.WriteByte(byte(LittleEndian)) // Add Byte Order Marker
		geoTypeBytes := encodeGeoType(MultiPolygonType, false, mp.Dimensions)
		buf.Write(geoTypeBytes)
	}

	countBytes := binary.LittleEndian.AppendUint32([]byte{}, uint32(len(mp.Polygons)))

	buf.Write(countBytes)

	for _, p := range mp.Polygons {

		buf.WriteByte(byte(LittleEndian))
		geoTypeBytes := encodeGeoType(PolygonType, false, mp.Dimensions)
		buf.Write(geoTypeBytes)
		pb := p.GetEWKB(false)
		buf.Write(pb.Bytes())
	}
	return *buf
}
