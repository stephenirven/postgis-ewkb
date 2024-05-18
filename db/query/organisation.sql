-- name: CreateOrganisation :one
INSERT INTO organisation (
  country_code, merchant_name
) VALUES (
  $1, $2
)
RETURNING *;

-- name: GetOrganisation :one
SELECT * FROM organisation
WHERE id = $1 LIMIT 1;

-- name: GetOrganisationForUpdate :one
SELECT * FROM organisation
WHERE id = $1 LIMIT 1
FOR NO KEY UPDATE;


-- -- name: GetDetails :many
-- SELECT sqlc.embed(accounts), sqlc.embed(transfers) 
-- FROM accounts JOIN transfers on (transfers.from_account_id = accounts.id OR transfers.to_account_id  = accounts.id);

-- name: ListOrganisations :many
SELECT * FROM organisation
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateOrganisation :one
UPDATE organisation
  set country_code = $2, 
  merchant_name = $3 
WHERE id = $1
RETURNING *;


-- name: DeleteOrganisation :exec
DELETE FROM organisation
WHERE id = $1;