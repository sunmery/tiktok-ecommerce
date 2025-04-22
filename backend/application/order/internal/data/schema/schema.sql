CREATE SCHEMA IF NOT EXISTS orders;
SET SEARCH_PATH TO orders;

-- 创建主订单表（用于汇总用户的所有子订单）
CREATE TABLE orders.orders
(
    id             SERIAL PRIMARY KEY,
    user_id        UUID                      NOT NULL, -- 关联用户ID
    currency       VARCHAR(3)  DEFAULT 'CNY' NOT NULL, -- 用户下单时使用的货币类型（ISO 4217）
    street_address TEXT                      NOT NULL, -- 反范式化存储地址信息，避免关联查询
    city           VARCHAR(100)              NOT NULL,
    state          VARCHAR(100)              NOT NULL,
    country        VARCHAR(100)              NOT NULL,
    zip_code       VARCHAR(10)               NOT NULL,
    payment_status   VARCHAR(20)    NOT NULL DEFAULT 'pending'
        CHECK (payment_status IN ('pending', 'paid', 'cancelled', 'failed', 'cancelled')),
    email          VARCHAR(320)              NOT NULL, -- 支持最大邮箱长度
    created_at     timestamptz DEFAULT now() NOT NULL, -- Unix时间戳，避免时区问题
    updated_at     timestamptz DEFAULT now() NOT NULL
);

-- 创建子订单表（按商家分单）
CREATE TABLE orders.sub_orders
(
    id               SERIAL PRIMARY KEY,                                 -- 子订单
    order_id         BIGINT         NOT NULL,                            -- 关联主订单ID（程序级外键）
    merchant_id      UUID           NOT NULL,                            -- 商家ID（来自商家服务）
    total_amount     NUMERIC(12, 2) NOT NULL,                            -- 精确金额计算（整数部分10位，小数2位）
    currency         VARCHAR(3)     NOT NULL,                            -- 实际结算货币
    status           VARCHAR(20)    NOT NULL,                            -- 订单状态
    items            JSONB          NOT NULL,                            -- 订单项快照（包含商品详情和当时价格）
    shipping_status VARCHAR(20) NOT NULL DEFAULT 'PENDING_SHIPMENT',
        CHECK (shipping_status IN ('PENDING_SHIPMENT', 'SHIPPED', 'IN_TRANSIT', 'DELIVERED', 'CONFIRMED', 'CANCELLED')),
    merchant_address JSONB,
    created_at       timestamptz             DEFAULT now() NOT NULL,
    updated_at       timestamptz             DEFAULT now() NOT NULL,
    tracking_number  VARCHAR(100),
    carrier          VARCHAR(100)
);
