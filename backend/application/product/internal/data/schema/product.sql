-- 创建分片表前先创建 schema
CREATE SCHEMA IF NOT EXISTS products;
SET SEARCH_PATH TO products;

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
-- 商品主表（分布式分片表）
-----------------------------
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

-----------------------------
-- 库存表, 按商家区分的商品库存表
-----------------------------
CREATE TABLE products.inventory
(
    product_id  UUID NOT NULL,                    -- 商品ID（关联商品表）
    merchant_id UUID NOT NULL,                    -- 商家ID
    stock       INT  NOT NULL CHECK (stock >= 0), -- 当前库存（不允许负数）
    PRIMARY KEY (product_id, merchant_id)         -- 联合主键（商品+商家唯一）
);

-- 配置分布式表
-- SELECT create_distributed_table('products.products', 'merchant_id');

-----------------------------
-- 商品图片表（共置分片表）
-----------------------------
CREATE TABLE products.product_images
(
    id          UUID                  DEFAULT uuidv7_sub_ms(),
    merchant_id UUID         NOT NULL, -- 分片键（必须）
    product_id  UUID         NOT NULL,
    url         VARCHAR(512) NOT NULL,
    is_primary  BOOLEAN      NOT NULL DEFAULT false,
    sort_order  SMALLINT              DEFAULT 0,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    PRIMARY KEY (merchant_id, id)
);

-----------------------------
-- 商品属性表（共置分片表）
-----------------------------
CREATE TABLE products.product_attributes
(
    merchant_id UUID        NOT NULL,     -- 分片键（必须）
    product_id  UUID        NOT NULL,
    attributes  JSONB       NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (merchant_id, product_id) -- 联合主键
);

-- SELECT create_distributed_table('products.product_attributes', 'merchant_id',
--     colocate_with => 'products');

-----------------------------
-- 商品审核记录表（共置分片表）
-----------------------------
CREATE TABLE products.product_audits
(
    id          UUID                 DEFAULT uuidv7_sub_ms(),
    merchant_id UUID        NOT NULL, -- 分片键（必须）
    product_id  UUID        NOT NULL,
    old_status  SMALLINT    NOT NULL,
    new_status  SMALLINT    NOT NULL,
    reason      TEXT,
    operator_id UUID        NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (merchant_id, id)
);
