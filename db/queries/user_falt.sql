-- Create
-- name: CreateUserFlat :one
INSERT INTO user_flats (
  username, flat_id, is_admin, balance
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- Read 
-- name: GetUserFlat :one
SELECT * FROM user_flats
WHERE id = $1 AND deleted_at IS NULL;

-- Update
-- name: UpdateUserFlat :one
UPDATE user_flats
SET 
  is_admin = $1, 
  balance = $2
WHERE id = $3 AND deleted_at IS NULL
RETURNING *;

-- Delete (Soft Delete)
-- name: SoftDeleteUserFlat :exec
UPDATE user_flats
SET 
  deleted_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- Hard Delete
-- name: HardDeleteUserFlat :exec
DELETE FROM user_flats
WHERE id = $1;

