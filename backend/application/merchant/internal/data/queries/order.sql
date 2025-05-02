CREATE SCHEMA IF NOT EXISTS orders;
SET SEARCH_PATH to orders;

-- name: UpdateOrderShippingStatus :one
-- WITH update_shipping_info_ship_status AS (
--     UPDATE orders.shipping_info
--         SET shipping_status = @shipping_status,
--             updated_at = now()
--         WHERE sub_order_id = @sub_order_id)
-- UPDATE orders.sub_orders
-- SET shipping_status = @shipping_status,
--     updated_at      = now()
-- WHERE id = @sub_order_id;
WITH update_shipping_info_ship_status AS (
    UPDATE orders.shipping_info
        SET
            merchant_id = @merchant_id,
            shipping_status = @shipping_status,
            tracking_number = @tracking_number,
            carrier = @carrier,
            receiver_address = @receiver_address,
            shipping_address = @shipping_address,
            shipping_fee = @shipping_fee,
            updated_at = now()
        WHERE sub_order_id = @sub_order_id)
UPDATE orders.sub_orders
SET shipping_status = @shipping_status,
    updated_at      = now()
WHERE id = @sub_order_id
RETURNING id, updated_at;


-- name: GetMerchantByOrderId :one
SELECT os.merchant_id
FROM orders.sub_orders os
         JOIN orders.orders o on os.order_id = o.id
WHERE o.id = @id;

-- name: GetMerchantOrders :many
SELECT os.order_id,
       oo.created_at,
       json_agg(
               item::jsonb ||
               jsonb_build_object(
                       'subOrderId', os.id,
                       'userId', oo.user_id,
                       'email', oo.email,
                       'totalAmount', os.total_amount,
                       'createdAt', os.created_at,
                       'updatedAt', os.updated_at,
                       'paymentStatus', os.status,
                       'shippingStatus', os.shipping_status,
                       'currency', oo.currency,
                       'address', json_build_object(
                               'streetAddress', oo.street_address,
                               'city', oo.city,
                               'state', oo.state,
                               'country', oo.country,
                               'zipCode', oo.zip_code
                                  )
               )
       ) AS items
FROM orders.sub_orders os
         JOIN orders.orders oo ON os.order_id = oo.id,
     json_array_elements(os.items::json) AS item
WHERE os.merchant_id = @merchant_id
GROUP BY os.order_id, os.merchant_id, oo.created_at
ORDER BY oo.created_at DESC
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
-- 创建货运信息
INSERT INTO orders.shipping_info(id, merchant_id, sub_order_id, shipping_status, tracking_number, carrier, delivery,
                                 shipping_address, receiver_address, shipping_fee)
VALUES (@id, @merchant_id, @sub_order_id, @shipping_status, @tracking_number, @carrier, @delivery,
        @shipping_address, @receiver_address, @shipping_fee)
RETURNING id, created_at;
