package geo

import (
	"bytes"
	"database/sql/driver"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strings"
)

// PostGIS License https://postgis.net/workshops/postgis-intro/license.html
//
// https://github.com/postgis/postgis/blob/master/doc/ZMSgeoms.txt
//

const (
	defaultSRID      = 4326         // default SRID for constructed geometry
	defaultByteOrder = LittleEndian // default byte order for constructed geometry (little endian)
)

type Dimensions byte

const (
	UNSET Dimensions = 0
	XY    Dimensions = 2
	XYZ   Dimensions = 3
	XYZM  Dimensions = 4
	XYM   Dimensions = 5
)

func (c Dimensions) String() string {
	switch c {
	case XY:
		return "XY"
	case XYZ:
		return "XYZ"
	case XYM:
		return "XYM"
	case XYZM:
		return "XYZM"
	default:
		return fmt.Sprintf("unknown: %" + string(byte(c)))
	}
}

type GeometrySubtype interface {
	GetEWKB(bool) bytes.Buffer
	GetDimensions() Dimensions
	GetGISGeometryType() GISGeometryType
	String() string
}

type GISGeometry struct {
	ByteOrder  ByteOrder
	GeoType    GISGeometryType
	Dimensions Dimensions
	SRIDFlag   bool
	SRID       uint32
	Geometry   GeometrySubtype
}

// Stringer interface
func (g GISGeometry) String() string {
	var sb strings.Builder
	sb.WriteString("(GISGeometry ")
	sb.WriteString("GeoType: " + g.GeoType.String() + ", ")
	sb.WriteString("Dimensions: " + g.Dimensions.String() + ", ")
	sb.WriteString(fmt.Sprintf("SRID: %v ", g.SRID) + ", ")

	sb.WriteString("Geometry: [")
	sb.WriteString(g.Geometry.String())

	sb.WriteString("])")
	return sb.String()
}

type GISGeometryType uint32

const (
	UNKNOWN                GISGeometryType = 0
	PointType              GISGeometryType = 1
	LineStringType         GISGeometryType = 2
	PolygonType            GISGeometryType = 3
	MultiPointType         GISGeometryType = 4
	MultiLineStringType    GISGeometryType = 5
	MultiPolygonType       GISGeometryType = 6
	GeometryCollectionType GISGeometryType = 7

	CircularStringType GISGeometryType = 8
	CompoundCurveType  GISGeometryType = 9
	CurvePolygonType   GISGeometryType = 10
	MultiCurveType     GISGeometryType = 11
	MultiSurfaceType   GISGeometryType = 12

	PolyHedralSurfaceType GISGeometryType = 15
	TINType               GISGeometryType = 16
	TriangleType          GISGeometryType = 17
)

func (g GISGeometryType) String() string {
	switch g {
	case PointType:
		return "PointType"
	case LineStringType:
		return "LineStringType"
	case PolygonType:
		return "PolygonType"
	case MultiPointType:
		return "MultiPointType"
	case MultiLineStringType:
		return "MultiLineStringType"
	case MultiPolygonType:
		return "MultiPolygonType"
	case GeometryCollectionType:
		return "GeometryCollectionType"
	case CircularStringType:
		return "CircularStringType"
	case CompoundCurveType:
		return "CompoundCurveType"
	case CurvePolygonType:
		return "CurvePolygonType"
	case MultiCurveType:
		return "MultiCurveType"
	case MultiSurfaceType:
		return "MultiSurfaceType"
	case PolyHedralSurfaceType:
		return "PolyhedralSurfaceType"
	case TINType:
		return "TINType"
	case TriangleType:
		return "TriangleType"
	case UNKNOWN:
		return "UNKNOWN"
	default:
		return "UNKNOWN"
	}
}

type ByteOrder byte

const (
	LittleEndian ByteOrder = 1
	BigEndian    ByteOrder = 2
)

func (b ByteOrder) String() string {
	switch b {
	case LittleEndian:
		return "LittleEndian"
	case BigEndian:
		return "BigEndian"
	default:
		return "UNKNOWN"
	}
}

// Flags for EWKB data - should be applied to most significant byte of geometry type
type Flag byte

const (
	wkbSRID Flag = 0x20 // 0010000 - SRID Presence flag
	wkbM    Flag = 0x40 // 0100000 - M-presence flag (wkbM)
	wkbZ    Flag = 0x80 // 1000000 - Z Coordinate presence flag
)

