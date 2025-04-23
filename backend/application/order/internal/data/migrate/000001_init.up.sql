CREATE SCHEMA IF NOT EXISTS orders;
SET SEARCH_PATH TO orders;

-- 创建主订单表（用于汇总用户的所有子订单）
-- 支付状态: 等待支付, 完成支付, 取消支付, 支付异常
CREATE TYPE orders.payment_status AS ENUM ('PENDING', 'PAID', 'CANCELLED', 'FAILED');
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
    payment_status payment_status            NOT NULL DEFAULT 'PENDING',
    created_at     timestamptz DEFAULT now() NOT NULL, -- Unix时间戳，避免时区问题
    updated_at     timestamptz DEFAULT now() NOT NULL
);
COMMENT
    ON TABLE orders.orders IS '主订单表，记录订单汇总信息';

-- 创建子订单表（按商家分单）
-- 货运状态: 等待操作, 等待发货, 发已货, 运输中, 已送达, 确认收货, 取消发货
CREATE TYPE orders.shipping_status AS ENUM ('WAIT_COMMAND','PENDING_SHIPMENT', 'SHIPPED', 'IN_TRANSIT', 'DELIVERED', 'CONFIRMED', 'CANCELLED');
CREATE TABLE orders.sub_orders
(
    id              BIGINT PRIMARY KEY,                                        -- 子订单
    order_id        BIGINT                    NOT NULL,                        -- 关联主订单ID（程序级外键）
    merchant_id     UUID                      NOT NULL,                        -- 商家ID（来自商家服务）
    total_amount    NUMERIC(12, 2)            NOT NULL,                        -- 精确金额计算（整数部分10位，小数2位）
    currency        VARCHAR(3)                NOT NULL,                        -- 实际结算货币
    status          VARCHAR(20)               NOT NULL,                        -- 订单状态：
    items           JSONB                     NOT NULL,                        -- 订单项快照（包含商品详情和当时价格）
    shipping_status orders.shipping_status    NOT NULL DEFAULT 'WAIT_COMMAND', -- 物流状态
    created_at      timestamptz DEFAULT now() NOT NULL,
    updated_at      timestamptz DEFAULT now() NOT NULL
);
COMMENT
    ON TABLE orders.sub_orders IS '子订单表，按商家分单存储';

-- 创建索引优化查询性能
CREATE INDEX idx_orders_user ON orders.orders (user_id, created_at DESC); -- 用户订单列表查询
CREATE INDEX idx_sub_orders_merchant ON orders.sub_orders (merchant_id, status);
