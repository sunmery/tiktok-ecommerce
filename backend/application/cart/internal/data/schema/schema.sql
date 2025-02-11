CREATE SCHEMA IF NOT EXISTS cart_schema;

-- 在 cart_schema 下创建 cart 表
CREATE TABLE IF NOT EXISTS cart_schema.cart (
    cart_id SERIAL PRIMARY KEY,                    -- 购物车唯一ID
    user_id INT NOT NULL,                          -- 用户ID
    status VARCHAR(50) NOT NULL DEFAULT 'active',  -- 购物车状态
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,  -- 创建时间
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP   -- 更新时间
);

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
-- 使用迁移工具可以不加事务但是如果使用脚本的话需要加事务