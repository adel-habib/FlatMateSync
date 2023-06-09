// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: flats.sql

package db

import (
	"context"
)

const createFlat = `-- name: CreateFlat :one
INSERT INTO flats (
  name
) VALUES (
  $1
)
RETURNING id, name, deleted_at, created_at
`

// Create
func (q *Queries) CreateFlat(ctx context.Context, name string) (Flat, error) {
	row := q.db.QueryRowContext(ctx, createFlat, name)
	var i Flat
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.DeletedAt,
		&i.CreatedAt,
	)
	return i, err
}

const getFlat = `-- name: GetFlat :one
SELECT id, name, deleted_at, created_at FROM flats
WHERE id = $1 AND deleted_at IS NULL
`

// Read
func (q *Queries) GetFlat(ctx context.Context, id int32) (Flat, error) {
	row := q.db.QueryRowContext(ctx, getFlat, id)
	var i Flat
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.DeletedAt,
		&i.CreatedAt,
	)
	return i, err
}

const getFlatWithUsers = `-- name: GetFlatWithUsers :many
SELECT 
  f.id, f.name, f.deleted_at, f.created_at,
  u.username, u.email, u.password_hash, u.oidc_id, u.oidc_provider, u.deleted_at, u.created_at
FROM
  flats f
JOIN 
  user_flats uf ON f.id = uf.flat_id
JOIN 
  users u ON uf.username = u.username
WHERE 
  f.id = $1
`

type GetFlatWithUsersRow struct {
	Flat Flat
	User User
}

func (q *Queries) GetFlatWithUsers(ctx context.Context, id int32) ([]GetFlatWithUsersRow, error) {
	rows, err := q.db.QueryContext(ctx, getFlatWithUsers, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetFlatWithUsersRow
	for rows.Next() {
		var i GetFlatWithUsersRow
		if err := rows.Scan(
			&i.Flat.ID,
			&i.Flat.Name,
			&i.Flat.DeletedAt,
			&i.Flat.CreatedAt,
			&i.User.Username,
			&i.User.Email,
			&i.User.PasswordHash,
			&i.User.OidcID,
			&i.User.OidcProvider,
			&i.User.DeletedAt,
			&i.User.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const hardDeleteFlat = `-- name: HardDeleteFlat :exec
DELETE FROM flats
WHERE id = $1
`

// Hard Delete
func (q *Queries) HardDeleteFlat(ctx context.Context, id int32) error {
	_, err := q.db.ExecContext(ctx, hardDeleteFlat, id)
	return err
}

const listFlats = `-- name: ListFlats :many
SELECT id, name, deleted_at, created_at FROM flats
WHERE deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $1 OFFSET $2
`

type ListFlatsParams struct {
	Limit  int32
	Offset int32
}

// List with pagination
func (q *Queries) ListFlats(ctx context.Context, arg ListFlatsParams) ([]Flat, error) {
	rows, err := q.db.QueryContext(ctx, listFlats, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Flat
	for rows.Next() {
		var i Flat
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.DeletedAt,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const softDeleteFlat = `-- name: SoftDeleteFlat :exec
UPDATE flats
SET 
  deleted_at = CURRENT_TIMESTAMP
WHERE id = $1
`

// Delete (Soft Delete)
func (q *Queries) SoftDeleteFlat(ctx context.Context, id int32) error {
	_, err := q.db.ExecContext(ctx, softDeleteFlat, id)
	return err
}

const updateFlat = `-- name: UpdateFlat :one
UPDATE flats
SET 
  name = $1
WHERE id = $2 AND deleted_at IS NULL
RETURNING id, name, deleted_at, created_at
`

type UpdateFlatParams struct {
	Name string
	ID   int32
}

// Update
func (q *Queries) UpdateFlat(ctx context.Context, arg UpdateFlatParams) (Flat, error) {
	row := q.db.QueryRowContext(ctx, updateFlat, arg.Name, arg.ID)
	var i Flat
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.DeletedAt,
		&i.CreatedAt,
	)
	return i, err
}
