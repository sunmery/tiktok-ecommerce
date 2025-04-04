-- 创建分片表前先创建 schema
CREATE SCHEMA IF NOT EXISTS merchant;
SET SEARCH_PATH TO merchant,products,orders;

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

CREATE TABLE products.products
(
    id               UUID                  DEFAULT uuidv7_sub_ms(),
    merchant_id      UUID         NOT NULL, -- 分片键（必须）
    name             VARCHAR(255) NOT NULL,
    description      TEXT,
    price            NUMERIC(15, 2) CHECK (price >= 0),
    status           SMALLINT     NOT NULL DEFAULT 1,
    current_audit_id UUID,
    category_id      int8         NOT NULL, -- 商品分类 ID
    created_at       TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at       TIMESTAMPTZ,
    PRIMARY KEY (merchant_id, id)           -- Citus 分片表需要包含分片键在PK中
);


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
    email          VARCHAR(320)              NOT NULL, -- 支持最大邮箱长度
    created_at     timestamptz DEFAULT now() NOT NULL, -- Unix时间戳，避免时区问题
    updated_at     timestamptz DEFAULT now() NOT NULL
);

-- 创建子订单表（按商家分单）
CREATE TABLE orders.sub_orders
(
    id           SERIAL PRIMARY KEY,                 -- 子订单
    order_id     BIGINT                    NOT NULL, -- 关联主订单ID（程序级外键）
    merchant_id  UUID                      NOT NULL, -- 商家ID（来自商家服务）
    total_amount NUMERIC(12, 2)            NOT NULL, -- 精确金额计算（整数部分10位，小数2位）
    currency     VARCHAR(3)                NOT NULL, -- 实际结算货币
    status       VARCHAR(20)               NOT NULL, -- 订单状态：
    items        JSONB                     NOT NULL, -- 订单项快照（包含商品详情和当时价格）
    created_at   timestamptz DEFAULT now() NOT NULL,
    updated_at   timestamptz DEFAULT now() NOT NULL
);
-- 主订单表
ALTER TABLE orders.orders
    ADD COLUMN payment_status VARCHAR(20) NOT NULL DEFAULT 'pending'
        CHECK (payment_status IN ('pending', 'paid', 'cancelled', 'failed', 'cancelled'));

-- 子订单表
ALTER TABLE orders.sub_orders
    ADD COLUMN payment_status VARCHAR(20) NOT NULL DEFAULT 'pending';
