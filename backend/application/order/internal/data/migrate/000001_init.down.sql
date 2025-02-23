-- 删除统计函数
DROP FUNCTION IF EXISTS orders.get_date_range_stats(UUID, TIMESTAMPTZ, TIMESTAMPTZ);

-- 删除索引
DROP INDEX IF EXISTS orders.idx_sub_orders_merchant;
DROP INDEX IF EXISTS orders.idx_orders_user;

-- 删除子订单表
DROP TABLE IF EXISTS orders.sub_orders;

-- 删除主订单表
DROP TABLE IF EXISTS orders.orders;

-- 删除 UUIDv7 生成函数
DROP FUNCTION IF EXISTS orders.uuidv7_sub_ms();

-- 重置搜索路径
RESET SEARCH_PATH;

-- 删除 orders 模式
DROP SCHEMA IF EXISTS orders CASCADE;
