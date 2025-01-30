-- 切换到正确的 schema
SET search_path TO addresses;

-- 删除外键约束
ALTER TABLE addresses
    DROP CONSTRAINT IF EXISTS addresses_user_id_fkey;

-- 删除索引
DROP INDEX IF EXISTS idx_addresses_user_id;

-- 删除表
DROP TABLE IF EXISTS addresses;

-- 如果需要，也可以删除 schema
DROP SCHEMA IF EXISTS addresses;
