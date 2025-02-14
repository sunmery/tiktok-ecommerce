CREATE SCHEMA IF NOT EXISTS cart_schema;

-- 在 cart_schema 下创建 cart 表
CREATE TABLE IF NOT EXISTS cart_schema.cart (
    cart_id SERIAL PRIMARY KEY,                    -- 购物车唯一ID
    owner VARCHAR(100) NOT NULL,                   -- 用户组织
    name VARCHAR(100) NOT NULL,                    -- 用户地址
    cart_name VARCHAR(100) NOT NULL DEFAULT 'cart',  -- 购物车名称，非空且默认值为"cart"
    status VARCHAR(50) NOT NULL DEFAULT 'active',  -- 购物车状态
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,  -- 创建时间
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,  -- 更新时间
    CONSTRAINT unique_owner_name_cart_name UNIQUE(owner, name, cart_name)  -- 保证 owner、name 和 cart_name 的组合唯一
);

-- 为 owner、name 和 cart_name 字段创建联合索引
CREATE INDEX idx_owner_name_cart_name ON cart_schema.cart (owner, name, cart_name);

-- 在 cart_schema 下创建 cart_items 表
CREATE TABLE IF NOT EXISTS cart_schema.cart_items (
    cart_item_id SERIAL PRIMARY KEY,              -- 购物车商品项唯一ID
    cart_id INT NOT NULL,                         -- 购物车ID
    product_id INT NOT NULL,                      -- 商品ID
    quantity INT NOT NULL CHECK (quantity > 0),    -- 商品数量
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,  -- 创建时间
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,  -- 更新时间
    CONSTRAINT unique_cart_product UNIQUE(cart_id, product_id)  -- 保证每个购物车商品的唯一性
);

-- 为 cart_id 和 product_id 字段创建联合索引
CREATE INDEX idx_cart_id_product_id ON cart_schema.cart_items (cart_id, product_id);