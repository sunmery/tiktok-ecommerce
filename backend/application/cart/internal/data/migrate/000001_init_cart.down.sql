-- 删除 cart_items 表
DROP TABLE IF EXISTS cart_schema.cart_items;

-- 删除 cart 表
DROP TABLE IF EXISTS cart_schema.cart;

-- 删除 cart_schema schema
DROP SCHEMA IF EXISTS cart_schema;

-- 使用迁移工具可以不加事务但是如果使用脚本的话需要加事务