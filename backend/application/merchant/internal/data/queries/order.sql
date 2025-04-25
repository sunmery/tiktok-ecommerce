CREATE SCHEMA IF NOT EXISTS orders;
SET SEARCH_PATH to orders;

-- name: UpdateOrderShippingStatus :exec
WITH update_shipping_info_ship_status AS (
    UPDATE orders.shipping_info
        SET shipping_status = @shipping_status,
            updated_at = now()
        WHERE sub_order_id = @sub_order_id)
UPDATE orders.sub_orders
SET shipping_status = @shipping_status,
    updated_at      = now()
WHERE id = @sub_order_id;

-- name: GetMerchantOrders :many
SELECT oo.id,
       oo.payment_status,
       oo.user_id,
       oo.currency,
       oo.street_address,
       oo.city,
       oo.state,
       oo.country,
       oo.zip_code,
       oo.email,
       os.order_id        AS order_id,
       os.merchant_id,
       os.total_amount,
       os.items,
       os.shipping_status AS shipping_status,
       os.created_at,
       os.updated_at
FROM orders.sub_orders os
         JOIN orders.orders oo on os.order_id = oo.id
WHERE os.merchant_id = @merchant_id
ORDER BY os.created_at DESC
LIMIT @page_size OFFSET @page;

-- 通过子订单 ID 去查询主订单的地址, 因为用户可能下单多个商品, 分别属于不同商家, 但地址并不会变化
-- name: GetConsumerAddress :one
SELECT o.id,
       user_id,
       street_address,
       city,
       state,
       country,
       zip_code,
       email,
       payment_status,
       s.created_at,
       s.updated_at,
       s.shipping_status
FROM orders.sub_orders s
         Join orders.orders o on s.order_id = o.id
WHERE s.id = @id;

-- name: CreateShip :one
INSERT INTO orders.shipping_info(id, merchant_id, sub_order_id, shipping_status, tracking_number, carrier, delivery,
                                 shipping_address, receiver_address, shipping_fee)
VALUES (@id, @merchant_id, @sub_order_id, @shipping_status, @tracking_number, @carrier, @delivery,
        @shipping_address, @receiver_address, @shipping_fee)
RETURNING id, created_at;
