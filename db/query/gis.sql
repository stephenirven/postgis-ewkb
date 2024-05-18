-- name: CreateGISData :one
INSERT INTO gisdata (
  geo
) VALUES (
  sqlc.arg(geo)::geometry
)
RETURNING *;

-- name: GetAllGISData :many
SELECT geo::geometry
FROM gisdata;

-- name: GetGISData :one
SELECT geo::geometry 
FROM gisdata
WHERE id = $1 LIMIT 1;

-- name: ListGISDataWithinDistance :many
SELECT geo::geometry 
FROM gisdata
WHERE
  ST_DWithin(
    geo::geometry,    
    sqlc.arg(geo)::geometry,
    sqlc.arg(range)::INT)
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: DeleteGISData :exec
DELETE FROM gisdata
WHERE id = $1;