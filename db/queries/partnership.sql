-- name: CreatePartnership :one
INSERT INTO partnerships (
    user_id_1, user_id_2
) VALUES (
    $1, $2
) RETURNING *;

-- name: UpdatePartnershipStatus :exec
UPDATE partnerships
SET status = $1
WHERE user_id_1 = $2 AND user_id_2 = $3;

-- name: ListPartners :many
SELECT u.id, u.username, u.display_name, u.created_at as register_date, u.avatar_url, u.is_ai
FROM users u
JOIN partnerships p ON (p.user_id_1 = u.id OR p.user_id_2 = u.id)
WHERE
    status = 'accepted'
    AND (p.user_id_1 = @user_id OR p.user_id_2 = @user_id)
    AND u.id != @user_id;

-- name: ListPotentialPartners :many
SELECT u.id, u.username, u.display_name, u.created_at as register_date, u.avatar_url, u.is_ai
FROM users u
LEFT JOIN partnerships p -- uses the "left join null filtering" technique, to select ones that don't satisfy the ON condition
    ON (p.user_id_1 = u.id OR p.user_id_2 = u.id)
    AND (
        p.user_id_1 = @user_id -- exclude users that $1 has liked, but not the other way around
        OR (p.user_id_2 = @user_id AND p.status = 'accepted') -- unless a partnership has established
    )
WHERE
    u.id != @user_id AND p.user_id_1 IS NULL;
