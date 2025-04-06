-- name: CreateUser :one
INSERT INTO users (email, username, display_name, avatar_url, is_ai)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetUserById :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;
