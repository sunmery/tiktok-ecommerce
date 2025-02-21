-- name: CreateOrder :one
INSERT INTO orders.orders (id, user_id, currency, street_address,
                           city, state, country, zip_code, email,
                           created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING *;

-- name: GetOrderByID :one
SELECT *
FROM orders.orders
WHERE id = $1;

-- name: ListOrdersByUser :many
SELECT *
FROM orders.orders
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: CreateSubOrder :one
INSERT INTO orders.sub_orders (id, order_id, merchant_id, total_amount,
                               currency, status, items, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: UpdateSubOrderStatus :exec
UPDATE orders.sub_orders
SET status     = $2,
    updated_at = $3
WHERE id = $1;
