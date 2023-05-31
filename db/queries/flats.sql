-- Create
-- name: CreateFlat :one
INSERT INTO flats (
  name
) VALUES (
  $1
)
RETURNING *;

-- Read 
-- name: GetFlat :one
SELECT * FROM flats
WHERE id = $1 AND deleted_at IS NULL;

-- Update
-- name: UpdateFlat :one
UPDATE flats
SET 
  name = $1
WHERE id = $2 AND deleted_at IS NULL
RETURNING *;

-- Delete (Soft Delete)
-- name: SoftDeleteFlat :exec
UPDATE flats
SET 
  deleted_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- Hard Delete
-- name: HardDeleteFlat :exec
DELETE FROM flats
WHERE id = $1;

-- List with pagination
-- name: ListFlats :many
SELECT * FROM flats
WHERE deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;
