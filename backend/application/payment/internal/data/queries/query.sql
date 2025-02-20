-- name: Creat :one
INSERT INTO table_name(owner, name)
VALUES ($1, $2) RETURNING *;

-- name: Get :many
SELECT *
FROM table_name
WHERE owner = @owner
  AND name = @name;

-- name: Update :one
UPDATE table_name
SET col1 = coalesce(sqlc.narg(col1), col1),
    col2 = coalesce(sqlc.narg(col2), col2)
WHERE name = @name RETURNING *;

-- name: Delete :one
DELETE
FROM table_name
WHERE name = @name RETURNING *;
