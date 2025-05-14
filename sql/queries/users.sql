-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: UpdateUser :one
UPDATE users
SET hashed_password = $1,
    email = $2,
    updated_at = NOW()
WHERE id = $3
RETURNING *;

-- name: GetUser :one
SELECT *
    FROM users
    WHERE email = $1;

-- name: GetUsers :many
SELECT email
    FROM users
    ORDER BY email;

-- name: DeleteUsers :exec
DELETE FROM users;

