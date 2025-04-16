SET SEARCH_PATH to merchant,orders;

-- name: ListOrdersByUser :many
SELECT id,
       order_id,
       merchant_id,
       total_amount,
       currency,
       status,
       items,
       created_at,
       updated_at
FROM orders.sub_orders
WHERE merchant_id = @merchant_id
ORDER BY created_at DESC
LIMIT @page_size OFFSET @page;

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
WHERE order_id = @order_id
ORDER BY created_at;