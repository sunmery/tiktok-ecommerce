SET search_path TO carts;
CREATE SCHEMA IF NOT EXISTS carts;

CREATE FUNCTION uuidv7_sub_ms() RETURNS uuid
AS
$$
select encode(
               substring(int8send(floor(t_ms)::int8) from 3) ||
               int2send((7 << 12)::int2 | ((t_ms - floor(t_ms)) * 4096)::int2) ||
               substring(uuid_send(gen_random_uuid()) from 9 for 8)
           , 'hex')::uuid
from (select extract(epoch from clock_timestamp()) * 1000 as t_ms) s
$$ LANGUAGE sql volatile;

CREATE TABLE IF NOT EXISTS carts.cart
(
    cart_id    SERIAL PRIMARY KEY,                               -- 购物车唯一ID
    user_id    uuid         NOT NULL,                            -- 用户ID
    cart_name  VARCHAR(100) NOT NULL DEFAULT 'cart',             -- 购物车名称，非空且默认值为"cart"
    created_at TIMESTAMP             DEFAULT now(),              -- 创建时间
    updated_at TIMESTAMP             DEFAULT now(),              -- 更新时间
    CONSTRAINT unique_user_cart_name UNIQUE (user_id, cart_name) -- 保证 user_id 和 cart_name 的组合唯一
);
COMMENT ON TABLE carts.cart IS '购物车表';

-- 为 user_id 和 cart_name 字段创建联合索引
CREATE INDEX idx_user_id_cart_name ON carts.cart (user_id, cart_name);

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
COMMENT ON TABLE carts.cart IS '购物车商品项表';

-- 为 cart_id、merchant_id 和 product_id 字段创建联合索引
CREATE INDEX idx_cart_id_merchant_id_product_id ON carts.cart_items (cart_id, merchant_id, product_id);
