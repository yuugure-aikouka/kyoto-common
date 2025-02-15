-- name: CreateUser :one
INSERT INTO users (
    username, display_name, avatar_url, is_ai
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: GetUser :one
SELECT * FROM users WHERE id = $1;
