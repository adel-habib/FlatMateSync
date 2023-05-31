-- Create
-- name: CreateUser :one
INSERT INTO users (
  username, email, password_hash, oidc_id, oidc_provider
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING *;

-- Read 
-- name: GetUser :one
SELECT * FROM users
WHERE username = $1 AND deleted_at IS NULL;

-- Update
-- name: UpdateUser :one
UPDATE users
SET 
  email = $1, 
  password_hash = $2, 
  oidc_id = $3, 
  oidc_provider = $4
WHERE username = $5 AND deleted_at IS NULL
RETURNING *;

-- Delete (Soft Delete)
-- name: SoftDeleteUser :exec
UPDATE users
SET 
  deleted_at = CURRENT_TIMESTAMP
WHERE username = $1;

-- Hard Delete
-- name: HardDeleteUser :exec
DELETE FROM users
WHERE username = $1;

-- List with pagination
-- name: ListUsers :many
SELECT * FROM users
WHERE deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;