func encodeGeoType(geoType GISGeometryType, srid bool, dimensions Dimensions) []byte {

	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf[:], uint32(geoType))

	switch dimensions {
	case XYZM:
		buf[3] = buf[3] | byte(wkbZ|wkbM) // apply ZM flags
	case XYZ:
		buf[3] = buf[3] | byte(wkbZ) // apply Z flag
	case XYM:
		buf[3] = buf[3] | byte(wkbM) // apply M flag
	case XY:
	default:

	}

	if srid {
		buf[3] = buf[3] | byte(wkbSRID) // apply SRID presence flag
	}
	return buf
}

func NewGISGeometry(geo GeometrySubtype) GISGeometry {

	return GISGeometry{
		ByteOrder:  defaultByteOrder,
		Geometry:   geo,
		GeoType:    geo.GetGISGeometryType(),
		Dimensions: geo.GetDimensions(),
	}
}
func (gis *GISGeometry) SetSRID(srid uint32) {
	gis.SRID = srid
	gis.SRIDFlag = true
}

func (gis *GISGeometry) SetGeometry(geo GeometrySubtype) {
	switch gis.GeoType {
	case PointType:
		gis.Geometry = geo
		gis.Dimensions = geo.GetDimensions()
		gis.GeoType = PointType
	case LineStringType:
		gis.Geometry = geo
		gis.Dimensions = geo.GetDimensions()
		gis.GeoType = LineStringType
	case PolygonType:
		gis.Geometry = geo
		gis.Dimensions = geo.GetDimensions()
		gis.GeoType = PolygonType
	case MultiPointType:
		gis.Geometry = geo
		gis.Dimensions = geo.GetDimensions()
		gis.GeoType = MultiPointType
	case MultiLineStringType:
		gis.Geometry = geo
		gis.Dimensions = geo.GetDimensions()
		gis.GeoType = MultiLineStringType
	case MultiPolygonType:
		gis.Geometry = geo
		gis.Dimensions = geo.GetDimensions()
		gis.GeoType = MultiPolygonType
	case GeometryCollectionType:
		gis.Geometry = geo
		gis.Dimensions = geo.GetDimensions()
		gis.GeoType = GeometryCollectionType
	case CircularStringType:
		gis.Geometry = geo
		gis.Dimensions = geo.GetDimensions()
		gis.GeoType = CircularStringType
	case CompoundCurveType:
		gis.Geometry = geo
		gis.Dimensions = geo.GetDimensions()
		gis.GeoType = CompoundCurveType
	case CurvePolygonType:
		gis.Geometry = geo
		gis.Dimensions = geo.GetDimensions()
		gis.GeoType = CurvePolygonType
	case MultiCurveType:
		gis.Geometry = geo
		gis.Dimensions = geo.GetDimensions()
		gis.GeoType = MultiCurveType
	case MultiSurfaceType:
		gis.Geometry = geo
		gis.Dimensions = geo.GetDimensions()
		gis.GeoType = MultiSurfaceType
	case PolyHedralSurfaceType:
		gis.Geometry = geo
		gis.Dimensions = geo.GetDimensions()
		gis.GeoType = PolyHedralSurfaceType
	case TINType:
		gis.Geometry = geo
		gis.Dimensions = geo.GetDimensions()
		gis.GeoType = TINType
	case TriangleType:
		gis.Geometry = geo
		gis.Dimensions = geo.GetDimensions()
		gis.GeoType = TriangleType
	default:
	}

}

func decodeGeotype(bytes []byte) (geoType GISGeometryType, SRIDFlag bool, dimensions Dimensions) {

	// Use 16 bit integer to ignore remaining flag bits
	geoType = GISGeometryType(binary.LittleEndian.Uint16(bytes))

	// Get the flags from the most significant byte of geometry type
	flags := bytes[3]

	// SRID flag indicates if the geometry has embedded SRID
	if flags&byte(wkbSRID) == byte(wkbSRID) {
		SRIDFlag = true
	}

	// Flags indicate if geometry has Z or M dimensional data
	if flags&byte(wkbZ|wkbM) == byte(wkbZ|wkbM) {
		dims := XYZM
		dimensions = dims
	} else if flags&byte(wkbM) == byte(wkbM) {
		dims := XYM
		dimensions = dims
	} else if flags&byte(wkbZ) == byte(wkbZ) {
		dims := XYZ
		dimensions = dims
	} else {
		dims := XY
		dimensions = dims
	}

	return
}

