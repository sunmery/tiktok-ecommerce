-- name: SetSensitiveWords :execrows
INSERT INTO admin.sensitive_words (created_by, category, word, level, is_active)
VALUES (@created_by::uuid, @category, @word, @level, @is_active)
ON CONFLICT (word) DO UPDATE
    SET category   = EXCLUDED.category,
        level     = EXCLUDED.level,
        is_active = EXCLUDED.is_active,
        updated_at = NOW();

-- name: GetSensitiveWords :many
SELECT *
FROM admin.sensitive_words
WHERE (@created_by::uuid IS NULL OR created_by = @created_by::uuid)
ORDER BY created_at DESC
LIMIT @page_size OFFSET @page;

-- name: GetSensitiveWordByID :one
SELECT *
FROM admin.sensitive_words
WHERE id = @id;

-- name: DeleteSensitiveWord :exec
DELETE
FROM admin.sensitive_words
WHERE id = @id
  AND created_by = @created_by::uuid;