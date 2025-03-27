-- name: CreatePaymentQuery :one
INSERT INTO payments.payments (id, order_id, user_id, amount, currency, method, status,
                               subject, trade_no, gateway_tx_id, pay_url, metadata)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
RETURNING *;

-- name: UpdateStatusQuery :one
UPDATE payments.payments
SET status        = $2,
    gateway_tx_id = $3,
    updated_at    = now()
WHERE id = $1
RETURNING *;

-- name: GetByIDQuery :one
SELECT *
FROM payments.payments
WHERE id = $1;

-- name: GetByOrderIDQuery :one
SELECT *
FROM payments.payments
WHERE order_id = $1;

-- name: GetPayment :one
SELECT *
FROM payments.payments
WHERE id = $1;

-- name: GetPaymentByOrderID :one
SELECT *
FROM payments.payments
WHERE order_id = $1;

-- name: GetPaymentByTradeNo :one
SELECT *
FROM payments.payments
WHERE trade_no = $1;

-- name: UpdatePaymentStatus :one
UPDATE payments.payments
SET 
    status = $4,
    gateway_tx_id = $3,
    updated_at = now()
WHERE (id = $1 AND $1 != 0) OR (order_id = $2 AND $2 != 0)
RETURNING *;

-- name: UpdatePaymentStatusByID :one
UPDATE payments.payments
SET status        = $2,
    gateway_tx_id = $3,
    updated_at    = now()
WHERE id = $1
RETURNING *;

-- name: UpdatePaymentStatusByOrderID :one
UPDATE payments.payments
SET status        = $2,
    gateway_tx_id = $3,
    updated_at    = now()
WHERE order_id = $1
RETURNING *;
