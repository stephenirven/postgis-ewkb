package geo

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
)

/*
	https://postgis.net/docs/using_postgis_dbmanagement.html#MultiLineString

A MultiLineString is a collection of LineStrings. A MultiLineString is closed
if each of its elements is closed.
*/
type MultiLineString struct {
	LineStrings []LineString
	Dimensions  Dimensions
}

func (g MultiLineString) GetGISGeometryType() GISGeometryType {
	return MultiLineStringType
}

// Stringer interface
func (mls MultiLineString) String() string {
	var sb strings.Builder
	sb.WriteString("(MultiLineString ")
	sb.WriteString(mls.Dimensions.String())
	sb.WriteString(" [")
	for _, g := range mls.LineStrings {
		sb.WriteString(g.String())
		sb.WriteString(" ")
	}

	sb.WriteString("])")
	return sb.String()
}

// Get the dimensions of the geometry
func (mls MultiLineString) GetDimensions() Dimensions {
	return mls.Dimensions
}

// Create a MultiLineString from a slice of LineStrings of the same dimensions.
// Length must be at least 1
func NewMultiLineString(l []LineString) (*MultiLineString, error) {
	m := MultiLineString{}
	m.LineStrings = l
	if len(l) < 1 {
		return nil, fmt.Errorf("multilinestring must have length of at least 1")
	}
	m.Dimensions = l[0].GetDimensions()

	return &m, nil
}

// Create a new MultiLineString from input byte buffer in EWKB format and dimensions.
func MultiLineStringFromEWKB(buffer *bytes.Buffer, dimensions Dimensions) (*MultiLineString, error) {

	mls := MultiLineString{}
	mls.Dimensions = dimensions
	if buffer.Len() < 4 {
		return nil, fmt.Errorf("byte array for multilinestring must be at least length 4")
	}
	count := binary.LittleEndian.Uint32(buffer.Next(4))

	for i := 0; i < int(count); i++ {

		// we don't need these, as they are passed from parent
		//bom := data[0]
		buffer.Next(1)
		//geometryType := data[1:5]
		buffer.Next(4)

		lineString, err := LineStringFromEWKB(buffer, dimensions)
		if err != nil {
			return nil, fmt.Errorf(err.Error())
		}

		mls.LineStrings = append(mls.LineStrings, *lineString)
	}

	return &mls, nil
}

// Get a byte slice containing the EKWB representation of the geometry
func (ml MultiLineString) GetEWKB(includeGeoType bool) bytes.Buffer {
	buf := new(bytes.Buffer)

	// Include geotype encoding if requested
	if includeGeoType {
		buf.WriteByte(byte(LittleEndian)) // Add Byte Order Marker
		geoTypeBytes := encodeGeoType(MultiLineStringType, false, ml.Dimensions)
		buf.Write(geoTypeBytes)
	}

	lenBytes := binary.LittleEndian.AppendUint32([]byte{}, uint32(len(ml.LineStrings)))

	buf.Write(lenBytes)

	for _, l := range ml.LineStrings {

		buf.WriteByte(byte(LittleEndian))
		geoTypeBytes := encodeGeoType(LineStringType, false, ml.Dimensions)
		buf.Write(geoTypeBytes)
		lb := l.GetEWKB(false)

		lb.WriteTo(buf)
	}
	return *buf
}
