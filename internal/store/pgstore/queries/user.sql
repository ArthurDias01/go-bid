-- name: CreateUser :one
INSERT INTO users (user_name, email, password_hash, bio)
VALUES ($1, $2, $3, $4)
RETURNING id;

-- name: GetUserByEmail :one
SELECT id, user_name, email, bio, created_at, updated_at, password_hash
FROM users
WHERE email = $1;

-- name: GetUserByID :one
SELECT id, user_name, email, bio, created_at, updated_at
FROM users
WHERE id = $1;

-- name: UpdateUser :one
UPDATE users
SET user_name = $2, email = $3, password_hash = $4, bio = $5
WHERE id = $1
RETURNING id;

-- name: DeleteUser :one
DELETE FROM users
WHERE id = $1
RETURNING id;
