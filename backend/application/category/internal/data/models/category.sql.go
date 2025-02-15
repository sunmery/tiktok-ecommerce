// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: category.sql

package models

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const CreateCategory = `-- name: CreateCategory :one
INSERT INTO categories.categories (name, level)
VALUES ($1, $2)
RETURNING id, name, created_at, updated_at
`

type CreateCategoryParams struct {
	Name  string `json:"name"`
	Level int32  `json:"level"`
}

type CreateCategoryRow struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// CreateCategory
//
//	INSERT INTO categories.categories (name, level)
//	VALUES ($1, $2)
//	RETURNING id, name, created_at, updated_at
func (q *Queries) CreateCategory(ctx context.Context, arg CreateCategoryParams) (CreateCategoryRow, error) {
	row := q.db.QueryRow(ctx, CreateCategory, arg.Name, arg.Level)
	var i CreateCategoryRow
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const FindOrCreateCategory = `-- name: FindOrCreateCategory :one
WITH s AS (
    SELECT id FROM categories.categories WHERE name = $1
), i AS (
    INSERT INTO categories.categories (name)
        SELECT $1
        WHERE NOT EXISTS (SELECT 1 FROM s)
        RETURNING id
)
SELECT id FROM i
UNION ALL
SELECT id FROM s
`

// FindOrCreateCategory
//
//	WITH s AS (
//	    SELECT id FROM categories.categories WHERE name = $1
//	), i AS (
//	    INSERT INTO categories.categories (name)
//	        SELECT $1
//	        WHERE NOT EXISTS (SELECT 1 FROM s)
//	        RETURNING id
//	)
//	SELECT id FROM i
//	UNION ALL
//	SELECT id FROM s
func (q *Queries) FindOrCreateCategory(ctx context.Context, dollar_1 *string) (uuid.UUID, error) {
	row := q.db.QueryRow(ctx, FindOrCreateCategory, dollar_1)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const GetCategoryByName = `-- name: GetCategoryByName :one
SELECT id, name, created_at, updated_at
FROM categories.categories
WHERE name = $1 LIMIT 1
`

type GetCategoryByNameRow struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// GetCategoryByName
//
//	SELECT id, name, created_at, updated_at
//	FROM categories.categories
//	WHERE name = $1 LIMIT 1
func (q *Queries) GetCategoryByName(ctx context.Context, name string) (GetCategoryByNameRow, error) {
	row := q.db.QueryRow(ctx, GetCategoryByName, name)
	var i GetCategoryByNameRow
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
