-- name: CreateRefreshToken :exec
INSERT INTO refresh_tokens (token) VALUES ($1);
