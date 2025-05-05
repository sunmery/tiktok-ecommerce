CREATE SCHEMA IF NOT EXISTS orders;
SET SEARCH_PATH TO orders;

-- 创建主订单表（用于汇总用户的所有子订单）
-- 支付状态: 等待支付, 完成支付, 取消支付, 支付异常
-- CREATE TYPE orders.payment_status AS ENUM ('PENDING', 'PAID', 'CANCELLED', 'FAILED');
CREATE TABLE orders.orders
(
    id             BIGINT PRIMARY KEY,
    user_id        UUID                      NOT NULL, -- 关联用户ID
    currency       VARCHAR(3)  DEFAULT 'CNY' NOT NULL, -- 用户下单时使用的货币类型（ISO 4217）
    street_address TEXT                      NOT NULL, -- 反范式化存储地址信息，避免关联查询
    city           VARCHAR(100)              NOT NULL,
    state          VARCHAR(100)              NOT NULL,
    country        VARCHAR(100)              NOT NULL,
    zip_code       VARCHAR(10)               NOT NULL, -- 支付状态
    email          VARCHAR(320)              NOT NULL, -- 支持最大邮箱长度
    payment_status VARCHAR(15)     NOT NULL DEFAULT 'PENDING' CHECK (
        payment_status IN ('PENDING', 'PAID', 'CANCELLED', 'FAILED')),
    created_at     timestamptz DEFAULT now() NOT NULL, -- Unix时间戳，避免时区问题
    updated_at     timestamptz DEFAULT now() NOT NULL
);
COMMENT
    ON TABLE orders.orders IS '主订单表，记录订单汇总信息';

-- 创建子订单表（按商家分单）
-- 货运状态: 等待操作, 等待发货, 发已货, 运输中, 已送达, 确认收货, 取消发货
-- CREATE TYPE orders.shipping_status AS ENUM ('WAIT_COMMAND','PENDING_SHIPMENT', 'SHIPPED', 'IN_TRANSIT', 'DELIVERED', 'CONFIRMED', 'CANCELLED');
CREATE TABLE orders.sub_orders
(
    id              BIGINT PRIMARY KEY,                 -- 子订单
    order_id        BIGINT                    NOT NULL, -- 关联主订单ID（程序级外键）
    merchant_id     UUID                      NOT NULL, -- 商家ID（来自商家服务）
    total_amount    NUMERIC(12, 2)            NOT NULL, -- 精确金额计算（整数部分10位，小数2位）
    currency        VARCHAR(3)                NOT NULL, -- 实际结算货币
    status          VARCHAR(20)               NOT NULL, -- 订单状态：
    items           JSONB                     NOT NULL, -- 订单项快照（包含商品详情和当时价格）
    shipping_status VARCHAR(15)               NOT NULL DEFAULT 'WAIT_COMMAND'
        CHECK ( shipping_status IN
                ('WAIT_COMMAND', 'PENDING_SHIPMENT', 'SHIPPED', 'IN_TRANSIT', 'DELIVERED', 'CONFIRMED',
                 'CANCELLED')),                         -- 物流状态
    created_at      timestamptz DEFAULT now() NOT NULL,
    updated_at      timestamptz DEFAULT now() NOT NULL
);
COMMENT
    ON TABLE orders.sub_orders IS '子订单表，按商家分单存储';

-- 创建物流表
CREATE TABLE orders.shipping_info
(
    id               BIGINT PRIMARY KEY,
    sub_order_id     BIGINT         NOT NULL,              -- 关联子订单ID
    merchant_id      UUID           NOT NULL,              -- 商家 id
    tracking_number  VARCHAR(100)   NOT NULL,              -- 物流单号
    carrier          VARCHAR(100)   NOT NULL,              -- 承运商
    shipping_status  VARCHAR(15)    NOT NULL DEFAULT 'WAIT_COMMAND'
        CHECK ( shipping_status IN
                ('WAIT_COMMAND', 'PENDING_SHIPMENT', 'SHIPPED', 'IN_TRANSIT', 'DELIVERED', 'CONFIRMED',
                 'CANCELLED')),                            -- 物流状态
    delivery         TIMESTAMPTZ    NULL,                  -- 送达时间
    shipping_address JSONB          NOT NULL DEFAULT '{}', -- 发货地址信息
    receiver_address JSONB          NOT NULL DEFAULT '{}', -- 收货地址信息
    shipping_fee     NUMERIC(10, 2) NOT NULL DEFAULT 0.00, -- 运费
    created_at       TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ    NOT NULL DEFAULT NOW()
);
