package geo_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stephenirven/go-postgis/geo"
)

func TestCircularString(t *testing.T) {

	dims := geo.XY

	cs := makeTestCircularString(t)

	ewkbBuffer := cs.GetEWKB(false)

	cs2, err := geo.CircularStringFromEWKB(&ewkbBuffer, dims)
	if err != nil {
		t.Error(err)
	}
	if ewkbBuffer.Len() > 0 {
		t.Error("recreated circularstring did not consume all data")
	}

	if !cmp.Equal(cs, cs2) {
		t.Errorf("circularstring %v was not equal to circularstring %v", cs, cs2)
	}

}

func makeTestCircularString(t *testing.T) *geo.CircularString {
	dims := geo.XY
	csp01, err := geo.NewPoint([]float64{0, 0}, dims)
	if err != nil {
		t.Error(err)
	}
	csp02, err := geo.NewPoint([]float64{1, 1}, dims)
	if err != nil {
		t.Error(err)
	}
	csp03, err := geo.NewPoint([]float64{1, 0}, dims)
	if err != nil {
		t.Error(err)
	}

	csp := []geo.Point{*csp01, *csp02, *csp03}

	cs, err := geo.NewCircularString(csp)
	if err != nil {
		t.Error(err)
	}
	return cs
}
