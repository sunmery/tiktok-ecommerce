// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: query.sql

package models

import (
	"context"

	"github.com/google/uuid"
)

const CreateComment = `-- name: CreateComment :one
INSERT INTO comments.comments (id, product_id, merchant_id, user_id, score, content)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, product_id, merchant_id, user_id, score, content, created_at, updated_at
`

type CreateCommentParams struct {
	ID         int64     `json:"id"`
	ProductID  uuid.UUID `json:"productID"`
	MerchantID uuid.UUID `json:"merchantID"`
	UserID     uuid.UUID `json:"userID"`
	Score      int32     `json:"score"`
	Content    string    `json:"content"`
}

// CreateComment
//
//	INSERT INTO comments.comments (id, product_id, merchant_id, user_id, score, content)
//	VALUES ($1, $2, $3, $4, $5, $6)
//	RETURNING id, product_id, merchant_id, user_id, score, content, created_at, updated_at
func (q *Queries) CreateComment(ctx context.Context, arg CreateCommentParams) (CommentsComments, error) {
	row := q.db.QueryRow(ctx, CreateComment,
		arg.ID,
		arg.ProductID,
		arg.MerchantID,
		arg.UserID,
		arg.Score,
		arg.Content,
	)
	var i CommentsComments
	err := row.Scan(
		&i.ID,
		&i.ProductID,
		&i.MerchantID,
		&i.UserID,
		&i.Score,
		&i.Content,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const DeleteComment = `-- name: DeleteComment :exec
DELETE
FROM comments.comments
WHERE id = $1
  AND user_id = $2
`

type DeleteCommentParams struct {
	ID     int64     `json:"id"`
	UserID uuid.UUID `json:"userID"`
}

// DeleteComment
//
//	DELETE
//	FROM comments.comments
//	WHERE id = $1
//	  AND user_id = $2
func (q *Queries) DeleteComment(ctx context.Context, arg DeleteCommentParams) error {
	_, err := q.db.Exec(ctx, DeleteComment, arg.ID, arg.UserID)
	return err
}

const GetCommentCount = `-- name: GetCommentCount :one
SELECT COUNT(*)
FROM comments.comments
WHERE product_id = $1
  AND merchant_id = $2
`

type GetCommentCountParams struct {
	ProductID  uuid.UUID `json:"productID"`
	MerchantID uuid.UUID `json:"merchantID"`
}

// GetCommentCount
//
//	SELECT COUNT(*)
//	FROM comments.comments
//	WHERE product_id = $1
//	  AND merchant_id = $2
func (q *Queries) GetCommentCount(ctx context.Context, arg GetCommentCountParams) (int64, error) {
	row := q.db.QueryRow(ctx, GetCommentCount, arg.ProductID, arg.MerchantID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const GetCommentsByProduct = `-- name: GetCommentsByProduct :many
SELECT id, product_id, merchant_id, user_id, score, content, created_at, updated_at
FROM comments.comments
WHERE product_id = $1
  AND merchant_id = $2
ORDER BY created_at DESC
LIMIT $4 OFFSET $3
`

type GetCommentsByProductParams struct {
	ProductID  uuid.UUID `json:"productID"`
	MerchantID uuid.UUID `json:"merchantID"`
	Page       int64     `json:"page"`
	PageSize   int64     `json:"pageSize"`
}

// GetCommentsByProduct
//
//	SELECT id, product_id, merchant_id, user_id, score, content, created_at, updated_at
//	FROM comments.comments
//	WHERE product_id = $1
//	  AND merchant_id = $2
//	ORDER BY created_at DESC
//	LIMIT $4 OFFSET $3
func (q *Queries) GetCommentsByProduct(ctx context.Context, arg GetCommentsByProductParams) ([]CommentsComments, error) {
	rows, err := q.db.Query(ctx, GetCommentsByProduct,
		arg.ProductID,
		arg.MerchantID,
		arg.Page,
		arg.PageSize,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []CommentsComments
	for rows.Next() {
		var i CommentsComments
		if err := rows.Scan(
			&i.ID,
			&i.ProductID,
			&i.MerchantID,
			&i.UserID,
			&i.Score,
			&i.Content,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const UpdateComment = `-- name: UpdateComment :one
UPDATE comments.comments
SET score      = COALESCE($1, score),
    content    = COALESCE($2, content),
    updated_at = NOW()
WHERE id = $3
  AND user_id = $4
RETURNING id, product_id, merchant_id, user_id, score, content, created_at, updated_at
`

type UpdateCommentParams struct {
	Score   int32     `json:"score"`
	Content string    `json:"content"`
	ID      int64     `json:"id"`
	UserID  uuid.UUID `json:"userID"`
}

// UpdateComment
//
//	UPDATE comments.comments
//	SET score      = COALESCE($1, score),
//	    content    = COALESCE($2, content),
//	    updated_at = NOW()
//	WHERE id = $3
//	  AND user_id = $4
//	RETURNING id, product_id, merchant_id, user_id, score, content, created_at, updated_at
func (q *Queries) UpdateComment(ctx context.Context, arg UpdateCommentParams) (CommentsComments, error) {
	row := q.db.QueryRow(ctx, UpdateComment,
		arg.Score,
		arg.Content,
		arg.ID,
		arg.UserID,
	)
	var i CommentsComments
	err := row.Scan(
		&i.ID,
		&i.ProductID,
		&i.MerchantID,
		&i.UserID,
		&i.Score,
		&i.Content,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
