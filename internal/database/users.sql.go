// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: users.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING id, created_at, updated_at, email, hashed_password, is_chirpy_red
`

type CreateUserParams struct {
	Email          string
	HashedPassword string
}

// email TEXT
func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser, arg.Email, arg.HashedPassword)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.HashedPassword,
		&i.IsChirpyRed,
	)
	return i, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, created_at, updated_at, email, hashed_password, is_chirpy_red FROM users WHERE email = $1
`

// email TEXT
func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.HashedPassword,
		&i.IsChirpyRed,
	)
	return i, err
}

const updateUser = `-- name: UpdateUser :one
UPDATE users
SET
    email = $2,
    hashed_password = $3,
    updated_at = NOW()
WHERE
    id = $1
RETURNING id, created_at, updated_at, email, hashed_password, is_chirpy_red
`

type UpdateUserParams struct {
	ID             uuid.UUID
	Email          string
	HashedPassword string
}

// id UUID
// email TEXT
// hashed_password TEXT
func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, updateUser, arg.ID, arg.Email, arg.HashedPassword)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.HashedPassword,
		&i.IsChirpyRed,
	)
	return i, err
}

const upgradeUserChirpyRed = `-- name: UpgradeUserChirpyRed :exec
UPDATE users
SET
    is_chirpy_red = $2,
    updated_at = NOW()
WHERE
    id = $1
`

type UpgradeUserChirpyRedParams struct {
	ID          uuid.UUID
	IsChirpyRed bool
}

// id UUID
// is_chirpy_red BOOLEAN
func (q *Queries) UpgradeUserChirpyRed(ctx context.Context, arg UpgradeUserChirpyRedParams) error {
	_, err := q.db.ExecContext(ctx, upgradeUserChirpyRed, arg.ID, arg.IsChirpyRed)
	return err
}
