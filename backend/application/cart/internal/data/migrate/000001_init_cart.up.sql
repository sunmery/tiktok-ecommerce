CREATE SCHEMA IF NOT EXISTS cart_schema;

-- 在 cart_schema 下创建 cart 表
CREATE TABLE IF NOT EXISTS cart_schema.cart (
    cart_id SERIAL PRIMARY KEY,                    -- 购物车唯一ID
    user_id VARCHAR(100) NOT NULL,                 -- 用户ID
    cart_name VARCHAR(100) NOT NULL DEFAULT 'cart',  -- 购物车名称，非空且默认值为"cart"
    status VARCHAR(50) NOT NULL DEFAULT 'active',  -- 购物车状态
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,  -- 创建时间
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,  -- 更新时间
    CONSTRAINT unique_user_cart_name UNIQUE(user_id, cart_name)  -- 保证 user_id 和 cart_name 的组合唯一
);

-- 为 user_id 和 cart_name 字段创建联合索引
CREATE INDEX idx_user_id_cart_name ON cart_schema.cart (user_id, cart_name);

-- 在 cart_schema 下创建 cart_items 表
CREATE TABLE IF NOT EXISTS cart_schema.cart_items (
    cart_item_id SERIAL PRIMARY KEY,              -- 购物车商品项唯一ID
    cart_id INT NOT NULL,                         -- 购物车ID
    merchant_id VARCHAR(100) NOT NULL,            -- 商家ID
    product_id INT NOT NULL,                      -- 商品ID
    quantity INT NOT NULL CHECK (quantity > 0),   -- 商品数量
    selected BOOLEAN NOT NULL DEFAULT TRUE,       -- 是否选中，默认为选中，且不能为空
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,  -- 创建时间
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,  -- 更新时间
    CONSTRAINT unique_cart_merchant_product UNIQUE(cart_id, merchant_id, product_id)  -- 保证每个购物车商品的唯一性
);

-- 为 cart_id、merchant_id 和 product_id 字段创建联合索引
CREATE INDEX idx_cart_id_merchant_id_product_id ON cart_schema.cart_items (cart_id, merchant_id, product_id);

-- 使用迁移工具可以不加事务但是如果使用脚本的话需要加事务