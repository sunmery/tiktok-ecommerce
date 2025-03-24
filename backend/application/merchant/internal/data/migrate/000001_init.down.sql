-- 删除索引
DROP INDEX IF EXISTS merchant.idx_stock_adjustments_created_at;
DROP INDEX IF EXISTS merchant.idx_stock_adjustments_merchant;
DROP INDEX IF EXISTS merchant.idx_stock_adjustments_product;
DROP INDEX IF EXISTS merchant.idx_stock_alerts_product_merchant;

-- 删除表
DROP TABLE IF EXISTS merchant.stock_adjustments;
DROP TABLE IF EXISTS merchant.stock_alerts;

-- 删除函数
DROP FUNCTION IF EXISTS merchant.uuidv7_sub_ms;

-- 删除 schema
DROP SCHEMA IF EXISTS merchant CASCADE;