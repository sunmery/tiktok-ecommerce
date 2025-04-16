-- name: CreateInventory :one
INSERT INTO products.inventory (product_id, merchant_id, stock)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpdateInventory :one
UPDATE products.inventory
SET stock = stock + @delta
WHERE product_id = @product_id
  AND merchant_id = @merchant_id
  AND stock + @delta >= 0 -- 防止负数库存
RETURNING *;

-- name: GetInventory :one
SELECT stock
FROM products.inventory
WHERE product_id = @product_id
  AND merchant_id = @merchant_id;