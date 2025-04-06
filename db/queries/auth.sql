-- name: CreateRefreshToken :exec
INSERT INTO refresh_tokens (token) VALUES ($1);

-- name: CheckRefreshTokenExist :one
SELECT COUNT(*) > 0 FROM refresh_tokens WHERE token = $1;
