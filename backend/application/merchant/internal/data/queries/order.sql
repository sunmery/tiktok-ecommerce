SET SEARCH_PATH to merchant,orders;

-- name: ListOrdersByUser :many
SELECT os.id AS sub_order_id,
       oo.id AS order_id,
       merchant_id,
       total_amount,
       os.currency,
       oo.payment_status,
       os.shipping_status,
       items,
       os.created_at,
       os.updated_at
FROM orders.sub_orders os
         JOIN orders.orders oo on os.order_id = oo.id
WHERE merchant_id = @merchant_id
ORDER BY created_at DESC
LIMIT @page_size OFFSET @page;

-- name: QuerySubOrders :many
SELECT s.order_id sub_orders_id,
       s.merchant_id,
       s.total_amount,
       s.currency,
       o.payment_status,
       s.shipping_status,
       s.items,
       s.created_at,
       s.updated_at
FROM orders.sub_orders s
         Join orders.orders o on s.order_id = o.id
WHERE order_id = @order_id
ORDER BY created_at;

-- name: CreateShip :one
INSERT INTO orders.shipping_info(id, merchant_id,sub_order_id, tracking_number, carrier, delivery,
                                 shipping_address, receiver_address, shipping_fee)
VALUES (@id, @merchant_id,@sub_order_id, @tracking_number, @carrier, @delivery,
        @shipping_address, @receiver_address, @shipping_fee)
RETURNING id, created_at;


-- name: UpdateOrderShippingStatus :exec
UPDATE orders.shipping_info
SET shipping_status = @shipping_status,
    updated_at      = now()
WHERE id = @id;

-- name: UpdateOrderShippingInfo :exec
UPDATE orders.sub_orders
SET shipping_status  = COALESCE(@shipping_status, shipping_status),
    tracking_number  = COALESCE(@tracking_number, tracking_number),
    carrier          = COALESCE(@carrier, carrier),
    merchant_address = COALESCE(@merchant_address, merchant_address),
    updated_at       = NOW()
WHERE id = @id;
