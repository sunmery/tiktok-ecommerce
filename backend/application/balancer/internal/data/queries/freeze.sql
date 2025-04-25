-- name: CreateFreeze :one
-- 创建冻结记录
INSERT INTO balances.balance_freezes (user_id, order_id, currency, amount, status, expires_at, created_at, updated_at)
VALUES ($1, $2, $3, $4, 'FROZEN', $5, NOW(), NOW())
RETURNING id; -- 返回冻结记录的 ID

-- name: GetFreeze :one
-- 根据 ID 获取冻结记录
SELECT * FROM balances.balance_freezes
WHERE id = $1;

-- name: GetFreezeByOrderForUser :one
-- 根据用户 ID 和订单 ID 获取冻结记录 (假设一个订单只有一个冻结记录)
SELECT * FROM balances.balance_freezes
WHERE user_id = $1 AND order_id = $2;

-- name: UpdateFreezeStatus :execrows
-- 更新冻结记录状态 (例如: FROZEN -> CONFIRMED 或 FROZEN -> CANCELED)
UPDATE balances.balance_freezes
SET status = @status,
    updated_at = NOW()
WHERE id = @id
  AND status = @current_status; -- 确保当前状态是预期的状态 (例如 'FROZEN')

-- name: GetExpiredFreezes :many
-- 获取所有已过期但仍处于冻结状态的记录 (用于定时任务处理)
SELECT * FROM balances.balance_freezes
WHERE status = 'FROZEN'
  AND expires_at < NOW();