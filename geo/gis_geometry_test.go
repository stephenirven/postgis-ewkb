package geo_test

import (
	"context"
	"database/sql"
	"log"
	"testing"

	"github.com/google/go-cmp/cmp"
	_ "github.com/lib/pq"
	db "github.com/stephenirven/go-postgis/db/sqlc"
	"github.com/stephenirven/go-postgis/geo"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://postgres:secret@localhost:5432/late?sslmode=disable"
)

var testQueries *db.Queries

func init() {
	testDB, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("failed to connect to database", err)
	}
	testQueries = db.New(testDB)
}
func TestGISPointGeometry(t *testing.T) {

	dims := geo.XY

	point := makeTestPoint(t, dims)
	gisGeometry := geo.NewGISGeometry(point)

	gisGeometry.SetSRID(4326)

	gisData, err := testQueries.CreateGISData(context.Background(), gisGeometry)

	if err != nil {
		t.Error(err)
	}

	geo2, err := testQueries.GetGISData(context.Background(), gisData.ID)
	if err != nil {
		t.Error(err)
	}

	if !cmp.Equal(gisGeometry, geo2) {
		t.Errorf("point %v was not equal to %v", &geo2, gisGeometry)
	}

}

func TestGISLineStringGeometry(t *testing.T) {

	dims := geo.XY

	lineString := makeTestLineString(t, 10, dims)
	gisGeometry := geo.NewGISGeometry(lineString)

	gisGeometry.SetSRID(4326)

	gisData, err := testQueries.CreateGISData(context.Background(), gisGeometry)

	if err != nil {
		t.Error(err)
	}

	geo2, err := testQueries.GetGISData(context.Background(), gisData.ID)
	if err != nil {
		t.Error(err)
	}

	if !cmp.Equal(gisGeometry, geo2) {
		t.Errorf("linestring %v was not equal to %v", &geo2, gisGeometry)
	}
}

func TestGISPolygonGeometry(t *testing.T) {

	polygon := makeTestPolygon(t)
	gisGeometry := geo.NewGISGeometry(polygon)

	gisGeometry.SetSRID(4326)

	gisData, err := testQueries.CreateGISData(context.Background(), gisGeometry)

	if err != nil {
		t.Error(err)
	}

	geo2, err := testQueries.GetGISData(context.Background(), gisData.ID)
	if err != nil {
		t.Error(err)
	}

	if !cmp.Equal(gisGeometry, geo2) {
		t.Errorf("polygon %v was not equal to %v", &geo2, gisGeometry)
	}

}

func TestGISMultiPointGeometry(t *testing.T) {

	dims := geo.XY

	multiPoint := makeTestMultiPoint(t, 1, dims)
	gisGeometry := geo.NewGISGeometry(multiPoint)

	gisGeometry.SetSRID(4326)

	gisData, err := testQueries.CreateGISData(context.Background(), gisGeometry)

	if err != nil {
		t.Error(err)
	}

	geo2, err := testQueries.GetGISData(context.Background(), gisData.ID)
	if err != nil {
		t.Error(err)
	}

	if !cmp.Equal(gisGeometry, geo2) {
		t.Errorf("multipoint %v was not equal to %v", &geo2, gisGeometry)
	}

}

func TestGISMultiLineStringGeometry(t *testing.T) {

	dims := geo.XY

	multiLineString := makeTestMultiLineString(t, 10, dims)
	gisGeometry := geo.NewGISGeometry(multiLineString)

	gisGeometry.SetSRID(4326)

	gisData, err := testQueries.CreateGISData(context.Background(), gisGeometry)

	if err != nil {
		t.Error(err)
	}

	geo2, err := testQueries.GetGISData(context.Background(), gisData.ID)
	if err != nil {
		t.Error(err)
	}

	if !cmp.Equal(gisGeometry, geo2) {
		t.Errorf("multilinestring %v was not equal to %v", &geo2, gisGeometry)
	}

}

func TestGISMultiPolygonGeometry(t *testing.T) {

	multiPolygon := makeTestMultiPolygon(t, 10)
	gisGeometry := geo.NewGISGeometry(multiPolygon)

	gisGeometry.SetSRID(4326)

	gisData, err := testQueries.CreateGISData(context.Background(), gisGeometry)

	if err != nil {
		t.Error(err)
	}

	geo2, err := testQueries.GetGISData(context.Background(), gisData.ID)
	if err != nil {
		t.Error(err)
	}

	if !cmp.Equal(gisGeometry, geo2) {
		t.Errorf("multipolygon %v was not equal to %v", &geo2, gisGeometry)
	}

}

func TestGISCurvePolygonGeometry(t *testing.T) {

	curvepolygon := makeTestCurvePolygon(t)
	gisGeometry := geo.NewGISGeometry(curvepolygon)

	gisGeometry.SetSRID(4326)

	gisData, err := testQueries.CreateGISData(context.Background(), gisGeometry)

	if err != nil {
		t.Error(err)
	}

	geo2, err := testQueries.GetGISData(context.Background(), gisData.ID)
	if err != nil {
		t.Error(err)
	}

	if !cmp.Equal(gisGeometry, geo2) {
		t.Errorf("curvepolygon %v was not equal to %v", &geo2, gisGeometry)
	}
}

