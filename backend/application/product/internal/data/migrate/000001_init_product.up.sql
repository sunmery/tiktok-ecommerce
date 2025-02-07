-- 手动在数据库中创建
CREATE SCHEMA IF NOT EXISTS products;

CREATE TABLE IF NOT EXISTS products.products
(
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(50)  NOT NULL,
    description TEXT         NOT NULL,
    picture     TEXT         NOT NULL,
    price       REAL         NOT NULL,
    categories  TEXT[]  NOT NULL
);
--
-- CREATE TABLE IF NOT EXISTS products.categories(
--    categories  jsonb[]
-- );
--
-- CREATE TABLE IF NOT EXISTS products.category_relate(
--     id serial primary key,
--     product_id int not null ,
--     category_id  int not null ,
--     created_at timestamptz default now() not null
-- );

-- 加快通过名称查询的商品
CREATE INDEX idx_products_name ON products.products (name);

-- GIN 索引会显著加速数组查询
CREATE INDEX idx_products_categories ON products.products USING GIN (categories);

-- 查询某类别的商品
-- SELECT * FROM products.products WHERE categories @> ARRAY['electronics'];
