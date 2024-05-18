package geo_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stephenirven/go-postgis/geo"
)

func TestMultiLineString(t *testing.T) {

	dims := geo.XY

	multiLineString := makeTestMultiLineString(t, 10, dims)

	ewkbBuffer := multiLineString.GetEWKB(false)

	multiLineString2, err := geo.MultiLineStringFromEWKB(&ewkbBuffer, dims)
	if err != nil {
		t.Error(err)
	}
	if ewkbBuffer.Len() > 0 {
		t.Error("recreated multilinestring did not consume all data")
	}

	if !cmp.Equal(multiLineString, multiLineString2) {
		t.Errorf("multilinestring %v was not equal to multilinestring %v", multiLineString, multiLineString2)
	}

}

func makeTestMultiLineString(t *testing.T, length int, dimensions geo.Dimensions) *geo.MultiLineString {
	lineStrings := []geo.LineString{}
	for i := 0; i < length; i++ {
		lineStrings = append(lineStrings, *makeTestLineString(t, 8, dimensions))
	}

	multiLineString, err := geo.NewMultiLineString(lineStrings)
	if err != nil {
		t.Error(err)
	}
	return multiLineString
}
