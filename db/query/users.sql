-- name: GetUserByPhone :one
SELECT * FROM users
WHERE phone = $1;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1;


-- name: CreateUser :one
INSERT INTO users (
  id,
  name,
  phone,
  image_url,
  status_text
) VALUES (
 $1, $2, $3,$4,$5
) RETURNING *;

-- name: UpdateUser :one
UPDATE users
SET
  role = COALESCE(sqlc.narg(role), role),
  name = COALESCE(sqlc.narg(name), name),
  phone = COALESCE(sqlc.narg(phone), phone),
  image_url = COALESCE(sqlc.narg(image_url), image_url),
  status_text = COALESCE(sqlc.narg(status_text), status_text)
WHERE id = sqlc.arg(id)
RETURNING *;

