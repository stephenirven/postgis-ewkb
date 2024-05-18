package geo

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
)

/*
	https://postgis.net/docs/using_postgis_dbmanagement.html#GeometryCollection

A GeometryCollection is a heterogeneous (mixed) collection of geometries.
*/
type GeometryCollection struct {
	Geometry   []GeometrySubtype
	Dimensions Dimensions
}

func (g GeometryCollection) GetGISGeometryType() GISGeometryType {
	return GeometryCollectionType
}

// Stringer interface
func (g GeometryCollection) String() string {
	var sb strings.Builder
	sb.WriteString("(GeometryCollection ")
	sb.WriteString(g.Dimensions.String())
	sb.WriteString(" [")
	for _, geom := range g.Geometry {
		sb.WriteString(geom.String())
		sb.WriteString(" ")
	}

	sb.WriteString("])")
	return sb.String()
}

// Get the dimensions of the geometry
func (c GeometryCollection) GetDimensions() Dimensions {
	return c.Dimensions
}

// Create a new GeometryCollection from input slice of GeometrySubTypes
// Geometry slice must be at least length 1
func NewGeometryCollection(g []GeometrySubtype) (*GeometryCollection, error) {
	gc := GeometryCollection{}
	gc.Geometry = g
	if len(g) == 0 {
		return nil, fmt.Errorf("error creating geometrycollection, no geometry provided")
	}

	gc.Dimensions = g[0].GetDimensions()

	return &gc, nil
}

// Create a new GeometryCollection from input byte buffer in EWKB format and dimensions.
func GeometryCollectionFromEWKB(b *bytes.Buffer, dimensions Dimensions) (*GeometryCollection, error) {

	geometryCollection := GeometryCollection{}
	geometryCollection.Dimensions = dimensions
	if b.Len() < 4 {
		return nil, fmt.Errorf("byte array for geometrycollection must be at least length 4")
	}
	count := binary.LittleEndian.Uint32(b.Next(4))

	for i := 0; i < int(count); i++ {

		b.Next(1) // move past BOM
		geoTypeBytes := b.Next(4)
		geoType := GISGeometryType(binary.LittleEndian.Uint16(geoTypeBytes)) // 4 bytes for geometry type
		flags := geoTypeBytes[3]                                             // flags are stored on the 3rd byte

		//var SRID uint16
		if flags&byte(wkbSRID) == byte(wkbSRID) {
			_ = binary.LittleEndian.Uint16(b.Next(4))
		}

		var geometrySubType GeometrySubtype
		var err error

		switch geoType {
		case PointType:
			geometrySubType, err = PointFromEWKB(b, dimensions)
			if err != nil {
				return nil, err
			}
		case LineStringType:
			geometrySubType, err = LineStringFromEWKB(b, dimensions)
			if err != nil {
				return nil, err
			}
		case PolygonType:
			geometrySubType, err = PolygonFromEWKB(b, dimensions)
			if err != nil {
				return nil, err
			}
		case MultiPointType:
			geometrySubType, err = MultiPointFromEWKB(b, dimensions)
			if err != nil {
				return nil, err
			}
		case MultiLineStringType:
			geometrySubType, err = MultiLineStringFromEWKB(b, dimensions)
			if err != nil {
				return nil, err
			}
		case MultiPolygonType:
			geometrySubType, err = MultiPolygonFromEWKB(b, dimensions)
			if err != nil {
				return nil, err
			}
		case CircularStringType:
			geometrySubType, err = CircularStringFromEWKB(b, dimensions)
			if err != nil {
				return nil, err
			}
		case CompoundCurveType:
			geometrySubType, err = CompoundCurveFromEWKB(b, dimensions)
			if err != nil {
				return nil, err
			}
		case CurvePolygonType:
			geometrySubType, err = CurvePolygonFromEWKB(b, dimensions)
			if err != nil {
				return nil, err
			}
		case MultiCurveType:
			geometrySubType, err = MultiCurveFromEWKB(b, dimensions)
			if err != nil {
				return nil, err
			}
		case MultiSurfaceType:
			geometrySubType, err = MultiSurfaceFromEWKB(b, dimensions)
			if err != nil {
				return nil, err
			}
		case PolyHedralSurfaceType:
			geometrySubType, err = PolyhedralSurfaceFromEWKB(b, dimensions)
			if err != nil {
				return nil, err
			}
		case TINType:
			geometrySubType, err = TINFromEWKB(b, dimensions)
			if err != nil {
				return nil, err
			}
		case TriangleType:
			geometrySubType, err = TriangleFromEWKB(b, dimensions)
			if err != nil {
				return nil, err
			}
		}

		geometryCollection.Geometry = append(geometryCollection.Geometry, geometrySubType)
	}

	return &geometryCollection, nil
}

// Get a byte slice containing the EKWB representation of the geometry
func (gc GeometryCollection) GetEWKB(includeGeoType bool) bytes.Buffer {
	buf := new(bytes.Buffer)

	// Include geotype encoding if requested
	if includeGeoType {
		buf.WriteByte(byte(LittleEndian)) // Add Byte Order Marker
		geoTypeBytes := encodeGeoType(GeometryCollectionType, false, gc.Dimensions)
		buf.Write(geoTypeBytes)
	}

	// Encode the length of the geometry elements
	lenBytes := binary.LittleEndian.AppendUint32([]byte{}, uint32(len(gc.Geometry)))
	buf.Write(lenBytes)

	// Add encoded geometry
	for _, g := range gc.Geometry {
		gb := g.GetEWKB(true)
		gb.WriteTo(buf)
	}
	return *buf
}
