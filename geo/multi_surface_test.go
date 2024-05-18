package geo_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stephenirven/go-postgis/geo"
)

func TestMultiSurface(t *testing.T) {

	dims := geo.XY

	multiSurface := makeTestMultiSurface(t)

	ewkbBuffer := multiSurface.GetEWKB(false)

	multiSurface2, err := geo.MultiSurfaceFromEWKB(&ewkbBuffer, dims)
	if err != nil {
		t.Error(err)
	}
	if ewkbBuffer.Len() > 0 {
		t.Error("recreated multisurface did not consume all data")
	}

	if !cmp.Equal(multiSurface, multiSurface2) {
		t.Errorf("multisurface %v was not equal to multisurface %v", multiSurface, multiSurface2)
	}

}

func makeTestMultiSurface(t *testing.T) *geo.MultiSurface {

	multiSurface, err := geo.NewMultiSurface()
	if err != nil {
		t.Error(err)
	}
	for i := 0; i < 10; i++ {
		err = multiSurface.AddPolygon(makeTestPolygon(t))
		if err != nil {
			t.Error(err)
		}
	}
	for i := 0; i < 10; i++ {
		err = multiSurface.AddCurvePolygon(makeTestCurvePolygon(t))
		if err != nil {
			t.Error(err)
		}
	}
	return multiSurface
}
