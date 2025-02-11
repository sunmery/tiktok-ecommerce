-- 手动在数据库中创建
CREATE SCHEMA IF NOT EXISTS products;

-- set search_path TO products;

-- 创建商品主表
CREATE TABLE IF NOT EXISTS products.products
(
    id                  SERIAL PRIMARY KEY,                                                -- 商品唯一标识（自增主键）

    -- 基础信息
    name                VARCHAR(100) NOT NULL CHECK (LENGTH(name) >= 2),                   -- 商品名称（至少2个字符）
    description         TEXT         NOT NULL CHECK (LENGTH(description) >= 10),           -- 商品描述（至少10个字符）
    picture             TEXT         NOT NULL CHECK (picture ~ '^https?://'),              -- 商品图片URL（必须为HTTP/HTTPS链接）

    -- 价格信息
    price               REAL         NOT NULL CHECK (price >= 0),                          -- 商品价格（精确到小数点后2位，非负数）

    -- 库存管理
    total_stock         INT          NOT NULL DEFAULT 0 CHECK (total_stock >= 0),          -- 总库存（物理库存，非负数）
    available_stock     INT GENERATED ALWAYS AS (total_stock - reserved_stock) STORED,     -- 可用库存（计算字段 = 总库存 - 预留库存）
    reserved_stock      INT          NOT NULL DEFAULT 0 CHECK (reserved_stock >= 0),       -- 预留库存（已下单未支付的库存，非负数）
    low_stock_threshold INT          NOT NULL DEFAULT 10 CHECK (low_stock_threshold >= 0), -- 低库存预警阈值（非负数）
    allow_negative      BOOLEAN      NOT NULL DEFAULT false,                               -- 是否允许超卖（true时可用库存可为负数）

    -- 审计字段
    created_at          TIMESTAMPTZ  NOT NULL DEFAULT NOW(),                               -- 创建时间（自动记录）
    updated_at          TIMESTAMPTZ  NOT NULL DEFAULT NOW(),                               -- 更新时间（自动记录）
    version             INT          NOT NULL DEFAULT 1                                    -- 乐观锁版本号（用于并发控制）
);

-- 创建分类表
CREATE TABLE IF NOT EXISTS products.categories
(
    id         SERIAL PRIMARY KEY,                       -- 分类唯一标识（自增主键）
    name       VARCHAR(50) UNIQUE NOT NULL,              -- 分类名称（唯一，非空）
    parent_id  INT REFERENCES products.categories (id),  -- 父分类ID（支持分类层级，外键引用自身）
    is_active  BOOLEAN            NOT NULL DEFAULT true, -- 分类是否启用（默认启用）
    created_at TIMESTAMPTZ        NOT NULL DEFAULT NOW() -- 创建时间（自动记录）
);

-- 创建商品-分类关联表（多对多关系）
CREATE TABLE IF NOT EXISTS products.product_categories
(
    product_id  INT REFERENCES products.products (id) ON DELETE CASCADE,   -- 商品ID（外键，级联删除）
    category_id INT REFERENCES products.categories (id) ON DELETE CASCADE, -- 分类ID（外键，级联删除）
    PRIMARY KEY (product_id, category_id)                                  -- 联合主键（确保唯一性）
);

-- 创建库存变更记录表（用于审计）
CREATE TABLE IF NOT EXISTS products.inventory_history
(
    id            BIGSERIAL PRIMARY KEY,                                            -- 记录唯一标识（自增主键）
    product_id    INT         NOT NULL REFERENCES products.products (id),           -- 商品ID（外键）
    old_stock     INT         NOT NULL,                                             -- 变更前库存
    new_stock     INT         NOT NULL CHECK (new_stock >= 0),                      -- 变更后库存（非负数）
    change_reason VARCHAR(20) NOT NULL CHECK (change_reason IN ('PURCHASE', 'ADJUSTMENT', 'RETURN', 'ORDER_RESERVED',
                                                                'ORDER_RELEASED')), -- 变更原因（枚举值）
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()                                -- 创建时间（自动记录）
);
