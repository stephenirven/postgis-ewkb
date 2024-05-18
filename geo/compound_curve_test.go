package geo_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stephenirven/go-postgis/geo"
)

func TestCompoundCurve(t *testing.T) {

	dims := geo.XY

	compoundCurve := makeTestCompoundCurve(t)

	ewkbBuffer := compoundCurve.GetEWKB(false)

	compoundCurve2, err := geo.CompoundCurveFromEWKB(&ewkbBuffer, dims)
	if err != nil {
		t.Error(err)
	}
	if ewkbBuffer.Len() > 0 {
		t.Error("recreated compoundcurve did not consume all data")
	}

	if !cmp.Equal(compoundCurve, compoundCurve2) {
		t.Errorf("compoundcurve %v was not equal to compoundcurve %v", compoundCurve, compoundCurve2)
	}

}

func makeTestCompoundCurve(t *testing.T) *geo.CompoundCurve {

	dimensions := geo.XY
	p1, err := geo.NewPoint([]float64{1, 0}, dimensions)
	if err != nil {
		t.Error(err)
	}
	p2, err := geo.NewPoint([]float64{1, -1}, dimensions)
	if err != nil {
		t.Error(err)
	}
	p3, err := geo.NewPoint([]float64{1, -3}, dimensions)
	if err != nil {
		t.Error(err)
	}
	p4, err := geo.NewPoint([]float64{1, -5}, dimensions)
	if err != nil {
		t.Error(err)
	}
	p5, err := geo.NewPoint([]float64{1, -7}, dimensions)
	if err != nil {
		t.Error(err)
	}
	ps1 := []geo.Point{*p1, *p2, *p3, *p4, *p5}

	l1, err := geo.NewLineString(ps1)
	if err != nil {
		t.Error(err)
	}

	csp01, err := geo.NewPoint([]float64{0, 0}, dimensions)
	if err != nil {
		t.Error(err)
	}
	csp02, err := geo.NewPoint([]float64{1, 1}, dimensions)
	if err != nil {
		t.Error(err)
	}
	csp03, err := geo.NewPoint([]float64{1, 0}, dimensions)
	if err != nil {
		t.Error(err)
	}

	csp := []geo.Point{*csp01, *csp02, *csp03}

	cs, err := geo.NewCircularString(csp)
	if err != nil {
		t.Error(err)
	}

	compoundCurve, err := geo.NewCompoundCurve()
	if err != nil {
		t.Error(err)
	}
	err = compoundCurve.AddCircularString(cs)
	if err != nil {
		t.Error(err)
	}
	err = compoundCurve.AddLineString(l1)
	if err != nil {
		t.Error(err)
	}

	return compoundCurve

}
