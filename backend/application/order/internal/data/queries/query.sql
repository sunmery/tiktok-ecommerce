-- name: CreateOrder :one
INSERT INTO orders.orders ( user_id, currency, street_address,
                           city, state, country, zip_code, email)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetOrderByID :one
SELECT *
FROM orders.orders
WHERE id = $1;

-- -- name: ListOrdersByUser :many
-- SELECT *
-- FROM orders.orders
-- WHERE user_id = $1
-- ORDER BY created_at DESC
-- LIMIT $2 OFFSET $3;

-- name: GetDateRangeStats :one
SELECT *
FROM orders.get_date_range_stats(
        p_user_id => $1,
        p_start => $2,
        p_end => $3
     );

-- 带日期过滤的查询
-- name: ListOrdersByUserWithDate :many
SELECT o.*,
       json_agg(so.*) AS sub_orders
FROM orders.orders o
         LEFT JOIN orders.sub_orders so ON o.id = so.order_id
WHERE o.user_id = @user_id::uuid
  AND o.created_at BETWEEN @start_time::timestamptz AND @end_time::timestamptz
GROUP BY o.id
ORDER BY o.created_at DESC
LIMIT @limits OFFSET @offsets;

-- name: CreateSubOrder :one
INSERT INTO orders.sub_orders (order_id, merchant_id, total_amount,
                               currency, status, items)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: UpdateSubOrderStatus :exec
UPDATE orders.sub_orders
SET status     = $2,
    updated_at = $3
WHERE id = $1;
