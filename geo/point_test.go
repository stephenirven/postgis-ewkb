package geo_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stephenirven/go-postgis/geo"
)

func TestPoint(t *testing.T) {

	dims := geo.XY

	p1 := makeTestPoint(t, dims)

	ewkbBuffer := p1.GetEWKB(false)

	p1b, err := geo.PointFromEWKB(&ewkbBuffer, dims)
	if err != nil {
		t.Error(err)
	}
	if ewkbBuffer.Len() > 0 {
		t.Errorf("recreated point did not consume all data: %v", ewkbBuffer.Bytes())
	}

	if !cmp.Equal(p1, p1b) {
		t.Errorf("Point %v was not equal to Point %v", p1, p1b)
	}

	dims = geo.XYZ
	p2 := makeTestPoint(t, dims)

	ewkbBuffer = p2.GetEWKB(false)

	p2b, err := geo.PointFromEWKB(&ewkbBuffer, dims)
	if err != nil {
		t.Error(err)
	}
	if ewkbBuffer.Len() > 0 {
		t.Error("recreated point did not consume all data")
	}

	if !cmp.Equal(p2, p2b) {
		t.Errorf("Point %v was not equal to Point %v", p2, p2b)
	}

	dims = geo.XYM
	p3 := makeTestPoint(t, dims)

	ewkbBuffer = p3.GetEWKB(false)

	p3b, err := geo.PointFromEWKB(&ewkbBuffer, dims)
	if err != nil {
		t.Error(err)
	}
	if ewkbBuffer.Len() > 0 {
		t.Error("recreated point did not consume all data")
	}

	if !cmp.Equal(p3, p3b) {
		t.Errorf("Point %v was not equal to Point %v", p3, p3b)
	}

	dims = geo.XYZM
	p4 := makeTestPoint(t, dims)

	ewkbBuffer = p4.GetEWKB(false)

	p4b, err := geo.PointFromEWKB(&ewkbBuffer, dims)
	if err != nil {
		t.Error(err)
	}
	if ewkbBuffer.Len() > 0 {
		t.Error("recreated point did not consume all data")
	}

	if !cmp.Equal(p4, p4b) {
		t.Errorf("Point %v was not equal to Point %v", p4, p4b)
	}

}

func makeTestPoint(t *testing.T, dimensions geo.Dimensions) *geo.Point {

	coords := []float64{}
	for i := 0; i < len(fmt.Sprintf("%v", dimensions)); i++ {
		coords = append(coords, rand.Float64()*90)
	}
	point, err := geo.NewPoint(coords, dimensions)
	if err != nil {
		t.Error(err)
	}

	return point

}
