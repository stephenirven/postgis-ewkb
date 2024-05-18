package geo

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
)

/*
	https://postgis.net/docs/using_postgis_dbmanagement.html#Polygon

A Polygon is a 2-dimensional planar region, delimited by an exterior boundary
(the shell) and zero or more interior boundaries (holes). Each boundary
is a LinearRing.
*/
type Polygon struct {
	LinearRings []LinearRing
	Dimensions  Dimensions
}

func (g Polygon) GetGISGeometryType() GISGeometryType {
	return PolygonType
}

// Stringer interface
func (p Polygon) String() string {
	var sb strings.Builder
	sb.WriteString("(Polygon ")
	sb.WriteString(p.Dimensions.String())
	sb.WriteString(" [")
	for _, l := range p.LinearRings {
		sb.WriteString(l.String())
		sb.WriteString(" ")
	}

	sb.WriteString("])")
	return sb.String()
}

// Get the dimensions of the geometry
func (p Polygon) GetDimensions() Dimensions {
	return p.Dimensions
}

// Create a Polygon from a slice of LinearRings of the same dimensions.
func NewPolygon(l []LinearRing) (*Polygon, error) {
	p := Polygon{}
	p.LinearRings = l
	p.Dimensions = l[0].GetDimensions()
	return &p, nil
}

// Create a new Polygon from input byte buffer in EWKB format and dimensions.
func PolygonFromEWKB(buffer *bytes.Buffer, dimensions Dimensions) (*Polygon, error) {
	poly := Polygon{}
	poly.Dimensions = dimensions
	if buffer.Len() < 4 {
		return nil, fmt.Errorf("byte array for polygon must be at least length 4")
	}
	count := binary.LittleEndian.Uint32(buffer.Next(4))

	if count < 1 {
		return nil, fmt.Errorf("polygons must have at least one linearring, %v provided", count)
	}

	for i := 0; i < int(count); i++ {
		linearRing, err := LinearRingFromEWKB(buffer, dimensions)
		if err != nil {
			return nil, fmt.Errorf(err.Error())
		}
		poly.LinearRings = append(poly.LinearRings, *linearRing)
	}

	return &poly, nil
}

// Get a byte slice containing the EKWB representation of the geometry
func (p Polygon) GetEWKB(includeGeoType bool) bytes.Buffer {
	buf := new(bytes.Buffer)

	// Include geotype encoding if requested
	if includeGeoType {
		buf.WriteByte(byte(LittleEndian)) // Add Byte Order Marker
		geoTypeBytes := encodeGeoType(PolygonType, false, p.Dimensions)
		buf.Write(geoTypeBytes)
	}

	// Encode the length
	lenBytes := binary.LittleEndian.AppendUint32([]byte{}, uint32(len(p.LinearRings)))
	buf.Write(lenBytes)

	for _, l := range p.LinearRings {
		// no BOM or geotype
		lr := l.GetEWKB(false)

		lr.WriteTo(buf)

	}
	return *buf
}
