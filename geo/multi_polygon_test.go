package geo_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stephenirven/go-postgis/geo"
)

func TestMultiPolygon(t *testing.T) {

	dims := geo.XY

	multiPolygon := makeTestMultiPolygon(t, 10)

	ewkbBuffer := multiPolygon.GetEWKB(false)

	multiPolygon2, err := geo.MultiPolygonFromEWKB(&ewkbBuffer, dims)
	if err != nil {
		t.Error(err)
	}
	if ewkbBuffer.Len() > 0 {
		t.Error("recreated multipolygon did not consume all data")
	}

	if !cmp.Equal(multiPolygon, multiPolygon2) {
		t.Errorf("multipolygon %v was not equal to multipolygon %v", multiPolygon, multiPolygon2)
	}

}

func makeTestMultiPolygon(t *testing.T, length int) *geo.MultiPolygon {
	polygons := []geo.Polygon{}
	for i := 0; i < length; i++ {
		polygons = append(polygons, *makeTestPolygon(t))
	}

	multiPolygon, err := geo.NewMultiPolygon(polygons)
	if err != nil {
		t.Error(err)
	}
	return multiPolygon
}
