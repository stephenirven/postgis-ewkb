package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	db "github.com/stephenirven/go-postgis/db/sqlc"
	"github.com/stephenirven/go-postgis/geo"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://postgres:secret@localhost:5432/late?sslmode=disable"
)

var testQueries *db.Queries
var testDB *sql.DB

func init() {
	testDB, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("failed to connect to database", err)
	}
	testQueries = db.New(testDB)
}

func main() {

	point, err := geo.NewPoint([]float64{0, 1, 2}, geo.XYZ)
	if err != nil {
		fmt.Print(err.Error())
	}
	fmt.Printf("Geo: %v\n", point)
	gisGeometry := geo.NewGISGeometry(point)

	gisGeometry.SetSRID(4326)

	fmt.Printf("Geo: %v\n", gisGeometry)

}

func getLocation(id int64) {

	row, err := testQueries.GetLocation(context.Background(), id)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Print(row.ID, " - ", row.Geo, "\n")
}

func getNear(point geo.GISGeometry, distance int32) {

	args := db.ListLocationsWithinDistanceParams{Limit: 10, Offset: 0, Location: point, Range: distance}

	rows, err := testQueries.ListLocationsWithinDistance(context.Background(), args)

	if err != nil {
		fmt.Println(err)
	}

	for _, loc := range rows {

		fmt.Print(loc.ID, " - ", loc.Geo)

		break
	}
}

func insert(geometry geo.GISGeometry) {
	crargs := db.CreateLocationParams{FullName: sql.NullString{String: "boo"}, Geo: geometry}

	_, err := testQueries.CreateLocation(context.Background(), crargs)

	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}

func getLocations() []db.GetLocationsRow {

	rows, err := testQueries.GetLocations(context.Background())
	if err != nil {
		fmt.Print(err.Error())
	}
	return rows

}