func TestGISCircularStringGeometry(t *testing.T) {

	circularString := makeTestCircularString(t)
	gisGeometry := geo.NewGISGeometry(circularString)

	gisGeometry.SetSRID(4326)

	gisData, err := testQueries.CreateGISData(context.Background(), gisGeometry)

	if err != nil {
		t.Error(err)
	}

	geo2, err := testQueries.GetGISData(context.Background(), gisData.ID)
	if err != nil {
		t.Error(err)
	}

	if !cmp.Equal(gisGeometry, geo2) {
		t.Errorf("circularstring %v was not equal to %v", &geo2, gisGeometry)
	}

}

func TestGISCompoundCurveGeometry(t *testing.T) {

	dims := geo.XY

	compoundCurve := makeTestPoint(t, dims)
	gisGeometry := geo.NewGISGeometry(compoundCurve)

	gisGeometry.SetSRID(4326)

	gisData, err := testQueries.CreateGISData(context.Background(), gisGeometry)

	if err != nil {
		t.Error(err)
	}

	geo2, err := testQueries.GetGISData(context.Background(), gisData.ID)
	if err != nil {
		t.Error(err)
	}

	if !cmp.Equal(gisGeometry, geo2) {
		t.Errorf("compoundcurve %v was not equal to %v", &geo2, gisGeometry)
	}

}

func TestGISMultiCurveGeometry(t *testing.T) {

	multiCurve := makeTestMultiCurve(t)
	gisGeometry := geo.NewGISGeometry(multiCurve)

	gisGeometry.SetSRID(4326)

	gisData, err := testQueries.CreateGISData(context.Background(), gisGeometry)

	if err != nil {
		t.Error(err)
	}

	geo2, err := testQueries.GetGISData(context.Background(), gisData.ID)
	if err != nil {
		t.Error(err)
	}

	if !cmp.Equal(gisGeometry, geo2) {
		t.Errorf("multicurve %v was not equal to %v", &geo2, gisGeometry)
	}

}

func TestGISMultiSurfaceGeometry(t *testing.T) {

	multiSurface := makeTestMultiSurface(t)
	gisGeometry := geo.NewGISGeometry(multiSurface)

	gisGeometry.SetSRID(4326)

	gisData, err := testQueries.CreateGISData(context.Background(), gisGeometry)

	if err != nil {
		t.Error(err)
	}

	geo2, err := testQueries.GetGISData(context.Background(), gisData.ID)
	if err != nil {
		t.Error(err)
	}

	if !cmp.Equal(gisGeometry, geo2) {
		t.Errorf("multisurface %v was not equal to %v", &geo2, gisGeometry)
	}

}

func TestGISPolyhedralSurfaceGeometry(t *testing.T) {

	polyhedralSurface := makeTestPolyhedralSurface(t)
	gisGeometry := geo.NewGISGeometry(polyhedralSurface)

	gisGeometry.SetSRID(4326)

	gisData, err := testQueries.CreateGISData(context.Background(), gisGeometry)

	if err != nil {
		t.Error(err)
	}

	geo2, err := testQueries.GetGISData(context.Background(), gisData.ID)
	if err != nil {
		t.Error(err)
	}

	if !cmp.Equal(gisGeometry, geo2) {
		t.Errorf("polyhedralsurface %v was not equal to %v", &geo2, gisGeometry)
	}

}

func TestGISTINGeometry(t *testing.T) {

	tin := makeTestTIN(t)
	gisGeometry := geo.NewGISGeometry(tin)

	gisGeometry.SetSRID(4326)

	gisData, err := testQueries.CreateGISData(context.Background(), gisGeometry)

	if err != nil {
		t.Error(err)
	}

	geo2, err := testQueries.GetGISData(context.Background(), gisData.ID)
	if err != nil {
		t.Error(err)
	}

	if !cmp.Equal(gisGeometry, geo2) {
		t.Errorf("tin %v was not equal to %v", &geo2, gisGeometry)
	}
}

func TestGISTriangleGeometry(t *testing.T) {

	triangle := makeTestTriangle(t)
	gisGeometry := geo.NewGISGeometry(triangle)

	gisGeometry.SetSRID(4326)

	gisData, err := testQueries.CreateGISData(context.Background(), gisGeometry)

	if err != nil {
		t.Error(err)
	}

	geo2, err := testQueries.GetGISData(context.Background(), gisData.ID)
	if err != nil {
		t.Error(err)
	}

	if !cmp.Equal(gisGeometry, geo2) {
		t.Errorf("triangle %v was not equal to %v", &geo2, gisGeometry)
	}
}

func TestGISGeometryCollection(t *testing.T) {

	geometryCollection := makeTestGeometryCollection(t)
	gisGeometry := geo.NewGISGeometry(geometryCollection)

	gisGeometry.SetSRID(4326)

	gisData, err := testQueries.CreateGISData(context.Background(), gisGeometry)

	if err != nil {
		t.Error(err)
	}

	geo2, err := testQueries.GetGISData(context.Background(), gisData.ID)
	if err != nil {
		t.Error(err)
	}

	if !cmp.Equal(gisGeometry, geo2) {
		t.Errorf("geometrycollection %v was not equal to %v", &geo2, gisGeometry)
	}

}
