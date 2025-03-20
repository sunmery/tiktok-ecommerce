CREATE TABLE IF NOT EXISTS carts.cart
(
    cart_id    SERIAL PRIMARY KEY,                               -- 购物车唯一ID
    user_id    uuid         NOT NULL,                            -- 用户ID
    cart_name  VARCHAR(100) NOT NULL DEFAULT 'cart',             -- 购物车名称，非空且默认值为"cart"
    status     VARCHAR(50)  NOT NULL DEFAULT 'active',           -- 购物车状态
    created_at TIMESTAMP             DEFAULT now(),              -- 创建时间
    updated_at TIMESTAMP             DEFAULT now(),              -- 更新时间
    CONSTRAINT unique_user_cart_name UNIQUE (user_id, cart_name) -- 保证 user_id 和 cart_name 的组合唯一
);

CREATE TABLE IF NOT EXISTS carts.cart_items
(
    cart_item_id uuid PRIMARY KEY DEFAULT uuidv7_sub_ms(),                            -- 购物车商品项唯一ID
    cart_id      INT            NOT NULL,                                             -- 购物车ID
    merchant_id  uuid           NOT NULL,                                             -- 商家ID
    product_id   uuid           NOT NULL,                                             -- 商品ID
    quantity     INT            NOT NULL CHECK (quantity > 0),                        -- 商品数量
    created_at   TIMESTAMP        DEFAULT now(),                                      -- 创建时间
    updated_at   TIMESTAMP        DEFAULT now(),                                      -- 更新时间
    CONSTRAINT unique_cart_merchant_product UNIQUE (cart_id, merchant_id, product_id) -- 保证每个购物车商品的唯一性
);
