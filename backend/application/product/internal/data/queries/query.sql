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

-- name: UpdateProduct :one
UPDATE products.products
SET name = $1, description = $2, picture = $3, price = $4, categories = $5
WHERE id = $6
RETURNING *;

-- name: DeleteProduct :exec
DELETE FROM products.products
WHERE id = @id
RETURNING *;
```

