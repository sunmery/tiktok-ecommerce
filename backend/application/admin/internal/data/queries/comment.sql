-- name: CreateBulkSensitiveWords :execrows
-- 跳过重复的敏感词
INSERT
INTO admin.sensitive_words(created_by, category, word, level, is_active)
SELECT unnest(@created_by::UUID[]),
       unnest(@categories::VARCHAR[]),
       unnest(@words::VARCHAR[]),
       unnest(@level::INT[]),
       unnest(@is_active::BOOL[])
ON CONFLICT (word) DO NOTHING;

-- name: GetSensitiveWords :many
SELECT *
FROM admin.sensitive_words
ORDER BY created_at DESC
LIMIT @page_size::INT OFFSET @page::INT;

-- name: UpdateSensitiveWord :execrows
UPDATE admin.sensitive_words
SET category   = COALESCE(@category, category),
    created_by = COALESCE(@created_by, created_by),
    word       = COALESCE(@word, word),
    level      = COALESCE(@level, level),
    is_active  = COALESCE(@is_active, is_active)
WHERE id = @id;

-- name: DeleteSensitiveWords :execrows
DELETE
FROM admin.sensitive_words
WHERE id = @id;
