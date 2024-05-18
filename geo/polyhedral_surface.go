package geo

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
)

/*
	https://postgis.net/docs/using_postgis_dbmanagement.html#PolyhedralSurface

A PolyhedralSurface is a contiguous collection of patches or facets which share
some edges. Each patch is a planar Polygon. If the Polygon coordinates have Z
ordinates then the surface is 3-dimensional.
*/
type PolyHedralSurface struct {
	Polygons   []Polygon
	Dimensions Dimensions
}

func (g PolyHedralSurface) GetGISGeometryType() GISGeometryType {
	return PolyHedralSurfaceType
}

// Stringer interface
func (p PolyHedralSurface) String() string {
	var sb strings.Builder
	sb.WriteString("(PolyhedralSurface ")
	sb.WriteString(p.Dimensions.String())
	sb.WriteString(" [")
	for _, g := range p.Polygons {
		sb.WriteString(g.String())
		sb.WriteString(" ")
	}

	sb.WriteString("])")
	return sb.String()
}

// Get the dimensions of the geometry
func (p PolyHedralSurface) GetDimensions() Dimensions {
	return p.Dimensions
}

// Create a PolyhedralSurface from a slice of Polygons of the same dimensions.
// Length must be at least 1
func NewPolyhedralSurface(p []Polygon) (*PolyHedralSurface, error) {
	poly := PolyHedralSurface{}
	poly.Polygons = p
	if len(p) == 0 {
		return nil, fmt.Errorf("error creating polyhedral surface, no polygons provided")
	}
	poly.Dimensions = p[0].GetDimensions()
	return &poly, nil
}

// Create a new PolyhedralSurface from input byte buffer in EWKB format and dimensions.
func PolyhedralSurfaceFromEWKB(b *bytes.Buffer, dimensions Dimensions) (*PolyHedralSurface, error) {

	polyhedralSurface := PolyHedralSurface{}
	polyhedralSurface.Dimensions = dimensions
	if b.Len() < 4 {
		return nil, fmt.Errorf("byte array for polyhedralsurface must be at least length 4")
	}
	count := binary.LittleEndian.Uint32(b.Next(4))

	for i := 0; i < int(count); i++ {

		// We don't need these
		//bom := data[0]
		b.Next(1)
		//geometryType := GISGeometryType(binary.LittleEndian.Uint16(data[1:5]))
		b.Next(4)

		polygon, err := PolygonFromEWKB(b, dimensions)
		if err != nil {
			return nil, fmt.Errorf(err.Error())
		}
		polyhedralSurface.Polygons = append(polyhedralSurface.Polygons, *polygon)
	}

	return &polyhedralSurface, nil
}

// Get a byte slice containing the EKWB representation of the geometry
func (ps PolyHedralSurface) GetEWKB(includeGeoType bool) bytes.Buffer {
	buf := new(bytes.Buffer)

	// Include geotype encoding if requested
	if includeGeoType {
		buf.WriteByte(byte(LittleEndian)) // Add Byte Order Marker
		geoTypeBytes := encodeGeoType(PolyHedralSurfaceType, false, ps.Dimensions)
		buf.Write(geoTypeBytes)
	}

	lenBytes := binary.LittleEndian.AppendUint32([]byte{}, uint32(len(ps.Polygons)))
	buf.Write(lenBytes)

	for _, p := range ps.Polygons {
		buf.WriteByte(byte(LittleEndian))
		geoTypeBytes := encodeGeoType(PolygonType, false, ps.Dimensions)
		buf.Write(geoTypeBytes)
		pb := p.GetEWKB(false)
		pb.WriteTo(buf)
	}
	return *buf
}
