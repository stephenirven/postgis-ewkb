package geo_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stephenirven/go-postgis/geo"
)

func TestGeometryCollection(t *testing.T) {

	geometryCollection := makeTestGeometryCollection(t)

	dims := geo.XY

	ewkbBuffer := geometryCollection.GetEWKB(false)

	geometryCollection2, err := geo.GeometryCollectionFromEWKB(&ewkbBuffer, dims)
	if err != nil {
		t.Error(err)
	}
	if ewkbBuffer.Len() > 0 {
		t.Error("recreated geometrycollection did not consume all data")
	}

	if !cmp.Equal(geometryCollection, geometryCollection2) {
		t.Errorf("geometrycollection %v was not equal to geometrycollection %v", geometryCollection, geometryCollection2)
	}

}

func makeTestGeometryCollection(t *testing.T) *geo.GeometryCollection {
	dims := geo.XY

	geometry := []geo.GeometrySubtype{}

	geometry = append(geometry, makeTestPoint(t, dims))
	geometry = append(geometry, makeTestLineString(t, 10, dims))
	geometry = append(geometry, makeTestPolygon(t))
	geometry = append(geometry, makeTestMultiPoint(t, 10, dims))
	geometry = append(geometry, makeTestMultiLineString(t, 10, dims))
	geometry = append(geometry, makeTestMultiPolygon(t, 10))
	geometry = append(geometry, makeTestCurvePolygon(t))

	// Can Geometry collections contain LinearRings?
	//geometry = append(geometry, makeTestLinearXYRing(t, *makeTestPoint(t, dims),10))

	geometry = append(geometry, makeTestCircularString(t))
	geometry = append(geometry, makeTestCompoundCurve(t))
	geometry = append(geometry, makeTestMultiCurve(t))
	geometry = append(geometry, makeTestMultiSurface(t))
	geometry = append(geometry, makeTestPolyhedralSurface(t))
	geometry = append(geometry, makeTestTIN(t))
	geometry = append(geometry, makeTestTriangle(t))

	geometryCollection, err := geo.NewGeometryCollection(geometry)
	if err != nil {
		t.Error(err)
	}
	return geometryCollection

}
