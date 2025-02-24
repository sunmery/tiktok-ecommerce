-- name: CreateOrder :one
INSERT INTO orders.orders (user_id, currency, street_address,
                           city, state, country, zip_code, email)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetOrderByID :one
SELECT *
FROM orders.orders
WHERE id = $1;

-- -- name: ListOrdersByUser :many
-- SELECT *
-- FROM orders.orders
-- WHERE user_id = $1
-- ORDER BY created_at DESC
-- LIMIT $2 OFFSET $3;

-- name: GetUserOrdersWithSuborders :many
SELECT o.id         AS order_id,
       o.currency   AS order_currency,
       o.street_address,
       o.city,
       o.state,
       o.country,
       o.zip_code,
       o.email,
       o.created_at AS order_created,
       jsonb_agg(
               jsonb_build_object(
                       'suborder_id', so.id,
                       'merchant_id', so.merchant_id,
                       'total_amount', so.total_amount,
                       'currency', so.currency,
                       'status', so.status,
                       'items', so.items,
                       'created_at', so.created_at,
                       'updated_at', so.updated_at
               ) ORDER BY so.created_at
       )            AS suborders
FROM orders.orders o
         LEFT JOIN orders.sub_orders so ON o.id = so.order_id
WHERE o.user_id = $1::uuid
GROUP BY o.id
ORDER BY o.created_at DESC;

-- name: CreateSubOrder :one
INSERT INTO orders.sub_orders (order_id, merchant_id, total_amount,
                               currency, status, items)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: UpdateSubOrderStatus :exec
UPDATE orders.sub_orders
SET status     = $2,
    updated_at = $3
WHERE id = $1;
