-- name: CreateComment :one
INSERT INTO comments.comments (product_id, merchant_id, user_id, score, content)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetCommentsByProduct :many
SELECT *
FROM comments.comments
WHERE product_id = @product_id
  AND merchant_id = @merchant_id
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetCommentCount :one
SELECT COUNT(*)
FROM comments.comments
WHERE product_id = @product_id
  AND merchant_id = @merchant_id;

-- name: UpdateComment :one
UPDATE comments.comments
SET score      = COALESCE($2, score),
    content    = COALESCE($3, content),
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
  AND user_id = $4
RETURNING *;

-- name: DeleteComment :exec
DELETE
FROM comments.comments
WHERE id = $1
  AND user_id = $2;