// Used to generate a database/sql/driver.Value to write
func (g GISGeometry) Value() (driver.Value, error) {
	ewkb := []byte{} // byte array to hold EWKB data

	ewkb = append(ewkb, byte(defaultByteOrder)) // Byte Order Marker

	// Encode the geotype and flags
	geoTypeBytes := encodeGeoType(g.GeoType, g.SRIDFlag, g.Dimensions)
	ewkb = append(ewkb, geoTypeBytes...) // Append geotype & flags

	if g.SRIDFlag || g.SRID != 0 { // Append SRID if supplied
		ewkb = binary.LittleEndian.AppendUint32(ewkb, uint32(g.SRID))
	}

	// Append the EWKB data for the geometry
	geometry := g.Geometry.GetEWKB(false)
	ewkb = append(ewkb, geometry.Bytes()...)

	// Encode and return HEX EWKB data
	hexewkb := make([]byte, hex.EncodedLen(len(ewkb)))
	hex.Encode(hexewkb, ewkb)

	return hexewkb, nil

}

// Used to map GISGeometry values into structs when read by the database driver
func (g *GISGeometry) Scan(value interface{}) error {

	// Format from PostGIS is Hex encoded EWKB - Extended Well Known Binary
	// Should start as a byte array
	hexewkb, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("scan expected []byte, got %T (%v)", value, value)
	}

	// Decode into EWKB byte array
	ewkb := make([]byte, hex.DecodedLen(len(hexewkb)))
	_, err := hex.Decode(ewkb, hexewkb)
	if err != nil {
		return err
	}

	if len(ewkb) < 9 {
		return fmt.Errorf("ewkb must be at least 9 bytes to contain byte order, type, and srid")
	}

	buffer := bytes.NewBuffer(ewkb)

	//The first byte indicates the byte order. 00 for big endian, or 01 for little endian.
	bom, err := buffer.ReadByte()
	if err != nil {
		return fmt.Errorf("unable to read byte order marker")
	}
	g.ByteOrder = ByteOrder(bom)

	if g.ByteOrder == LittleEndian {

		g.GeoType, g.SRIDFlag, g.Dimensions = decodeGeotype(buffer.Next(4))

		// Get the SRID if present
		if g.SRIDFlag {
			g.SRID = binary.LittleEndian.Uint32(buffer.Next(4))
		}

		// Get the geometry from the remaining data
		var geometry GeometrySubtype

		switch g.GeoType {
		case PointType:
			geometry, err = PointFromEWKB(buffer, g.Dimensions)
			if err != nil {
				return fmt.Errorf(err.Error())
			}

		case LineStringType:
			geometry, err = LineStringFromEWKB(buffer, g.Dimensions)
			if err != nil {
				return fmt.Errorf(err.Error())
			}

		case PolygonType:
			geometry, err = PolygonFromEWKB(buffer, g.Dimensions)
			if err != nil {
				return err
			}

		case MultiPointType:
			geometry, err = MultiPointFromEWKB(buffer, g.Dimensions)
			if err != nil {
				return err
			}

		case MultiLineStringType:
			geometry, err = MultiLineStringFromEWKB(buffer, g.Dimensions)
			if err != nil {
				return err
			}

		case MultiPolygonType:
			geometry, err = MultiPolygonFromEWKB(buffer, g.Dimensions)
			if err != nil {
				return err
			}

		case GeometryCollectionType:
			geometry, err = GeometryCollectionFromEWKB(buffer, g.Dimensions)
			if err != nil {
				return err
			}

		case CircularStringType:
			geometry, err = CircularStringFromEWKB(buffer, g.Dimensions)
			if err != nil {
				return err
			}

		case CompoundCurveType:
			geometry, err = CompoundCurveFromEWKB(buffer, g.Dimensions)
			if err != nil {
				return err
			}

		case CurvePolygonType:
			geometry, err = CurvePolygonFromEWKB(buffer, g.Dimensions)
			if err != nil {
				return err
			}

		case MultiCurveType:
			geometry, err = MultiCurveFromEWKB(buffer, g.Dimensions)
			if err != nil {
				return err
			}

		case MultiSurfaceType:
			geometry, err = MultiSurfaceFromEWKB(buffer, g.Dimensions)
			if err != nil {
				return err
			}

		case PolyHedralSurfaceType:
			geometry, err = PolyhedralSurfaceFromEWKB(buffer, g.Dimensions)
			if err != nil {
				return err
			}

		case TINType:
			geometry, err = TINFromEWKB(buffer, g.Dimensions)
			if err != nil {
				return err
			}

		case TriangleType:
			geometry, err = TriangleFromEWKB(buffer, g.Dimensions)
			if err != nil {
				return err
			}

		default:
			return fmt.Errorf("unknown geometry type: %v", g.GeoType)

		}

		g.Geometry = geometry

	} else {
		return fmt.Errorf("big endian is currently unsupported")

	}

	return nil

}
