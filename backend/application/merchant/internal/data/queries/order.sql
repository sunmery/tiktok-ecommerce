SET SEARCH_PATH to merchant,orders;

-- name: ListOrdersByUser :many
SELECT oo.id AS order_id,
       os.merchant_id,
       total_amount,
       os.currency,
       os.shipping_status,
       oo.payment_status,
       si.sub_order_id,
       si.tracking_number,
       si.carrier,
       si.shipping_status,
       si.delivery,
       si.shipping_address,
       si.receiver_address,
       si.shipping_fee,
       si.created_at,
       items
FROM orders.sub_orders os
         JOIN orders.orders oo on os.order_id = oo.id
         LEFT JOIN orders.shipping_info si on os.id = si.sub_order_id
WHERE os.merchant_id = @merchant_id
ORDER BY si.created_at DESC
LIMIT @page_size OFFSET @page;

-- name: QuerySubOrders :many
SELECT os.id AS sub_order_id,
       merchant_id,
       total_amount,
       oo.currency,
       status,
       items,
       oo.created_at,
       oo.updated_at,
       oo.payment_status,
       os.shipping_status
FROM orders.sub_orders os
         Join orders.orders oo on os.order_id = oo.id
WHERE order_id = @order_id
ORDER BY created_at;

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
