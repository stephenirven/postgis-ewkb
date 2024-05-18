package geo_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stephenirven/go-postgis/geo"
)

func TestTriangle(t *testing.T) {

	dims := geo.XY

	triangle := makeTestTriangle(t)

	ewkbBuffer := triangle.GetEWKB(false)

	triangle2, err := geo.TriangleFromEWKB(&ewkbBuffer, dims)

	if err != nil {
		t.Error(err)
	}
	if ewkbBuffer.Len() > 0 {
		t.Error("recreated triangle did not consume all data")
	}

	if !cmp.Equal(triangle, triangle2) {
		t.Errorf("triangle %v was not equal to triangle %v", triangle, triangle2)
	}

}

func makeTestTriangle(t *testing.T) *geo.Triangle {

	dims := geo.XY

	var points [4]geo.Point
	points[0] = *makeTestPoint(t, dims)
	points[1] = *makeTestPoint(t, dims)
	points[2] = *makeTestPoint(t, dims)
	points[3] = points[0]

	triangle, err := geo.NewTriangle(points)
	if err != nil {
		t.Error(err)
	}

	return triangle
}
