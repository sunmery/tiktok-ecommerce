-- 删除索引
DROP INDEX IF EXISTS payments.idx_payments_status;
DROP INDEX IF EXISTS payments.idx_payments_order;

-- 删除 payments 表
DROP TABLE IF EXISTS payments.payments;

-- 删除 uuidv7_sub_ms 函数
DROP FUNCTION IF EXISTS payments.uuidv7_sub_ms();

-- 删除 payments 模式
DROP SCHEMA IF EXISTS payments CASCADE;
