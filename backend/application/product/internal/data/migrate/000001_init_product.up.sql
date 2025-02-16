CREATE SCHEMA IF NOT EXISTS products;

-- 创建商品表
CREATE TABLE products.products
(
    -- 使用UUID作为主键，防止ID猜测
    id          UUID PRIMARY KEY      DEFAULT gen_random_uuid(),

    -- 商品基本信息
    name        TEXT         NOT NULL CHECK (length(name) <= 200),  -- 商品名称，限制200字符
    description VARCHAR(200) NOT NULL,                              -- 商品描述
    picture     VARCHAR(255) NOT NULL,                              -- 商品图片
    price       REAL         NOT NULL CHECK (price > 0),            -- 价格

    -- 库存管理
    stock       INT          NOT NULL DEFAULT 0 CHECK (stock >= 0), -- 可用库存

    -- 分类
    category_id UUID         NOT NULL,

    -- 时间戳
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),                -- 创建时间
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()                 -- 更新时间
);

-- 创建索引优化常用查询
CREATE INDEX idx_products_created_at ON products.products (created_at);
CREATE INDEX idx_products_price ON products.products (price);

-- 自动更新updated_at的触发器
CREATE OR REPLACE FUNCTION update_modified_column()
    RETURNS TRIGGER AS
$$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_products_updated_at
    BEFORE UPDATE
    ON products.products
    FOR EACH ROW
EXECUTE FUNCTION update_modified_column();
