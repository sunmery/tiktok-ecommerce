CREATE SCHEMA IF NOT EXISTS categories;

-- 创建分类表
CREATE TABLE categories.categories
(
    id         UUID PRIMARY KEY      DEFAULT gen_random_uuid(),
    name       VARCHAR(100) NOT NULL CHECK (length(name) >= 2),
    -- 分类层级, 最大三层, 例如: 电子产品 - 手机 - 安卓手机
    level      INT          NOT NULL CHECK (level >= 1 AND level <= 3),
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),

    -- 唯一性约束保证分类名称不重复
    CONSTRAINT uniq_category_name UNIQUE (name)
);

-- 闭包表
CREATE TABLE categories.category_closure
(
    ancestor   UUID NOT NULL REFERENCES categories.categories (id),
    descendant UUID NOT NULL REFERENCES categories.categories (id),
    depth      INT  NOT NULL CHECK (depth >= 0),
    PRIMARY KEY (ancestor, descendant)
);

-- 将分类表设为 Citus 参考表（所有节点全量存储）
-- SELECT create_reference_table('categories');

-- 创建商品表分布式策略（按 category_id 分片）
-- SELECT create_distributed_table('products', 'category_id');

-- 索引优化
CREATE INDEX idx_category_closure_ancestor ON categories.category_closure(ancestor);
CREATE INDEX idx_category_closure_descendant ON categories.category_closure(descendant);

-- 更新分类表触发器（与商品表相同）
-- CREATE TRIGGER update_categories_updated_at
--     BEFORE UPDATE
--     ON categories
--     FOR EACH ROW
-- EXECUTE FUNCTION update_modified_column();
