package geo_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stephenirven/go-postgis/geo"
)

func TestCurvePolygon(t *testing.T) {

	dims := geo.XY
	cp := makeTestCurvePolygon(t)

	ewkbBuffer := cp.GetEWKB(false)

	cp2, err := geo.CurvePolygonFromEWKB(&ewkbBuffer, dims)
	if err != nil {
		t.Error(err)
	}
	if ewkbBuffer.Len() > 0 {
		t.Error("recreated curvepolygon did not consume all data")
	}

	if !cmp.Equal(cp, cp2) {
		t.Errorf("curvepolygon %v was not equal to curvepolygon %v", cp, cp2)
	}

}

func makeTestCurvePolygon(t *testing.T) *geo.CurvePolygon {

	dims := geo.XY

	p1, err := geo.NewPoint([]float64{1, 0}, dims)
	if err != nil {
		t.Error(err)
	}
	p2, err := geo.NewPoint([]float64{1, -1}, dims)
	if err != nil {
		t.Error(err)
	}
	p3, err := geo.NewPoint([]float64{1, -3}, dims)
	if err != nil {
		t.Error(err)
	}
	p4, err := geo.NewPoint([]float64{1, -5}, dims)
	if err != nil {
		t.Error(err)
	}
	p5, err := geo.NewPoint([]float64{1, -7}, dims)
	if err != nil {
		t.Error(err)
	}
	ps1 := []geo.Point{*p1, *p2, *p3, *p4, *p5}

	lineString, err := geo.NewLineString(ps1)
	if err != nil {
		t.Error(err)
	}

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

	circularString, err := geo.NewCircularString(csp)
	if err != nil {
		t.Error(err)
	}

	CompoundCurve, err := geo.NewCompoundCurve()
	if err != nil {
		t.Error(err)
	}
	err = CompoundCurve.AddCircularString(circularString)
	if err != nil {
		t.Error(err)
	}
	err = CompoundCurve.AddLineString(lineString)
	if err != nil {
		t.Error(err)
	}

	csp1, err := geo.NewPoint([]float64{0, 0}, dims)
	if err != nil {
		t.Error(err)
	}
	csp2, err := geo.NewPoint([]float64{1, 1}, dims)
	if err != nil {
		t.Error(err)
	}
	csp3, err := geo.NewPoint([]float64{1, 0}, dims)
	if err != nil {
		t.Error(err)
	}

	circularString2, err := geo.NewCircularString([]geo.Point{*csp1, *csp2, *csp3})
	if err != nil {
		t.Error(err)
	}

	curvePolygon := geo.NewCurvePolygon()
	err = curvePolygon.AddCircularString(circularString2)
	if err != nil {
		t.Error(err)
	}
	err = curvePolygon.AddCompoundCurve(CompoundCurve)
	if err != nil {
		t.Error(err)
	}
	err = curvePolygon.AddLineString(lineString)
	if err != nil {
		t.Error(err)
	}
	return curvePolygon

}
