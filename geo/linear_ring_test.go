package geo_test

import (
	"math/rand"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stephenirven/go-postgis/geo"
)

func TestLinearRing(t *testing.T) {

	dims := geo.XY

	linearRing := makeTestLinearXYRing(t, 10)

	ewkbBuffer := linearRing.GetEWKB(false)

	linearRing2, err := geo.LinearRingFromEWKB(&ewkbBuffer, dims)
	if err != nil {
		t.Error(err)
	}
	if ewkbBuffer.Len() > 0 {
		t.Error("recreated linearring did not consume all data")
	}

	if !cmp.Equal(linearRing, linearRing2) {
		t.Errorf("linearring %v was not equal to linearring %v", linearRing, linearRing2)
	}

}

func makeTestLinearXYRing(t *testing.T, offset float64) *geo.LinearRing {

	x := rand.Float64() * 45
	y := rand.Float64() * 45
	dims := geo.XY

	p1, err := geo.NewPoint([]float64{x - offset, y - offset}, dims)
	if err != nil {
		t.Error(err)
	}
	p2, err := geo.NewPoint([]float64{x - offset, y + offset}, dims)
	if err != nil {
		t.Error(err)
	}
	p3, err := geo.NewPoint([]float64{x + offset, y + offset}, dims)
	if err != nil {
		t.Error(err)
	}
	p4, err := geo.NewPoint([]float64{x + 10, y - 10}, dims)
	if err != nil {
		t.Error(err)
	}
	p5, err := geo.NewPoint([]float64{x - offset, y - offset}, dims)
	if err != nil {
		t.Error(err)
	}
	ps1 := []geo.Point{*p1, *p2, *p3, *p4, *p5}
	if err != nil {
		t.Error(err)
	}

	linearRing, err := geo.NewLinearRing(ps1)
	if err != nil {
		t.Error(err)
	}
	return linearRing

}
