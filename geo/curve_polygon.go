package geo

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
)

/*
	https://postgis.net/docs/using_postgis_dbmanagement.html#CurvePolygon

A CurvePolygon is like a polygon, with an outer ring and zero or more inner
rings. The difference is that a ring can be a CircularString or CompoundCurve
as well as a LineString.
*/
type CurvePolygon struct {
	Geometry   []GeometrySubtype
	Dimensions Dimensions
}

func (g CurvePolygon) GetGISGeometryType() GISGeometryType {
	return CurvePolygonType
}

// Get the dimensions of the geometry
func (c CurvePolygon) GetDimensions() Dimensions {
	return c.Dimensions
}

// Stringer interface
func (c CurvePolygon) String() string {
	var sb strings.Builder
	sb.WriteString("(CurvePolygon ")
	sb.WriteString(c.Dimensions.String())
	sb.WriteString(" [")
	for _, g := range c.Geometry {
		sb.WriteString(g.String())
		sb.WriteString(" ")
	}

	sb.WriteString("])")
	return sb.String()
}

// Create a new CurvePolygon
func NewCurvePolygon() *CurvePolygon {
	cp := CurvePolygon{}
	return &cp
}

func (c *CurvePolygon) AddCircularString(cs *CircularString) error {
	if c.Dimensions != UNSET && c.Dimensions != cs.Dimensions {
		return fmt.Errorf("dimensions %v do not match %v", cs.Dimensions, c.Dimensions)
	}
	if c.Dimensions == UNSET {
		c.Dimensions = cs.Dimensions
	}
	c.Geometry = append(c.Geometry, cs)

	return nil
}

func (c *CurvePolygon) AddCompoundCurve(cc *CompoundCurve) error {
	if c.Dimensions != UNSET && c.Dimensions != cc.Dimensions {
		return fmt.Errorf("dimensions %v do not match %v", cc.Dimensions, c.Dimensions)
	}
	if c.Dimensions == UNSET {
		c.Dimensions = cc.Dimensions
	}
	c.Geometry = append(c.Geometry, cc)
	return nil
}

func (c *CurvePolygon) AddLineString(ls *LineString) error {
	if c.Dimensions != UNSET && c.Dimensions != ls.Dimensions {
		return fmt.Errorf("dimensions %v do not match %v", ls.Dimensions, c.Dimensions)
	}
	if c.Dimensions == UNSET {
		c.Dimensions = ls.Dimensions
	}
	c.Geometry = append(c.Geometry, ls)
	return nil
}

// Create CurvePolygon from input byte buffer in EWKB format and dimensions.
func CurvePolygonFromEWKB(b *bytes.Buffer, dimensions Dimensions) (*CurvePolygon, error) {

	curvePolygon := CurvePolygon{}
	curvePolygon.Dimensions = dimensions
	if b.Len() < 4 {
		return nil, fmt.Errorf("byte array for curvepolygon must be at least length 4")
	}
	count := binary.LittleEndian.Uint32(b.Next(4))

	for i := 0; i < int(count); i++ {

		//bom := ByteOrder(data[0]) 	// First byte is Byte Order Mark
		b.Next(1)                                                         // move past BOM
		geoType := GISGeometryType(binary.LittleEndian.Uint16(b.Next(4))) // 4 bytes for geometry type
		//flags := data[3] // flags are stored on the 3rd byte

		var err error
		var geometry GeometrySubtype

		switch geoType {

		case CircularStringType:
			geometry, err = CircularStringFromEWKB(b, dimensions)
			if err != nil {
				return nil, err
			}

		case CompoundCurveType:
			geometry, err = CompoundCurveFromEWKB(b, dimensions)
			if err != nil {
				return nil, err
			}

		case LineStringType:
			geometry, err = LineStringFromEWKB(b, dimensions)
			if err != nil {
				return nil, err
			}

		}
		curvePolygon.Geometry = append(curvePolygon.Geometry, geometry)

	}

	return &curvePolygon, nil
}

// Get a byte slice containing the EKWB representation of the geometry
func (cp CurvePolygon) GetEWKB(includeGeoType bool) bytes.Buffer {
	buf := new(bytes.Buffer)

	// Include geotype encoding if requested
	if includeGeoType {
		buf.WriteByte(byte(LittleEndian)) // Add Byte Order Marker
		geoTypeBytes := encodeGeoType(CurvePolygonType, false, cp.Dimensions)
		buf.Write(geoTypeBytes)
	}

	lenBytes := binary.LittleEndian.AppendUint32([]byte{}, uint32(len(cp.Geometry)))
	buf.Write(lenBytes)

	for _, g := range cp.Geometry {

		buf.WriteByte(byte(LittleEndian))

		var geoTypeBytes []byte

		switch t := g.(type) {
		case *CircularString:
			geoTypeBytes = encodeGeoType(CircularStringType, false, cp.Dimensions)
		case *CompoundCurve:
			geoTypeBytes = encodeGeoType(CompoundCurveType, false, cp.Dimensions)
		case *LineString:
			geoTypeBytes = encodeGeoType(LineStringType, false, cp.Dimensions)
		default:
			fmt.Printf("curvepolygon must only contain circularstring, compoundcurve and linestring: %T\n", t)
		}
		buf.Write(geoTypeBytes)
		gb := g.GetEWKB(false)
		gb.WriteTo(buf)
	}

	return *buf
}
