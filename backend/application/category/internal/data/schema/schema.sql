-- 创建 Schema 并设置搜索路径
CREATE SCHEMA IF NOT EXISTS categories;
SET search_path TO categories, public;

-- 启用 ltree 扩展（需超级用户权限）
CREATE EXTENSION IF NOT EXISTS ltree;

-- 核心分类表（树形结构核心）
CREATE TABLE categories.categories
(
    id         UUID PRIMARY KEY,                   -- 分类ID（UUID生成）
    parent_id  UUID,                               -- 父分类ID（NULL表示根节点）
    level      SMALLINT    NOT NULL DEFAULT 1
        CHECK (level BETWEEN 1 AND 4),             -- 层级深度（限制四级）
    path       LTREE       NOT NULL,               -- 层级路径（使用ltree类型）
    name       VARCHAR(50) NOT NULL,               -- 分类名称
    sort_order SMALLINT    NOT NULL DEFAULT 0
        CHECK (sort_order >= 0),                   -- 同级排序（0-32767）
    is_leaf    BOOLEAN     NOT NULL DEFAULT FALSE, -- 是否为叶子节点
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 创建索引
CREATE INDEX idx_categories_path_gist ON categories.categories USING GIST (path);
CREATE INDEX idx_categories_parent ON categories.categories (parent_id);
CREATE INDEX idx_categories_leaf ON categories.categories (is_leaf) WHERE is_leaf;

COMMENT ON TABLE categories.categories IS '商品分类主表（ltree路径+闭包表双重优化）';

-- 闭包关系表（存储所有层级关系）
CREATE TABLE categories.category_closure
(
    ancestor   UUID     NOT NULL REFERENCES categories.categories (id), -- 祖先节点ID
    descendant UUID     NOT NULL REFERENCES categories.categories (id), -- 后代节点ID
    depth      SMALLINT NOT NULL CHECK (depth BETWEEN 0 AND 3),         -- 层级间隔
    PRIMARY KEY (ancestor, descendant)
);

-- 闭包表索引
CREATE INDEX idx_closure_descendant ON categories.category_closure (descendant);

-- 初始化根节点示例（可选）
INSERT INTO categories.categories (id, parent_id, path, name, level)
VALUES ('00000000-0000-0000-0000-000000000000', -- 根节点特殊UUID
        NULL,
        'root',
        'Root Category',
        1);
