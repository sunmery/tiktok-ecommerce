-- name: CreateProduct :one
INSERT INTO products.products(name, description, picture, price, categories)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: ListProducts :many
SELECT *
FROM products.products
ORDER BY id
OFFSET @page LIMIT @page_size;

-- name: GetProduct :one
SELECT *
FROM products.products
WHERE id = @id
LIMIT 1;

-- name: SearchProducts :many
SELECT *
FROM products.products
WHERE name ILIKE '%' || @name || '%';
