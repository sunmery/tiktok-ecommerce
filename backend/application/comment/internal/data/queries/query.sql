-- name: CreateComment :one
WITH check_sensitive_word AS (
    -- 检查 content 是否包含敏感词
    SELECT EXISTS (SELECT 1
                   FROM admin.sensitive_words sw
                   WHERE @content::text ILIKE '%' || sw.word || '%'
                     AND sw.is_active = TRUE) AS has_sensitive_word),
     insert_comment AS (
         -- 如果没有检测到敏感词，则执行插入操作
         INSERT INTO comments.comments (id, product_id, merchant_id, user_id, score, content)
             SELECT @id::bigint, @product_id::uuid, @merchant_id::uuid, @user_id::uuid, @score, @content
             WHERE NOT (SELECT has_sensitive_word FROM check_sensitive_word)
             RETURNING id, product_id, merchant_id, user_id, score, content, created_at, updated_at)
-- 返回最终结果：是否命中敏感词
SELECT has_sensitive_word AS is_sensitive
FROM check_sensitive_word;

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
