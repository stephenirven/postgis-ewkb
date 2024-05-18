package geo_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stephenirven/go-postgis/geo"
)

func TestTIN(t *testing.T) {

	dims := geo.XY

	tin := makeTestTIN(t)

	ewkbBuffer := tin.GetEWKB(false)

	tin2, err := geo.TINFromEWKB(&ewkbBuffer, dims)
	if err != nil {
		t.Error(err)
	}
	if ewkbBuffer.Len() > 0 {
		t.Error("recreated tin did not consume all data")
	}

	if !cmp.Equal(tin, tin2) {
		t.Errorf("tin %v was not equal to tin %v", tin, tin2)
	}

}

func makeTestTIN(t *testing.T) *geo.TIN {
	triangles := []geo.Triangle{}
	for i := 0; i < 10; i++ {
		triangles = append(triangles, *makeTestTriangle(t))
	}

	tin, err := geo.NewTIN(triangles)
	if err != nil {
		t.Error(err)
	}
	return tin
}
