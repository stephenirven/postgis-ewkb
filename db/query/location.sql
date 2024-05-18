-- name: CreateLocation :one
INSERT INTO location (
  full_name, geo
) VALUES (
  $1, sqlc.arg(geo)::geometry
)
RETURNING *;

-- name: GetLocations :many
SELECT id, organisation_id, user_id, full_name,
		line1
		line2,
		city,
		county,
		country_code,
		geo::geometry,
		created_at
FROM location;

-- name: GetLocation :one
SELECT 
		id, organisation_id, user_id, full_name,
		line1
		line2,
		city,
		county,
		country_code,
		geo::geometry,
		created_at
FROM location
WHERE id = $1 LIMIT 1;


-- name: GetLocationForUpdate :one
SELECT 
		id, organisation_id, user_id, full_name,
		line1
		line2,
		city,
		county,
		country_code,
		geo::geometry,
		created_at
FROM location
WHERE id = $1 LIMIT 1
FOR NO KEY UPDATE;

-- name: ListLocationsWithinDistance :many
SELECT 
		id, organisation_id, user_id, full_name,
		line1
		line2,
		city,
		county,
		country_code,
		geo::geometry,
		created_at
FROM location
WHERE
  ST_DWithin(
    geo::geometry,    
    sqlc.arg(location)::geometry,
    sqlc.arg(range)::INT)
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: DeleteLocation :exec
DELETE FROM location
WHERE id = $1;