-- name: CreatePaymentQuery :one
INSERT INTO payments.payments (payment_id, order_id, amount, currency, method, status,
                               gateway_tx_id, metadata)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *;

-- name: UpdateStatusQuery :one
UPDATE payments.payments
SET status        = $2,
    gateway_tx_id = $3,
    updated_at    = $4
WHERE payment_id = $1 RETURNING *;

-- name: GetByIDQuery :one
SELECT *
FROM payments.payments
WHERE payment_id = $1;

-- name: GetByOrderIDQuery :one
SELECT *
FROM payments.payments
WHERE order_id = $1;
