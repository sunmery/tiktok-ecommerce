-- 删除 cart_items 表
DROP TABLE IF EXISTS carts.cart_items;

-- 删除 cart 表
DROP TABLE IF EXISTS carts.cart;

-- 删除 carts schema
DROP SCHEMA IF EXISTS carts CASCADE;


-- 使用迁移工具可以不加事务但是如果使用脚本的话需要加事务
