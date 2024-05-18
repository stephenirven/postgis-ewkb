package geo_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stephenirven/go-postgis/geo"
)

func TestPolyhedralSurface(t *testing.T) {

	dims := geo.XY

	polyhedralSurface := makeTestPolyhedralSurface(t)

	ewkbBuffer := polyhedralSurface.GetEWKB(false)

	polyhedralSurface2, err := geo.PolyhedralSurfaceFromEWKB(&ewkbBuffer, dims)
	if err != nil {
		t.Error(err)
	}
	if ewkbBuffer.Len() > 0 {
		t.Error("recreated polyhedralsurface did not consume all data")
	}

	if !cmp.Equal(polyhedralSurface, polyhedralSurface2) {
		t.Errorf("polyhedralsurface %v was not equal to polyhedralsurface %v", polyhedralSurface, polyhedralSurface2)
	}
}

func makeTestPolyhedralSurface(t *testing.T) *geo.PolyHedralSurface {

	polygons := []geo.Polygon{}
	for i := 0; i < 10; i++ {
		polygons = append(polygons, *makeTestPolygon(t))
	}
	polyhedralSurface, err := geo.NewPolyhedralSurface(polygons)
	if err != nil {
		t.Error(err)
	}
	return polyhedralSurface
}
