package geo_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stephenirven/go-postgis/geo"
)

func TestPolygon(t *testing.T) {

	dims := geo.XY

	poly := makeTestPolygon(t)

	ewkbBuffer := poly.GetEWKB(false)

	poly2, err := geo.PolygonFromEWKB(&ewkbBuffer, dims)
	if err != nil {
		t.Error(err)
	}
	if ewkbBuffer.Len() > 0 {
		t.Error("recreated polygon did not consume all data")
	}

	if !cmp.Equal(poly, poly2) {
		t.Errorf("polygon %v was not equal to polygon %v", poly, poly2)
	}

}

func makeTestPolygon(t *testing.T) *geo.Polygon {

	linearRing1 := makeTestLinearXYRing(t, 10)

	linearRing2 := makeTestLinearXYRing(t, 5)

	polygon, err := geo.NewPolygon([]geo.LinearRing{*linearRing1, *linearRing2})
	if err != nil {
		t.Error(err)
	}
	return polygon
}
