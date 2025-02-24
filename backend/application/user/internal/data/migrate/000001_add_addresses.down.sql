-- 切换到正确的 schema
SET search_path TO addresses;

-- 删除索引
DROP INDEX IF EXISTS idx_addresses_user_id;

-- 删除表
DROP TABLE IF EXISTS addresses CASCADE;
