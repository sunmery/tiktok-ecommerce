-- name: CreateProduct :one
INSERT INTO products.products (name,
                               description,
                               picture,
                               price,
                               stock,
                               category_id)
VALUES ($1, $2, $3, $4, $5,$6)
RETURNING *;

-- name: GetProduct :one
SELECT *
FROM products.products
WHERE id = $1
LIMIT 1;

-- name: ListProducts :many
SELECT *
FROM products.products
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: UpdateProduct :one
UPDATE products.products
SET name        = COALESCE($2, name),
    description       = COALESCE($3, description),
    picture       = COALESCE($4, picture),
    price       = COALESCE($5, price),
    stock       = COALESCE($6, stock),
    category_id = COALESCE($7, category_id)
WHERE id = $1
RETURNING *;

-- name: DeleteProduct :exec
DELETE
FROM products.products
WHERE id = $1;

-- name: LockProductStock :one
SELECT stock
FROM products.products
WHERE id = $1
    FOR UPDATE;

-- name: UpdateProductStock :exec
UPDATE products.products
SET stock = $2
WHERE id = $1;
