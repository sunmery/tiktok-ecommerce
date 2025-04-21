-- name: CreateOrder :one
INSERT INTO orders.orders (id, user_id, currency, street_address,
                           city, state, country, zip_code, email)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: GetOrderByID :one
SELECT o.*,
       json_agg(
               json_build_object(
                       'id', os.id,
                       'merchant_id', os.merchant_id,
                       'total_amount', os.total_amount,
                       'currency', os.currency,
                       'status', os.status,
                       'items', os.items,
                       'created_at', os.created_at,
                       'updated_at', os.updated_at
               )
       ) AS sub_orders
FROM orders.orders o
         LEFT JOIN orders.sub_orders os ON o.id = os.order_id
WHERE o.user_id = @user_id
  AND o.id = @order_id
GROUP BY o.id;

-- name: GetOrderByUserID :one
SELECT os.*,
       oo.payment_status
FROM orders.sub_orders os
         JOIN orders.orders oo
              ON os.order_id = oo.id
WHERE user_id = @user_id;

-- name: ListOrders :many
SELECT os.*,
       oo.payment_status
FROM orders.sub_orders os
         JOIN orders.orders oo
              ON os.order_id = oo.id
ORDER BY os.created_at DESC
LIMIT @page_size OFFSET @page;

-- name: GetConsumerOrders :many
SELECT oo.*,
       json_agg(
               json_build_object(
                       'id', os.id,
                       'merchant_id', os.merchant_id,
                       'total_amount', os.total_amount,
                       'currency', os.currency,
                       'status', os.status,
                       'shipping_status', os.shipping_status,
                       'items', os.items,
                       'created_at', os.created_at,
                       'updated_at', os.updated_at
               )
       ) AS sub_orders
FROM orders.orders oo
         LEFT JOIN orders.sub_orders os ON oo.id = os.order_id
WHERE oo.user_id = @user_id
GROUP BY oo.id
LIMIT @page_size OFFSET @page;

-- name: GetUserOrdersWithSuborders :many
SELECT o.id,
       o.street_address,
       o.city,
       o.state,
       o.country,
       o.zip_code,
       o.email,
       o.payment_status,

       o.created_at AS order_created,
       jsonb_agg(
               jsonb_build_object(
                       'sub_order_id', os.id,
                       'merchant_id', os.merchant_id,
                       'total_amount', os.total_amount,
                       'currency', os.currency,
                       'status', os.status,
                       'shipping_status', os.shipping_status,
                       'items', os.items,
                       'created_at', os.created_at,
                       'updated_at', os.updated_at
               ) ORDER BY os.created_at
       )            AS suborders
FROM orders.orders o
         LEFT JOIN orders.sub_orders os ON o.id = os.order_id
WHERE o.user_id = @user_id::uuid
GROUP BY o.id, o.currency, o.street_address, o.city, o.state, o.country, o.zip_code, o.email, o.created_at
ORDER BY o.created_at DESC;

-- name: QuerySubOrders :many
SELECT os.*,
       oo.payment_status
FROM orders.sub_orders os
         Join orders.orders oo on os.order_id = oo.id
WHERE order_id = @order_id
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
SET status     = $1,
    updated_at = now()
WHERE order_id = $2
RETURNING *;

-- name: GetOrderStatus :one
SELECT o.payment_status,
       s.shipping_status,
       s.id sub_order_id,
       o.id order_id
FROM orders.sub_orders s
         JOIN orders.orders o
              ON s.order_id = o.id
WHERE s.id = @id;

-- name: UpdateOrderPaymentStatus :exec
UPDATE orders.orders
SET payment_status = $2,
    updated_at     = now()
WHERE id = $1;

-- name: UpdateOrderShippingStatus :exec
UPDATE orders.sub_orders
SET shipping_status = $2,
    updated_at      = now()
WHERE id = $1;

-- name: UpdateOrderShippingInfo :one
UPDATE orders.sub_orders
SET shipping_status = @shipping_status,
    tracking_number = @tracking_number,
    carrier         = @carrier,
    updated_at      = NOW()
WHERE id = @id
RETURNING *;
