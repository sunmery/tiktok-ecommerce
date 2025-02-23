-- 创建分类表
CREATE TABLE categories.categories
(
    id         UUID PRIMARY KEY      DEFAULT gen_random_uuid(),
    name       VARCHAR(100) NOT NULL CHECK (length(name) >= 2),
    -- 分类层级, 最大三层, 例如: 电子产品 - 手机 - 安卓手机
    level      INT          NOT NULL CHECK (level >= 1 AND level <= 3),
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    parent_id UUID REFERENCES categories.categories(id),

    -- 唯一性约束保证分类名称不重复
    CONSTRAINT uniq_category_parent UNIQUE (name, parent_id)
);

-- 闭包表
CREATE TABLE categories.category_closure
(
    ancestor   UUID NOT NULL REFERENCES categories.categories (id),
    descendant UUID NOT NULL REFERENCES categories.categories (id),
    depth      INT  NOT NULL CHECK (depth >= 0),
    PRIMARY KEY (ancestor, descendant)
);
