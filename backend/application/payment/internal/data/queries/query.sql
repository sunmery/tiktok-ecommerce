-- name: CreatePaymentQuery :one
INSERT INTO payments.payments (id, order_id, consumer_id, amount, currency, method, status,
                               subject, trade_no, freeze_id,consumer_version, merchant_version)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
RETURNING *;

-- name: GetPaymentByTradeNo :one
-- 根据商户订单号查询支付记录
SELECT *
FROM payments.payments
WHERE trade_no = $1;

-- name: UpdateStatusQuery :one
UPDATE payments.payments
SET status     = $2,
    id         = $3,
    updated_at = now()
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

-- name: UpdatePaymentStatus :one
UPDATE payments.payments
SET status     = @status,
    updated_at = now()
WHERE (id = @id AND @id != 0)
   OR (order_id = @order_id AND @order_id != 0)
RETURNING *;
