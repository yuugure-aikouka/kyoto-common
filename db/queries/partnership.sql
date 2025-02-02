
-- name: ListPartners :many
SELECT u.id, u.username, u.display_name, u.created_at as register_date, u.avatar_url, u.is_ai
FROM users u
JOIN partnerships p ON (p.user_id_1 = u.id OR p.user_id_2 = u.id)
WHERE
    status = 'accepted'
    AND (p.user_id_1 = @user_id OR p.user_id_2 = @user_id)
    AND u.id != @user_id;
