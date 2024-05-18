package geo

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
)

/* https://postgis.net/docs/using_postgis_dbmanagement.html#MultiSurface
A MultiSurface is a collection of surfaces, which can be (linear) Polygons
or CurvePolygons.
*/

type MultiSurface struct {
	Geometry   []GeometrySubtype
	Dimensions Dimensions
}

func (g MultiSurface) GetGISGeometryType() GISGeometryType {
	return MultiSurfaceType
}

// Stringer interface
func (ms MultiSurface) String() string {
	var sb strings.Builder
	sb.WriteString("(MultiSurface ")
	sb.WriteString(ms.Dimensions.String())
	sb.WriteString(" [")
	for _, g := range ms.Geometry {
		sb.WriteString(g.String())
		sb.WriteString(" ")
	}

	sb.WriteString("])")
	return sb.String()
}

// Get the dimensions of the geometry
func (ms MultiSurface) GetDimensions() Dimensions {
	return ms.Dimensions
}

// Create a MultiSurface from a slice of Geometry of the same dimensions.
// [Must contain ]
// Length must be at least 1
func NewMultiSurface() (*MultiSurface, error) {
	m := MultiSurface{}
	return &m, nil
}

func (ms *MultiSurface) AddPolygon(p *Polygon) error {
	if ms.Dimensions != UNSET && ms.Dimensions != p.Dimensions {
		return fmt.Errorf("dimensions %v do not match %v", p.Dimensions, ms.Dimensions)
	}
	if ms.Dimensions == UNSET {
		ms.Dimensions = p.Dimensions
	}
	ms.Geometry = append(ms.Geometry, p)
	return nil
}

func (ms *MultiSurface) AddCurvePolygon(cp *CurvePolygon) error {
	if ms.Dimensions != UNSET && ms.Dimensions != cp.Dimensions {
		return fmt.Errorf("dimensions %v do not match %v", cp.Dimensions, ms.Dimensions)
	}
	if ms.Dimensions == UNSET {
		ms.Dimensions = cp.Dimensions
	}
	ms.Geometry = append(ms.Geometry, cp)
	return nil
}

// Create a new MultiSurface from input byte buffer in EWKB format and dimensions.
func MultiSurfaceFromEWKB(buffer *bytes.Buffer, dimensions Dimensions) (*MultiSurface, error) {

	multiSurface := MultiSurface{}
	multiSurface.Dimensions = dimensions
	if buffer.Len() < 4 {
		return nil, fmt.Errorf("byte array for multisurface must be at least length 4")
	}
	count := binary.LittleEndian.Uint32(buffer.Next(4))

	for i := 0; i < int(count); i++ {
		buffer.Next(1) // move past BOM
		geoTypeBytes := buffer.Next(4)
		geoType := GISGeometryType(binary.LittleEndian.Uint16(geoTypeBytes)) // 4 bytes for geometry type
		flags := geoTypeBytes[3]                                             // flags are stored on the 3rd byte

		if flags&byte(wkbSRID) == byte(wkbSRID) {
			//SRID := binary.LittleEndian.Uint16(data[:4])
			buffer.Next(4) // move past SRID
		}

		var geometry GeometrySubtype
		var err error

		switch geoType {
		case PolygonType:
			geometry, err = PolygonFromEWKB(buffer, dimensions)
			if err != nil {
				return nil, err
			}
		case CurvePolygonType:
			geometry, err = CurvePolygonFromEWKB(buffer, dimensions)
			if err != nil {
				return nil, err
			}
		default:
			fmt.Printf("multisurface must only contain polygons / curvepolygon: %v", geoType)

		}
		multiSurface.Geometry = append(multiSurface.Geometry, geometry)

	}

	return &multiSurface, nil
}

// Get a byte slice containing the EKWB representation of the geometry
func (ms MultiSurface) GetEWKB(includeGeoType bool) bytes.Buffer {
	buf := new(bytes.Buffer)

	// Include geotype encoding if requested
	if includeGeoType {
		buf.WriteByte(byte(LittleEndian)) // Add Byte Order Marker
		geoTypeBytes := encodeGeoType(MultiSurfaceType, false, ms.Dimensions)
		buf.Write(geoTypeBytes)
	}

	lenBytes := binary.LittleEndian.AppendUint32([]byte{}, uint32(len(ms.Geometry)))
	buf.Write(lenBytes)

	for _, p := range ms.Geometry {

		buf.WriteByte(byte(LittleEndian))
		var geoTypeBytes []byte
		switch shape := p.(type) {
		case *CurvePolygon:
			geoTypeBytes = encodeGeoType(CurvePolygonType, false, ms.Dimensions)
		case *Polygon:
			geoTypeBytes = encodeGeoType(PolygonType, false, ms.Dimensions)
		default:
			fmt.Printf("multisurface must only contain polygon/curvepolygon: %T", shape)
		}
		buf.Write(geoTypeBytes)
		pb := p.GetEWKB(false)
		pb.WriteTo(buf)
	}
	return *buf
}
