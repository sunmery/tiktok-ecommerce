-- name: CreateProduct :one
INSERT INTO products.item (
  sku, name, description, category_paths,
  main_image_url, image_gallery, price, status
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING *;

-- name: GetProduct :one
SELECT * FROM products.item
WHERE id = $1 LIMIT 1;

-- name: UpdateProduct :one
UPDATE products.item SET
  name = $2,
  description = $3,
  category_paths = $4,
  main_image_url = $5,
  image_gallery = $6,
  price = $7,
  status = $8,
  updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: SearchCategory :many
SELECT * FROM products.item
WHERE category_paths @> '["5f8d04b3"]'::jsonb;

-- name: DeleteProduct :exec
DELETE FROM products.item
WHERE id = $1;
