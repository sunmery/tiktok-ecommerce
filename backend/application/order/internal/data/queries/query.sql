-- name: UpdatePaymentStatus :one
SELECT id, user_id, payment_status
FROM orders.orders
WHERE id = @id
    FOR UPDATE;

-- name: MarkOrderAsPaid :one
UPDATE orders.orders
SET payment_status = @payment_status,
    updated_at     = now()
WHERE id = $2
RETURNING *;

-- name: MarkSubOrderAsPaid :one
UPDATE orders.sub_orders
SET status     = @status,
    updated_at = now()
WHERE order_id = @order_id
RETURNING *;
