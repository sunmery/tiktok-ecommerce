-----------------------------
-- 商品主表（分区表基准表）
-----------------------------
CREATE TABLE products.products (
    id BIGSERIAL PRIMARY KEY,         -- 使用bigserial防止溢出
    name VARCHAR(255) NOT NULL,       -- 商品名称
    description TEXT,                 -- 商品描述
    price NUMERIC(15,2) CHECK (price >= 0), -- 精确金额计算

    status SMALLINT NOT NULL DEFAULT 1, -- 状态: 1=草稿 2=待审 3=已上架 4=已驳回
    merchant_id BIGINT NOT NULL,       -- 商家ID（业务关联）
    current_audit_id BIGINT,           -- 当前审核记录ID（逻辑外键）
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ             -- 软删除标记
);

COMMENT ON COLUMN products.products.status IS '1: Draft, 2: Pending, 3: Approved, 4: Rejected';
COMMENT ON COLUMN products.products.current_audit_id IS '最新审核记录ID';

-----------------------------
-- 商品图片表（无主图标记）
-----------------------------
CREATE TABLE products.product_images (
    id BIGSERIAL PRIMARY KEY,
    product_id BIGINT NOT NULL,        -- 逻辑关联products.id
    url VARCHAR(512) NOT NULL,         -- 图片存储地址
    is_primary BOOLEAN NOT NULL DEFAULT false, -- 是否主图
    sort_order SMALLINT DEFAULT 0,     -- 排序字段
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-----------------------------
-- 商品属性表（JSONB方案）
-----------------------------
CREATE TABLE products.product_attributes (
    product_id BIGINT NOT NULL,        -- 逻辑关联products.id
    attributes JSONB NOT NULL,         -- 存储结构: {"color": "red", "sizes": ["S", "M"]}
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (product_id)          -- 一个商品对应一组属性
);

COMMENT ON TABLE products.product_attributes IS '使用JSONB存储灵活属性结构';
-----------------------------
-- 商品审核记录表（审计日志）
-----------------------------
CREATE TABLE products.product_audits (
    id BIGSERIAL PRIMARY KEY,
    product_id BIGINT NOT NULL,        -- 逻辑关联products.id
    old_status SMALLINT NOT NULL,      -- 操作前状态
    new_status SMALLINT NOT NULL,      -- 操作后状态
    reason TEXT,                       -- 驳回原因
    operator_id BIGINT NOT NULL,       -- 操作人ID
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
