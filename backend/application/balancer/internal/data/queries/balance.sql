-- name: GetUserBalance :one
-- 获取指定用户和币种的余额信息
SELECT available, frozen, version, currency
FROM balances.user_balances
WHERE user_id = $1
  AND currency = $2;

-- name: CreateUserBalance :one
-- 为用户创建指定币种的初始余额记录 (通常在用户注册或首次涉及该币种时调用)
INSERT INTO balances.user_balances (user_id, currency, available, frozen, version)
VALUES ($1, $2, 0, 0, 0)
RETURNING *;

-- name: GetMerchantBalance :one
-- 获取指定商家和币种的余额信息
SELECT available, version, currency
FROM balances.merchant_balances
WHERE merchant_id = $1
  AND currency = $2;

-- name: CreateMerchantBalance :one
-- 为商家创建指定币种的初始余额记录
INSERT INTO balances.merchant_balances (merchant_id, currency, available, version, created_at, updated_at)
VALUES ($1, $2, 0, 0, NOW(), NOW())
RETURNING *;

-- name: IncreaseUserAvailableBalance :execrows
-- 增加用户可用余额 (用于充值成功, 取消提现, 取消冻结成功后资金退回) - 使用乐观锁
UPDATE balances.user_balances
SET available  = available + sqlc.arg(amount), -- 金额参数 (分)
    version    = version + 1,
    updated_at = NOW()
WHERE user_id = $1
  AND currency = $2
  AND version = sqlc.arg(expected_version);
-- 乐观锁检查

-- name: DecreaseUserAvailableBalance :execrows
-- 减少用户可用余额 (用于发起提现) - 使用乐观锁
UPDATE balances.user_balances
SET available  = available - sqlc.arg(amount), -- 金额参数 (分)
    version    = version + 1,
    updated_at = NOW()
WHERE user_id = $1
  AND currency = $2
  AND available >= sqlc.arg(amount) -- 确保余额充足
  AND version = sqlc.arg(expected_version);
-- 乐观锁检查

-- name: FreezeUserBalance :execrows
-- 冻结用户余额 (减少可用，增加冻结) - 使用乐观锁
UPDATE balances.user_balances
SET available  = available - sqlc.arg(amount), -- 金额参数 (分)
    frozen     = frozen + sqlc.arg(amount),
    version    = version + 1,
    updated_at = NOW()
WHERE user_id = $1
  AND currency = $2
  AND available >= sqlc.arg(amount) -- 确保可用余额充足
  AND version = sqlc.arg(expected_version);
-- 乐观锁检查

-- name: UnfreezeUserBalance :execrows
-- 取消冻结 (增加可用，减少冻结) - 使用乐观锁
UPDATE balances.user_balances
SET available  = available + sqlc.arg(amount), -- 金额参数 (分)
    frozen     = frozen - sqlc.arg(amount),
    version    = version + 1,
    updated_at = NOW()
WHERE user_id = $1
  AND currency = $2
  AND frozen >= sqlc.arg(amount) -- 确保冻结余额充足
  AND version = sqlc.arg(expected_version);
-- 乐观锁检查

-- name: ConfirmUserFreeze :execrows
-- 确认冻结 (仅减少冻结金额，资金将流向商家) - 使用乐观锁
UPDATE balances.user_balances
SET frozen     = frozen - sqlc.arg(amount), -- 金额参数 (分)
    version    = version + 1,
    updated_at = NOW()
WHERE user_id = $1
  AND currency = $2
  AND frozen >= sqlc.arg(amount) -- 确保冻结余额充足
  AND version = sqlc.arg(expected_version);
-- 乐观锁检查

-- name: IncreaseMerchantAvailableBalance :execrows
-- 增加商家可用余额 (用于确认转账成功) - 使用乐观锁
UPDATE balances.merchant_balances
SET available  = available + sqlc.arg(amount), -- 金额参数 (分)
    version    = version + 1,
    updated_at = NOW()
WHERE merchant_id = $1
  AND currency = $2
  AND version = sqlc.arg(expected_version); -- 乐观锁检查