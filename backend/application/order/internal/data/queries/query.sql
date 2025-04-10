-- name: CreateOrder :one
INSERT INTO orders.orders (id, user_id, currency, street_address,
                           city, state, country, zip_code, email)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: GetOrderByID :one
SELECT o.*,
       json_agg(
               json_build_object(
                       'id', so.id,
                       'merchant_id', so.merchant_id,
                       'total_amount', so.total_amount,
                       'currency', so.currency,
                       'status', so.status,
                       'items', so.items,
                       'created_at', so.created_at,
                       'updated_at', so.updated_at
               )
       ) AS sub_orders
FROM orders.orders o
         LEFT JOIN orders.sub_orders so ON o.id = so.order_id
WHERE o.user_id = @user_id
  AND o.id = @order_id
GROUP BY o.id;

-- name: GetOrderByUserID :one
SELECT *
FROM orders.orders
WHERE user_id = @user_id;

-- name: ListOrders :many
SELECT *
FROM orders.sub_orders
ORDER BY created_at DESC
LIMIT @page_size OFFSET @page;

-- name: GetConsumerOrders :many
SELECT o.*,
       json_agg(
               json_build_object(
                       'id', so.id,
                       'merchant_id', so.merchant_id,
                       'total_amount', so.total_amount,
                       'currency', so.currency,
                       'status', so.status,
                       'items', so.items,
                       'created_at', so.created_at,
                       'updated_at', so.updated_at
               )
       ) AS sub_orders
FROM orders.orders o
         LEFT JOIN orders.sub_orders so ON o.id = so.order_id
WHERE o.user_id = @user_id
GROUP BY o.id
LIMIT @page_size OFFSET @page;

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
GROUP BY o.id, o.currency, o.street_address, o.city, o.state, o.country, o.zip_code, o.email, o.created_at
ORDER BY o.created_at DESC;

-- name: QuerySubOrders :many
SELECT id,
       merchant_id,
       total_amount,
       currency,
       status,
       items,
       created_at,
       updated_at
FROM orders.sub_orders
WHERE order_id = $1
ORDER BY created_at;

-- name: CreateSubOrder :one
INSERT INTO orders.sub_orders (id, order_id, merchant_id, total_amount,
                               currency, status, items)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: UpdateSubOrderStatus :exec
UPDATE orders.sub_orders
SET status     = $2,
    updated_at = $3
WHERE id = $1;

-- name: UpdatePaymentStatus :one
SELECT id, user_id, payment_status
FROM orders.orders
WHERE id = $1
    FOR UPDATE;

-- name: MarkOrderAsPaid :one
UPDATE orders.orders
SET payment_status = $1,
    updated_at     = now()
WHERE id = $2
RETURNING *;

-- name: MarkSubOrderAsPaid :one
UPDATE orders.sub_orders
SET payment_status = $1,
    updated_at     = now()
WHERE order_id = $2
RETURNING *;

-- name: UpdateOrderPaymentStatus :exec
UPDATE orders.orders
SET payment_status = $2,
    updated_at     = now()
WHERE id = $1;
