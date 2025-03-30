SET SEARCH_PATH to merchant,orders;

-- name: ListOrdersByUser :many
SELECT id,
       user_id,
       currency,
       street_address,
       city,
       state,
       country,
       zip_code,
       email,
       created_at,
       updated_at,
       payment_status
FROM orders.orders
WHERE user_id = @user_id
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