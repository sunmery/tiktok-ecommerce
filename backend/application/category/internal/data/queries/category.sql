-- name: GetCategoryByName :one
SELECT id, name, created_at, updated_at
FROM categories.categories
WHERE name = $1 LIMIT 1;

-- name: CreateCategory :one
INSERT INTO categories.categories (name, level)
VALUES ($1, $2)
RETURNING id, name, created_at, updated_at;

-- name: FindOrCreateCategory :one
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
SELECT id FROM s;
