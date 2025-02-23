SET search_path TO products;

-- 删除分片兼容索引
DROP INDEX IF EXISTS idx_audits_product;

-- 删除商品审核记录表
DROP TABLE IF EXISTS product_audits;

-- 删除商品属性表
DROP TABLE IF EXISTS product_attributes;

-- 删除商品图片表的唯一约束索引
DROP INDEX IF EXISTS idx_unique_primary_image;

-- 删除商品图片表
DROP TABLE IF EXISTS product_images;

-- 删除库存表
DROP TABLE IF EXISTS inventory;

-- 删除商品主表
DROP TABLE IF EXISTS products CASCADE;
