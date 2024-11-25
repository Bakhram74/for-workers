-- name: CreateVehicle :one
INSERT INTO vehicle (
  id,
  user_id,
  type,
  brand,
  number,
  region,
  country
) VALUES (
 $1, $2, $3, $4, $5, $6,$7
) RETURNING *;

-- name: GetVehicleByID :one
SELECT * FROM vehicle
WHERE id = $1 AND user_id = $2;

-- name: DeleteVehicle :exec
DELETE FROM vehicle
WHERE id = $1 AND user_id = $2;


-- name: UpdateVehicle :one
UPDATE vehicle
SET
  type = COALESCE(sqlc.narg(type), type),
  brand = COALESCE(sqlc.narg(brand), brand),
  number = COALESCE(sqlc.narg(number), number),
  region = COALESCE(sqlc.narg(region), region),
  country = COALESCE(sqlc.narg(country), country)
WHERE id = sqlc.arg(id) AND user_id = sqlc.arg(user_id)
RETURNING *;


-- name: GetVehicleByRegion :many
SELECT * FROM vehicle
WHERE number LIKE $1 AND region = $2
LIMIT $3 OFFSET $4;

-- name: CountVehicleByRegion :one
SELECT COUNT(id) FROM vehicle
WHERE number LIKE $1 AND region = $2;


-- name: CountVehicleByNumber :one
SELECT COUNT(id) FROM vehicle
 WHERE number LIKE $1;

-- name: GetVehicleByNumber :many
SELECT * FROM vehicle
WHERE number LIKE $1 LIMIT $2 OFFSET $3;

