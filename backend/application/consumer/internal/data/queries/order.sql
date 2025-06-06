-- name: GetConsumerOrders :many
SELECT oo.id     AS order_id,
       os.id     AS sub_order_id,
       os.total_amount,
       os.currency,
       os.status AS payment_status,
       os.shipping_status,
       os.items,
       oo.email,
       oo.street_address,
       oo.city,
       oo.state,
       oo.country,
       oo.zip_code,
       os.created_at,
       os.updated_at
FROM orders.orders oo
         LEFT JOIN orders.sub_orders os
                   ON oo.id = os.order_id
WHERE oo.user_id = @user_id
GROUP BY oo.id, os.id, os.total_amount, os.currency, os.status, os.shipping_status, os.items, oo.email
LIMIT @page_size OFFSET @page;

-- name: CreateOrder :one
INSERT INTO orders.orders (id, user_id, currency, street_address,
                           city, state, country, zip_code, email)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING id, user_id, currency, street_address,
    city, state, country, zip_code, email;

-- name: GetOrderByID :one
SELECT oo.id     AS order_id,
       os.id     AS sub_order_id,
       os.total_amount,
       os.currency,
       os.status AS payment_status,
       os.shipping_status,
       os.items,
       oo.email,
       oo.street_address,
       oo.city,
       oo.state,
       oo.country,
       oo.zip_code,
       os.created_at,
       os.updated_at
FROM orders.orders oo
         JOIN orders.sub_orders os
              ON oo.id = os.order_id
WHERE oo.user_id = @user_id
  AND os.id = @sub_order_id
GROUP BY oo.id, os.id, os.total_amount, os.currency, os.status, os.shipping_status, os.items, oo.email;

-- name: GetSubOrderByID :one
SELECT o.*,
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
FROM orders.orders o
         LEFT JOIN orders.sub_orders os ON o.id = os.order_id
WHERE o.user_id = @user_id::uuid
  AND os.id = @order_id
GROUP BY o.id;

-- name: GetOrderByUserID :one
SELECT os.id,
       os.order_id,
       os.merchant_id,
       os.total_amount,
       os.currency,
       os.status,
       os.items,
       os.created_at,
       os.updated_at,
       oo.payment_status,
       os.shipping_status
FROM orders.sub_orders os
         JOIN orders.orders oo
              ON os.order_id = oo.id
WHERE user_id = @user_id;

-- name: GetConsumerOrdersWithSuborders :many
SELECT o.id,
       o.street_address,
       o.city,
       o.state,
       o.country,
       o.zip_code,
       o.email,
       o.payment_status,
       os.shipping_status,
       os.merchant_id,
       os.id AS sub_order_id,
       os.total_amount,
       os.currency,
       os.items,
       o.created_at,
       o.updated_at
FROM orders.orders o
         LEFT JOIN orders.sub_orders os ON o.id = os.order_id
WHERE o.user_id = @user_id
  AND o.id = @order_id
GROUP BY o.id, os.id, o.currency, os.merchant_id, o.street_address, o.city, o.state, o.country, o.zip_code, o.email,
         o.created_at,
         os.shipping_status
ORDER BY o.created_at DESC;

-- name: QuerySubOrders :many
SELECT os.id,
       os.order_id,
       os.merchant_id,
       os.total_amount,
       os.currency,
       os.status,
       os.items,
       os.created_at,
       os.updated_at,
       oo.payment_status,
       os.shipping_status
FROM orders.sub_orders os
         Join orders.orders oo on os.order_id = oo.id
WHERE order_id = @order_id
ORDER BY created_at;

-- name: CreateSubOrder :one
INSERT INTO orders.sub_orders (id, order_id, merchant_id, total_amount,
                               currency, status, items)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, order_id, merchant_id, total_amount,
    currency, status, items;

-- -- name: UpdateSubOrderStatus :exec
-- UPDATE orders.sub_orders
-- SET status     = @status,
--     updated_at = NOW()
-- WHERE id = @id;

-- name: CreateOrderShipping :exec
INSERT INTO orders.shipping_info(id, merchant_id, sub_order_id, shipping_status, tracking_number, carrier, delivery,
                                 shipping_address, receiver_address, shipping_fee)
VALUES (@id, @merchant_id, @sub_order_id, @shipping_status, @tracking_number, @carrier, @delivery,
        @shipping_address, @receiver_address, @shipping_fee);

-- name: GetShipOrderStatus :one
SELECT id,
       sub_order_id,
       tracking_number,
       carrier,
       shipping_status,
       delivery,
       shipping_address,
       receiver_address,
       shipping_fee,
       created_at,
       updated_at
FROM orders.shipping_info
WHERE sub_order_id = @id;

-- name: UpdateOrderShippingStatus :exec
WITH update_shipinfo_status AS (
    UPDATE orders.shipping_info
        SET shipping_status = @shipping_status
        WHERE sub_order_id = @sub_order_id)
UPDATE orders.sub_orders
SET shipping_status = @shipping_status
WHERE id = @sub_order_id;
