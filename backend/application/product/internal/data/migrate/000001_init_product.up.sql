-- 创建分片表前先创建 schema
CREATE SCHEMA IF NOT EXISTS products;

-----------------------------
-- 商品主表（分布式分片表）
-----------------------------
CREATE TABLE products.products (
    id BIGSERIAL,                     -- 分布式环境下建议使用 snowflake 等分布式ID
    merchant_id BIGINT NOT NULL,       -- 分片键（必须）
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price NUMERIC(15,2) CHECK (price >= 0),
    stock INT CHECK (stock >= 0),
    status SMALLINT NOT NULL DEFAULT 1,
    current_audit_id BIGINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    PRIMARY KEY (merchant_id, id)     -- Citus 分片表需要包含分片键在PK中
);

-- 配置分布式表
-- SELECT create_distributed_table('products.products', 'merchant_id');

-----------------------------
-- 商品图片表（共置分片表）
-----------------------------
CREATE TABLE products.product_images (
    id BIGSERIAL,
    merchant_id BIGINT NOT NULL,       -- 分片键（必须）
    product_id BIGINT NOT NULL,
    url VARCHAR(512) NOT NULL,
    is_primary BOOLEAN NOT NULL DEFAULT false,
    sort_order SMALLINT DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (merchant_id, id)
);

-- 创建共置分片表
-- SELECT create_distributed_table('products.product_images', 'merchant_id',
--     colocate_with => 'products');

-- 唯一约束需要包含分片键
CREATE UNIQUE INDEX idx_unique_primary_image
ON products.product_images(merchant_id, product_id, is_primary)
WHERE is_primary = true;

-----------------------------
-- 商品属性表（共置分片表）
-----------------------------
CREATE TABLE products.product_attributes (
    merchant_id BIGINT NOT NULL,       -- 分片键（必须）
    product_id BIGINT NOT NULL,
    attributes JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (merchant_id, product_id)  -- 联合主键
);

-- SELECT create_distributed_table('products.product_attributes', 'merchant_id',
--     colocate_with => 'products');

-----------------------------
-- 商品审核记录表（共置分片表）
-----------------------------
CREATE TABLE products.product_audits (
    id BIGSERIAL,
    merchant_id BIGINT NOT NULL,       -- 分片键（必须）
    product_id BIGINT NOT NULL,
    old_status SMALLINT NOT NULL,
    new_status SMALLINT NOT NULL,
    reason TEXT,
    operator_id BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (merchant_id, id)
);

-- SELECT create_distributed_table('products.product_audits', 'merchant_id',
--     colocate_with => 'products');

-- 创建分片兼容索引
CREATE INDEX idx_audits_product ON products.product_audits
USING BTREE (merchant_id, product_id, created_at DESC);