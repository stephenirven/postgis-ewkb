package geo_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stephenirven/go-postgis/geo"
)

func TestLineString(t *testing.T) {

	dims := geo.XY

	l := makeTestLineString(t, 10, dims)

	ewkbBuffer := l.GetEWKB(false)

	m, err := geo.LineStringFromEWKB(&ewkbBuffer, dims)
	if err != nil {
		t.Error(err)
	}
	if ewkbBuffer.Len() > 0 {
		t.Error("recreated point did not consume all data")
	}

	if !cmp.Equal(l, m) {
		t.Errorf("linestring %v was not equal to linestring %v", l, m)
	}

}

func makeTestLineString(t *testing.T, length int, dimensions geo.Dimensions) *geo.LineString {

	points := []geo.Point{}
	for i := 0; i < length; i++ {
		points = append(points, *makeTestPoint(t, dimensions))

	}

	lineString, err := geo.NewLineString(points)

	if err != nil {
		t.Error(err)
	}
	return lineString

}
