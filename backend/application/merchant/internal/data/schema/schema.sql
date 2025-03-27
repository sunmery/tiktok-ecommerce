-- 创建分片表前先创建 schema
CREATE SCHEMA IF NOT EXISTS merchant;
SET SEARCH_PATH TO merchant,products;

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
