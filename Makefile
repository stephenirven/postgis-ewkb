postgres:
	docker run --name postgis -p 5432:5432 -e POSTGRES_PASSWORD=secret -d postgis/postgis 
createdb:
	docker exec -it postgis createdb --username=postgres --owner=postgres late 
dropdb: 
	docker exec -it postgis dropdb --username=postgres late
migrateup:
	migrate -path db/migration -database "postgresql://postgres:secret@localhost:5432/late?sslmode=disable" -verbose up
migratedown:
	migrate -path db/migration -database "postgresql://postgres:secret@localhost:5432/late?sslmode=disable" -verbose down
sqlc:
	sqlc generate
test:
	go test -v -cover ./... -count 1
cover: 
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test cover