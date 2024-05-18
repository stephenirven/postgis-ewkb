package geo_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stephenirven/go-postgis/geo"
)

func TestMultiCurve(t *testing.T) {

	dims := geo.XY

	multiCurve := makeTestMultiCurve(t)

	ewkbBuffer := multiCurve.GetEWKB(false)

	multiCurve2, err := geo.MultiCurveFromEWKB(&ewkbBuffer, dims)
	if err != nil {
		t.Error(err)
	}
	if ewkbBuffer.Len() > 0 {
		t.Error("recreated multicurve did not consume all data")
	}

	if !cmp.Equal(multiCurve, multiCurve2) {
		t.Errorf("multicurve %v was not equal to multicurve %v", multiCurve, multiCurve2)
	}
}

func makeTestMultiCurve(t *testing.T) *geo.MultiCurve {

	multiCurve, err := geo.NewMultiCurve()
	if err != nil {
		t.Error(err)
	}

	circularString := makeTestCircularString(t)
	err = multiCurve.AddCircularString(circularString)
	if err != nil {
		t.Error(err)
	}

	compoundCurve := makeTestCompoundCurve(t)
	err = multiCurve.AddCompoundCurve(compoundCurve)
	if err != nil {
		t.Error(err)
	}

	lineString := makeTestLineString(t, 10, geo.XY)
	err = multiCurve.AddLineString(lineString)
	if err != nil {
		t.Error(err)
	}

	return multiCurve
}
