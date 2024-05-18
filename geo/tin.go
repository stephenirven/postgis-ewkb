package geo

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
)

/*
	https://postgis.net/docs/using_postgis_dbmanagement.html#TIN

A TIN is a collection of non-overlapping Triangles representing a Triangulated
Irregular Network.
*/
type TIN struct {
	Triangles  []Triangle
	Dimensions Dimensions
}

func (g TIN) GetGISGeometryType() GISGeometryType {
	return TINType
}

// Stringer interface
func (t TIN) String() string {
	var sb strings.Builder
	sb.WriteString("(TIN ")
	sb.WriteString(t.Dimensions.String())
	sb.WriteString(" [")
	for _, tri := range t.Triangles {
		sb.WriteString(tri.String())
		sb.WriteString(" ")
	}

	sb.WriteString("])")
	return sb.String()
}

// Get the dimensions of the geometry
func (t TIN) GetDimensions() Dimensions {
	return t.Dimensions
}

// Create a TIN from a slice of Triangles of the same dimensions.
// Length must be at least 1
func NewTIN(t []Triangle) (*TIN, error) {
	ti := TIN{}
	ti.Triangles = t
	if len(t) < 1 {
		return nil, fmt.Errorf("error creating tin. no triangles provided")
	}
	ti.Dimensions = t[0].GetDimensions()
	return &ti, nil
}

// Create a new TIN from input byte buffer in EWKB format and dimensions.
func TINFromEWKB(b *bytes.Buffer, dimensions Dimensions) (*TIN, error) {

	tin := TIN{}
	tin.Dimensions = dimensions
	if b.Len() < 4 {
		return nil, fmt.Errorf("byte array for tin must be at least length 4")
	}

	count := binary.LittleEndian.Uint32(b.Next(4))

	for i := 0; i < int(count); i++ {
		//bom := data[0]
		b.Next(1)
		// geometryType := data[1:5]
		b.Next(4)

		triangle, err := TriangleFromEWKB(b, dimensions)
		if err != nil {
			return nil, fmt.Errorf(err.Error())
		}

		tin.Triangles = append(tin.Triangles, *triangle)
	}

	return &tin, nil
}

// Get a byte slice containing the EKWB representation of the geometry
func (t TIN) GetEWKB(includeGeoType bool) bytes.Buffer {
	buf := new(bytes.Buffer)

	// Include geotype encoding if requested
	if includeGeoType {
		buf.WriteByte(byte(LittleEndian)) // Add Byte Order Marker
		geoTypeBytes := encodeGeoType(TINType, false, t.Dimensions)
		buf.Write(geoTypeBytes)
	}

	lenBytes := binary.LittleEndian.AppendUint32([]byte{}, uint32(len(t.Triangles)))
	buf.Write(lenBytes)

	for _, t := range t.Triangles {
		buf.WriteByte(byte(LittleEndian))

		geoTypeBytes := encodeGeoType(TriangleType, false, t.Dimensions)
		buf.Write(geoTypeBytes)
		tb := t.GetEWKB(false)
		tb.WriteTo(buf)
	}
	return *buf
}
