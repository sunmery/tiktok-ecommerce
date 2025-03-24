-- name: UpdateOrderPaymentStatus :exec
UPDATE orders.orders
SET payment_status = $2,
    updated_at = now()
WHERE id = $1;