-- name: CreateComment :one
INSERT INTO comments.comments (id, product_id, merchant_id, user_id, score, content)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetCommentsByProduct :many
SELECT *
FROM comments.comments
WHERE product_id = @product_id
  AND merchant_id = @merchant_id
ORDER BY created_at DESC
LIMIT @page_size OFFSET @page;

-- name: GetCommentCount :one
SELECT COUNT(*)
FROM comments.comments
WHERE product_id = @product_id
  AND merchant_id = @merchant_id;

-- name: UpdateComment :one
UPDATE comments.comments
SET score      = COALESCE(@score, score),
    content    = COALESCE(@content, content),
    updated_at = NOW()
WHERE id = @id
  AND user_id = @user_id
RETURNING *;

-- name: DeleteComment :exec
DELETE
FROM comments.comments
WHERE id = @id
  AND user_id = @user_id;