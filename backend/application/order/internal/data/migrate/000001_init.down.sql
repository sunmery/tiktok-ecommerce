-- migrations/000001_create_orders_tables.down.sql

-- 逆向操作需要按依赖顺序删除表
DROP TABLE IF EXISTS orders.sub_orders;
DROP TABLE IF EXISTS orders.orders;
