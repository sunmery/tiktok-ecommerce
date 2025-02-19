CREATE SCHEMA IF NOT EXISTS orders;

CREATE TABLE orders.orders
(
    id             SERIAL PRIMARY KEY,
    owner          VARCHAR(100)    NOT NULL,    -- 订单所有者
    name           VARCHAR(100)    NOT NULL,    -- 订单名称
    email          VARCHAR(100)    NOT NULL,    -- 用户邮箱
    street_address VARCHAR(255)    NOT NULL,    -- 街道地址
    city           VARCHAR(100)    NOT NULL,    -- 城市
    state          VARCHAR(100)    NOT NULL,    -- 州/省
    country        VARCHAR(100)    NOT NULL,    -- 国家
    zip_code       INT             NOT NULL,    -- 邮政编码（数字类型）
    currency       CHAR(3)         NOT NULL,    -- 货币类型, ISO 4217 标准
    status         VARCHAR(10)     DEFAULT 'pending' NOT NULL, -- 订单状态: pending, paid, cancelled
    created_at     timestamptz     NOT NULL,    -- 创建时间
    updated_at     timestamptz     DEFAULT now() NOT NULL     -- 更新时间
);

CREATE TABLE orders.order_items
(
    id         SERIAL PRIMARY KEY,
    order_id   INT           NOT NULL REFERENCES orders.orders (id) ON DELETE CASCADE,
    product_id INT           NOT NULL,         -- 商品ID
    name       VARCHAR(100)  NOT NULL,         -- 商品名称
    quantity   INT           NOT NULL,         -- 商品数量
    price      REAL NOT NULL          -- 商品单价
);
