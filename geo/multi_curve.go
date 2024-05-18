package geo

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
)

/*
	https://postgis.net/docs/using_postgis_dbmanagement.html#MultiCurve

A MultiCurve is a collection of curves which can include LineStrings,
CircularStrings or CompoundCurves.
*/
type MultiCurve struct {
	Geometry   []GeometrySubtype
	Dimensions Dimensions
}

func (g MultiCurve) GetGISGeometryType() GISGeometryType {
	return MultiCurveType
}

// Stringer interface
func (mc MultiCurve) String() string {
	var sb strings.Builder
	sb.WriteString("(MultiCurve ")
	sb.WriteString(mc.Dimensions.String())
	sb.WriteString(" [")
	for _, g := range mc.Geometry {
		sb.WriteString(g.String())
		sb.WriteString(" ")
	}

	sb.WriteString("])")
	return sb.String()
}

// Get the dimensions of the geometry
func (mc MultiCurve) GetDimensions() Dimensions {
	return mc.Dimensions
}

// Create a MultiCurve from a slice of Geometry of the same dimensions.
// [Containing only Polygons and CurvePolygons]
// Length of geometry must be at least 1.
func NewMultiCurve() (*MultiCurve, error) {
	m := MultiCurve{}
	return &m, nil
}

func (mc *MultiCurve) AddLineString(l *LineString) error {
	if mc.Dimensions != UNSET && mc.Dimensions != l.Dimensions {
		return fmt.Errorf("dimensions %v do not match %v", l.Dimensions, mc.Dimensions)
	}
	if mc.Dimensions == UNSET {
		mc.Dimensions = l.Dimensions
	}
	mc.Geometry = append(mc.Geometry, l)
	return nil
}
func (mc *MultiCurve) AddCircularString(cs *CircularString) error {
	if mc.Dimensions != 0 && mc.Dimensions != cs.Dimensions {
		return fmt.Errorf("dimensions %v do not match %v", cs.Dimensions, mc.Dimensions)
	}
	if mc.Dimensions == UNSET {
		mc.Dimensions = cs.Dimensions
	}
	mc.Geometry = append(mc.Geometry, cs)
	return nil
}
func (mc *MultiCurve) AddCompoundCurve(cc *CompoundCurve) error {
	if mc.Dimensions != 0 && mc.Dimensions != cc.Dimensions {
		return fmt.Errorf("dimensions %v do not match %v", cc.Dimensions, mc.Dimensions)
	}
	if mc.Dimensions == UNSET {
		mc.Dimensions = cc.Dimensions
	}
	mc.Geometry = append(mc.Geometry, cc)
	return nil
}

// Create a new MultiCurve from input byte buffer in EWKB format and dimensions.
func MultiCurveFromEWKB(b *bytes.Buffer, dimensions Dimensions) (*MultiCurve, error) {

	multiCurve := MultiCurve{}
	multiCurve.Dimensions = dimensions
	if b.Len() < 4 {
		return nil, fmt.Errorf("byte array for multicurve must be at least length 4")
	}
	count := binary.LittleEndian.Uint32(b.Next(4))

	for i := 0; i < int(count); i++ {

		//bom := ByteOrder(data[0]) 	// First byte is Byte Order Mark
		b.Next(1)
		geoType := GISGeometryType(binary.LittleEndian.Uint16(b.Next(4))) // 4 bytes for geometry type
		//flags := data[3] // flags are stored on the 3rd byte

		var geometry GeometrySubtype
		var err error

		switch geoType {

		case LineStringType:
			geometry, err = LineStringFromEWKB(b, dimensions)
			if err != nil {
				return nil, err
			}
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

		}
		multiCurve.Geometry = append(multiCurve.Geometry, geometry)

	}

	return &multiCurve, nil
}

// Get a byte slice containing the EKWB representation of the geometry
func (mc MultiCurve) GetEWKB(includeGeoType bool) bytes.Buffer {
	buf := new(bytes.Buffer)

	// Include geotype encoding if requested
	if includeGeoType {
		buf.WriteByte(byte(LittleEndian)) // Add Byte Order Marker
		geoTypeBytes := encodeGeoType(MultiCurveType, false, mc.Dimensions)
		buf.Write(geoTypeBytes)
	}

	lenBytes := binary.LittleEndian.AppendUint32([]byte{}, uint32(len(mc.Geometry)))
	buf.Write(lenBytes)

	for _, g := range mc.Geometry {

		buf.WriteByte(byte(LittleEndian))

		var geoTypeBytes []byte
		switch t := g.(type) {
		case *CircularString:
			geoTypeBytes = encodeGeoType(CircularStringType, false, mc.Dimensions)
		case *CompoundCurve:
			geoTypeBytes = encodeGeoType(CompoundCurveType, false, mc.Dimensions)
		case *LineString:
			geoTypeBytes = encodeGeoType(LineStringType, false, mc.Dimensions)
		default:
			fmt.Printf("multicurve must only contain circularstring, compound curve, linestring - %T ", t)
		}
		buf.Write(geoTypeBytes)
		gb := g.GetEWKB(false)
		gb.WriteTo(buf)

	}
	return *buf
}
