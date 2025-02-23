-- name: CreatePayRecord :one
INSERT INTO pay_record (
    user_id,
    order_id,
    transcation_id,
    amount,
    pay_at,
    status
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetPayRecordsByUserId :many
SELECT *
FROM pay_record
WHERE user_id = $1 AND deleted_at IS NULL
ORDER BY created_at DESC;

-- name: GetPayRecordByOrderId :one
SELECT *
FROM pay_record
WHERE order_id = $1 AND deleted_at IS NULL
LIMIT 1;

-- name: GetPayRecordByTransactionId :one
SELECT *
FROM pay_record
WHERE transcation_id = $1 AND deleted_at IS NULL
LIMIT 1;
