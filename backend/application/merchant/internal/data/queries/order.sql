SET SEARCH_PATH to merchant,orders;

-- name: ListOrdersByUser :many
SELECT s.id,
       order_id,
       merchant_id,
       total_amount,
       s.currency,
       o.payment_status,
       s.shipping_status,
       items,
       s.created_at,
       s.updated_at
FROM orders.sub_orders s
         JOIN orders.orders o on s.order_id = o.id
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