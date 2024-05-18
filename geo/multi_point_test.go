package geo_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stephenirven/go-postgis/geo"
)

func TestMultiPoint(t *testing.T) {

	dims := geo.XY

	mp := makeTestMultiPoint(t, 10, dims)

	ewkbBuffer := mp.GetEWKB(false)

	mp2, err := geo.MultiPointFromEWKB(&ewkbBuffer, dims)
	if err != nil {
		t.Error(err)
	}
	if ewkbBuffer.Len() > 0 {
		t.Errorf("recreated multipoint did not consume all data: %v", ewkbBuffer.Bytes())
	}

	if !cmp.Equal(mp, mp2) {
		t.Errorf("multipoint %v was not equal to multipoint %v", mp, mp2)
	}

}

func makeTestMultiPoint(t *testing.T, length int, dimensions geo.Dimensions) *geo.MultiPoint {

	points := []geo.Point{}
	for i := 0; i < length; i++ {
		points = append(points, *makeTestPoint(t, dimensions))
	}

	multiPoint, err := geo.NewMultiPoint(points)
	if err != nil {
		t.Error(err)
	}

	return multiPoint
}
