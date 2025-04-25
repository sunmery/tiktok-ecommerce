-- name: CreateTransaction :one
-- 创建交易流水记录
INSERT INTO balances.transactions (id,
                                   type, amount, currency, from_user_id, to_merchant_id,
                                   payment_method_type, payment_account, payment_extra, status,
                                   created_at, updated_at
    -- freeze_id 可以在创建支付流水时关联
    -- , freeze_id
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW(), NOW()
           -- , $10
       )
RETURNING id;
-- 返回交易记录的 ID

-- name: UpdateTransactionStatus :execrows
-- 更新交易流水状态
UPDATE balances.transactions
SET status     = $1,
    updated_at = NOW()
WHERE id = $2;

-- name: GetMerchantTransactions :many
-- 根据 商家ID 获取交易流水记录
SELECT *
FROM balances.transactions
WHERE to_merchant_id = @merchant_id
  AND currency = COALESCE(@currency, currency)
  AND status = COALESCE(@status, status)
LIMIT @page_size OFFSET @page;

-- name: GetConsumerTransactions :many
-- 根据 用户ID 获取交易流水记录
SELECT *
FROM balances.transactions
WHERE from_user_id = @user_id
  AND currency = COALESCE(@currency, currency)
  AND status = COALESCE(@status, status)
LIMIT @page_size OFFSET @page;


-- name: GetUserPaymentMethod :one
-- 获取用户支付方式详情 (可能在提现时需要)
SELECT *
FROM balances.user_payment_methods
WHERE id = $1
  AND user_id = $2;

-- name: GetMerchantPaymentMethod :one
-- 获取商家支付方式详情 (未来可能用于商家提现)
SELECT *
FROM balances.merchant_payment_methods
WHERE id = $1
  AND merchant_id = $2;