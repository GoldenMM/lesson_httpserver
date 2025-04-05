-- name: CreateRefreshToken :one
-- -- -- user_id UUID
-- -- -- expires_at TIMESTAMP
-- -- -- revoked_at TIMESTAMP
INSERT INTO refresh_tokens (token, created_at, updated_at, expires_at, revoked_at, user_id)
VALUES (
    $1,
    NOW(),
    NOW(),
    $2,
    NULL,
    $3
)
RETURNING *;

-- name: GetRefreshToken :one
SELECT * FROM refresh_tokens
WHERE token = $1 AND (expires_at > NOW()) AND (revoked_at IS NULL);

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens
SET 
    revoked_at = NOW(),
    updated_at = NOW()
WHERE
    token = $1;

-- name: GetAllRefreshTokens :many
SELECT * FROM refresh_tokens;