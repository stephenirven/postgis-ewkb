package geo

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
)

/* https://postgis.net/docs/using_postgis_dbmanagement.html#CompoundCurve
A CompoundCurve is a single continuous curve that may contain both circular arc segments and
linear segments. That means that in addition to having well-formed components, the end point
of every component (except the last) must be coincident with the start point of the following
component.
*/

type CompoundCurve struct {
	Geometry   []GeometrySubtype
	Dimensions Dimensions
}

func (g CompoundCurve) GetGISGeometryType() GISGeometryType {
	return CompoundCurveType
}

// Stringer interface
func (c CompoundCurve) String() string {
	var sb strings.Builder
	sb.WriteString("(CompoundCurve ")
	sb.WriteString(c.Dimensions.String())
	sb.WriteString(" [")
	for _, g := range c.Geometry {
		sb.WriteString(g.String())
		sb.WriteString(" ")
	}

	sb.WriteString("])")
	return sb.String()
}

func (*CompoundCurve) Equals(*CompoundCurve) bool {
	return true
}

// Get the dimensions of the geometry
func (c CompoundCurve) GetDimensions() Dimensions {
	return c.Dimensions
}

// Create a new CompoundCurve
func NewCompoundCurve() (*CompoundCurve, error) {
	cc := CompoundCurve{}
	return &cc, nil
}

func (c *CompoundCurve) AddLineString(ls *LineString) error {
	if c.Dimensions != UNSET && c.Dimensions != ls.Dimensions {
		return fmt.Errorf("dimensions %v do not match %v", ls.Dimensions, c.Dimensions)
	}
	if c.Dimensions == UNSET {
		c.Dimensions = ls.Dimensions
	}
	c.Geometry = append(c.Geometry, ls)
	return nil
}

func (c *CompoundCurve) AddCircularString(cs *CircularString) error {
	if c.Dimensions != UNSET && c.Dimensions != cs.Dimensions {
		return fmt.Errorf("dimensions %v do not match %v", cs.Dimensions, c.Dimensions)
	}
	if c.Dimensions == UNSET {
		c.Dimensions = cs.Dimensions
	}
	c.Geometry = append(c.Geometry, cs)

	return nil
}

// Create CompoundCurve from input byte buffer in EWKB format and dimensions.
func CompoundCurveFromEWKB(buffer *bytes.Buffer, dimensions Dimensions) (*CompoundCurve, error) {

	compoundCurve := CompoundCurve{}
	compoundCurve.Dimensions = dimensions
	if buffer.Len() < 4 {
		return nil, fmt.Errorf("byte array for compoundcurve must be at least length 4")
	}
	count := binary.LittleEndian.Uint32(buffer.Next(4))

	for i := 0; i < int(count); i++ {

		//bom := ByteOrder(data[0]) 	// First byte is Byte Order Mark
		buffer.Next(1)                                                         // move past BOM
		geoType := GISGeometryType(binary.LittleEndian.Uint16(buffer.Next(4))) // 4 bytes for geometry type
		//flags := data[3] // flags are stored on the 3rd byte

		var err error
		var geometry GeometrySubtype
		switch geoType {

		case LineStringType:
			geometry, err = LineStringFromEWKB(buffer, dimensions)
			if err != nil {
				return nil, err
			}

		case CircularStringType:
			geometry, err = CircularStringFromEWKB(buffer, dimensions)
			if err != nil {
				return nil, err
			}
		}
		compoundCurve.Geometry = append(compoundCurve.Geometry, geometry)

	}
	return &compoundCurve, nil
}

// Get a byte slice containing the EKWB representation of the geometry
func (cc CompoundCurve) GetEWKB(includeGeoType bool) bytes.Buffer {
	buf := new(bytes.Buffer)

	// Include geotype encoding if requested
	if includeGeoType {
		buf.WriteByte(byte(LittleEndian)) // Add Byte Order Marker
		geoTypeBytes := encodeGeoType(CompoundCurveType, false, cc.Dimensions)
		buf.Write(geoTypeBytes)
	}

	lenBytes := binary.LittleEndian.AppendUint32([]byte{}, uint32(len(cc.Geometry)))
	buf.Write(lenBytes)

	for _, g := range cc.Geometry {

		var geoTypeBytes []byte
		buf.WriteByte(byte(LittleEndian))
		switch t := g.(type) {
		case *CircularString:
			geoTypeBytes = encodeGeoType(CircularStringType, false, cc.Dimensions)
		case *LineString:
			geoTypeBytes = encodeGeoType(LineStringType, false, cc.Dimensions)
		default:
			fmt.Printf("compoundcurve must only contain circularstring and linestring: %T", t)
		}
		buf.Write(geoTypeBytes)
		gb := g.GetEWKB(false)
		gb.WriteTo(buf)
	}

	return *buf
}
