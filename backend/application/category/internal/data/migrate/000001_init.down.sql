-- 删除闭包表
DROP TABLE IF EXISTS categories.category_closure;

-- 删除分类表
DROP TABLE IF EXISTS categories.categories;

-- 删除序列
DROP SEQUENCE IF EXISTS categories.categories_id_seq;

-- 删除 ltree 扩展
DROP EXTENSION IF EXISTS ltree;

-- 删除模式（如果为空）
DROP SCHEMA IF EXISTS categories cascade;
