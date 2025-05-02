-- 创建分片表前先创建 schema
CREATE SCHEMA IF NOT EXISTS merchant;
SET SEARCH_PATH TO merchant;

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

-----------------------------
-- 库存警报表（存储库存警报阈值配置）
-----------------------------
CREATE TABLE merchant.stock_alerts
(
    id          UUID                 DEFAULT uuidv7_sub_ms() PRIMARY KEY,
    product_id  UUID        NOT NULL,                       -- 产品ID
    merchant_id UUID        NOT NULL,                       -- 商家ID
    threshold   INT         NOT NULL CHECK (threshold > 0), -- 警报阈值（必须大于0）
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),         -- 创建时间
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),         -- 更新时间
    UNIQUE (product_id, merchant_id)                        -- 每个产品每个商家只能有一个警报配置
);

-----------------------------
-- 库存调整记录表（记录库存变更历史）
-----------------------------
CREATE TABLE merchant.stock_adjustments
(
    id          UUID                 DEFAULT uuidv7_sub_ms() PRIMARY KEY,
    product_id  UUID        NOT NULL,              -- 产品ID
    merchant_id UUID        NOT NULL,              -- 商家ID
    quantity    INT         NOT NULL,              -- 调整数量（正数增加，负数减少）
    reason      TEXT        NULL,                  -- 调整原因(选填)
    operator_id UUID        NOT NULL,              -- 操作人ID
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW() -- 创建时间
);


-- 商家地址表设计
CREATE TYPE merchant.address_type AS ENUM ('WAREHOUSE', 'RETURN', 'STORE', 'BILLING');
CREATE TABLE merchants.addresses
(
    id             BIGINT PRIMARY KEY,
    merchant_id    UUID                  NOT NULL,
    address_type   merchant.address_type NOT NULL DEFAULT 'WAREHOUSE', -- 地址类型扩展
    contact_person VARCHAR(255)          NOT NULL,                     -- 联系人姓名
    contact_phone  VARCHAR(20)           NOT NULL,                     -- 联系电话
    street_address TEXT                  NOT NULL,
    city           VARCHAR(255)          NOT NULL,
    state          VARCHAR(20)           NOT NULL,
    country        VARCHAR(100)          NOT NULL,
    zip_code       VARCHAR(20)           NOT NULL,
    is_default     BOOLEAN               NOT NULL DEFAULT FALSE,
    created_at     TIMESTAMPTZ           NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ           NOT NULL DEFAULT NOW()
);
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
    payment_status orders.payment_status     NOT NULL DEFAULT 'PENDING',
    created_at     timestamptz DEFAULT now() NOT NULL, -- Unix时间戳，避免时区问题
    updated_at     timestamptz DEFAULT now() NOT NULL
);

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

-- 创建物流表
CREATE TABLE IF NOT EXISTS orders.shipping_info
(
    id               BIGINT PRIMARY KEY,
    sub_order_id     BIGINT                 NOT NULL,              -- 关联子订单ID
    merchant_id      UUID                   NOT NULL,              -- 商家 id
    tracking_number  VARCHAR(100)           NOT NULL,              -- 物流单号
    carrier          VARCHAR(100)           NOT NULL,              -- 承运商
    shipping_status  orders.shipping_status NOT NULL DEFAULT 'WAIT_COMMAND',
    delivery         TIMESTAMPTZ            NOT NULL,              -- 送达时间
    shipping_address JSONB                  NOT NULL,              -- 发货地址信息
    receiver_address JSONB                  NOT NULL,              -- 收货地址信息
    shipping_fee     NUMERIC(10, 2)         NOT NULL DEFAULT 0.00, -- 运费
    created_at       TIMESTAMPTZ            NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ            NOT NULL DEFAULT NOW()
);
