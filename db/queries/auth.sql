
-- name: InvalidateActiveMagicLink :exec
DELETE FROM magic_links WHERE email = $1 AND expires_at > NOW();

-- name: CreateMagicLink :exec
INSERT INTO magic_links (
    email, token, expires_at
) VALUES (
    $1, $2, NOW() + INTERVAL '15 minutes'
);
